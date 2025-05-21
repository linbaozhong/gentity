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
	"github.com/linbaozhong/gentity/pkg/types"
	"strings"
	"time"
)

var (
	Err_Expression_Empty_Param = errors.New("Expression parameter must have one value")
	Err_Condition_Empty_Param  = errors.New("Condition parameter must have one value")
)

type (
	FieldType interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
			~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
			~float32 | ~float64 |
			~string |
			time.Time | types.Time
	}
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
	if len(val) > 0 && val[0] != nil {
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
	}
}

// Decr 自减
// val 默认为1
func (f Field) Decr(val ...any) ExprSetter {
	var v any
	if len(val) > 0 && val[0] != nil {
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
	}
}

// Replace 替换
func (f Field) Replace(old, new string) ExprSetter {
	return func() (string, any) {
		// 参数校验
		if old == "" {
			return "1 = 0", Err_Expression_Empty_Param
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
	}
}

// Expr 其它表达式
func (f Field) Expr(expr string) ExprSetter {
	return func() (string, any) {
		if expr == "" {
			return "", Err_Expression_Empty_Param
		}
		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" = ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" = ")
		sb.WriteString(Placeholder)

		return sb.String(), expr
	}
}

// Eq 等于
func (f Field) Eq(val any) Condition {
	return func() (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" = ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" = ")
		sb.WriteString(Placeholder)

		return sb.String(), val
	}
}

// NotEq 不等于
func (f Field) NotEq(val any) Condition {
	return func() (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" != ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" != ")
		sb.WriteString(Placeholder)

		return sb.String(), val
	}
}

// Gt 大于
func (f Field) Gt(val any) Condition {
	return func() (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" > ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" > ")
		sb.WriteString(Placeholder)

		return sb.String(), val
	}
}

// Gte 大于或等于
func (f Field) Gte(val any) Condition {
	return func() (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" >= ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" >= ")
		sb.WriteString(Placeholder)

		return sb.String(), val
	}
}

// Lt 小于
func (f Field) Lt(val any) Condition {
	return func() (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" < ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" < ")
		sb.WriteString(Placeholder)

		return sb.String(), val
	}
}

// Lte 小于或等于
func (f Field) Lte(val any) Condition {
	return func() (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}

		// 使用 strings.Builder 进行字符串拼接
		var sb strings.Builder
		// 预先计算并分配足够的内存空间
		sb.Grow(len(f.Quote()) + len(" <= ") + len(Placeholder))
		sb.WriteString(f.Quote())
		sb.WriteString(" <= ")
		sb.WriteString(Placeholder)

		return sb.String(), val
	}
}

func checkSlice(vals ...any) error {
	for _, val := range vals {
		// // 使用反射判断是否为切片类型
		// if reflect.TypeOf(val).Kind() == reflect.Slice {
		// 	return errors.New("params cannot be slices")
		// }
		switch val.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string, bool, time.Time:
			continue
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Uint, types.Uint8, types.Uint16,
			types.Uint32, types.Uint64, types.Float32, types.Float64, types.String, types.Bool, types.Time:
			continue
		default:
			return errors.New("Parameter type error")
		}
	}
	return nil
}

// In 包含
func (f Field) In(vals ...any) Condition {
	return func() (string, any) {
		l := len(vals)
		if l == 0 {
			return "1 = 0", Err_Condition_Empty_Param
		}
		if err := checkSlice(vals...); err != nil {
			return "1 = 0", err
		}

		var sb strings.Builder
		sb.Grow(len(f.Quote()) + len(" In (") + (len(Placeholder)+1)*l)
		sb.WriteString(f.Quote())
		sb.WriteString(" In (")
		sb.WriteString(strings.Repeat(Placeholder+",", l)[:(len(Placeholder)+1)*l-1])
		sb.WriteString(")")
		return sb.String(), vals
	}
}

// NotIn 不包含
func (f Field) NotIn(vals ...any) Condition {
	return func() (string, any) {
		l := len(vals)
		if l == 0 {
			return "1 = 0", Err_Condition_Empty_Param
		}
		if err := checkSlice(vals...); err != nil {
			return "1 = 0", err
		}

		var sb strings.Builder
		sb.Grow(len(f.Quote()) + len(" Not In (") + (len(Placeholder)+1)*l)
		sb.WriteString(f.Quote())
		sb.WriteString(" Not In (")
		sb.WriteString(strings.Repeat(Placeholder+",", l)[:(len(Placeholder)+1)*l-1])
		sb.WriteString(")")
		return sb.String(), vals
	}
}

// Between 在区间
func (f Field) Between(vals ...any) Condition {
	return func() (string, any) {
		if len(vals) != 2 {
			return "1 = 0", errors.New("Between condition must have two value")
		}
		if err := checkSlice(vals...); err != nil {
			return "1 = 0", err
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
			return "1 = 0", Err_Condition_Empty_Param
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
			return "1 = 0", Err_Condition_Empty_Param
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
	}
}

// Rlike 右匹配
func (f Field) Rlike(val any) Condition {
	return func() (string, any) {
		// 检查 val 是否为空
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
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
	}
}
