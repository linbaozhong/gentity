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
)

type (
	Field struct {
		Name      string
		Json      string
		OmitEmpty bool
		Table     string
		Type      string
	}

	Function   func() string
	Condition  func() (string, any)
	Order      func() (string, []Field)
	Setter     func() (Field, any)
	ExprSetter func() (string, any)
)

// Quote 为字段添加引号
func (f Field) Quote() string {
	return f.TableName() + "." + f.FieldName()
}

// TableName 为表名添加引号
func (f Field) TableName() string {
	return Quote_Char + f.Table + Quote_Char
}

// FieldName 为字段名添加引号
func (f Field) FieldName() string {
	return Quote_Char + f.Name + Quote_Char
}

// Set 为字段设置值
func (f Field) Set(val any) Setter {
	return func() (Field, any) {
		return f, val
	}
}

// Incr 自增
// val 默认为1
func (f Field) Incr(val ...any) ExprSetter {
	var v any
	if len(val) > 0 {
		v = val[0]
	} else {
		v = 1
	}
	return func() (string, any) {
		return f.Quote() + " = " + f.Quote() + " + " + Placeholder, v
	}
}

// Decr 自减
// val 默认为1
func (f Field) Decr(val ...any) ExprSetter {
	var v any
	if len(val) > 0 {
		v = val[0]
	} else {
		v = 1
	}
	return func() (string, any) {
		return f.Quote() + " = " + f.Quote() + " - " + Placeholder, v
	}
}

// Replace 替换
func (f Field) Replace(old, new string) ExprSetter {
	return func() (string, any) {
		return f.Quote() + " = REPLACE(" + f.Quote() + ",'" + old + "','" + new + "')", nil
	}
}

// Expr 其它表达式
func (f Field) Expr(expr string) ExprSetter {
	return func() (string, any) {
		return f.Quote() + " = " + expr, nil
	}
}

// Eq 等于
func (f Field) Eq(val any) Condition {
	return func() (string, any) {
		return f.Quote() + " = " + Placeholder, val
	}
}

// NotEq 不等于
func (f Field) NotEq(val any) Condition {
	return func() (string, any) {
		return f.Quote() + " != " + Placeholder, val
	}
}

// Gt 大于
func (f Field) Gt(val any) Condition {
	return func() (string, any) {
		return f.Quote() + " > " + Placeholder, val
	}
}

// Gte 大于或等于
func (f Field) Gte(val any) Condition {
	return func() (string, any) {
		return f.Quote() + " >= " + Placeholder, val
	}
}

// Lt 小于
func (f Field) Lt(val any) Condition {
	return func() (string, any) {
		return f.Quote() + " < " + Placeholder, val
	}
}

// Lte 小于或等于
func (f Field) Lte(val any) Condition {
	return func() (string, any) {
		return f.Quote() + " <= " + Placeholder, val
	}
}

// In 包含
func (f Field) In(vals ...any) Condition {
	return func() (string, any) {
		l := len(vals)
		return f.Quote() + " In (" + strings.Repeat(Placeholder+",", l)[:2*l-1] + ")", vals
	}
}

// NotIn 不包含
func (f Field) NotIn(vals ...any) Condition {
	return func() (string, any) {
		l := len(vals)
		return f.Quote() + " Not In (" + strings.Repeat(Placeholder+",", l)[:2*l-1] + ")", vals
	}
}

// Between 在区间
func (f Field) Between(vals ...any) Condition {
	return func() (string, any) {
		return f.Quote() + " BETWEEN " + Placeholder + " AND " + Placeholder, vals
	}
}

// Like 匹配
func (f Field) Like(val any) Condition {
	return func() (string, any) {
		return "(" + f.Quote() + " LIKE CONCAT('%'," + Placeholder + ",'%'))", val
	}
}

// Llike 左匹配
func (f Field) Llike(val any) Condition {
	return func() (string, any) {
		return "(" + f.Quote() + " LIKE CONCAT('%'," + Placeholder + "))", val
	}
}

// Rlike 右匹配
func (f Field) Rlike(val any) Condition {
	return func() (string, any) {
		return "(" + f.Quote() + " LIKE CONCAT(" + Placeholder + ",'%'))", val
	}
}

// Null 为空
func (f Field) Null(val any) Condition {
	return func() (string, any) {
		return " ISNULL(" + Placeholder + ")", val
	}
}

// NotNull 不为空
func (f Field) NotNull(val any) Condition {
	return func() (string, any) {
		return " NOT ISNULL(" + Placeholder + ")", val
	}
}

// AsName 别名
func (f Field) AsName(name string) string {
	return f.Quote() + " AS " + name
}

// Sum 合计
func (f Field) Sum(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		return "IFNULL(Sum(" + f.Quote() + "),0) AS " + a
	}
}

// Avg 平均
func (f Field) Avg(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		return "IFNULL(Avg(" + f.Quote() + "),0) AS " + a
	}
}

// Count 计数
func (f Field) Count(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		return "IFNULL(Count(" + f.Quote() + "),0) AS " + a
	}
}

// Max 最大值
func (f Field) Max(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		return "IFNULL(Max(" + f.Quote() + "),0) AS " + a
	}
}

// Min 最小值
func (f Field) Min(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		return "IFNULL(Min(" + f.Quote() + "),0) AS " + a
	}
}

// /////////////////////
// OR
//
// // OrEq 等于
// func (f Field) OrEq(val any) Condition {
// 	return func() (string, any) {
// 		return "OR " + f.Quote() + " = " + Placeholder, val
// 	}
// }
//
// // OrNotEq 不等于
// func (f Field) OrNotEq(val any) Condition {
// 	return func() (string, any) {
// 		return "OR " + f.Quote() + " != " + Placeholder, val
// 	}
// }
//
// // OrGt 大于
// func (f Field) OrGt(val any) Condition {
// 	return func() (string, any) {
// 		return "OR " + f.Quote() + " > " + Placeholder, val
// 	}
// }
//
// // OrGte 大于或等于
// func (f Field) OrGte(val any) Condition {
// 	return func() (string, any) {
// 		return "OR " + f.Quote() + " >= " + Placeholder, val
// 	}
// }
//
// // OrLt 小于
// func (f Field) OrLt(val any) Condition {
// 	return func() (string, any) {
// 		return "OR " + f.Quote() + " < " + Placeholder, val
// 	}
// }
//
// // OrLte 小于或等于
// func (f Field) OrLte(val any) Condition {
// 	return func() (string, any) {
// 		return "OR " + f.Quote() + " <= " + Placeholder, val
// 	}
// }
//
// // OrIn 包含
// func (f Field) OrIn(vals ...any) Condition {
// 	return func() (string, any) {
// 		l := len(vals)
// 		return "OR " + f.Quote() + " In (" + strings.Repeat(Placeholder+",", l)[:2*l-1] + ")", vals
// 	}
// }
//
// // OrNotIn 不包含
// func (f Field) OrNotIn(vals ...any) Condition {
// 	return func() (string, any) {
// 		l := len(vals)
// 		return "OR " + f.Quote() + " Not In (" + strings.Repeat(Placeholder+",", l)[:2*l-1] + ")", vals
// 	}
// }
//
// // OrBetween 在区间
// func (f Field) OrBetween(vals ...any) Condition {
// 	return func() (string, any) {
// 		return "OR " + f.Quote() + " BETWEEN " + Placeholder + " AND " + Placeholder, vals
// 	}
// }
//
// // OrLike 匹配
// func (f Field) OrLike(val any) Condition {
// 	return func() (string, any) {
// 		return "OR " + f.Quote() + " LIKE CONCAT('%'," + Placeholder + ",'%')", val
// 	}
// }
//
// // OrLlike 左匹配
// func (f Field) OrLlike(val any) Condition {
// 	return func() (string, any) {
// 		return "OR " + f.Quote() + " LIKE CONCAT('%'," + Placeholder + ")", val
// 	}
// }
//
// // OrRlike 右匹配
// func (f Field) OrRlike(val any) Condition {
// 	return func() (string, any) {
// 		return "OR " + f.Quote() + " LIKE CONCAT(" + Placeholder + ",'%')", val
// 	}
// }
//
// // OrNull 或为空
// func (f Field) OrNull(val any) Condition {
// 	return func() (string, any) {
// 		return "OR ISNULL(" + Placeholder + ")", val
// 	}
// }
//
// // OrNotNull 或不为空
// func (f Field) OrNotNull(val any) Condition {
// 	return func() (string, any) {
// 		return "OR NOT ISNULL(" + Placeholder + ")", val
// 	}
// }
