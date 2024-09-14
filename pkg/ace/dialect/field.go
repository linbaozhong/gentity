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

package dialect

import (
	"strings"
	"time"
)

type (
	// 基本数据类型
	BaseType interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | string | ~bool | time.Time
	}
	Field[T BaseType] struct {
		Name  string
		Table string
		Type  string
	}
	Function               func() string
	Condition[T BaseType]  func() (string, any)
	Setter[T BaseType]     func() (Field[T], T)
	ExprSetter[T BaseType] func() (string, any)
)

// Quote 为字段添加引号
func (f Field[T]) Quote() string {
	return f.TableName() + "." + f.FieldName()
}

// TableName 为表名添加引号
func (f Field[T]) TableName() string {
	return Quote_Char + f.Table + Quote_Char
}

// FieldName 为字段名添加引号
func (f Field[T]) FieldName() string {
	return Quote_Char + f.Name + Quote_Char
}

// Set 为字段设置值
func (f Field[T]) Set(val T) Setter[T] {
	return func() (Field[T], T) {
		return f, val
	}
}

// Incr 自增
// val 默认为1
func (f Field[T]) Incr(val ...T) ExprSetter[T] {
	var v T
	if len(val) > 0 {
		v = val[0]
	} else {
		v = T(1)
	}
	return func() (string, any) {
		return f.Quote() + " = " + f.Quote() + " + " + Placeholder, v
	}
}

// Decr 自减
// val 默认为1
func (f Field[T]) Decr(val ...T) ExprSetter[T] {
	var v T
	if len(val) > 0 {
		v = val[0]
	} else {
		v = T(1)
	}
	return func() (string, any) {
		return f.Quote() + " = " + f.Quote() + " - " + Placeholder, v
	}
}

// Replace 替换
func (f Field[T]) Replace(old, new string) ExprSetter[T] {
	return func() (string, any) {
		return f.Quote() + " = REPLACE(" + f.Quote() + ",'" + old + "','" + new + "')", nil
	}
}

// Expr 其它表达式
func (f Field[T]) Expr(expr string) ExprSetter[T] {
	return func() (string, any) {
		return f.Quote() + " = " + expr, nil
	}
}

// Eq 等于
func (f Field[T]) Eq(val T) Condition[T] {
	return func() (string, any) {
		return f.Quote() + " = " + Placeholder, val
	}
}

// NotEq 不等于
func (f Field[T]) NotEq(val T) Condition[T] {
	return func() (string, any) {
		return f.Quote() + " != " + Placeholder, val
	}
}

// Gt 大于
func (f Field[T]) Gt(val T) Condition[T] {
	return func() (string, any) {
		return f.Quote() + " > " + Placeholder, val
	}
}

// Gte 大于或等于
func (f Field[T]) Gte(val T) Condition[T] {
	return func() (string, any) {
		return f.Quote() + " >= " + Placeholder, val
	}
}

// Lt 小于
func (f Field[T]) Lt(val T) Condition[T] {
	return func() (string, any) {
		return f.Quote() + " < " + Placeholder, val
	}
}

// Lte 小于或等于
func (f Field[T]) Lte(val T) Condition[T] {
	return func() (string, any) {
		return f.Quote() + " <= " + Placeholder, val
	}
}

// In 包含
func (f Field[T]) In(vals ...T) Condition[T] {
	return func() (string, any) {
		l := len(vals)
		return f.Quote() + " In (" + strings.Repeat(Placeholder+",", l)[:2*l-1] + ") ", vals
	}
}

// NotIn 不包含
func (f Field[T]) NotIn(vals ...T) Condition[T] {
	return func() (string, any) {
		l := len(vals)
		return f.Quote() + " Not In (" + strings.Repeat(Placeholder+",", l)[:2*l-1] + ") ", vals
	}
}

// Between 在区间
func (f Field[T]) Between(vals ...T) Condition[T] {
	return func() (string, any) {
		return f.Quote() + " BETWEEN " + Placeholder + " AND " + Placeholder, vals
	}
}

// Like 匹配
func (f Field[T]) Like(val T) Condition[T] {
	return func() (string, any) {
		return "(" + f.Quote() + " LIKE CONCAT('%'," + Placeholder + ",'%'))", val
	}
}

// Llike 左匹配
func (f Field[T]) Llike(val T) Condition[T] {
	return func() (string, any) {
		return "(" + f.Quote() + " LIKE CONCAT('%'," + Placeholder + "))", val
	}
}

// Rlike 右匹配
func (f Field[T]) Rlike(val T) Condition[T] {
	return func() (string, any) {
		return "(" + f.Quote() + " LIKE CONCAT(" + Placeholder + ",'%'))", val
	}
}

// Null 为空
func (f Field[T]) Null(val T) Condition[T] {
	return func() (string, any) {
		return " ISNULL(" + Placeholder + ")", val
	}
}

// NotNull 不为空
func (f Field[T]) NotNull(val T) Condition[T] {
	return func() (string, any) {
		return " NOT ISNULL(" + Placeholder + ")", val
	}
}

// Sum 合计
func (f Field[T]) Sum(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		return "IFNULL(Sum(" + f.Quote() + "),0) AS " + a
	}
}

// Avg 平均
func (f Field[T]) Avg(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		return "IFNULL(Avg(" + f.Quote() + "),0) AS " + a
	}
}

// Count 计数
func (f Field[T]) Count(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		return "IFNULL(Count(" + f.Quote() + "),0) AS " + a
	}
}

// Max 最大值
func (f Field[T]) Max(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		return "IFNULL(Max(" + f.Quote() + "),0) AS " + a
	}
}

// Min 最小值
func (f Field[T]) Min(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		return "IFNULL(Min(" + f.Quote() + "),0) AS " + a
	}
}
