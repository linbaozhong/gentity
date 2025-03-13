package service

import (
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/log"
	"sync"
)

type service struct{}

var (
	db       *ace.DB
	openOnce sync.Once
	serve    *service
)

func init() {
	// 注册服务启动器
	app.RegisterServiceLauncher(serve)
}

func (s *service) Launch() error {
	// 启动服务
	var e error
	openOnce.Do(func() {
		db, e = ace.Connect("mysql",
			"user:password@tcp(0.0.0.0:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
		if e != nil {
			return
		}
		db.SetMaxOpenConns(50)
		db.SetMaxIdleConns(25)
		db.SetDebug(true)
	})
	if e != nil {
		log.Fatal(e)
	}
	return e
}
