package orm

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/linbaozhong/gentity/pkg/log"
)

type (
	Session *sqlx.Tx
)

var (
	xdb        *sqlx.DB
	Err_NoRows = sql.ErrNoRows
)

func Connect(driverName, dns string) (*sqlx.DB, error) {
	var e error
	xdb, e = sqlx.Connect(driverName, dns)
	if e != nil {
		log.Panic(e)
	}
	return xdb, e
}

func Close() bool {
	if xdb != nil {
		if e := xdb.Close(); e != nil {
			log.Error(e)
			return false
		}
	}
	return true
}

// Transaction 事务处理
func Transaction(ctx context.Context, f func(tx *sqlx.Tx) (interface{}, error)) (interface{}, error) {
	tx, e := xdb.BeginTxx(ctx, nil)
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
