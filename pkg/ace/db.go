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
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		Debug() bool
		Insert(tableName string) *Creator
		Delete(tableName string) *Deleter
		Update(tableName string) *Updater
		Select(tableName string) *Selector
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
func (s *DB) Transaction(ctx context.Context, f func(tx *sql.Tx) (interface{}, error)) (interface{}, error) {
	tx, e := s.BeginTx(ctx, nil)
	if e != nil {
		return nil, e
	}

	var result interface{}
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

func (s *DB) Insert(tableName string) *Creator {
	return newCreate(s, tableName)
}

func (s *DB) Update(tableName string) *Updater {
	return NewUpdate(s, tableName)
}

func (s *DB) Delete(tableName string) *Deleter {
	return newDelete(s, tableName)
}

func (s *DB) Select(tableName string) *Selector {
	return newSelect(s, tableName)
}
