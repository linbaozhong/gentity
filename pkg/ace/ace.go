// Copyright © 2023 Linbaozhong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ace

import (
	"database/sql"
	"errors"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/ace/dialect/mysql"
	"github.com/linbaozhong/gentity/pkg/ace/dialect/postgres"
	"github.com/linbaozhong/gentity/pkg/ace/dialect/sqlite"
	"github.com/linbaozhong/gentity/pkg/ace/dialect/sqlserver"
	"github.com/linbaozhong/gentity/pkg/app"
)

const (
	CacheTypeMemory  cacheType = "memory"
	CacheTypeRedis   cacheType = "redis"
	CacheTypeSyncMap cacheType = "sync"
)

// Method
const (
	Method_Get    MethodType = iota // 读取符合条件的首行数据
	Method_List                     // 读取符合条件的所有数据
	Method_Post                     // 插入数据
	Method_Patch                    // 更新数据
	Method_Delete                   // 删除数据
)

type (
	cacheType string
	// MethodType DAL方法类型
	MethodType uint8
)

var (
	Err_ToSql = errors.New("the SQL command terminates execution after the Debug() method")

	// _obj  *DB
	// _once sync.Once
)

// Connect
func Connect(driverName, dns string) (*DB, error) {
	// var err error
	// _once.Do(func() {
	// 根据驱动自动选择方言
	var d dialect.Dialect = Dialect(driverName)

	db, e := sql.Open(driverName, dns)
	if e != nil {
		// err = e
		return nil, e
	}
	if e = db.Ping(); e != nil {
		db.Close()
		// err = e
		return nil, e
	}

	_obj := &DB{
		DB:         db,
		driverName: driverName,
		mapper:     mapper(),
		dialect:    d,
		debug:      false,
		cacheType:  CacheTypeSyncMap,
		cacheOpts:  nil,
	}

	app.RegisterServiceCloser(_obj)
	// // 注册事务方法
	// Transaction = _obj.Transaction

	// })

	return _obj, e
}

func Dialect(driverName string) dialect.Dialect {
	switch driverName {
	case "postgres", "postgresql":
		return &postgres.PostgreSQL{}
	case "sqlite", "sqlite3":
		return &sqlite.SQLite{}
	case "sqlserver":
		return &sqlserver.SQLServer{}
	}
	return &mysql.MySQL{}
}

//
// // GetDB
// // 调用该方法前，确保已经调用过 Connect 方法并确保没有 error 产生
// func GetDB() *DB {
// 	if _obj == nil {
// 		log.Panic("db not init")
// 	}
// 	return _obj
// }

// // //////////////////////////
// func GetExec(exec ...Executer) Executer {
// 	if len(exec) > 0 {
// 		return exec[0]
// 	}
// 	return GetDB()
// }
