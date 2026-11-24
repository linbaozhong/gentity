package deploy

import (
	"context"
	"fmt"

	"lets/internal/config"
	"lets/internal/model"
)

// CertPayload 待部署的证书内容
type CertPayload struct {
	Domain      string
	Certificate string // fullchain PEM
	PrivateKey  string // privkey PEM
	CASCertID   string
}

// Deployer 部署器接口
type Deployer interface {
	// Type 返回部署类型
	Type() string
	// Deploy 部署证书到目标资源
	Deploy(ctx context.Context, target *model.DeployTarget, cert *CertPayload) error
}

// Manager 部署管理器
type Manager struct {
	deployers map[string]Deployer
}

// NewManager 创建部署管理器
func NewManager(cfg *config.Config) (*Manager, error) {
	m := &Manager{deployers: make(map[string]Deployer)}

	cdn, err := NewCDNDeployer(cfg)
	if err != nil {
		return nil, err
	}
	m.deployers[model.DeployCDN] = cdn

	slb, err := NewSLBDeployer(cfg)
	if err != nil {
		return nil, err
	}
	m.deployers[model.DeploySLB] = slb

	nginx, err := NewNginxDeployer(cfg)
	if err != nil {
		return nil, err
	}
	m.deployers[model.DeployNginx] = nginx

	return m, nil
}

// Deploy 根据目标类型分发部署
func (m *Manager) Deploy(ctx context.Context, target *model.DeployTarget, cert *CertPayload) error {
	d, ok := m.deployers[target.Type]
	if !ok {
		return fmt.Errorf("unsupported deploy type: %s", target.Type)
	}
	return d.Deploy(ctx, target, cert)
}
