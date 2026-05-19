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
	// 允许：字母、数字、下划线、空格、常见操作符、括号
	// 禁止：引号、分号、注释符号等可能导致注入的字符
	if cond == "" {
		return false
	}
	// 检查是否包含危险字符
	dangerous := []string{"'", "\"", ";", "--", "/*", "*/", "xp_", "EXEC", "UNION"}
	lower := strings.ToLower(cond)
	for _, d := range dangerous {
		if strings.Contains(lower, d) {
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
	o.cond = append(o.cond, cond{
		op: dialect.Operator_and,
		conditions: append([]dialect.Condition{}, func(*uint8, dialect.Dialect) (string, any) {
			return cnd, append([]any{}, params...)
		}),
	})
	return o
}

// And 添加 AND 查询条件，子条件之间为 and 关系。
func (o *orm) And(fns ...dialect.Condition) Builder {
	o.cond = append(o.cond, cond{
		op:         dialect.Operator_and,
		conditions: fns,
	})
	return o
}

// AndOr 添加 AND 查询条件，所有子条件之间为 or 关系。
func (o *orm) AndOr(fns ...dialect.Condition) Builder {
	o.cond = append(o.cond, cond{
		op: dialect.Operator_and,
		children: append([]cond{}, cond{
			op:         dialect.Operator_or,
			conditions: fns,
		}),
	})
	return o
}

// Or 添加 OR 查询条件，子条件之间为 or 关系。
func (o *orm) Or(fns ...dialect.Condition) Builder {
	// return o.buildWhereSimpleCondition(fns, dialect.Operator_or)
	o.cond = append(o.cond, cond{
		op:         dialect.Operator_or,
		conditions: fns,
	})
	return o
}

// OrAnd 添加 OR 查询条件，子条件之间为 and 关系。
func (o *orm) OrAnd(fns ...dialect.Condition) Builder {
	// return o.buildWhereBracketsCondition(fns, dialect.Operator_or, dialect.Operator_and)
	o.cond = append(o.cond, cond{
		op: dialect.Operator_or,
		children: append([]cond{}, cond{
			op:         dialect.Operator_and,
			conditions: fns,
		}),
	})
	return o
}

func (o *orm) parseCond(d []cond) (where strings.Builder, params []any, e error) {
	for i, c := range d {
		if i > 0 {
			where.WriteString(c.op.String())
		}
		if len(c.conditions) > 0 {
			where.WriteString("(")
		}
		for i2, cd := range c.conditions {
			s, v := cd(&o.paramIndex, o.db.Dialect())
			// s, v := cd(&o.paramIndex, o.db.Dialect())
			if e = parseWhereParams(v, &params); e != nil {
				o.err = e
				return
			}
			if i2 > 0 {
				where.WriteString(c.op.String())
			}
			if len(c.conditions) > 1 {
				where.WriteString("(" + s + ")")
			} else {
				where.WriteString(s)
			}
		}
		if len(c.conditions) > 0 {
			where.WriteString(")")
		}
		if len(c.children) > 0 {
			if s, v, er := o.parseCond(c.children); e != nil {
				o.err = er
				return
			} else {
				if e = parseWhereParams(v, &params); e != nil {
					o.err = e
					return
				}
				where.WriteString(s.String())
			}
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

//
// // buildWhereSimpleCondition 构建简单WHERE条件（每个条件前都加操作符）
// // 用于 Or 和 And 方法
// // fns: 条件函数数组
// // innerOperator: 条件之间的连接操作符
// func (o *orm) buildWhereSimpleCondition(fns []dialect.Condition, innerOperator dialect.LogicalOperator) Builder {
// 	if len(fns) == 0 || o.err != nil {
// 		return o
// 	}
//
// 	// if o.where.Len() > 0 {
// 	// 	o.where.WriteString(string(innerOperator))
// 	// }
// 	//
// 	// tmpWhereParams := make([]any, len(o.whereParams), len(o.whereParams)+len(fns))
// 	// copy(tmpWhereParams, o.whereParams)
// 	//
// 	// for i, fn := range fns {
// 	// 	conditions, val := fn(&o.paramIndex, o.db.Dialect())
// 	//
// 	// 	// 空值检查：跳过空条件
// 	// 	if conditions == "" {
// 	// 		continue
// 	// 	}
// 	// 	if i > 0 {
// 	// 		o.where.WriteString(string(innerOperator))
// 	// 	}
// 	// 	o.where.WriteString(conditions)
// 	// 	if err := parseWhereParams(val, &tmpWhereParams); err != nil {
// 	// 		o.err = err
// 	// 		return o
// 	// 	}
// 	// }
// 	// o.whereParams = tmpWhereParams
//
// 	//
// 	o.cond = append(o.cond, cond{
// 		op:         innerOperator,
// 		conditions: fns,
// 	})
// 	return o
// }
//
// // buildWhereBracketsCondition 构建带括号的WHERE条件（第一个条件前不加操作符）
// // 用于 AndOr 和 OrAnd 方法
// // fns: 条件函数数组
// // prefixOperator: 与已有条件的连接操作符
// // innerOperator: 条件之间的连接操作符
// func (o *orm) buildWhereBracketsCondition(fns []dialect.Condition, prefixOperator, innerOperator dialect.LogicalOperator) Builder {
// 	if len(fns) == 0 || o.err != nil {
// 		return o
// 	}
//
// 	// if o.where.Len() == 0 {
// 	// 	o.where.WriteString("(")
// 	// } else {
// 	// 	o.where.WriteString(string(prefixOperator) + "(")
// 	// }
// 	//
// 	// tmpWhereParams := make([]any, len(o.whereParams), len(o.whereParams)+len(fns))
// 	// copy(tmpWhereParams, o.whereParams)
// 	//
// 	// for i, fn := range fns {
// 	// 	conditions, val := fn(&o.paramIndex, o.db.Dialect())
// 	//
// 	// 	// 空值检查：跳过空条件
// 	// 	if conditions == "" {
// 	// 		continue
// 	// 	}
// 	//
// 	// 	if i > 0 {
// 	// 		// if strings.HasPrefix(conditions, dialect.Operator_or) || strings.HasPrefix(conditions, dialect.Operator_and) {
// 	// 		// 	o.where.WriteString(" ")
// 	// 		// } else {
// 	// 		o.where.WriteString(string(innerOperator))
// 	// 		// }
// 	// 	}
// 	// 	o.where.WriteString(conditions)
// 	// 	if err := parseWhereParams(val, &tmpWhereParams); err != nil {
// 	// 		o.err = err
// 	// 		return o
// 	// 	}
// 	// }
// 	// o.whereParams = tmpWhereParams
// 	// o.where.WriteString(")")
//
// 	//
// 	o.cond = append(o.cond, cond{
// 		op:         prefixOperator,
// 		conditions: fns,
// 	})
// 	return o
// }
