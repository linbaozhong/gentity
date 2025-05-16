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
	"strings"
)

type (
	conditions []dialect.Condition
	orders     []dialect.Order
	sets       []dialect.Setter
	fields     []dialect.Field
)

// Conds 函数用于创建一个条件列表。它接收可变数量的 Condition 类型的参数，
// 返回一个 conditions 类型的切片，该切片包含了所有传入的条件。
// 该函数可用于构建复杂的查询条件。
func Conds(fns ...dialect.Condition) *conditions {
	r := conditions(fns)
	return &r
}

// Conds 函数用于添加条件到条件列表中。它接收可变数量的 Condition 类型的参数，
// 将这些条件添加到 conditions 类型的切片中，并返回更新后的条件列表。
func (c *conditions) Conds(fns ...dialect.Condition) *conditions {
	if len(fns) > 0 {
		*c = append(*c, fns...)
	}
	return c
}

// And 函数用于将多个条件组合成一个逻辑与条件。它接收可变数量的 Condition 类型的参数，
// 返回一个新的 conditions 类型的切片，该切片包含了所有传入的条件。
func (c *conditions) And(fns ...dialect.Condition) *conditions {
	if len(fns) > 0 {
		*c = append(*c, and(fns...))
	}
	return c
}

// Or 函数用于将多个条件组合成一个逻辑或条件。它接收可变数量的 Condition 类型的参数，
// 返回一个新的 conditions 类型的切片，该切片包含了所有传入的条件。
func (c *conditions) Or(fns ...dialect.Condition) *conditions {
	if len(fns) > 0 {
		*c = append(*c, or(fns...))
	}
	return c
}

// ToSlice 函数用于获取 conditions 类型的切片。它返回一个包含所有条件的切片。
// 该函数可用于将 conditions 类型的切片转换为其他类型的切片。
func (c *conditions) ToSlice() []dialect.Condition {
	return *c
}

// Len 函数用于获取 conditions 类型的切片的长度。它返回 conditions 类型的切片的长度。
func (c *conditions) Len() int {
	return len(*c)
}

func or(fns ...dialect.Condition) dialect.Condition {
	return func() (string, any) {
		if len(fns) == 0 {
			return "", nil
		}
		var (
			buf    strings.Builder
			params = make([]any, 0, len(fns))
		)

		buf.WriteString(dialect.Operator_or + "(")
		for i, fn := range fns {
			cond, val := fn()
			if i > 0 {
				if strings.HasPrefix(cond, dialect.Operator_or) || strings.HasPrefix(cond, dialect.Operator_and) {
					buf.WriteString(" ")
				} else {
					buf.WriteString(dialect.Operator_and)
				}
			}
			buf.WriteString(cond)
			if vals, ok := val.([]any); ok {
				params = append(params, vals...)
			} else {
				params = append(params, val)
			}
		}
		buf.WriteString(")")

		return buf.String(), params
	}
}

func and(fns ...dialect.Condition) dialect.Condition {
	return func() (string, any) {
		if len(fns) == 0 {
			return "", nil
		}
		var (
			buf    strings.Builder
			params = make([]any, 0, len(fns))
		)

		buf.WriteString(dialect.Operator_and + "(")
		for i, fn := range fns {
			cond, val := fn()
			if i > 0 {
				if strings.HasPrefix(cond, dialect.Operator_or) || strings.HasPrefix(cond, dialect.Operator_and) {
					buf.WriteString(" ")
				} else {
					buf.WriteString(dialect.Operator_or)
				}
			}
			buf.WriteString(cond)
			if vals, ok := val.([]any); ok {
				params = append(params, vals...)
			} else {
				params = append(params, val)
			}
		}
		buf.WriteString(")")

		return buf.String(), params
	}
}

// -------------------

// Fields 函数用于创建一个字段列表。它接收可变数量的 dialect.Field 类型的参数，
// 如果传入的参数不为空，则返回包含这些参数的切片；否则返回一个空的字段切片。
func Fields(fs ...dialect.Field) *fields {
	r := fields(fs)
	return &r
}

// Fields 方法用于添加字段到字段列表，它接收可变数量的 dialect.Field 类型的参数，
// 将这些字段添加到字段列表中，并返回更新后的字段列表
func (f *fields) Fields(fs ...dialect.Field) *fields {
	if len(fs) > 0 {
		*f = append(*f, fs...)
	}
	return f
}

// ToSlice 方法用于获取 fields 类型的切片，返回包含所有字段的切片
func (f *fields) ToSlice() []dialect.Field {
	return *f
}

// Len 方法用于获取 fields 类型的切片的长度，返回切片的长度
func (f *fields) Len() int {
	return len(*f)
}

// -------------------

// Sets 函数用于创建一个设置器列表。它接收可变数量的 Setter 类型的参数，
// 返回一个 sets 类型的切片，该切片包含了所有传入的设置器。
// 该函数可用于构建复杂的更新语句。
func Sets(fns ...dialect.Setter) *sets {
	r := sets(fns)
	return &r
}

// Sets 函数用于添加设置器到设置器列表中。它接收可变数量的 Setter 类型的参数，
// 将这些设置器添加到 sets 类型的切片中，并返回更新后的设置器列表。
func (s *sets) Sets(fns ...dialect.Setter) *sets {
	if len(fns) > 0 {
		*s = append(*s, fns...)
	}
	return s
}

// ToSlice 函数用于获取 sets 类型的切片。它返回一个包含所有设置器的切片。
// 该函数可用于将 sets 类型的切片转换为其他类型的切片。
func (s *sets) ToSlice() []dialect.Setter {
	return *s
}

// Len 函数用于获取 sets 类型的切片的长度。它返回 sets 类型的切片的长度。
func (s *sets) Len() int {
	return len(*s)
}

// -------------------

// Orders 函数用于创建一个排序规则列表。它接收可变数量的 Field 类型的参数，
// 返回一个 orders 类型的切片，该切片包含了所有传入的排序规则。
// 该函数可用于指定查询结果的排序方式。
func Orders(fns ...dialect.Field) *orders {
	if len(fns) > 0 {
		r := append(orders{}, asc(fns...))
		return &r
	}
	return &orders{}
}

// Asc 函数用于添加升序排序规则到排序规则列表中。它接收可变数量的 Field 类型的参数，
// 将这些规则添加到 orders 类型的切片中，并返回更新后的排序规则列表。
// 该函数可用于指定查询结果按指定字段进行升序排序。
func (o *orders) Asc(fns ...dialect.Field) *orders {
	if len(fns) > 0 {
		*o = append(*o, asc(fns...))
	}
	return o
}

// Desc 函数用于添加降序排序规则到排序规则列表中。它接收可变数量的 Field 类型的参数，
// 将这些规则添加到 orders 类型的切片中，并返回更新后的排序规则列表。
// 该函数可用于指定查询结果按指定字段进行降序排序。
func (o *orders) Desc(fns ...dialect.Field) *orders {
	if len(fns) > 0 {
		*o = append(*o, desc(fns...))
	}
	return o
}

// ToSlice 函数用于获取 orders 类型的切片。它返回一个包含所有排序规则的切片。
// 该函数可用于将 orders 类型的切片转换为其他类型的切片。
func (o *orders) ToSlice() []dialect.Order {
	return *o
}

// Len 函数用于获取 orders 类型的切片的长度。它返回 orders 类型的切片的长度。
func (o *orders) Len() int {
	return len(*o)
}

// asc 函数用于创建一个升序排序的规则。它接收可变数量的 dialect.Field 类型的参数，
// 返回一个实现了 dialect.Order 接口的函数，该函数会返回排序操作符 "ASC" 和指定的字段列表。
// 该函数可用于指定查询结果按指定字段进行升序排序。
func asc(fs ...dialect.Field) dialect.Order {
	// 返回一个匿名函数，该函数实现了 dialect.Order 接口，返回排序操作符 "ASC" 和字段列表
	return func() (string, []dialect.Field) {
		// 返回升序操作符
		return dialect.Operator_Asc, fs
	}
}

// desc 函数用于创建一个降序排序的规则。它接收可变数量的 dialect.Field 类型的参数，
// 返回一个实现了 dialect.Order 接口的函数，该函数会返回排序操作符 "DESC" 和指定的字段列表。
// 该函数可用于指定查询结果按指定字段进行降序排序。
func desc(fs ...dialect.Field) dialect.Order {
	// 返回一个匿名函数，该函数实现了 dialect.Order 接口，返回排序操作符 "DESC" 和字段列表
	return func() (string, []dialect.Field) {
		// 返回降序操作符
		return dialect.Operator_Desc, fs
	}
}
