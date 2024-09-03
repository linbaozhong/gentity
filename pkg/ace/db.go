package ace

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/linbaozhong/gentity/pkg/log"
)

type (
	Cruder interface {
		C() *Creator
		R() *Selector
		U() *Updater
		D() *Deleter
	}
	DB struct {
		*sql.DB
		debug bool // 如果是调试模式，则打印sql命令及错误
	}
)

func Connect(driverName, dns string) (*DB, error) {
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
