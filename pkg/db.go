package orm

import (
	"context"
	"database/sql"
	"fmt"
	"ganji/pkg/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type (
	Session *sqlx.Tx
)

var (
	db         *sqlx.DB
	Err_NoRows = sql.ErrNoRows
)

func Db() *sqlx.DB {
	if db == nil {
		var e error
		db, e = sqlx.Connect("mysql",
			fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
				"lbz",
				"p@ssw0rd",
				"127.0.0.1:33061",
				"6lime",
				"charset=utf8mb4&parseTime=true&loc=Local&readTimeout=30s",
			))
		if e != nil {
			log.Panic(e)
		}
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(time.Second * 60)
	}
	return db
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
func Transaction(f func(tx *sqlx.Tx) (interface{}, error)) (interface{}, error) {
	tx, e := Db().BeginTxx(context.Background(), nil)
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
