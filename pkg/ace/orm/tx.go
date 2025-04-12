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

package orm

import (
	"context"
	"database/sql"
	"github.com/linbaozhong/gentity/pkg/ace/reflectx"
	"github.com/linbaozhong/gentity/pkg/cachego"
)

type (
	Tx struct {
		*sql.Tx
		mapper      *reflectx.Mapper
		cache       func(name string) cachego.Cache
		transaction func(ctx context.Context, f func(tx *Tx) (any, error)) (any, error)
		debug       bool // 如果是调试模式，则打印sql命令及错误
	}
)

func (t *Tx) Mapper() *reflectx.Mapper {
	return t.mapper
}

//func (t *Tx) Create() CreateBuilder {
//	return newCreate(t)
//}
//func (t *Tx) Update() UpdateBuilder {
//	return newUpdate(t)
//}
//func (t *Tx) Delete() DeleteBuilder {
//	return newDelete(t)
//}
//func (t *Tx) Select() SelectBuilder {
//	return newSelect(t)
//}

func (t *Tx) Cache(name string) cachego.Cache {
	return t.cache(name)
}
func (t *Tx) Debug() bool {
	return t.debug
}
func (t *Tx) IsDB() bool {
	return false
}
func (t *Tx) Transaction(ctx context.Context, f func(tx *Tx) (any, error)) (any, error) {
	return t.transaction(ctx, f)
}
