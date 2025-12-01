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
	"context"
	"database/sql"
	"errors"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/log"
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
	Err_ToSql = errors.New("the SQL command terminates execution after the ToSql() method")

	_obj        *DB
	Transaction func(ctx context.Context, f func(tx *Tx) (any, error)) (any, error)
	// // Db 实验
	// Db *sql.DB
)

// // Open 实验
// func Open(driverName, dns string) (*sql.DB, error) {
// 	dialect.Register(driverName)
// 	var e error
// 	Db, e = sql.Open(driverName, dns)
// 	if e != nil {
// 		return nil, e
// 	}
// 	if e = Db.Ping(); e != nil {
// 		Db.Close()
// 		return nil, e
// 	}
// 	return Db, e
// }

// Connect
func Connect(driverName, dns string) (*DB, error) {
	dialect.Register(driverName)
	db, e := sql.Open(driverName, dns)
	if e != nil {
		return nil, e
	}
	if e = db.Ping(); e != nil {
		db.Close()
		return nil, e
	}

	_obj = &DB{}
	_obj.DB = db
	_obj.driverName = driverName
	_obj.mapper = mapper()
	_obj.debug = false

	app.RegisterServiceCloser(_obj)
	// 注册事务方法
	Transaction = _obj.Transaction

	return _obj, e
}

// GetDB
// 调用该方法前，确保已经调用过 Connect 方法并确保没有 error 产生
func GetDB() *DB {
	if _obj == nil {
		log.Panic("db not init")
	}
	return _obj
}

// //////////////////////////
func GetExec(exec ...Executer) Executer {
	if len(exec) > 0 {
		return exec[0]
	}
	return GetDB()
}
