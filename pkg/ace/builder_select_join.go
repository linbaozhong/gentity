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

// Join 添加连接查询条件
func (o *orm) Join(joinType dialect.JoinType, left, right dialect.Field, fns ...dialect.Condition) Builder {
	var on strings.Builder

	if len(fns) > 0 {
		tmpJoinParams := make([]any, len(o.joinParams), len(o.joinParams)+len(fns))
		copy(tmpJoinParams, o.joinParams)

		for _, fn := range fns {
			on.WriteString(dialect.Operator_and)
			cond, val := fn()

			on.WriteString(cond)
			if err := parseWhereParams(val, &tmpJoinParams); err != nil {
				o.err = err
				return o
			}
			// if vals, ok := val.([]any); ok {
			// 	o.joinParams = append(o.joinParams, vals...)
			// } else {
			// 	o.joinParams = append(o.joinParams, val)
			// }
		}
		o.joinParams = tmpJoinParams
	}

	o.join = append(o.join, [3]string{
		string(joinType),
		right.TableName(),
		left.Quote() + "=" + right.Quote() + on.String(),
	})
	return o
}

// LeftJoin 添加左连接查询条件。
func (o *orm) LeftJoin(left, right dialect.Field, fns ...dialect.Condition) Builder {
	return o.Join(dialect.Left_Join, left, right, fns...)
}

// RightJoin 添加右连接查询条件。
func (o *orm) RightJoin(left, right dialect.Field, fns ...dialect.Condition) Builder {
	return o.Join(dialect.Right_Join, left, right, fns...)
}
