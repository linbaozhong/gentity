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

type Orderer interface {
	OrderFunc(ords ...dialect.Order) Builder
	Order(cols ...dialect.Field) Builder
	Asc(cols ...dialect.Field) Builder
	Desc(cols ...dialect.Field) Builder
}

// OrderFunc 方法用于根据传入的排序规则函数设置排序规则
// 它会遍历传入的排序规则函数，根据规则函数的返回值调用 Asc 或 Desc 方法
func (o *orm) OrderFunc(ords ...dialect.Order) Builder {
	for _, ord := range ords {
		sord, fs := ord()
		if sord == dialect.Operator_Desc {
			o.Desc(fs...)
		} else {
			o.Asc(fs...)
		}
	}
	return o
}

// OrderField
// Deprecated: 此方法后续版本可能会被移除，建议使用 OrderFunc 方法
func (o *orm) OrderField(ords ...dialect.Order) Builder {
	return o.OrderFunc(ords...)
}

// Order 指定查询结果的排序字段，默认升序。
func (o *orm) Order(cols ...dialect.Field) Builder {
	return o.Asc(cols...)
}

// Asc 指定查询结果按指定列升序排序。
func (o *orm) Asc(cols ...dialect.Field) Builder {
	if len(cols) == 0 {
		return o
	}
	for _, col := range cols {
		if o.orderBy.Len() > 0 {
			o.orderBy.WriteByte(',')
		}
		o.orderBy.WriteString(col.Quote())
	}
	return o
}

// Desc 指定查询结果按指定列降序排序
func (o *orm) Desc(cols ...dialect.Field) Builder {
	if len(cols) == 0 {
		return o
	}
	for _, col := range cols {
		if o.orderBy.Len() > 0 {
			o.orderBy.WriteByte(',')
		}
		o.orderBy.WriteString(col.Quote() + dialect.Operator_Desc)
	}
	return o
}
