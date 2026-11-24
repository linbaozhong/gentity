package deploy

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"

	"lets/internal/config"
	"lets/internal/model"
)

// NginxDeployer ECS 上的 Nginx 部署：OSS 私有 bucket 中转 + 云助手远程执行
type NginxDeployer struct {
	cfg       *config.Config
	ecsClient *sdk.Client
	ossClient *oss.Client
}

func NewNginxDeployer(cfg *config.Config) (*NginxDeployer, error) {
	ecsClient, err := sdk.NewClientWithAccessKey(cfg.Aliyun.Region, cfg.Aliyun.AccessKeyID, cfg.Aliyun.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("create ECS sdk client: %w", err)
	}
	ossCfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.Aliyun.AccessKeyID, cfg.Aliyun.AccessKeySecret, "")).
		WithRegion(cfg.Aliyun.Region)
	return &NginxDeployer{
		cfg:       cfg,
		ecsClient: ecsClient,
		ossClient: oss.NewClient(ossCfg),
	}, nil
}

func (d *NginxDeployer) Type() string { return model.DeployNginx }

func (d *NginxDeployer) Deploy(ctx context.Context, target *model.DeployTarget, cert *CertPayload) error {
	regionID := target.RegionID
	if regionID == "" {
		regionID = d.cfg.Aliyun.Region
	}
	nginxPath := target.NginxPath
	if nginxPath == "" {
		nginxPath = fmt.Sprintf("/etc/nginx/ssl/%s", cert.Domain)
	}

	// 1. 生成 OSS 临时路径（每个域名独立，含时间戳）
	prefix := strings.Trim(d.cfg.Aliyun.OSS.PrivatePrefix, "/")
	baseKey := fmt.Sprintf("%s/%s/%d", prefix, cert.Domain, time.Now().UnixMilli())
	certKey := baseKey + "/fullchain.pem"
	keyKey := baseKey + "/privkey.pem"

	// 2. 上传证书链和私钥到私有 bucket
	if err := d.putObject(ctx, certKey, []byte(cert.Certificate)); err != nil {
		return err
	}
	if err := d.putObject(ctx, keyKey, []byte(cert.PrivateKey)); err != nil {
		// 清理已上传的证书
		_ = d.deleteObject(ctx, certKey)
		return err
	}
	// 函数结束无论如何都清理临时文件
	defer func() {
		_ = d.deleteObject(ctx, certKey)
		_ = d.deleteObject(ctx, keyKey)
	}()

	// 3. 通过 ECS 云助手下发部署脚本
	script := d.buildDeployScript(certKey, keyKey, nginxPath, regionID, cert.Domain)
	invID, err := d.runCommand(ctx, target.ResourceID, regionID, script)
	if err != nil {
		return err
	}

	// 4. 轮询执行结果
	return d.waitForInvocation(ctx, regionID, invID, target.ResourceID)
}

func (d *NginxDeployer) putObject(ctx context.Context, key string, data []byte) error {
	_, err := d.ossClient.PutObject(ctx, &oss.PutObjectRequest{
		Bucket:      oss.Ptr(d.cfg.Aliyun.OSS.PrivateBucket),
		Key:         oss.Ptr(key),
		Body:        bytesReader(data),
		Acl:         oss.ObjectACLPrivate,
		ContentType: oss.Ptr("application/octet-stream"),
	})
	if err != nil {
		return fmt.Errorf("upload cert to private OSS (%s): %w", key, err)
	}
	return nil
}

func (d *NginxDeployer) deleteObject(ctx context.Context, key string) error {
	_, err := d.ossClient.DeleteObject(ctx, &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(d.cfg.Aliyun.OSS.PrivateBucket),
		Key:    oss.Ptr(key),
	})
	return err
}

// buildDeployScript 生成在 ECS 上执行的部署脚本
// 注意：ECS 实例需绑定 RAM 角色并授权 oss:GetObject 读取私有 bucket
func (d *NginxDeployer) buildDeployScript(certKey, keyKey, nginxPath, region, domain string) string {
	bucket := d.cfg.Aliyun.OSS.PrivateBucket
	endpoint := fmt.Sprintf("oss-%s.aliyuncs.com", region)
	return fmt.Sprintf(`#!/bin/bash
set -euo pipefail

CERT_KEY="%s"
KEY_KEY="%s"
BUCKET="%s"
ENDPOINT="%s"
NGINX_PATH="%s"
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# 使用 ECS 实例 RAM 角色临时凭证从私有 OSS 下载（无需明文 AK）
TOKEN=$(curl -sf "http://100.100.100.200/latest/meta-data/Ram/security-credentials/")
ROLE_NAME=$(echo "$TOKEN" | head -1)
CREDS=$(curl -sf "http://100.100.100.200/latest/meta-data/Ram/security-credentials/${ROLE_NAME}")
AK=$(echo "$CREDS" | python3 -c 'import sys,json; print(json.load(sys.stdin)["AccessKeyId"])')
SK=$(echo "$CREDS" | python3 -c 'import sys,json; print(json.load(sys.stdin)["AccessKeySecret"])')
STS=$(echo "$CREDS" | python3 -c 'import sys,json; print(json.load(sys.stdin)["SecurityToken"])')

# 使用 ossutil 下载（ECS 需预装 ossutil）
ossutil -e "$ENDPOINT" -i "$AK" -k "$SK" --token "$STS" cp "oss://$BUCKET/$CERT_KEY" "$TMP_DIR/fullchain.pem" -f
ossutil -e "$ENDPOINT" -i "$AK" -k "$SK" --token "$STS" cp "oss://$BUCKET/$KEY_KEY" "$TMP_DIR/privkey.pem" -f

# 校验证书与私钥匹配
CERT_HASH=$(openssl x509 -noout -modulus -in "$TMP_DIR/fullchain.pem" | openssl md5)
KEY_HASH=$(openssl rsa -noout -modulus -in "$TMP_DIR/privkey.pem" 2>/dev/null | openssl md5)
if [ "$CERT_HASH" != "$KEY_HASH" ]; then
  echo "ERROR: certificate and private key mismatch"
  exit 1
fi

# 备份旧证书并原子替换
mkdir -p "$NGINX_PATH"
if [ -f "$NGINX_PATH/fullchain.pem" ]; then
  cp "$NGINX_PATH/fullchain.pem" "$NGINX_PATH/fullchain.pem.bak"
  cp "$NGINX_PATH/privkey.pem" "$NGINX_PATH/privkey.pem.bak"
fi
install -m 0644 "$TMP_DIR/fullchain.pem" "$NGINX_PATH/fullchain.pem"
install -m 0600 "$TMP_DIR/privkey.pem" "$NGINX_PATH/privkey.pem"

# 测试并 reload
if ! nginx -t 2>&1; then
  echo "ERROR: nginx config test failed, rolling back"
  if [ -f "$NGINX_PATH/fullchain.pem.bak" ]; then
    cp "$NGINX_PATH/fullchain.pem.bak" "$NGINX_PATH/fullchain.pem"
    cp "$NGINX_PATH/privkey.pem.bak" "$NGINX_PATH/privkey.pem"
  fi
  exit 1
fi
nginx -s reload || systemctl reload nginx

echo "SUCCESS: certificate deployed for %s"
`, certKey, keyKey, bucket, endpoint, nginxPath, domain)
}

func (d *NginxDeployer) runCommand(ctx context.Context, instanceID, regionID, script string) (string, error) {
	req := requests.NewCommonRequest()
	req.Method = "POST"
	req.Scheme = "HTTPS"
	req.Domain = "ecs.aliyuncs.com"
	req.Version = "2014-05-26"
	req.ApiName = "RunCommand"
	req.QueryParams["InstanceId.1"] = instanceID
	req.QueryParams["Type"] = "RunShellScript"
	req.QueryParams["CommandContent"] = base64.StdEncoding.EncodeToString([]byte(script))
	req.QueryParams["RegionId"] = regionID
	req.QueryParams["Timeout"] = "120"

	resp, err := d.ecsClient.ProcessCommonRequest(req)
	if err != nil {
		return "", fmt.Errorf("run ECS command: %w", err)
	}
	if !resp.IsSuccess() {
		return "", fmt.Errorf("run ECS command failed: HTTP %d - %s",
			resp.GetHttpStatus(), resp.GetHttpContentString())
	}
	var result struct {
		InvocationId string `json:"InvocationId"`
	}
	if err = jsonUnmarshalECS(resp.GetHttpContentBytes(), &result); err != nil {
		return "", err
	}
	return result.InvocationId, nil
}

func (d *NginxDeployer) waitForInvocation(ctx context.Context, regionID, invID, instanceID string) error {
	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		req := requests.NewCommonRequest()
		req.Method = "POST"
		req.Scheme = "HTTPS"
		req.Domain = "ecs.aliyuncs.com"
		req.Version = "2014-05-26"
		req.ApiName = "DescribeInvocationResults"
		req.QueryParams["InvocationId"] = invID
		req.QueryParams["InstanceId"] = instanceID
		req.QueryParams["RegionId"] = regionID

		resp, err := d.ecsClient.ProcessCommonRequest(req)
		if err != nil {
			return fmt.Errorf("describe ECS invocation: %w", err)
		}
		if !resp.IsSuccess() {
			return fmt.Errorf("describe ECS invocation failed: HTTP %d", resp.GetHttpStatus())
		}
		var result struct {
			InvocationResults struct {
				InvocationResult []struct {
					Status   string `json:"Status"`
					Output   string `json:"Output"`
					ExitCode int    `json:"ExitCode"`
				} `json:"InvocationResult"`
			} `json:"InvocationResults"`
		}
		if err = jsonUnmarshalECS(resp.GetHttpContentBytes(), &result); err != nil {
			return err
		}
		results := result.InvocationResults.InvocationResult
		if len(results) > 0 {
			r := results[0]
			switch r.Status {
			case "Success":
				if r.ExitCode != 0 || strings.Contains(r.Output, "ERROR") {
					return fmt.Errorf("nginx deploy failed: exit=%d output=%s", r.ExitCode, r.Output)
				}
				return nil
			case "Failed", "Cancelled":
				return fmt.Errorf("nginx deploy invocation %s: %s output=%s", r.Status, invID, r.Output)
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(3 * time.Second):
		}
	}
	return fmt.Errorf("timeout waiting for ECS invocation %s", invID)
}
