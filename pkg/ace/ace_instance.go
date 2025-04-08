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
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"reflect"
)

// Create 创建
func Create(x ...Executer) CreateBuilder {
	return newCreate(GetExec(x...))
}

// Select 查询
func Select(x ...Executer) SelectBuilder {
	return newSelect(GetExec(x...))
}

// Update 更新
func Update(x ...Executer) UpdateBuilder {
	return newUpdate(GetExec(x...))
}

// Delete 删除
func Delete(x ...Executer) DeleteBuilder {
	return newDelete(GetExec(x...))
}

// setTableName 设置表名
func setTableName(p *string, name any) {
	switch v := name.(type) {
	case string:
		*p = v
	case dialect.TableNamer:
		*p = v.TableName()
	default:
		// 避免多次调用 reflect.ValueOf 和 reflect.Indirect
		value := reflect.ValueOf(name)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		*p = value.Type().Name()
	}
}

// //////////////////

// Selector Sql查询构造器
func Selector() Selecter {
	return newSelecter()
}
