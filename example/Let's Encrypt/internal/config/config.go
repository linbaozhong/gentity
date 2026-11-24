package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config 应用全局配置
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	ACME      ACMEConfig      `yaml:"acme"`
	DNS       DNSConfig       `yaml:"dns"`
	Aliyun    AliyunConfig    `yaml:"aliyun"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
	Log       LogConfig       `yaml:"log"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
	Mode string `yaml:"mode"` // debug / release
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type ACMEConfig struct {
	// 使用 Let's Encrypt staging 环境测试，生产 switched 到 production
	CADir       string  `yaml:"ca_dir"` // staging / production
	Email       string  `yaml:"email"`
	AccountPath string  `yaml:"account_path"`
	RenewBefore float64 `yaml:"renew_before"` // 有效期比例，默认 0.33；ARI 优先
	KeyType     string  `yaml:"key_type"`     // ECDSA / RSA
}

type DNSConfig struct {
	Default   string                    `yaml:"default"`
	Providers map[string]DNSProviderCfg `yaml:"providers"`
	Domains   map[string]string         `yaml:"domains"` // 域名或 *.domain → provider 名
}

type DNSProviderCfg struct {
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Region          string `yaml:"region"`
}

type AliyunConfig struct {
	Region          string    `yaml:"region"`
	AccessKeyID     string    `yaml:"access_key_id"`
	AccessKeySecret string    `yaml:"access_key_secret"`
	CAS             CASConfig `yaml:"cas"`
	OSS             OSSConfig `yaml:"oss"`
	ECS             ECSConfig `yaml:"ecs"`
}

type CASConfig struct {
	Enabled bool `yaml:"enabled"`
}

type OSSConfig struct {
	// 私有 bucket，用于 Nginx 证书私钥中转
	PrivateBucket string `yaml:"private_bucket"`
	PrivatePrefix string `yaml:"private_prefix"`
}

type ECSConfig struct {
	Enabled bool `yaml:"enabled"`
}

type SchedulerConfig struct {
	Enabled  bool `yaml:"enabled"`
	Interval int  `yaml:"interval"` // 小时
}

type LogConfig struct {
	Level string `yaml:"level"`
	Dir   string `yaml:"dir"`
}

var (
	cfg  *Config
	once sync.Once
)

// Load 加载配置：优先读 YAML 文件，环境变量覆盖敏感字段
func Load(path string) (*Config, error) {
	var err error
	once.Do(func() {
		cfg = &Config{}
		data, e := os.ReadFile(path)
		if e != nil {
			err = fmt.Errorf("read config file %s: %w", path, e)
			return
		}
		if e = yaml.Unmarshal(data, cfg); e != nil {
			err = fmt.Errorf("parse config file: %w", e)
			return
		}
		cfg.applyDefaults()
		cfg.applyEnv()
	})
	return cfg, err
}

// Get 返回已加载的配置
func Get() *Config {
	return cfg
}

func (c *Config) applyDefaults() {
	if c.Server.Addr == "" {
		c.Server.Addr = ":8080"
	}
	if c.Database.Path == "" {
		c.Database.Path = "data/certs.db"
	}
	if c.ACME.CADir == "" {
		c.ACME.CADir = "staging" // 默认 staging，安全优先
	}
	if c.ACME.AccountPath == "" {
		c.ACME.AccountPath = "data/accounts"
	}
	if c.ACME.RenewBefore <= 0 {
		c.ACME.RenewBefore = 0.33
	}
	if c.ACME.KeyType == "" {
		c.ACME.KeyType = "ECDSA"
	}
	if c.DNS.Default == "" {
		c.DNS.Default = "alidns"
	}
	if c.Scheduler.Interval <= 0 {
		c.Scheduler.Interval = 6
	}
	if c.Aliyun.Region == "" {
		c.Aliyun.Region = "cn-hangzhou"
	}
	if c.Aliyun.OSS.PrivatePrefix == "" {
		c.Aliyun.OSS.PrivatePrefix = "cert-transfer"
	}
}

// applyEnv 环境变量覆盖（敏感信息不落配置文件）
func (c *Config) applyEnv() {
	if v := os.Getenv("LETS_ALIYUN_AK_ID"); v != "" {
		c.Aliyun.AccessKeyID = v
	}
	if v := os.Getenv("LETS_ALIYUN_AK_SECRET"); v != "" {
		c.Aliyun.AccessKeySecret = v
	}
	if v := os.Getenv("LETS_ACME_EMAIL"); v != "" {
		c.ACME.Email = v
	}
	if v := os.Getenv("LETS_CA_DIR"); v != "" {
		c.ACME.CADir = v
	}
	if v := os.Getenv("LETS_DB_PATH"); v != "" {
		c.Database.Path = v
	}
	if v := os.Getenv("LETS_DNS_ALIDNS_AK_ID"); v != "" {
		if p, ok := c.DNS.Providers["alidns"]; ok {
			p.AccessKeyID = v
			c.DNS.Providers["alidns"] = p
		}
	}
	if v := os.Getenv("LETS_DNS_ALIDNS_AK_SECRET"); v != "" {
		if p, ok := c.DNS.Providers["alidns"]; ok {
			p.AccessKeySecret = v
			c.DNS.Providers["alidns"] = p
		}
	}
}

// IsProduction 是否使用 LE 生产环境
func (c *Config) IsProduction() bool {
	return c.ACME.CADir == "production"
}
