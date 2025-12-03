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

type Columner interface {
	Table(a any) Builder
	GetTableName() string
	// GetCols 获取 orm 对象要查询的列。
	GetCols() []dialect.Field
	// Distinct 设置查询结果去重，并指定去重的列。
	Distinct(cols ...dialect.Field) Builder
	// Cols 指定要查询的列
	Cols(cols ...dialect.Field) Builder
	// Func 添加聚合函数到查询中
	Func(fns ...dialect.Function) Builder
	// Omits 忽略指定的列
	Omit(cols ...dialect.Field) Builder
	// PureCols 只包含指定的列，忽略其他列
	PureCols(cols ...dialect.Field) Builder
}

// GetCols 获取 orm 对象要查询的列。
func (s *orm) GetCols() []dialect.Field {
	return s.cols
}

// Distinct 设置查询结果去重，并指定去重的列。
func (o *orm) Distinct(cols ...dialect.Field) Builder {
	o.distinct = true
	o.cols = cols
	return o
}

// Cols 指定要查询的列
func (o *orm) Cols(cols ...dialect.Field) Builder {
	o.cols = append(o.cols, cols...)
	return o
}

// Omits 忽略指定的列
func (o *orm) Omit(cols ...dialect.Field) Builder {
	o.omits = append(o.omits, cols...)
	return o
}

// PureCols 只包含指定的列，忽略其他列
func (o *orm) PureCols(cols ...dialect.Field) Builder {
	o.cols = cols
	return o
}

// Func 添加聚合函数到查询中
func (o *orm) Func(fns ...dialect.Function) Builder {
	tmpFuncs := make([]string, len(o.funcs), len(fns)+len(o.funcs))
	copy(tmpFuncs, o.funcs)
	for _, fn := range fns {
		tmpFuncs = append(tmpFuncs, fn())
	}
	o.funcs = tmpFuncs
	return o
}
