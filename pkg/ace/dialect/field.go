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
	"fmt"
	"math"
	"strings"
)

var (
	Err_Expression_Empty_Param = errors.New("Expression parameter must have one value")
	Err_Condition_Empty_Param  = errors.New("Condition parameter must have one value")
)

type (
	SetOp int8 // 赋值运算符

	Field struct {
		Name      string
		Json      string
		OmitEmpty bool
		Table     string
		Type      string
	}

	Function  func(Dialect) string
	Condition func(*uint8, Dialect) (string, any)
	Order     func() (OrderType, []Field)
	Setter    func() (Field, any, SetOp)
)

const (
	Op_Normal    SetOp = iota // insert 赋值
	Op_Increment              // update 自增
	Op_Decrement              // update 自减
	Op_Replace                // update 替换
	Op_Expr                   // update 其它表达式
)

// Quote 为字段添加引号
func (f *Field) Quote(d Dialect) string {
	return d.Quote(f.Table) + "." + d.Quote(f.Name)
}

// TableName 为表名添加引号
func (f *Field) TableName(d Dialect) string {
	return d.Quote(f.Table)
}

// FieldName 为字段名添加引号
func (f *Field) FieldName(d Dialect) string {
	return d.Quote(f.Name)
}

// Set 赋值：为字段设置值
func (f *Field) Set(val any) Setter {
	return func() (Field, any, SetOp) {
		return *f, val, Op_Normal
	}
}

// Incr 赋值：自增
// val 默认为1
func (f *Field) Incr(val ...any) Setter {
	var v any
	if len(val) == 0 || val[0] == nil {
		v = 1
	} else {
		v = val[0]
	}
	return func() (Field, any, SetOp) {
		return *f, v, Op_Increment
	}
}

// Decr 赋值：自减
// val 默认为1
func (f *Field) Decr(val ...any) Setter {
	var v any
	if len(val) == 0 || val[0] == nil {
		v = 1
	} else {
		v = val[0]
	}
	return func() (Field, any, SetOp) {
		return *f, v, Op_Decrement
	}
}

// Replace 赋值：替换
func (f *Field) Replace(old, new string) Setter {
	// 参数校验
	if old == "" {
		return func() (Field, any, SetOp) {
			return *f, Err_Expression_Empty_Param, Op_Replace
		}
	}

	return func() (Field, any, SetOp) {
		return *f, []any{old, new}, Op_Replace
	}
}

// Expr 赋值：其它表达式
func (f *Field) Expr(expr string) Setter {
	// 参数校验
	if expr == "" {
		return func() (Field, any, SetOp) {
			return *f, Err_Expression_Empty_Param, Op_Replace
		}
	}
	return func() (Field, any, SetOp) {
		return *f, expr, Op_Expr
	}
}

func ParseSetter(set Setter, i *uint8, d Dialect) (string, any, error) {
	f, v, op := set()
	if e, ok := v.(error); ok {
		return "", nil, e
	}
	switch op {
	case Op_Increment:
		return f.Quote(d) + " = " + f.Quote(d) + " + " + d.Placeholder(i), v, nil
	case Op_Decrement:
		return f.Quote(d) + " = " + f.Quote(d) + " - " + d.Placeholder(i), v, nil
	case Op_Replace:
		return f.Quote(d) + " = REPLACE(" + f.Quote(d) + "," + d.Placeholder(i) + "," + d.Placeholder(i) + ")", v, nil
	case Op_Expr:
		return f.Quote(d) + " = " + d.Placeholder(i), v, nil
	default:
		return f.Quote(d), v, Err_Expression_Empty_Param
	}
}

// Eq 条件：等于
func (f *Field) Eq(val any) Condition {
	return func(i *uint8, d Dialect) (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}
		return f.Quote(d) + " = " + d.Placeholder(i), val
	}
}

// NotEq 条件：不等于
func (f *Field) NotEq(val any) Condition {
	return func(i *uint8, d Dialect) (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}
		return f.Quote(d) + " != " + d.Placeholder(i), val
	}
}

// Gt 条件：大于
func (f *Field) Gt(val any) Condition {
	return func(i *uint8, d Dialect) (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}
		return f.Quote(d) + " > " + d.Placeholder(i), val
	}
}

// Gte 条件：大于或等于
func (f *Field) Gte(val any) Condition {
	return func(i *uint8, d Dialect) (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}
		return f.Quote(d) + " >= " + d.Placeholder(i), val
	}
}

// Lt 条件：小于
func (f *Field) Lt(val any) Condition {
	return func(i *uint8, d Dialect) (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}
		return f.Quote(d) + " < " + d.Placeholder(i), val
	}
}

// Lte 条件：小于或等于
func (f *Field) Lte(val any) Condition {
	return func(i *uint8, d Dialect) (string, any) {
		// 空值检查，可根据实际需求决定是否保留
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}
		return f.Quote(d) + " <= " + d.Placeholder(i), val
	}
}

func checkSlice(vals ...any) error {
	for _, val := range vals {
		if val == nil {
			continue
		}
		switch val.(type) {
		case func():
			return errors.New("function type is not allowed")
		case chan any:
			return errors.New("channel type is not allowed")
		case complex64, complex128:
			return errors.New("complex type is not allowed")
		case struct{}:
			return errors.New("struct type is not allowed in IN clause")
		}
	}
	return nil
}

// In 条件：包含
func (f *Field) In(vals ...any) Condition {
	return func(i *uint8, d Dialect) (string, any) {
		l := len(vals)
		if l == 0 {
			return "1 = 0", Err_Condition_Empty_Param
		}
		if err := checkSlice(vals...); err != nil {
			return "1 = 0", err
		}

		var sb strings.Builder
		// 预分配足够的内存空间：len(f.Quote()) + len(" In (") + (len(Placeholder)+1)*l
		sb.Grow(2*l + len(f.Quote(d)) + 5)
		sb.WriteString(f.Quote(d))
		sb.WriteString(" In (")
		// 循环生成不同的占位符
		for j := 0; j < l; j++ {
			if j > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(d.Placeholder(i))
		}
		sb.WriteString(")")

		return sb.String(), vals
	}
}

// NotIn 条件：不包含
func (f *Field) NotIn(vals ...any) Condition {
	return func(i *uint8, d Dialect) (string, any) {
		l := len(vals)
		if l == 0 {
			return "1 = 0", Err_Condition_Empty_Param
		}
		if err := checkSlice(vals...); err != nil {
			return "1 = 0", err
		}

		var sb strings.Builder
		// 预分配足够的内存空间：len(f.Quote()) + len(" Not In (") + (len(Placeholder)+1)*l
		sb.Grow(2*l + len(f.Quote(d)) + 9)
		sb.WriteString(f.Quote(d))
		sb.WriteString(" Not In (")
		// 循环生成不同的占位符
		for j := 0; j < l; j++ {
			if j > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(d.Placeholder(i))
		}
		sb.WriteString(")")
		return sb.String(), vals
	}
}

// Between 条件：在区间
func (f *Field) Between(vals ...any) Condition {
	return func(i *uint8, d Dialect) (string, any) {
		if len(vals) != 2 {
			return "1 = 0", errors.New("Between condition must have two value")
		}
		if err := checkSlice(vals...); err != nil {
			return "1 = 0", err
		}
		return f.Quote(d) + " Between " + d.Placeholder(i) + " And " + d.Placeholder(i), vals
	}
}

// Like 条件：匹配
func (f *Field) Like(val any) Condition {
	return func(i *uint8, d Dialect) (string, any) {
		// 检查 val 是否为空
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}
		return f.Quote(d) + " Like CONCAT('%'," + d.Placeholder(i) + ",'%')", val
	}
}

// Llike 条件：左匹配
func (f *Field) Llike(val any) Condition {
	return func(i *uint8, d Dialect) (string, any) {
		// 检查 val 是否为空
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}
		return f.Quote(d) + " Like CONCAT('%'," + d.Placeholder(i) + ")", val
	}
}

// Rlike 条件：右匹配
func (f *Field) Rlike(val any) Condition {
	return func(i *uint8, d Dialect) (string, any) {
		// 检查 val 是否为空
		if val == nil {
			return "1 = 0", Err_Condition_Empty_Param
		}
		return f.Quote(d) + " Like CONCAT(" + d.Placeholder(i) + ",'%')", val
	}
}

// Null 条件：为空
func (f *Field) Null() Condition {
	return func(i *uint8, d Dialect) (string, any) {
		return f.Quote(d) + " IS NULL", nil
	}
}

// NotNull 条件：不为空
func (f *Field) NotNull() Condition {
	return func(i *uint8, d Dialect) (string, any) {
		return f.Quote(d) + " IS NOT NULL", nil
	}
}

// AsName 别名
func (f *Field) AsName(name string, d Dialect) string {
	if name == "" {
		return f.Quote(d)
	}
	return f.Quote(d) + " AS " + name
}

// Sum 聚合表达式：合计
func (f *Field) Sum(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func(d Dialect) string {
		// var sb strings.Builder
		// // 预先计算并分配足够的内存空间：len(" IFNULL(Sum(") + len(f.Quote()) + len("),0) AS ") + len(a)
		// sb.Grow(len(f.Quote(d)) + 20 + len(a))
		// sb.WriteString(" IFNULL(Sum(")
		// sb.WriteString(f.Quote(d))
		// sb.WriteString("),0) AS ")
		// sb.WriteString(a)
		// return sb.String()
		return d.Null("(Sum("+f.Quote(d)+"),0)") + " AS " + a
	}
}

// Avg 聚合表达式：平均
func (f *Field) Avg(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func(d Dialect) string {
		// var sb strings.Builder
		// // 预先计算并分配足够的内存空间：len(" IFNULL(Avg(") + len(f.Quote()) + len("),0) AS ") + len(a)
		// sb.Grow(len(f.Quote(d)) + 20 + len(a))
		// sb.WriteString(" IFNULL(Avg(")
		// sb.WriteString(f.Quote(d))
		// sb.WriteString("),0) AS ")
		// sb.WriteString(a)
		// return sb.String()
		return d.Null("(Avg("+f.Quote(d)+"),0)") + " AS " + a
	}
}

// Count 聚合表达式：计数
func (f *Field) Count(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func(d Dialect) string {
		// var sb strings.Builder
		// // 预先计算并分配足够的内存空间：len(" IFNULL(Count(") + len(f.Quote()) + len("),0) AS ") + len(a)
		// sb.Grow(len(f.Quote(d)) + 22 + len(a))
		// sb.WriteString(" IFNULL(Count(")
		// sb.WriteString(f.Quote(d))
		// sb.WriteString("),0) AS ")
		// sb.WriteString(a)
		// return sb.String()
		return d.Null("(Count("+f.Quote(d)+"),0)") + " AS " + a
	}
}

// Max 聚合表达式：最大值
func (f *Field) Max(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func(d Dialect) string {
		// var sb strings.Builder
		// // 预先计算并分配足够的内存空间：len(" IFNULL(Max(") + len(f.Quote()) + len("),0) AS ") + len(a)
		// sb.Grow(len(f.Quote(d)) + 20 + len(a))
		// sb.WriteString(" IFNULL(Max(")
		// sb.WriteString(f.Quote(d))
		// sb.WriteString("),0) AS ")
		// sb.WriteString(a)
		// return sb.String()
		return d.Null("(Max("+f.Quote(d)+"),0)") + " AS " + a
	}
}

// Min 聚合表达式：最小值
func (f *Field) Min(as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func(d Dialect) string {
		// var sb strings.Builder
		// // 预先计算并分配足够的内存空间：len("IFNULL(Min(") + len(f.Quote()) + len("),0) AS ") + len(a)
		// sb.Grow(len(f.Quote(d)) + 20 + len(a))
		// sb.WriteString(" IFNULL(Min(")
		// sb.WriteString(f.Quote(d))
		// sb.WriteString("),0) AS ")
		// sb.WriteString(a)
		// return sb.String()
		return d.Null("(Min("+f.Quote(d)+"),0)") + " AS " + a
	}
}

// Distance 列表达式：POINT类型字段到指定定位的距离
func (f *Field) Distance(lng, lat float64, as ...string) Function {
	var a = f.Name
	if len(as) > 0 {
		a = as[0]
	}
	return func(d Dialect) string {
		// var sb strings.Builder
		// sb.Grow(len(f.Quote(d)) + 80 + len(a))
		// sb.WriteString(
		// 	fmt.Sprintf("ST_Distance_Sphere(%s, ST_PointFromText(CONCAT('POINT(%f %f)')),4326)",
		// 		f.Quote(d), lat, lng))
		// sb.WriteString(" AS ")
		// sb.WriteString(a)
		// return sb.String()
		return fmt.Sprintf("ST_Distance_Sphere(%s, ST_PointFromText(CONCAT('POINT(%f %f)')),4326) AS %s",
			f.Quote(d), lat, lng, a)
	}
}

// MBRContains 条件：判断点是否在指定距离（米）的范围内
// @param lng 经度
// @param lat 纬度
// @param radius 半径(米)
func (f *Field) MBRContains(lng, lat, radius float64) Condition {
	lat_offset := radius / 111320
	lat1, lat2 := lat+lat_offset, lat-lat_offset

	lng_offset := radius / (111320 * math.Cos(lat*math.Pi/180))
	lng1, lng2 := lng+lng_offset, lng-lng_offset
	return func(i *uint8, d Dialect) (string, any) {
		return fmt.Sprintf("MBRContains(ST_GeomFromText(CONCAT('POLYGON((',%f,' ',%f,', ',%f,' ',%f,', ',%f,' ',%f,', ',%f,' ',%f,', ',%f,' ',%f,'))'),4326),%s)",
			lat2, lng2, lat1, lng2, lat1, lng1, lat2, lng1, lat2, lng2, f.Quote(d)), nil
	}
}
