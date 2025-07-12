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
	"github.com/linbaozhong/gentity/pkg/ace/reflectx"
)

type (
	Tx struct {
		*sql.Tx
		mapper *reflectx.Mapper
		//cache       func(name string) cachego.Cache
		transaction func(ctx context.Context, f func(tx *Tx) (any, error)) (any, error)
		debug       bool // 如果是调试模式，则打印sql命令及错误
	}
)

func (t *Tx) Mapper() *reflectx.Mapper {
	return t.mapper
}

//func (t *Tx) Cache(name string) cachego.Cache {
//	return t.cache(name)
//}

func (t *Tx) Debug() bool {
	return t.debug
}
func (t *Tx) IsDB() bool {
	return false
}
func (t *Tx) Transaction(ctx context.Context, f func(tx *Tx) (any, error)) (any, error) {
	return t.transaction(ctx, f)
}

// QueryContext 执行查询操作
func (s *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.Tx.QueryContext(ctx, query, args...)
}

// QueryRowContext 执行单行查询操作
func (s *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.Tx.QueryRowContext(ctx, query, args...)
}

// ExecContext 执行更新、插入、删除等操作
func (s *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.Tx.ExecContext(ctx, query, args...)
}

// PrepareContext 为以后的查询或执行创建一个准备好的语句。可以从返回的语句并发地运行多个查询或执行。调用者必须调用语句的Stmt。当不再需要语句时，关闭方法。
// 所提供的上下文用于语句的准备，而不是用于语句的执行。
func (s *Tx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.Tx.PrepareContext(ctx, query)
}
