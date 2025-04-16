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
	"bytes"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

type (
	conditions []dialect.Condition
	orders     []dialect.Order
	sets       []dialect.Setter
)

// Sets 函数用于创建一个设置器列表。它接收可变数量的 Setter 类型的参数，
// 返回一个 sets 类型的切片，该切片包含了所有传入的设置器。
// 该函数可用于构建复杂的更新语句。
func Sets(fns ...dialect.Setter) sets {
	return fns
}

// Conds 函数用于创建一个条件列表。它接收可变数量的 Condition 类型的参数，
// 返回一个 conditions 类型的切片，该切片包含了所有传入的条件。
// 该函数可用于构建复杂的查询条件。
func Conds(fns ...dialect.Condition) conditions {
	return fns
}

// Where 函数用于添加条件到条件列表中。它接收可变数量的 Condition 类型的参数，
// 将这些条件添加到 conditions 类型的切片中，并返回更新后的条件列表。
func (c conditions) Where(fns ...dialect.Condition) conditions {
	return append(c, fns...)
}

// And 函数用于将多个条件组合成一个逻辑与条件。它接收可变数量的 Condition 类型的参数，
// 返回一个新的 conditions 类型的切片，该切片包含了所有传入的条件。
func (c conditions) And(fns ...dialect.Condition) conditions {
	return append(c, and(fns...))
}

// Or 函数用于将多个条件组合成一个逻辑或条件。它接收可变数量的 Condition 类型的参数，
// 返回一个新的 conditions 类型的切片，该切片包含了所有传入的条件。
func (c conditions) Or(fns ...dialect.Condition) conditions {
	return append(c, or(fns...))
}

func or(fns ...dialect.Condition) dialect.Condition {
	return func() (string, any) {
		if len(fns) == 0 {
			return "", nil
		}
		var (
			buf    bytes.Buffer
			params = make([]any, 0, len(fns))
		)

		buf.WriteString(dialect.Operator_or + "(")
		for i, fn := range fns {
			cond, val := fn()
			if i > 0 {
				buf.WriteString(dialect.Operator_and)
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
			buf    bytes.Buffer
			params = make([]any, 0, len(fns))
		)

		buf.WriteString(dialect.Operator_and + "(")
		for i, fn := range fns {
			cond, val := fn()
			if i > 0 {
				buf.WriteString(dialect.Operator_or)
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

// Orders 函数用于创建一个排序规则列表。它接收可变数量的 Field 类型的参数，
// 返回一个 orders 类型的切片，该切片包含了所有传入的排序规则。
// 该函数可用于指定查询结果的排序方式。
func Orders(fns ...dialect.Field) orders {
	return append([]dialect.Order{nil}, asc(fns...))
}

// Asc 函数用于添加升序排序规则到排序规则列表中。它接收可变数量的 Field 类型的参数，
// 将这些规则添加到 orders 类型的切片中，并返回更新后的排序规则列表。
// 该函数可用于指定查询结果按指定字段进行升序排序。
func (o orders) Asc(fns ...dialect.Field) orders {
	return append(o, asc(fns...))
}

// Desc 函数用于添加降序排序规则到排序规则列表中。它接收可变数量的 Field 类型的参数，
// 将这些规则添加到 orders 类型的切片中，并返回更新后的排序规则列表。
// 该函数可用于指定查询结果按指定字段进行降序排序。
func (o orders) Desc(fns ...dialect.Field) orders {
	return append(o, desc(fns...))
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
