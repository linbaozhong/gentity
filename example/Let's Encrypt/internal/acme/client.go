package acme

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"

	"lets/internal/config"
)

const (
	caProduction = "https://acme-v02.api.letsencrypt.org/directory"
	caStaging    = "https://acme-staging-v02.api.letsencrypt.org/directory"
)

// acmeUser 实现 lego 的 registration.User 接口
type acmeUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *acmeUser) GetEmail() string                        { return u.Email }
func (u *acmeUser) GetRegistration() *registration.Resource { return u.Registration }
func (u *acmeUser) GetPrivateKey() crypto.PrivateKey        { return u.key }

// accountData 账户密钥持久化结构
type accountData struct {
	Email        string `json:"email"`
	PrivateKey   string `json:"private_key"` // PEM 编码
	Registration []byte `json:"registration"`
}

// Client ACME 客户端
type Client struct {
	cfg       *config.Config
	lego      *lego.Client
	user      *acmeUser
	keyType   certcrypto.KeyType
	mu        sync.Mutex
	certCache map[string]*certificate.Resource // domain → cert，用于续期
}

// CertResult 签发结果
type CertResult struct {
	Domain      string
	Domains     []string
	Certificate string // fullchain PEM
	PrivateKey  string // privkey PEM
	NotBefore   time.Time
	NotAfter    time.Time
}

// NewClient 创建 ACME 客户端（加载或注册账户）
func NewClient(cfg *config.Config) (*Client, error) {
	keyType := certcrypto.EC256
	if strings.EqualFold(cfg.ACME.KeyType, "RSA") {
		keyType = certcrypto.RSA2048
	}

	user, err := loadOrCreateAccount(cfg)
	if err != nil {
		return nil, fmt.Errorf("load ACME account: %w", err)
	}

	caURL := caStaging
	if cfg.IsProduction() {
		caURL = caProduction
	}
	legoCfg := lego.NewConfig(user)
	legoCfg.Certificate.KeyType = keyType
	legoCfg.CADirURL = caURL

	client, err := lego.NewClient(legoCfg)
	if err != nil {
		return nil, fmt.Errorf("create lego client: %w", err)
	}

	c := &Client{
		cfg:       cfg,
		lego:      client,
		user:      user,
		keyType:   keyType,
		certCache: make(map[string]*certificate.Resource),
	}

	// 若账户未注册，则注册
	if user.Registration == nil {
		reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		if err != nil {
			return nil, fmt.Errorf("register ACME account: %w", err)
		}
		user.Registration = reg
		if err = c.saveAccount(); err != nil {
			return nil, fmt.Errorf("save ACME account: %w", err)
		}
	}
	return c, nil
}

// SetDNSProvider 设置 DNS-01 挑战提供商
func (c *Client) SetDNSProvider(p challenge.Provider) error {
	return c.lego.Challenge.SetDNS01Provider(p)
}

// Obtain 签发新证书
func (c *Client) Obtain(domains []string) (*CertResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	req := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}
	cert, err := c.lego.Certificate.Obtain(req)
	if err != nil {
		return nil, fmt.Errorf("obtain certificate: %w", err)
	}
	c.certCache[domains[0]] = cert
	return c.toResult(domains, cert)
}

// Renew 续期证书
func (c *Client) Renew(domains []string, oldCert *certificate.Resource) (*CertResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if oldCert == nil {
		oldCert = c.certCache[domains[0]]
	}
	if oldCert == nil {
		return nil, fmt.Errorf("no previous certificate for %s", domains[0])
	}
	cert, err := c.lego.Certificate.Renew(*oldCert, true, false, "")
	if err != nil {
		return nil, fmt.Errorf("renew certificate: %w", err)
	}
	c.certCache[domains[0]] = cert
	return c.toResult(domains, cert)
}

// ARIInfo ARI 建议的续期信息
type ARIInfo struct {
	RenewAt    time.Time
	RetryAfter time.Duration
}

// GetARI 查询 ACME Renewal Information（若 CA 支持）
func (c *Client) GetARI(certPEM string) (*ARIInfo, error) {
	crt, err := parseCert(certPEM)
	if err != nil {
		return nil, err
	}
	info, err := c.lego.Certificate.GetRenewalInfo(certificate.RenewalInfoRequest{Cert: crt})
	if err != nil {
		return nil, err
	}
	ari := &ARIInfo{RetryAfter: info.RetryAfter}
	if renewAt := info.ShouldRenewAt(time.Now(), 0); renewAt != nil {
		ari.RenewAt = *renewAt
	}
	return ari, nil
}

func parseCert(certPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, fmt.Errorf("decode certificate PEM failed")
	}
	return x509.ParseCertificate(block.Bytes)
}

func (c *Client) toResult(domains []string, cert *certificate.Resource) (*CertResult, error) {
	notBefore, notAfter, err := parseCertDates(cert.Certificate)
	if err != nil {
		return nil, err
	}
	return &CertResult{
		Domain:      domains[0],
		Domains:     domains,
		Certificate: string(cert.Certificate),
		PrivateKey:  string(cert.PrivateKey),
		NotBefore:   notBefore,
		NotAfter:    notAfter,
	}, nil
}

func parseCertDates(certPEM []byte) (time.Time, time.Time, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return time.Time{}, time.Time{}, fmt.Errorf("decode certificate PEM failed")
	}
	crt, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return crt.NotBefore, crt.NotAfter, nil
}

func loadOrCreateAccount(cfg *config.Config) (*acmeUser, error) {
	dir := cfg.ACME.AccountPath
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "account.json")

	if data, err := os.ReadFile(path); err == nil {
		var ad accountData
		if err = json.Unmarshal(data, &ad); err != nil {
			return nil, fmt.Errorf("parse account file: %w", err)
		}
		key, err := parsePrivateKey([]byte(ad.PrivateKey))
		if err != nil {
			return nil, err
		}
		u := &acmeUser{Email: ad.Email, key: key}
		if len(ad.Registration) > 0 {
			var reg registration.Resource
			if err = json.Unmarshal(ad.Registration, &reg); err != nil {
				return nil, err
			}
			u.Registration = &reg
		}
		return u, nil
	}

	// 创建新账户密钥（ECDSA P-256 或 RSA 2048）
	var key crypto.PrivateKey
	var genErr error
	if strings.EqualFold(cfg.ACME.KeyType, "RSA") {
		key, genErr = rsa.GenerateKey(rand.Reader, 2048)
	} else {
		key, genErr = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}
	if genErr != nil {
		return nil, genErr
	}
	return &acmeUser{Email: cfg.ACME.Email, key: key}, nil
}

func (c *Client) saveAccount() error {
	dir := c.cfg.ACME.AccountPath
	path := filepath.Join(dir, "account.json")
	keyPEM, err := encodePrivateKey(c.user.key)
	if err != nil {
		return err
	}
	ad := accountData{Email: c.user.Email, PrivateKey: string(keyPEM)}
	if c.user.Registration != nil {
		ad.Registration, err = json.Marshal(c.user.Registration)
		if err != nil {
			return err
		}
	}
	data, err := json.MarshalIndent(ad, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func encodePrivateKey(key crypto.PrivateKey) ([]byte, error) {
	switch k := key.(type) {
	case *ecdsa.PrivateKey:
		der, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}
		return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der}), nil
	case *rsa.PrivateKey:
		der := x509.MarshalPKCS1PrivateKey(k)
		return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}), nil
	default:
		return nil, fmt.Errorf("unsupported key type %T", key)
	}
}

func parsePrivateKey(data []byte) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("decode private key PEM failed")
	}
	if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	return x509.ParsePKCS8PrivateKey(block.Bytes)
}
