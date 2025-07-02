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

import "github.com/linbaozhong/gentity/pkg/ace/dialect"

// Table 设置表名
func Table(t any) Builder {
	return newOrm().Table(t)
}

// // Distinct 设置去重字段
// func Distinct(cols ...dialect.Field) Builder {
// 	return newOrm().Distinct(cols...)
// }

// Cols 设置查询字段
func Cols(cols ...dialect.Field) Builder {
	return newOrm().Cols(cols...)
}

// Omits 设置忽略字段
func Omit(cols ...dialect.Field) Builder {
	return newOrm().Omit(cols...)
}

// Func 聚合函数查询
func Func(fns ...dialect.Function) Builder {
	return newOrm().Func(fns...)
}

// Where 设置条件
func Where(fns ...dialect.Condition) Builder {
	return newOrm().Where(fns...)
}

// Set
// 用于设置更新语句中的字段和值
// 例如：Set(dialect.F("name", "linbaozhong"))
func Set(fns ...dialect.Setter) Builder {
	return newOrm().Set(fns...)
}

// SetExpr
// 用于设置更新语句中的表达式
// 例如：SetExpr(dialect.Expr("age", "age + 1"))
func SetExpr(fns ...dialect.ExprSetter) Builder {
	return newOrm().SetExpr(fns...)
}
