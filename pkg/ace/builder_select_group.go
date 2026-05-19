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
	if len(cols) == 0 || o.err != nil {
		return o
	}

	o.groupBy = append(o.groupBy, cols...)
	return o
}

// Having 指定查询结果的分组条件
func (o *orm) Having(fns ...dialect.Condition) Builder {
	if len(fns) == 0 || o.err != nil {
		return o
	}

	o.having = append(o.having, cond{
		op:         dialect.Operator_and,
		conditions: fns,
	})
	return o
}
