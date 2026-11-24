package model

import "time"

// CertStatus 证书状态
const (
	CertStatusActive   = "active"
	CertStatusRenewing = "renewing"
	CertStatusFailed   = "failed"
)

// DeployTargetType 部署目标类型
const (
	DeployCDN   = "cdn"
	DeploySLB   = "slb"
	DeployNginx = "nginx"
)

// DeployStatus 部署状态
const (
	DeployStatusSuccess = "success"
	DeployStatusFailed  = "failed"
	DeployStatusPending = "pending"
)

// CertMeta 证书元数据（不含私钥，私钥交由 CAS 托管）
type CertMeta struct {
	ID            int64     `json:"id" db:"id"`
	Domain        string    `json:"domain" db:"domain"`
	Domains       string    `json:"domains" db:"domains"` // JSON 数组: ["example.com","*.example.com"]
	CASCertID     string    `json:"cas_cert_id" db:"cas_cert_id"`
	NotBefore     time.Time `json:"not_before" db:"not_before"`
	NotAfter      time.Time `json:"not_after" db:"not_after"`
	ARIRetryAfter time.Time `json:"ari_retry_after" db:"ari_retry_after"`
	DNSProvider   string    `json:"dns_provider" db:"dns_provider"`
	KeyType       string    `json:"key_type" db:"key_type"`
	Status        string    `json:"status" db:"status"`
	LastError     string    `json:"last_error" db:"last_error"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// DeployTarget 部署目标
type DeployTarget struct {
	ID         int64     `json:"id" db:"id"`
	CertMetaID int64     `json:"cert_meta_id" db:"cert_meta_id"`
	Type       string    `json:"type" db:"type"`               // cdn / slb / nginx
	ResourceID string    `json:"resource_id" db:"resource_id"` // CDN域名 / SLB实例ID / ECS实例ID
	RegionID   string    `json:"region_id" db:"region_id"`     // ECS 地域
	NginxPath  string    `json:"nginx_path" db:"nginx_path"`   // Nginx 证书目录
	Extra      string    `json:"extra" db:"extra"`             // JSON 扩展配置（如 SLB 监听器 ID）
	Status     string    `json:"status" db:"status"`
	LastError  string    `json:"last_error" db:"last_error"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// DeployLog 部署/续期操作日志
type DeployLog struct {
	ID         int64     `json:"id" db:"id"`
	CertMetaID int64     `json:"cert_meta_id" db:"cert_meta_id"`
	Action     string    `json:"action" db:"action"` // apply / renew / deploy / delete_old
	Target     string    `json:"target" db:"target"`
	Success    bool      `json:"success" db:"success"`
	Message    string    `json:"message" db:"message"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// DNSProviderConfig DNS 服务商配置（脱敏存储）
type DNSProviderConfig struct {
	Name            string `json:"name"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret,omitempty"`
	Region          string `json:"region,omitempty"`
}

// ApplyRequest 申请证书请求
type ApplyRequest struct {
	Domain       string        `json:"domain" binding:"required"`
	Wildcard     bool          `json:"wildcard"`
	DNSProvider  string        `json:"dns_provider"`
	DeployTarget []TargetInput `json:"deploy_targets"`
}

// TargetInput 部署目标输入
type TargetInput struct {
	Type       string `json:"type" binding:"required"`
	ResourceID string `json:"resource_id" binding:"required"`
	RegionID   string `json:"region_id"`
	NginxPath  string `json:"nginx_path"`
	Extra      string `json:"extra"`
}

// CertResponse 证书响应（不含私钥）
type CertResponse struct {
	ID          int64           `json:"id"`
	Domain      string          `json:"domain"`
	Domains     []string        `json:"domains"`
	CASCertID   string          `json:"cas_cert_id"`
	NotBefore   time.Time       `json:"not_before"`
	NotAfter    time.Time       `json:"not_after"`
	DaysLeft    int             `json:"days_left"`
	DNSProvider string          `json:"dns_provider"`
	Status      string          `json:"status"`
	Targets     []*DeployTarget `json:"targets"`
}
