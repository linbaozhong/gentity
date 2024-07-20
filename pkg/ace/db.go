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
)

var db *sqlx.DB

func Db() *sqlx.DB {
	if db == nil {
		log.Panic("db is nil")
	}
	return db
}
func Connect(driverName, dns string) (*sqlx.DB, error) {
	var e error
	db, e = sqlx.Connect(driverName, dns)
	if e != nil {
		log.Panic(e)
	}
	return db, e
}

func Close() bool {
	if db != nil {
		if e := db.Close(); e != nil {
			log.Error(e)
			return false
		}
	}
	return true
}

// Transaction 事务处理
func Transaction(ctx context.Context, f func(tx *sqlx.Tx) (interface{}, error)) (interface{}, error) {
	tx, e := db.BeginTxx(ctx, nil)
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
