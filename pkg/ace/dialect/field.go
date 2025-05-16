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
	"errors"
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
	var sb strings.Builder
	// 预先计算并分配足够的内存空间
	sb.Grow(len(f.TableName()) + len(".") + len(f.FieldName()))
	sb.WriteString(f.TableName())
	sb.WriteString(".")
	sb.WriteString(f.FieldName())
	return sb.String()
	// return f.TableName() + "." + f.FieldName()
}

// TableName 为表名添加引号
func (f Field) TableName() string {
	var sb strings.Builder
	// 预先计算并分配足够的内存空间
	sb.Grow(len(Quote_Char) + len(f.Table) + len(Quote_Char))
	sb.WriteString(Quote_Char)
	sb.WriteString(f.Table)
	sb.WriteString(Quote_Char)
	return sb.String()
	// return Quote_Char + f.Table + Quote_Char
}

// FieldName 为字段名添加引号
func (f Field) FieldName() string {
	var sb strings.Builder
	// 预先计算并分配足够的内存空间
	sb.Grow(len(Quote_Char) + len(f.Name) + len(Quote_Char))
	sb.WriteString(Quote_Char)
	sb.WriteString(f.Name)
	sb.WriteString(Quote_Char)
	return sb.String()
	// return Quote_Char + f.Name + Quote_Char
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
		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" = ") + len(f.Quote()) + len(" + ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" = ")
		sb.WriteString(f.Quote())
		sb.WriteString(" + ")
		sb.WriteString(Placeholder)
		return sb.String(), v
		// return f.Quote() + " = " + f.Quote() + " + " + Placeholder, v
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
		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" = ") + len(f.Quote()) + len(" - ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" = ")
		sb.WriteString(f.Quote())
		sb.WriteString(" - ")
		sb.WriteString(Placeholder)

		return sb.String(), v
		// return f.Quote() + " = " + f.Quote() + " - " + Placeholder, v
	}
}

// Replace 替换
func (f Field) Replace(old, new string) ExprSetter {
	return func() (string, any) {
		// 参数校验
		if old == "" {
			return "1 = 0", errors.New("Replace expression must have one value")
		}
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" = REPLACE(") + len(f.Quote()) + len(",?,?)"))
		sb.WriteString(f.Quote())
		sb.WriteString(" = REPLACE(")
		sb.WriteString(f.Quote())
		sb.WriteString(",")
		sb.WriteString(Placeholder)
		sb.WriteString(",")
		sb.WriteString(Placeholder)
		sb.WriteString(")")
		return sb.String(), []any{old, new}
		// return f.Quote() + " = REPLACE(" + f.Quote() + ",'" + old + "','" + new + "')", nil
	}
}

// Expr 其它表达式
func (f Field) Expr(expr string) ExprSetter {
	return func() (string, any) {
		if expr == "" {
			return "", errors.New("Expr expression must have one value")
		}
		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" = ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" = ")
		sb.WriteString(Placeholder)

		return sb.String(), expr
		// return f.Quote() + " = " + expr, nil
	}
}

// Eq 等于
func (f Field) Eq(val any) Condition {
	return func() (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", errors.New("Eq condition must have one value")
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" = ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" = ")
		sb.WriteString(Placeholder)

		return sb.String(), val
		// return f.Quote() + " = " + Placeholder, val
	}
}

// NotEq 不等于
func (f Field) NotEq(val any) Condition {
	return func() (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", errors.New("Not Eq condition must have one value")
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" != ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" != ")
		sb.WriteString(Placeholder)

		return sb.String(), val
		// return f.Quote() + " != " + Placeholder, val
	}
}

// Gt 大于
func (f Field) Gt(val any) Condition {
	return func() (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", errors.New("Gt condition must have one value")
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" > ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" > ")
		sb.WriteString(Placeholder)

		return sb.String(), val
		// return f.Quote() + " > " + Placeholder, val
	}
}

// Gte 大于或等于
func (f Field) Gte(val any) Condition {
	return func() (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", errors.New("Gte condition must have one value")
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" >= ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" >= ")
		sb.WriteString(Placeholder)

		return sb.String(), val
		// return f.Quote() + " >= " + Placeholder, val
	}
}

// Lt 小于
func (f Field) Lt(val any) Condition {
	return func() (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", errors.New("Lt condition must have one value")
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" < ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" < ")
		sb.WriteString(Placeholder)

		return sb.String(), val
		// return f.Quote() + " < " + Placeholder, val
	}
}

// Lte 小于或等于
func (f Field) Lte(val any) Condition {
	return func() (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", errors.New("Lte condition must have one value")
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" <= ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" <= ")
		sb.WriteString(Placeholder)

		return sb.String(), val
		// return f.Quote() + " <= " + Placeholder, val
	}
}

// In 包含
func (f Field) In(vals ...any) Condition {
	return func() (string, any) {
		l := len(vals)
		if l == 0 {
			return "1 = 0", errors.New("In condition must have at least one value")
		}
		_vals := make([]any, 0, l)
		for _, val := range vals {
			if _, ok := val.([]any); ok {
				// _vals = append(_vals, v...)
				// continue
				return "1 = 0", errors.New("params cannot be slices of []any")
			}
			_vals = append(_vals, val)
		}
		l = len(_vals)
		var sb strings.Builder
		sb.Grow(len(f.Quote()) + len(" In (") + (len(Placeholder)+1)*l)
		sb.WriteString(f.Quote())
		sb.WriteString(" In (")
		sb.WriteString(strings.Repeat(Placeholder+",", l)[:(len(Placeholder)+1)*l-1])
		sb.WriteString(")")
		return sb.String(), _vals
	}
}

// NotIn 不包含
func (f Field) NotIn(vals ...any) Condition {
	return func() (string, any) {
		l := len(vals)
		if l == 0 {
			return "1 = 0", errors.New("Not In condition must have at least one value")
		}
		_vals := make([]any, 0, l)
		for _, val := range vals {
			if _, ok := val.([]any); ok {
				// _vals = append(_vals, v...)
				// continue
				return "1 = 0", errors.New("params cannot be slices of []any")
			}
			_vals = append(_vals, val)
		}
		l = len(_vals)
		var sb strings.Builder
		sb.Grow(len(f.Quote()) + len(" Not In (") + (len(Placeholder)+1)*l)
		sb.WriteString(f.Quote())
		sb.WriteString(" Not In (")
		sb.WriteString(strings.Repeat(Placeholder+",", l)[:(len(Placeholder)+1)*l-1])
		sb.WriteString(")")
		return sb.String(), _vals
	}
}

// Between 在区间
func (f Field) Between(vals ...any) Condition {
	return func() (string, any) {
		if len(vals) != 2 {
			return "1 = 0", errors.New("Between condition must have two value")
		}
		var sb strings.Builder
		sb.Grow(len(f.Quote()) + len(" Between ") + len(Placeholder) + len(" And ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" Between ")
		sb.WriteString(Placeholder)
		sb.WriteString(" And ")
		sb.WriteString(Placeholder)
		return sb.String(), vals
	}
}

// Like 匹配
func (f Field) Like(val any) Condition {
	return func() (string, any) {
		// 检查 val 是否为空
		if val == nil {
			return "1 = 0", errors.New("Like condition must have one value")
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len("(") + len(f.Quote()) + len(" LIKE CONCAT('%',") + len(Placeholder) + len(",'%'))"))
		sb.WriteString("(")
		sb.WriteString(f.Quote())
		sb.WriteString(" LIKE CONCAT('%',")
		sb.WriteString(Placeholder)
		sb.WriteString(",'%'))")

		return sb.String(), val
	}
}

// Llike 左匹配
func (f Field) Llike(val any) Condition {
	return func() (string, any) {
		// 检查 val 是否为空
		if val == nil {
			return "1 = 0", errors.New("LeftLike condition must have one value")
		}
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len("(") + len(f.Quote()) + len(" LIKE CONCAT('%',") + len(Placeholder) + len("))"))
		sb.WriteString("(")
		sb.WriteString(f.Quote())
		sb.WriteString(" LIKE CONCAT('%',")
		sb.WriteString(Placeholder)
		sb.WriteString("))")

		return sb.String(), val
		// return "(" + f.Quote() + " LIKE CONCAT('%'," + Placeholder + "))", val
	}
}

// Rlike 右匹配
func (f Field) Rlike(val any) Condition {
	return func() (string, any) {
		// 检查 val 是否为空
		if val == nil {
			return "1 = 0", errors.New("RightLike condition must have one value")
		}
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len("(") + len(f.Quote()) + len(" LIKE CONCAT(") + len(Placeholder) + len(",'%'))"))
		sb.WriteString("(")
		sb.WriteString(f.Quote())
		sb.WriteString(" LIKE CONCAT(")
		sb.WriteString(Placeholder)
		sb.WriteString(",'%'))")
		return sb.String(), val
		// return "(" + f.Quote() + " LIKE CONCAT(" + Placeholder + ",'%'))", val
	}
}

// Null 为空
func (f Field) Null() Condition {
	return func() (string, any) {
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(" ISNULL(") + len(f.Quote()) + len(")"))
		sb.WriteString(" ISNULL(")
		sb.WriteString(f.Quote())
		sb.WriteString(")")
		return sb.String(), nil
		// return " ISNULL(" + f.Quote() + ")", nil
	}
}

// NotNull 不为空
func (f Field) NotNull() Condition {
	return func() (string, any) {
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(" NOT ISNULL(") + len(f.Quote()) + len(")"))
		sb.WriteString(" NOT ISNULL(")
		sb.WriteString(f.Quote())
		sb.WriteString(")")
		return sb.String(), nil
		// return " NOT ISNULL(" + f.Quote() + ")", nil
	}
}

// AsName 别名
func (f Field) AsName(name string) string {
	if name == "" {
		return f.Quote()
	}
	var sb strings.Builder
	// 预先计算并分配足够的内存空间
	sb.Grow(len(f.Quote()) + len(" AS ") + len(name))
	sb.WriteString(f.Quote())
	sb.WriteString(" AS ")
	sb.WriteString(name)
	return sb.String()
	// return f.Quote() + " AS " + name
}

// Sum 合计
func (f Field) Sum(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len("IFNULL(Sum(") + len(f.Quote()) + len("),0) AS ") + len(a))
		sb.WriteString("IFNULL(Sum(")
		sb.WriteString(f.Quote())
		sb.WriteString("),0) AS ")
		sb.WriteString(a)
		return sb.String()
		// return "IFNULL(Sum(" + f.Quote() + "),0) AS " + a
	}
}

// Avg 平均
func (f Field) Avg(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len("IFNULL(Avg(") + len(f.Quote()) + len("),0) AS ") + len(a))
		sb.WriteString("IFNULL(Avg(")
		sb.WriteString(f.Quote())
		sb.WriteString("),0) AS ")
		sb.WriteString(a)
		return sb.String()
		// return "IFNULL(Avg(" + f.Quote() + "),0) AS " + a
	}
}

// Count 计数
func (f Field) Count(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len("IFNULL(Count(") + len(f.Quote()) + len("),0) AS ") + len(a))
		sb.WriteString("IFNULL(Count(")
		sb.WriteString(f.Quote())
		sb.WriteString("),0) AS ")
		sb.WriteString(a)
		return sb.String()
		// return "IFNULL(Count(" + f.Quote() + "),0) AS " + a
	}
}

// Max 最大值
func (f Field) Max(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len("IFNULL(Max(") + len(f.Quote()) + len("),0) AS ") + len(a))
		sb.WriteString("IFNULL(Max(")
		sb.WriteString(f.Quote())
		sb.WriteString("),0) AS ")
		sb.WriteString(a)
		return sb.String()
		// return "IFNULL(Max(" + f.Quote() + "),0) AS " + a
	}
}

// Min 最小值
func (f Field) Min(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func() string {
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len("IFNULL(Min(") + len(f.Quote()) + len("),0) AS ") + len(a))
		sb.WriteString("IFNULL(Min(")
		sb.WriteString(f.Quote())
		sb.WriteString("),0) AS ")
		sb.WriteString(a)
		return sb.String()
		// return "IFNULL(Min(" + f.Quote() + "),0) AS " + a
	}
}
