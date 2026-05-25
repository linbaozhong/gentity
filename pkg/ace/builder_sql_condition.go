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
	"errors"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"strings"
)

type Wherer interface {
	Where(fns ...dialect.Condition) Builder
	RawWhere(cond string, params ...any) Builder
	And(fns ...dialect.Condition) Builder
	Or(fns ...dialect.Condition) Builder
	AndOr(fns ...dialect.Condition) Builder
	OrAnd(fns ...dialect.Condition) Builder
}

// Where 添加查询条件，子条件之间为 and 关系。
func (o *orm) Where(fns ...dialect.Condition) Builder {
	return o.And(fns...)
}

// RawWhere 添加原生查询条件
func (o *orm) RawWhere(cond string, params ...any) Builder {
	return o.RawWhereSafe(cond, params...)
}

// Err_RawWhere_Invalid_Condition 非法条件错误
var Err_RawWhere_Invalid_Condition = errors.New("RawWhere: invalid condition format")

// isValidCondition 验证条件格式
func isValidCondition(cond string) bool {
	if cond == "" {
		return false
	}
	// 白名单模式：仅允许 字母、数字、下划线、空格、括号、常见操作符、? 占位符、. 和 ,
	for _, r := range cond {
		isAlphaNum := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9')
		isSafeChar := r == '_' || r == ' ' || r == '(' || r == ')' || r == '?' ||
			r == '.' || r == ',' || r == '=' || r == '>' || r == '<' ||
			r == '!' || r == '+' || r == '-' || r == '*' || r == '/' ||
			r == '|' || r == '&'
		if !isAlphaNum && !isSafeChar {
			return false
		}
	}
	return true
}

// RawWhereSafe 安全地添加原生查询条件
// 支持格式：
//   - 单个条件: "name = ?" 或 "age > ?"
//   - 多个条件: "name = ? AND age > ?"
//   - 括号组: "(name = ? OR age > ?)"
//   - IN 子句: "status IN (?)"  (params 需传入 []any 类型)
func (o *orm) RawWhereSafe(cnd string, params ...any) Builder {
	// 验证条件格式
	if !isValidCondition(cnd) {
		o.err = Err_RawWhere_Invalid_Condition
		return o
	}
	// 在闭包外一次性拷贝，避免每次 Condition 求值时重复分配
	capturedParams := make([]any, len(params))
	copy(capturedParams, params)

	o.cond = append(o.cond, dialect.Condition{
		Op: dialect.Operator_and,
		Condition: func(*uint16, dialect.Dialect) (string, any) {
			return cnd, capturedParams

		},
	})
	return o
}

// And 添加 AND 查询条件，子条件之间为 and 关系。
func (o *orm) And(fns ...dialect.Condition) Builder {
	var l = len(fns)

	if l == 0 || o.err != nil {
		return o
	}

	if l == 1 {
		fns[0].Op = dialect.Operator_and
		o.cond = append(o.cond, fns[0])
	} else {
		for i := 0; i < l; i++ {
			fns[i].Op = dialect.Operator_and
		}
		o.cond = append(o.cond, dialect.Condition{
			Op:       dialect.Operator_and,
			Children: fns,
		})
	}
	return o
}

// AndOr 添加 AND 查询条件，所有子条件之间为 or 关系。
func (o *orm) AndOr(fns ...dialect.Condition) Builder {
	var l = len(fns)

	if l == 0 || o.err != nil {
		return o
	}

	if l == 1 {
		fns[0].Op = dialect.Operator_and
		o.cond = append(o.cond, fns[0])
	} else {
		for i := 0; i < l; i++ {
			fns[i].Op = dialect.Operator_or
		}
		o.cond = append(o.cond, dialect.Condition{
			Op:       dialect.Operator_and,
			Children: fns,
		})
	}
	return o
}

// Or 添加 OR 查询条件，子条件之间为 or 关系。
func (o *orm) Or(fns ...dialect.Condition) Builder {
	var l = len(fns)

	if l == 0 || o.err != nil {
		return o
	}

	if l == 1 {
		fns[0].Op = dialect.Operator_or
		o.cond = append(o.cond, fns[0])
	} else {
		for i := 0; i < l; i++ {
			fns[i].Op = dialect.Operator_or
		}
		o.cond = append(o.cond, dialect.Condition{
			Op:       dialect.Operator_or,
			Children: fns,
		})
	}
	return o
}

// OrAnd 添加 OR 查询条件，子条件之间为 and 关系。
func (o *orm) OrAnd(fns ...dialect.Condition) Builder {
	var l = len(fns)

	if l == 0 || o.err != nil {
		return o
	}

	if l == 1 {
		fns[0].Op = dialect.Operator_or
		o.cond = append(o.cond, fns[0])
	} else {
		for i := 0; i < l; i++ {
			fns[i].Op = dialect.Operator_and
		}
		o.cond = append(o.cond, dialect.Condition{
			Op:       dialect.Operator_or,
			Children: fns,
		})
	}
	return o
}

func (o *orm) parseCond(d []dialect.Condition) (where strings.Builder, params []any, e error) {
	for i, c := range d {
		if i > 0 {
			where.WriteString(c.Op.String())
		}
		if len(c.Children) > 0 {
			var (
				s strings.Builder
				v []any
			)
			if s, v, e = o.parseCond(c.Children); e != nil {
				o.err = e
				return
			} else {
				if e = parseWhereParams(v, &params); e != nil {
					o.err = e
					return
				}
				where.WriteString("(" + s.String() + ")")
			}
		} else {
			s, v := c.Condition(&o.paramIndex, o.db.Dialect())
			if e = parseWhereParams(v, &params); e != nil {
				o.err = e
				return
			}
			where.WriteString("(" + s + ")")
		}
	}

	return
}

func parseWhereParams(val any, params *[]any) error {
	switch v := val.(type) {
	case error:
		return v
	case []any:
		*params = append(*params, v...)
	default:
		if v != nil {
			*params = append(*params, v)
		}
	}
	return nil
}
