package dns

import (
	"time"

	"github.com/go-acme/lego/v4/providers/dns/alidns"

	"lets/internal/config"
)

// aliDNS 阿里云 DNS 服务商（包装 lego 的 alidns provider）
type aliDNS struct {
	p *alidns.DNSProvider
}

func newAliDNS(cfg config.DNSProviderCfg) (Provider, error) {
	// 用 lego 的 Config 显式构造，避免依赖环境变量（支持多服务商共存）
	legoCfg := alidns.NewDefaultConfig()
	legoCfg.APIKey = cfg.AccessKeyID
	legoCfg.SecretKey = cfg.AccessKeySecret
	if cfg.Region != "" {
		legoCfg.RegionID = cfg.Region
	}
	p, err := alidns.NewDNSProviderConfig(legoCfg)
	if err != nil {
		return nil, err
	}
	return &aliDNS{p: p}, nil
}

func (a *aliDNS) Present(domain, token, keyAuth string) error {
	return a.p.Present(domain, token, keyAuth)
}

func (a *aliDNS) CleanUp(domain, token, keyAuth string) error {
	return a.p.CleanUp(domain, token, keyAuth)
}

func (a *aliDNS) Timeout() (time.Duration, time.Duration) {
	return a.p.Timeout()
}
