package dns

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"lets/internal/config"
)

// Provider DNS 服务商接口，与 lego 的 DNS01Provider 对齐
type Provider interface {
	Present(domain, token, keyAuth string) error
	CleanUp(domain, token, keyAuth string) error
	Timeout() (time.Duration, time.Duration)
}

// Factory DNS 服务商工厂
type Factory struct {
	cfg       *config.Config
	mu        sync.Mutex
	providers map[string]Provider
}

// NewFactory 创建工厂
func NewFactory(cfg *config.Config) *Factory {
	return &Factory{
		cfg:       cfg,
		providers: make(map[string]Provider),
	}
}

// ForDomain 根据域名选择合适的 DNS 服务商
// 匹配规则：精确域名 → *.根域名 → 默认
func (f *Factory) ForDomain(domain string) (string, Provider, error) {
	name := f.matchProviderName(domain)
	p, err := f.getProvider(name)
	if err != nil {
		return "", nil, err
	}
	return name, p, nil
}

// Get 按名称获取服务商
func (f *Factory) Get(name string) (Provider, error) {
	return f.getProvider(name)
}

func (f *Factory) matchProviderName(domain string) string {
	domain = strings.ToLower(strings.TrimPrefix(domain, "*."))
	// 精确匹配
	if name, ok := f.cfg.DNS.Domains[domain]; ok {
		return name
	}
	// 通配根域名匹配（配置 key 为 *.example.com）
	wildKey := "*." + domain
	if name, ok := f.cfg.DNS.Domains[wildKey]; ok {
		return name
	}
	// 逐级向上查找父域名
	parts := strings.Split(domain, ".")
	for i := 1; i < len(parts)-1; i++ {
		parent := strings.Join(parts[i:], ".")
		if name, ok := f.cfg.DNS.Domains[parent]; ok {
			return name
		}
		if name, ok := f.cfg.DNS.Domains["*."+parent]; ok {
			return name
		}
	}
	return f.cfg.DNS.Default
}

func (f *Factory) getProvider(name string) (Provider, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if p, ok := f.providers[name]; ok {
		return p, nil
	}
	pcfg, ok := f.cfg.DNS.Providers[name]
	if !ok {
		return nil, fmt.Errorf("DNS provider %q not configured", name)
	}
	var p Provider
	var err error
	switch name {
	case "alidns":
		p, err = newAliDNS(pcfg)
	default:
		return nil, fmt.Errorf("DNS provider %q not supported yet", name)
	}
	if err != nil {
		return nil, fmt.Errorf("create DNS provider %q: %w", name, err)
	}
	f.providers[name] = p
	return p, nil
}

// ListNames 列出已配置的服务商名称
func (f *Factory) ListNames() []string {
	names := make([]string, 0, len(f.cfg.DNS.Providers))
	for name := range f.cfg.DNS.Providers {
		names = append(names, name)
	}
	return names
}
