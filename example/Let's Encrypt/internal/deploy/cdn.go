package deploy

import (
	"context"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"lets/internal/config"
	"lets/internal/model"
)

// CDNDeployer 阿里云 CDN 证书部署
type CDNDeployer struct {
	cfg    *config.Config
	client *sdk.Client
}

func NewCDNDeployer(cfg *config.Config) (*CDNDeployer, error) {
	client, err := sdk.NewClientWithAccessKey(cfg.Aliyun.Region, cfg.Aliyun.AccessKeyID, cfg.Aliyun.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("create CDN sdk client: %w", err)
	}
	return &CDNDeployer{cfg: cfg, client: client}, nil
}

func (d *CDNDeployer) Type() string { return model.DeployCDN }

func (d *CDNDeployer) Deploy(ctx context.Context, target *model.DeployTarget, cert *CertPayload) error {
	// 使用 upload 类型直接传证书内容（避免依赖 CAS 部署次数限制）
	req := requests.NewCommonRequest()
	req.Method = "POST"
	req.Scheme = "HTTPS"
	req.Domain = "cdn.aliyuncs.com"
	req.Version = "2018-05-10"
	req.ApiName = "SetCdnDomainSSLCertificate"
	req.QueryParams["DomainName"] = target.ResourceID
	req.QueryParams["CertName"] = fmt.Sprintf("%s-%d", cert.Domain, unixMilli())
	req.QueryParams["CertType"] = "upload"
	req.QueryParams["SSLPub"] = cert.Certificate
	req.QueryParams["SSLPri"] = cert.PrivateKey
	req.QueryParams["SSLProtocol"] = "on"

	resp, err := d.client.ProcessCommonRequest(req)
	if err != nil {
		return fmt.Errorf("deploy CDN certificate for %s: %w", target.ResourceID, err)
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("deploy CDN certificate for %s failed: HTTP %d - %s",
			target.ResourceID, resp.GetHttpStatus(), resp.GetHttpContentString())
	}
	return nil
}
