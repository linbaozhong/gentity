package ace

import (
	"context"
	"database/sql"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/log"
)

type (
	Cruder interface {
		// C Create 命令体
		C() *Creator
		// R Read 命令体
		R() *Selector
		// U Update 命令体
		U() *Updater
		// D Delete 命令体
		D() *Deleter
	}

	Executer interface {
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		Debug() bool
		C(tableName string) *Creator
		D(tableName string) *Deleter
		U(tableName string) *Updater
		R(tableName string) *Selector
	}

	DB struct {
		*sql.DB
		debug bool // 如果是调试模式，则打印sql命令及错误
	}
)

func Connect(driverName, dns string) (*DB, error) {
	dialect.Register(driverName)
	db, e := sql.Open(driverName, dns)
	if e != nil {
		log.Panic(e)
	}
	if e = db.Ping(); e != nil {
		log.Panic(e)
	}

	return &DB{db, false}, e
}

// SetDebug
func (s *DB) SetDebug(debug bool) {
	s.debug = debug
}

// Debug
func (s *DB) Debug() bool {
	return s.debug
}

// Transaction 事务处理
func (s *DB) Transaction(ctx context.Context, f func(tx *sql.Tx) (any, error)) (any, error) {
	tx, e := s.BeginTx(ctx, nil)
	if e != nil {
		return nil, e
	}

	var result any
	result, e = f(tx)
	if e != nil {
		if e = tx.Rollback(); e != nil {
			log.Error(e)
		}
		return result, e
	}

	if e = tx.Commit(); e != nil {
		return result, e
	}

	return result, nil
}

func (s *DB) C(tableName string) *Creator {
	return newCreate(s, tableName)
}

func (s *DB) U(tableName string) *Updater {
	return NewUpdate(s, tableName)
}

func (s *DB) D(tableName string) *Deleter {
	return newDelete(s, tableName)
}

func (s *DB) R(tableName string) *Selector {
	return newSelect(s, tableName)
}
