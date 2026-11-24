package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"lets/internal/model"

	_ "modernc.org/sqlite"
)

// Store SQLite 存储
type Store struct {
	db *sql.DB
}

// New 创建存储并初始化表结构
func New(path string) (*Store, error) {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create data dir: %w", err)
		}
	}
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	s := &Store{db: db}
	if err = s.initSchema(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initSchema() error {
	const schema = `
CREATE TABLE IF NOT EXISTS cert_meta (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL UNIQUE,
    domains TEXT NOT NULL,
    cas_cert_id TEXT NOT NULL DEFAULT '',
    not_before DATETIME,
    not_after DATETIME,
    ari_retry_after DATETIME,
    dns_provider TEXT NOT NULL DEFAULT '',
    key_type TEXT NOT NULL DEFAULT 'ECDSA',
    status TEXT NOT NULL DEFAULT 'active',
    last_error TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS deploy_target (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cert_meta_id INTEGER NOT NULL,
    type TEXT NOT NULL,
    resource_id TEXT NOT NULL,
    region_id TEXT NOT NULL DEFAULT '',
    nginx_path TEXT NOT NULL DEFAULT '',
    extra TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    last_error TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (cert_meta_id) REFERENCES cert_meta(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_deploy_target_cert ON deploy_target(cert_meta_id);

CREATE TABLE IF NOT EXISTS deploy_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cert_meta_id INTEGER NOT NULL,
    action TEXT NOT NULL,
    target TEXT NOT NULL DEFAULT '',
    success INTEGER NOT NULL DEFAULT 1,
    message TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_deploy_log_cert ON deploy_log(cert_meta_id);
`
	_, err := s.db.Exec(schema)
	return err
}

// Close 关闭数据库
func (s *Store) Close() error { return s.db.Close() }

// SaveCertMeta 保存或更新证书元数据（按 domain 唯一）
func (s *Store) SaveCertMeta(m *model.CertMeta) error {
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	m.UpdatedAt = now
	_, err := s.db.Exec(`
INSERT INTO cert_meta (domain, domains, cas_cert_id, not_before, not_after, ari_retry_after, dns_provider, key_type, status, last_error, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(domain) DO UPDATE SET
    domains=excluded.domains, cas_cert_id=excluded.cas_cert_id,
    not_before=excluded.not_before, not_after=excluded.not_after,
    ari_retry_after=excluded.ari_retry_after, dns_provider=excluded.dns_provider,
    key_type=excluded.key_type, status=excluded.status, last_error=excluded.last_error,
    updated_at=excluded.updated_at`,
		m.Domain, m.Domains, m.CASCertID, m.NotBefore, m.NotAfter, m.ARIRetryAfter,
		m.DNSProvider, m.KeyType, m.Status, m.LastError, m.CreatedAt, m.UpdatedAt)
	return err
}

// GetCertMetaByDomain 按域名查询
func (s *Store) GetCertMetaByDomain(domain string) (*model.CertMeta, error) {
	m := &model.CertMeta{}
	err := s.db.QueryRow(`SELECT id, domain, domains, cas_cert_id, not_before, not_after,
		ari_retry_after, dns_provider, key_type, status, last_error, created_at, updated_at
		FROM cert_meta WHERE domain=?`, domain).Scan(
		&m.ID, &m.Domain, &m.Domains, &m.CASCertID, &m.NotBefore, &m.NotAfter,
		&m.ARIRetryAfter, &m.DNSProvider, &m.KeyType, &m.Status, &m.LastError,
		&m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetCertMetaByID 按 ID 查询
func (s *Store) GetCertMetaByID(id int64) (*model.CertMeta, error) {
	m := &model.CertMeta{}
	err := s.db.QueryRow(`SELECT id, domain, domains, cas_cert_id, not_before, not_after,
		ari_retry_after, dns_provider, key_type, status, last_error, created_at, updated_at
		FROM cert_meta WHERE id=?`, id).Scan(
		&m.ID, &m.Domain, &m.Domains, &m.CASCertID, &m.NotBefore, &m.NotAfter,
		&m.ARIRetryAfter, &m.DNSProvider, &m.KeyType, &m.Status, &m.LastError,
		&m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// ListCertMeta 列出所有证书
func (s *Store) ListCertMeta() ([]*model.CertMeta, error) {
	rows, err := s.db.Query(`SELECT id, domain, domains, cas_cert_id, not_before, not_after,
		ari_retry_after, dns_provider, key_type, status, last_error, created_at, updated_at
		FROM cert_meta ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.CertMeta
	for rows.Next() {
		m := &model.CertMeta{}
		if err = rows.Scan(&m.ID, &m.Domain, &m.Domains, &m.CASCertID, &m.NotBefore,
			&m.NotAfter, &m.ARIRetryAfter, &m.DNSProvider, &m.KeyType, &m.Status,
			&m.LastError, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

// ListActiveCerts 列出需要检查续期的证书（active + renewing）
func (s *Store) ListActiveCerts() ([]*model.CertMeta, error) {
	rows, err := s.db.Query(`SELECT id, domain, domains, cas_cert_id, not_before, not_after,
		ari_retry_after, dns_provider, key_type, status, last_error, created_at, updated_at
		FROM cert_meta WHERE status IN (?, ?)`, model.CertStatusActive, model.CertStatusRenewing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.CertMeta
	for rows.Next() {
		m := &model.CertMeta{}
		if err = rows.Scan(&m.ID, &m.Domain, &m.Domains, &m.CASCertID, &m.NotBefore,
			&m.NotAfter, &m.ARIRetryAfter, &m.DNSProvider, &m.KeyType, &m.Status,
			&m.LastError, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

// DeleteCertMeta 删除证书及关联数据
func (s *Store) DeleteCertMeta(id int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.Exec(`DELETE FROM deploy_target WHERE cert_meta_id=?`, id); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM deploy_log WHERE cert_meta_id=?`, id); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM cert_meta WHERE id=?`, id); err != nil {
		return err
	}
	return tx.Commit()
}

// SaveDeployTarget 保存部署目标，返回 ID
func (s *Store) SaveDeployTarget(t *model.DeployTarget) (int64, error) {
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	t.UpdatedAt = now
	res, err := s.db.Exec(`
INSERT INTO deploy_target (cert_meta_id, type, resource_id, region_id, nginx_path, extra, status, last_error, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    status=excluded.status, last_error=excluded.last_error, updated_at=excluded.updated_at`,
		t.CertMetaID, t.Type, t.ResourceID, t.RegionID, t.NginxPath, t.Extra, t.Status, t.LastError, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return 0, err
	}
	if t.ID > 0 {
		return t.ID, nil
	}
	return res.LastInsertId()
}

// ListDeployTargets 列出证书的部署目标
func (s *Store) ListDeployTargets(certMetaID int64) ([]*model.DeployTarget, error) {
	rows, err := s.db.Query(`SELECT id, cert_meta_id, type, resource_id, region_id, nginx_path, extra, status, last_error, created_at, updated_at
		FROM deploy_target WHERE cert_meta_id=?`, certMetaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.DeployTarget
	for rows.Next() {
		t := &model.DeployTarget{}
		if err = rows.Scan(&t.ID, &t.CertMetaID, &t.Type, &t.ResourceID, &t.RegionID,
			&t.NginxPath, &t.Extra, &t.Status, &t.LastError, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

// UpdateDeployTargetStatus 更新部署目标状态
func (s *Store) UpdateDeployTargetStatus(id int64, status, errMsg string) error {
	_, err := s.db.Exec(`UPDATE deploy_target SET status=?, last_error=?, updated_at=? WHERE id=?`,
		status, errMsg, time.Now(), id)
	return err
}

// AddDeployLog 添加操作日志
func (s *Store) AddDeployLog(log *model.DeployLog) error {
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	_, err := s.db.Exec(`INSERT INTO deploy_log (cert_meta_id, action, target, success, message, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`, log.CertMetaID, log.Action, log.Target, log.Success, log.Message, log.CreatedAt)
	return err
}

// ListDeployLogs 列出操作日志
func (s *Store) ListDeployLogs(limit int) ([]*model.DeployLog, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.Query(`SELECT id, cert_meta_id, action, target, success, message, created_at
		FROM deploy_log ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.DeployLog
	for rows.Next() {
		l := &model.DeployLog{}
		if err = rows.Scan(&l.ID, &l.CertMetaID, &l.Action, &l.Target, &l.Success, &l.Message, &l.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, rows.Err()
}

// ParseDomains 解析 domains JSON 字段
func ParseDomains(s string) []string {
	var domains []string
	_ = json.Unmarshal([]byte(s), &domains)
	return domains
}

// MarshalDomains 序列化 domains
func MarshalDomains(domains []string) string {
	b, _ := json.Marshal(domains)
	return string(b)
}
