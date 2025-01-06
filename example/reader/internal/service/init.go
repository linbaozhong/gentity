package service

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/log"
	"sync"
)

var (
	db       *ace.DB
	openOnce sync.Once
)

func Open(ctx context.Context) error {
	openOnce.Do(func() {
		var err error
		db, err = ace.Connect(ctx, "mysql",
			"user:password@tcp(0.0.0.0:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
		if err != nil {
			log.Fatal(err)
		}
		db.SetMaxOpenConns(50)
		db.SetMaxIdleConns(25)
		db.SetDebug(true)
	})
	return nil
}

func GetDB() *ace.DB {
	return db
}
