package deploy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"lets/internal/config"
	"lets/internal/model"
)

// SLBDeployer 阿里云 SLB 证书部署
type SLBDeployer struct {
	cfg    *config.Config
	client *sdk.Client
}

// slbListenerExtra SLB 部署扩展配置（存于 DeployTarget.Extra）
type slbListenerExtra struct {
	ListenerIDs []string `json:"listener_ids"` // 需要更新证书的 HTTPS 监听器 ID
}

func NewSLBDeployer(cfg *config.Config) (*SLBDeployer, error) {
	client, err := sdk.NewClientWithAccessKey(cfg.Aliyun.Region, cfg.Aliyun.AccessKeyID, cfg.Aliyun.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("create SLB sdk client: %w", err)
	}
	return &SLBDeployer{cfg: cfg, client: client}, nil
}

func (d *SLBDeployer) Type() string { return model.DeploySLB }

func (d *SLBDeployer) Deploy(ctx context.Context, target *model.DeployTarget, cert *CertPayload) error {
	// 1. 上传服务端证书，获得新的 ServerCertificateId
	certID, err := d.uploadServerCert(cert)
	if err != nil {
		return err
	}

	// 2. 解析需要更新的监听器列表
	var extra slbListenerExtra
	if target.Extra != "" {
		if err = json.Unmarshal([]byte(target.Extra), &extra); err != nil {
			return fmt.Errorf("parse SLB extra config: %w", err)
		}
	}
	if len(extra.ListenerIDs) == 0 {
		// 未指定监听器时，只完成证书上传（需用户在控制台手动绑定）
		return nil
	}

	// 3. 逐个更新 HTTPS 监听器绑定的证书
	for _, listenerID := range extra.ListenerIDs {
		if err = d.updateListener(target.ResourceID, listenerID, certID); err != nil {
			return err
		}
	}
	return nil
}

func (d *SLBDeployer) uploadServerCert(cert *CertPayload) (string, error) {
	req := requests.NewCommonRequest()
	req.Method = "POST"
	req.Scheme = "HTTPS"
	req.Domain = "slb.aliyuncs.com"
	req.Version = "2014-05-15"
	req.ApiName = "SetServerCertificate"
	req.QueryParams["ServerCertificate"] = cert.Certificate
	req.QueryParams["PrivateKey"] = cert.PrivateKey
	req.QueryParams["ServerCertificateName"] = fmt.Sprintf("%s-%d", cert.Domain, unixMilli())
	req.QueryParams["RegionId"] = d.cfg.Aliyun.Region

	resp, err := d.client.ProcessCommonRequest(req)
	if err != nil {
		return "", fmt.Errorf("upload SLB certificate: %w", err)
	}
	if !resp.IsSuccess() {
		return "", fmt.Errorf("upload SLB certificate failed: HTTP %d - %s",
			resp.GetHttpStatus(), resp.GetHttpContentString())
	}
	var result struct {
		ServerCertificateId string `json:"ServerCertificateId"`
	}
	if err = json.Unmarshal(resp.GetHttpContentBytes(), &result); err != nil {
		return "", fmt.Errorf("parse SLB upload response: %w", err)
	}
	return result.ServerCertificateId, nil
}

func (d *SLBDeployer) updateListener(lbID, listenerID, certID string) error {
	req := requests.NewCommonRequest()
	req.Method = "POST"
	req.Scheme = "HTTPS"
	req.Domain = "slb.aliyuncs.com"
	req.Version = "2014-05-15"
	req.ApiName = "SetLoadBalancerHTTPSListenerAttribute"
	req.QueryParams["LoadBalancerId"] = lbID
	req.QueryParams["ListenerPort"] = listenerID
	req.QueryParams["ServerCertificateId"] = certID
	req.QueryParams["RegionId"] = d.cfg.Aliyun.Region

	resp, err := d.client.ProcessCommonRequest(req)
	if err != nil {
		return fmt.Errorf("update SLB listener %s: %w", listenerID, err)
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("update SLB listener %s failed: HTTP %d - %s",
			listenerID, resp.GetHttpStatus(), resp.GetHttpContentString())
	}
	return nil
}
