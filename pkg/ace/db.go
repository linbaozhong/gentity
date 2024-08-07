package ace

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/linbaozhong/gentity/pkg/log"
)

type (
	Executer interface {
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}
	DB struct {
		*sqlx.DB
	}
)

func Connect(driverName, dns string) (*DB, error) {
	db, e := sqlx.Connect(driverName, dns)
	if e != nil {
		log.Panic(e)
	}
	return &DB{db}, e
}

// Transaction 事务处理
func (s *DB) Transaction(ctx context.Context, f func(tx *sqlx.Tx) (interface{}, error)) (interface{}, error) {
	tx, e := s.BeginTxx(ctx, nil)
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
