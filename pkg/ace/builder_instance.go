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

// Select 设置查询字段
func Cols(cols ...dialect.Field) Builder {
	return newOrm().Cols(cols...)
}

// Distinct 设置去重字段
func Distinct(cols ...dialect.Field) Builder {
	return newOrm().Distinct(cols...)
}

// Omits 设置忽略字段
func Omits(cols ...dialect.Field) Builder {
	return newOrm().Omits(cols...)
}

// Funcs 聚合函数查询
func Funcs(fns ...dialect.Function) Builder {
	return newOrm().Funcs(fns...)
}

// Join 设置连接
func Join(joinType dialect.JoinType, left, right dialect.Field, fns ...dialect.Condition) Builder {
	return newOrm().Join(joinType, left, right, fns...)
}

// LeftJoin 设置左连接
func LeftJoin(left, right dialect.Field, fns ...dialect.Condition) Builder {
	return newOrm().LeftJoin(left, right, fns...)
}

// RightJoin 设置右连接
func RightJoin(left, right dialect.Field, fns ...dialect.Condition) Builder {
	return newOrm().RightJoin(left, right, fns...)
}

// Where 设置条件
func Where(fns ...dialect.Condition) Builder {
	return newOrm().Where(fns...)
}

// Order 设置排序，默认升序
func Order(cols ...dialect.Field) Builder {
	return newOrm().Order(cols...)
}

// Asc 指定查询结果按指定列升序排序。
func OrderAsc(cols ...dialect.Field) Builder {
	return newOrm().Asc(cols...)
}

// Desc 指定查询结果按指定列降序排序
func OrderDesc(cols ...dialect.Field) Builder {
	return newOrm().Desc(cols...)
}

// Group 指定查询结果的分组字段
func Group(cols ...dialect.Field) Builder {
	return newOrm().Group(cols...)
}

// Having 指定查询结果的分组条件
func Having(fns ...dialect.Condition) Builder {
	return newOrm().Having(fns...)
}

// Limit
// size 大小
// start 开始位置
func Limit(size uint, start ...uint) Builder {
	return newOrm().Limit(size, start...)
}

// Page
// pageIndex 页码
// pageSize 页大小
func Page(pageIndex, pageSize uint) Builder {
	return newOrm().Page(pageIndex, pageSize)
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
