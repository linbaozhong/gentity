package orm

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/linbaozhong/gentity/pkg/log"
)

type (
	session struct {
		*sqlx.DB
	}
	ExtContext interface {
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	} // sqlx.ExtContext
)

func Connect(driverName, dns string) (session, error) {
	xdb, e := sqlx.Connect(driverName, dns)
	if e != nil {
		log.Panic(e)
	}
	return session{
		xdb,
	}, e
}

func (s session) Close() bool {
	if s.DB != nil {
		if e := s.DB.Close(); e != nil {
			log.Error(e)
			return false
		}
	}
	return true
}

// Transaction 事务处理
func (s session) Transaction(ctx context.Context, f func(tx *sqlx.Tx) (any, error)) (any, error) {
	tx, e := s.BeginTxx(ctx, nil)
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
