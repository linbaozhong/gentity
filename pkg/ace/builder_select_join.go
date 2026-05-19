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
	"fmt"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"strings"
)

// Join 添加连接查询条件
func (o *orm) Join(joinType dialect.JoinType, left, right dialect.Field, fns ...dialect.Condition) Builder {
	if o.err != nil {
		return o
	}
	o.join = append(o.join, join{
		joinType:   joinType,
		table:      right,
		left:       left,
		right:      right,
		conditions: fns,
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

func (o *orm) parseJoin(d []join) (joinStr strings.Builder, params []any, e error) {
	for _, j := range d {
		joinStr.WriteString(fmt.Sprintf(" %s JOIN %s ON (%s = %s",
			j.joinType,
			j.table.TableName(o.db.Dialect()),
			j.left.Quote(o.db.Dialect()),
			j.right.Quote(o.db.Dialect())))
		for _, condition := range j.conditions {
			joinStr.WriteString(dialect.Operator_and.String())
			str, val := condition(&o.paramIndex, o.db.Dialect())
			joinStr.WriteString(str)
			if err := parseWhereParams(val, &params); err != nil {
				e = err
				return
			}
		}
		joinStr.WriteString(")")
	}
	return
}
