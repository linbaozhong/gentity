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

type (
	cacheType string
)

var (
	Err_ToSql = errors.New("the SQL command terminates execution after the ToSql() method")

	_obj        *DB
	Transaction func(ctx context.Context, f func(tx *Tx) (any, error)) (any, error)
)

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

//
// func Or(fns ...dialect.Condition) dialect.Condition {
//	return func() (string, any) {
//		if len(fns) == 0 {
//			return "", nil
//		}
//		var (
//			buf    bytes.Buffer
//			params = make([]any, 0, len(fns))
//		)
//
//		buf.WriteString(dialect.Operator_or + "(")
//		for i, fn := range fns {
//			cond, val := fn()
//			if i > 0 {
//				buf.WriteString(dialect.Operator_and)
//			}
//			buf.WriteString(cond)
//			if vals, ok := val.([]any); ok {
//				params = append(params, vals...)
//			} else {
//				params = append(params, val)
//			}
//		}
//		buf.WriteString(")")
//
//		return buf.String(), params
//	}
// }
//
// func And(fns ...dialect.Condition) dialect.Condition {
//	return func() (string, any) {
//		if len(fns) == 0 {
//			return "", nil
//		}
//		var (
//			buf    bytes.Buffer
//			params = make([]any, 0, len(fns))
//		)
//
//		buf.WriteString(dialect.Operator_and + "(")
//		for i, fn := range fns {
//			cond, val := fn()
//			if i > 0 {
//				buf.WriteString(dialect.Operator_or)
//			}
//			buf.WriteString(cond)
//			if vals, ok := val.([]any); ok {
//				params = append(params, vals...)
//			} else {
//				params = append(params, val)
//			}
//		}
//		buf.WriteString(")")
//
//		return buf.String(), params
//	}
// }
//
// // Asc 函数用于创建一个升序排序的规则。它接收可变数量的 dialect.Field 类型的参数，
// // 返回一个实现了 dialect.Order 接口的函数，该函数会返回排序操作符 "ASC" 和指定的字段列表。
// // 该函数可用于指定查询结果按指定字段进行升序排序。
// func Asc(fs ...dialect.Field) dialect.Order {
//	// 返回一个匿名函数，该函数实现了 dialect.Order 接口，返回排序操作符 "ASC" 和字段列表
//	return func() (string, []dialect.Field) {
//		// 返回升序操作符
//		return dialect.Operator_Asc, fs
//	}
// }
//
// // Desc 函数用于创建一个降序排序的规则。它接收可变数量的 dialect.Field 类型的参数，
// // 返回一个实现了 dialect.Order 接口的函数，该函数会返回排序操作符 "DESC" 和指定的字段列表。
// // 该函数可用于指定查询结果按指定字段进行降序排序。
// func Desc(fs ...dialect.Field) dialect.Order {
//	// 返回一个匿名函数，该函数实现了 dialect.Order 接口，返回排序操作符 "DESC" 和字段列表
//	return func() (string, []dialect.Field) {
//		// 返回降序操作符
//		return dialect.Operator_Desc, fs
//	}
// }
