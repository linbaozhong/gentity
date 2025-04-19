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

type Grouper interface {
	Group(cols ...dialect.Field) Builder
	Having(fns ...dialect.Condition) Builder
}

// Group 指定查询结果的分组字段
func (o *orm) Group(cols ...dialect.Field) Builder {
	if len(cols) == 0 {
		return o
	}
	for _, col := range cols {
		if o.groupBy.Len() > 0 {
			o.groupBy.WriteByte(',')
		}
		o.groupBy.WriteString(col.Quote())
	}
	return o
}

// Having 指定查询结果的分组条件
func (o *orm) Having(fns ...dialect.Condition) Builder {
	if len(fns) == 0 {
		return o
	}

	o.having.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			o.having.WriteString(dialect.Operator_and)
		}
		cond, val := fn()
		o.having.WriteString(cond)
		if vals, ok := val.([]any); ok {
			o.havingParams = append(o.havingParams, vals...)
		}
	}
	return o
}
