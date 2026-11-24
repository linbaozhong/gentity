package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"lets/internal/acme"
	"lets/internal/cas"
	"lets/internal/config"
	"lets/internal/deploy"
	"lets/internal/dns"
	"lets/internal/model"
	"lets/internal/storage"
)

var (
	ErrCertExists   = errors.New("certificate already exists")
	ErrCertNotFound = errors.New("certificate not found")
)

// CertService 证书服务：编排签发、托管、部署、续期全流程
type CertService struct {
	cfg        *config.Config
	store      *storage.Store
	acmeClient *acme.Client
	casClient  *cas.Client
	dnsFactory *dns.Factory
	deployer   *deploy.Manager

	// 防止同一域名并发申请/续期
	locks sync.Map // domain → *sync.Mutex
}

func NewCertService(
	cfg *config.Config,
	store *storage.Store,
	acmeClient *acme.Client,
	casClient *cas.Client,
	dnsFactory *dns.Factory,
	deployer *deploy.Manager,
) *CertService {
	return &CertService{
		cfg:        cfg,
		store:      store,
		acmeClient: acmeClient,
		casClient:  casClient,
		dnsFactory: dnsFactory,
		deployer:   deployer,
	}
}

// Apply 申请新证书并托管部署
func (s *CertService) Apply(ctx context.Context, req *model.ApplyRequest) (*model.CertResponse, error) {
	domain := normalizeDomain(req.Domain)
	if domain == "" || strings.ContainsAny(domain, " /\\") {
		return nil, fmt.Errorf("invalid domain: %s", req.Domain)
	}

	// 检查是否已存在
	if _, err := s.store.GetCertMetaByDomain(domain); err == nil {
		return nil, ErrCertExists
	}

	// 构建域名列表（通配符：同时包含根域和 *.根域）
	domains := []string{domain}
	if req.Wildcard {
		domains = append(domains, "*."+domain)
	}

	// 获取并锁定域名
	mu := s.lockDomain(domain)
	mu.Lock()
	defer mu.Unlock()

	// 选择 DNS 服务商
	providerName, provider, err := s.dnsFactory.ForDomain(domain)
	if err != nil {
		return nil, fmt.Errorf("select DNS provider: %w", err)
	}
	if req.DNSProvider != "" {
		providerName = req.DNSProvider
		provider, err = s.dnsFactory.Get(req.DNSProvider)
		if err != nil {
			return nil, err
		}
	}
	if err = s.acmeClient.SetDNSProvider(provider); err != nil {
		return nil, fmt.Errorf("set DNS provider: %w", err)
	}

	// 签发证书
	slog.Info("obtaining certificate", "domains", domains, "dns", providerName)
	result, err := s.acmeClient.Obtain(domains)
	if err != nil {
		return nil, err
	}

	// 上传到 CAS 托管
	upload, err := s.casClient.Upload(domain, result.Certificate, result.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("upload to CAS: %w", err)
	}

	// 保存元数据
	meta := &model.CertMeta{
		Domain:      domain,
		Domains:     storage.MarshalDomains(domains),
		CASCertID:   upload.CertID,
		NotBefore:   result.NotBefore,
		NotAfter:    result.NotAfter,
		DNSProvider: providerName,
		KeyType:     s.cfg.ACME.KeyType,
		Status:      model.CertStatusActive,
	}
	if err = s.store.SaveCertMeta(meta); err != nil {
		// 回滚：删除 CAS 证书
		_ = s.casClient.Delete(upload.CertID)
		return nil, fmt.Errorf("save cert meta: %w", err)
	}

	// 保存部署目标并执行部署
	certPayload := &deploy.CertPayload{
		Domain:      domain,
		Certificate: result.Certificate,
		PrivateKey:  result.PrivateKey,
		CASCertID:   upload.CertID,
	}
	saved, err := s.saveAndDeploy(ctx, meta.ID, req.DeployTarget, certPayload)
	if err != nil {
		_ = s.store.AddDeployLog(&model.DeployLog{
			CertMetaID: meta.ID, Action: "apply", Success: false, Message: err.Error(),
		})
		return nil, err
	}

	_ = s.store.AddDeployLog(&model.DeployLog{
		CertMetaID: meta.ID, Action: "apply", Success: true, Message: "certificate issued and deployed",
	})
	return s.buildResponse(meta, saved), nil
}

// Renew 续期指定域名的证书
func (s *CertService) Renew(ctx context.Context, domain string) error {
	domain = normalizeDomain(domain)
	meta, err := s.store.GetCertMetaByDomain(domain)
	if err != nil {
		return ErrCertNotFound
	}

	mu := s.lockDomain(domain)
	mu.Lock()
	defer mu.Unlock()

	// 标记为 renewing
	meta.Status = model.CertStatusRenewing
	if err = s.store.SaveCertMeta(meta); err != nil {
		return err
	}

	domains := storage.ParseDomains(meta.Domains)
	provider, err := s.dnsFactory.Get(meta.DNSProvider)
	if err != nil {
		return s.failRenew(meta, err)
	}
	if err = s.acmeClient.SetDNSProvider(provider); err != nil {
		return s.failRenew(meta, err)
	}

	slog.Info("renewing certificate", "domain", domain)
	result, err := s.acmeClient.Renew(domains, nil)
	if err != nil {
		return s.failRenew(meta, err)
	}

	// 上传新证书到 CAS
	newUpload, err := s.casClient.Upload(domain, result.Certificate, result.PrivateKey)
	if err != nil {
		return s.failRenew(meta, err)
	}
	oldCertID := meta.CASCertID

	// 更新元数据
	meta.CASCertID = newUpload.CertID
	meta.NotBefore = result.NotBefore
	meta.NotAfter = result.NotAfter
	meta.Status = model.CertStatusActive
	meta.LastError = ""
	meta.ARIRetryAfter = time.Time{}
	if err = s.store.SaveCertMeta(meta); err != nil {
		_ = s.casClient.Delete(newUpload.CertID)
		return err
	}

	// 重新部署到所有目标
	targets, err := s.store.ListDeployTargets(meta.ID)
	if err != nil {
		return err
	}
	certPayload := &deploy.CertPayload{
		Domain:      domain,
		Certificate: result.Certificate,
		PrivateKey:  result.PrivateKey,
		CASCertID:   newUpload.CertID,
	}
	var deployErrs []string
	for _, target := range targets {
		if err = s.deployer.Deploy(ctx, target, certPayload); err != nil {
			deployErrs = append(deployErrs, fmt.Sprintf("%s/%s: %v", target.Type, target.ResourceID, err))
			_ = s.store.UpdateDeployTargetStatus(target.ID, model.DeployStatusFailed, err.Error())
		} else {
			_ = s.store.UpdateDeployTargetStatus(target.ID, model.DeployStatusSuccess, "")
		}
	}

	if len(deployErrs) > 0 {
		msg := strings.Join(deployErrs, "; ")
		return s.failRenew(meta, fmt.Errorf("deploy failed: %s", msg))
	}

	// 全部部署成功 → 删除 CAS 旧证书
	if oldCertID != "" && oldCertID != newUpload.CertID {
		if err = s.casClient.Delete(oldCertID); err != nil {
			slog.Warn("failed to delete old CAS certificate", "cert_id", oldCertID, "error", err)
		}
	}
	_ = s.store.AddDeployLog(&model.DeployLog{
		CertMetaID: meta.ID, Action: "renew", Success: true,
		Message: fmt.Sprintf("renewed, old cert %s deleted", oldCertID),
	})
	slog.Info("certificate renewed", "domain", domain, "new_cert_id", newUpload.CertID)
	return nil
}

func (s *CertService) failRenew(meta *model.CertMeta, cause error) error {
	meta.Status = model.CertStatusFailed
	meta.LastError = cause.Error()
	_ = s.store.SaveCertMeta(meta)
	_ = s.store.AddDeployLog(&model.DeployLog{
		CertMetaID: meta.ID, Action: "renew", Success: false, Message: cause.Error(),
	})
	return cause
}

// CheckRenewal 检查单个证书是否需要续期（供定时任务调用）
func (s *CertService) CheckRenewal(ctx context.Context, meta *model.CertMeta) (bool, error) {
	// 优先使用 ARI 建议时间
	if s.cfg.ACME.KeyType != "" && !meta.ARIRetryAfter.IsZero() {
		return time.Now().After(meta.ARIRetryAfter), nil
	}
	// 降级：按有效期比例计算（应对证书生命周期缩短）
	lifetime := meta.NotAfter.Sub(meta.NotBefore)
	if lifetime <= 0 {
		return true, nil
	}
	renewAt := meta.NotAfter.Add(-time.Duration(float64(lifetime) * s.cfg.ACME.RenewBefore))
	return time.Now().After(renewAt), nil
}

// GetMetaForCheck 供定时任务获取完整的 CertMeta（含续期判断所需字段）
func (s *CertService) GetMetaForCheck(domain string) (*model.CertMeta, error) {
	return s.store.GetCertMetaByDomain(normalizeDomain(domain))
}

// ListLogs 列出操作日志
func (s *CertService) ListLogs(limit int) ([]*model.DeployLog, error) {
	return s.store.ListDeployLogs(limit)
}

// ListDNSProviders 列出已配置的 DNS 服务商
func (s *CertService) ListDNSProviders() []string {
	return s.dnsFactory.ListNames()
}

// List 列出所有证书
func (s *CertService) List(ctx context.Context) ([]*model.CertResponse, error) {
	metas, err := s.store.ListCertMeta()
	if err != nil {
		return nil, err
	}
	result := make([]*model.CertResponse, 0, len(metas))
	for _, m := range metas {
		targets, _ := s.store.ListDeployTargets(m.ID)
		result = append(result, s.buildResponse(m, targets))
	}
	return result, nil
}

// Get 获取单个证书详情
func (s *CertService) Get(ctx context.Context, domain string) (*model.CertResponse, error) {
	meta, err := s.store.GetCertMetaByDomain(normalizeDomain(domain))
	if err != nil {
		return nil, ErrCertNotFound
	}
	targets, _ := s.store.ListDeployTargets(meta.ID)
	return s.buildResponse(meta, targets), nil
}

// Delete 删除证书（含 CAS 清理）
func (s *CertService) Delete(ctx context.Context, domain string) error {
	meta, err := s.store.GetCertMetaByDomain(normalizeDomain(domain))
	if err != nil {
		return ErrCertNotFound
	}
	if meta.CASCertID != "" {
		if err = s.casClient.Delete(meta.CASCertID); err != nil {
			slog.Warn("delete CAS cert failed", "cert_id", meta.CASCertID, "error", err)
		}
	}
	return s.store.DeleteCertMeta(meta.ID)
}

func (s *CertService) saveAndDeploy(ctx context.Context, certID int64, inputs []model.TargetInput, cert *deploy.CertPayload) ([]*model.DeployTarget, error) {
	var saved []*model.DeployTarget
	for _, input := range inputs {
		t := &model.DeployTarget{
			CertMetaID: certID,
			Type:       input.Type,
			ResourceID: input.ResourceID,
			RegionID:   input.RegionID,
			NginxPath:  input.NginxPath,
			Extra:      input.Extra,
			Status:     model.DeployStatusPending,
		}
		id, err := s.store.SaveDeployTarget(t)
		if err != nil {
			return saved, err
		}
		t.ID = id
		if err = s.deployer.Deploy(ctx, t, cert); err != nil {
			_ = s.store.UpdateDeployTargetStatus(id, model.DeployStatusFailed, err.Error())
			_ = s.store.AddDeployLog(&model.DeployLog{
				CertMetaID: certID, Action: "deploy", Target: t.Type + "/" + t.ResourceID,
				Success: false, Message: err.Error(),
			})
			return saved, fmt.Errorf("deploy %s/%s: %w", t.Type, t.ResourceID, err)
		}
		_ = s.store.UpdateDeployTargetStatus(id, model.DeployStatusSuccess, "")
		saved = append(saved, t)
	}
	return saved, nil
}

func (s *CertService) buildResponse(m *model.CertMeta, targets []*model.DeployTarget) *model.CertResponse {
	if targets == nil {
		targets = []*model.DeployTarget{}
	}
	daysLeft := 0
	if !m.NotAfter.IsZero() {
		daysLeft = int(time.Until(m.NotAfter).Hours() / 24)
	}
	return &model.CertResponse{
		ID:          m.ID,
		Domain:      m.Domain,
		Domains:     storage.ParseDomains(m.Domains),
		CASCertID:   m.CASCertID,
		NotBefore:   m.NotBefore,
		NotAfter:    m.NotAfter,
		DaysLeft:    daysLeft,
		DNSProvider: m.DNSProvider,
		Status:      m.Status,
		Targets:     targets,
	}
}

func (s *CertService) lockDomain(domain string) *sync.Mutex {
	mu, _ := s.locks.LoadOrStore(domain, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func normalizeDomain(d string) string {
	d = strings.ToLower(strings.TrimSpace(d))
	d = strings.TrimPrefix(d, "*.")
	return strings.TrimSuffix(d, ".")
}
