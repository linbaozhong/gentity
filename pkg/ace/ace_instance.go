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
func (s *DB) Table(a any, as ...string) Builder {
	return New(s).Table(a, as...)
}

// Cols 指定要操作的列
func (s *DB) Cols(cols ...dialect.Field) Builder {
	return New(s).Cols(cols...)
}

// PureCols 只包含指定的列
func (s *DB) PureCols(cols ...dialect.Field) Builder {
	return New(s).PureCols(cols...)
}

// Omit 忽略指定的列
func (s *DB) Omit(cols ...dialect.Field) Builder {
	return New(s).Omit(cols...)
}

// Distinct 设置查询结果去重
func (s *DB) Distinct(cols ...dialect.Field) Builder {
	return New(s).Distinct(cols...)
}

// Func 添加聚合函数
func (s *DB) Func(fns ...dialect.Function) Builder {
	return New(s).Func(fns...)
}

// Where 设置查询条件
func (s *DB) Where(fns ...dialect.Condition) Builder {
	return New(s).Where(fns...)
}

// RawWhere 设置原始 SQL 条件
func (s *DB) RawWhere(cond string, params ...any) Builder {
	return New(s).RawWhere(cond, params...)
}

// And 添加 AND 条件
func (s *DB) And(fns ...dialect.Condition) Builder {
	return New(s).And(fns...)
}

// Or 添加 OR 条件
func (s *DB) Or(fns ...dialect.Condition) Builder {
	return New(s).Or(fns...)
}

// AndOr 添加 AND-OR 组合条件
func (s *DB) AndOr(fns ...dialect.Condition) Builder {
	return New(s).AndOr(fns...)
}

// OrAnd 添加 OR-AND 组合条件
func (s *DB) OrAnd(fns ...dialect.Condition) Builder {
	return New(s).OrAnd(fns...)
}

// Order 设置排序（默认 ASC）
func (s *DB) Order(cols ...dialect.Field) Builder {
	return New(s).Order(cols...)
}

// OrderFunc 使用自定义排序函数
func (s *DB) OrderFunc(ords ...dialect.Order) Builder {
	return New(s).OrderFunc(ords...)
}

// Asc 升序排序
func (s *DB) Asc(cols ...dialect.Field) Builder {
	return New(s).Asc(cols...)
}

// Desc 降序排序
func (s *DB) Desc(cols ...dialect.Field) Builder {
	return New(s).Desc(cols...)
}

// Group 设置分组
func (s *DB) Group(cols ...dialect.Field) Builder {
	return New(s).Group(cols...)
}

// Having 设置分组条件
func (s *DB) Having(fns ...dialect.Condition) Builder {
	return New(s).Having(fns...)
}

// Join 添加 JOIN
func (s *DB) Join(joinType dialect.JoinType, left, right dialect.Field, fns ...dialect.Condition) Builder {
	return New(s).Join(joinType, left, right, fns...)
}

// LeftJoin 添加 LEFT JOIN
func (s *DB) LeftJoin(left, right dialect.Field, fns ...dialect.Condition) Builder {
	return New(s).LeftJoin(left, right, fns...)
}

// RightJoin 添加 RIGHT JOIN
func (s *DB) RightJoin(left, right dialect.Field, fns ...dialect.Condition) Builder {
	return New(s).RightJoin(left, right, fns...)
}

// Page 设置分页
func (s *DB) Page(pageIndex, pageSize uint) Builder {
	return New(s).Page(pageIndex, pageSize)
}

// PageByBookmark 基于书签的分页
func (s *DB) PageByBookmark(size uint, bm dialect.Condition) Builder {
	return New(s).PageByBookmark(size, bm)
}

// Limit 设置限制
func (s *DB) Limit(size uint, start ...uint) Builder {
	return New(s).Limit(size, start...)
}

// Set 设置更新字段和值
func (s *DB) Set(fns ...dialect.Setter) Builder {
	return New(s).Set(fns...)
}

// SetExpr 设置更新表达式
func (s *DB) SetExpr(fns ...dialect.Setter) Builder {
	return New(s).SetExpr(fns...)
}
