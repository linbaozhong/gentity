package cas

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"lets/internal/config"
)

const (
	casEndpoint = "cas.aliyuncs.com"
	casVersion  = "2020-04-07"
)

// Client 阿里云数字证书管理服务客户端
type Client struct {
	cfg    *config.Config
	client *sdk.Client
}

// UploadResult 上传结果
type UploadResult struct {
	CertID   string
	CertName string
}

// CertDetail 证书详情
type CertDetail struct {
	CertID     string
	Common     string
	NotAfter   string
	Status     string
	Cert       string
	PrivateKey string
}

// NewClient 创建 CAS 客户端
func NewClient(cfg *config.Config) (*Client, error) {
	client, err := sdk.NewClientWithAccessKey(cfg.Aliyun.Region, cfg.Aliyun.AccessKeyID, cfg.Aliyun.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("create CAS sdk client: %w", err)
	}
	return &Client{cfg: cfg, client: client}, nil
}

// Upload 上传证书（含私钥）到 CAS，返回 CertId
func (c *Client) Upload(name, certPEM, keyPEM string) (*UploadResult, error) {
	req := requests.NewCommonRequest()
	req.Method = "POST"
	req.Scheme = "HTTPS"
	req.Domain = casEndpoint
	req.Version = casVersion
	req.ApiName = "UploadUserCertificate"
	req.QueryParams["CertName"] = name
	req.QueryParams["Cert"] = certPEM
	req.QueryParams["Key"] = keyPEM

	resp, err := c.client.ProcessCommonRequest(req)
	if err != nil {
		return nil, fmt.Errorf("upload certificate to CAS: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("upload certificate to CAS failed: HTTP %d", resp.GetHttpStatus())
	}
	certID, err := jsonGetString(resp.GetHttpContentString(), "CertId")
	if err != nil {
		return nil, err
	}
	return &UploadResult{CertID: certID, CertName: name}, nil
}

// Delete 删除 CAS 中的证书
func (c *Client) Delete(certID string) error {
	req := requests.NewCommonRequest()
	req.Method = "POST"
	req.Scheme = "HTTPS"
	req.Domain = casEndpoint
	req.Version = casVersion
	req.ApiName = "DeleteUserCertificate"
	req.QueryParams["CertId"] = certID

	resp, err := c.client.ProcessCommonRequest(req)
	if err != nil {
		return fmt.Errorf("delete certificate from CAS: %w", err)
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("delete certificate from CAS failed: HTTP %d", resp.GetHttpStatus())
	}
	return nil
}

// GetDetail 查询证书详情（含证书内容）
func (c *Client) GetDetail(certID string) (*CertDetail, error) {
	req := requests.NewCommonRequest()
	req.Method = "POST"
	req.Scheme = "HTTPS"
	req.Domain = casEndpoint
	req.Version = casVersion
	req.ApiName = "GetUserCertificateDetail"
	req.QueryParams["CertId"] = certID

	resp, err := c.client.ProcessCommonRequest(req)
	if err != nil {
		return nil, fmt.Errorf("get certificate detail: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("get certificate detail failed: HTTP %d", resp.GetHttpStatus())
	}
	body := resp.GetHttpContentString()
	d := &CertDetail{CertID: certID}
	d.Common, _ = jsonGetString(body, "CommonName")
	d.NotAfter, _ = jsonGetString(body, "NotAfter")
	d.Status, _ = jsonGetString(body, "Status")
	d.Cert, _ = jsonGetString(body, "Cert")
	d.PrivateKey, _ = jsonGetString(body, "Key")
	return d, nil
}

// ListByDomain 按域名查询已上传证书（用于去重）
func (c *Client) ListByDomain(domain string) (string, error) {
	req := requests.NewCommonRequest()
	req.Method = "POST"
	req.Scheme = "HTTPS"
	req.Domain = casEndpoint
	req.Version = casVersion
	req.ApiName = "ListUserCertificateOrder"
	req.QueryParams["OrderType"] = "UPLOAD"
	req.QueryParams["Keyword"] = strings.TrimPrefix(domain, "*.")
	req.QueryParams["ShowSize"] = "50"
	req.QueryParams["CurrentPage"] = "1"

	resp, err := c.client.ProcessCommonRequest(req)
	if err != nil {
		return "", fmt.Errorf("list CAS certificates: %w", err)
	}
	if !resp.IsSuccess() {
		return "", fmt.Errorf("list CAS certificates failed: HTTP %d", resp.GetHttpStatus())
	}
	return resp.GetHttpContentString(), nil
}

// jsonGetString 简单 JSON 取值（避免引入复杂结构体）
func jsonGetString(body, key string) (string, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		return "", fmt.Errorf("parse CAS response: %w", err)
	}
	data, ok := m[key]
	if !ok || data == nil {
		return "", fmt.Errorf("CAS response missing field %s: %s", key, truncate(body, 200))
	}
	switch v := data.(type) {
	case string:
		return v, nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}
