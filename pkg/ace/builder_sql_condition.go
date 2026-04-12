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
)

// buildWhereSimpleCondition 构建简单WHERE条件（每个条件前都加操作符）
// 用于 Or 和 And 方法
// fns: 条件函数数组
// prefixOperator: 与已有条件的连接操作符
// innerOperator: 条件之间的连接操作符
func (o *orm) buildWhereSimpleCondition(fns []dialect.Condition, innerOperator string) Builder {
	if len(fns) == 0 || o.err != nil {
		return o
	}

	if o.where.Len() > 0 {
		o.where.WriteString(innerOperator)
	}

	tmpWhereParams := make([]any, len(o.whereParams), len(o.whereParams)+len(fns))
	copy(tmpWhereParams, o.whereParams)

	for i, fn := range fns {
		cond, val := fn()

		// 空值检查：跳过空条件
		if cond == "" {
			continue
		}
		if i > 0 {
			o.where.WriteString(innerOperator)
		}
		o.where.WriteString(cond)
		if err := parseWhereParams(val, &tmpWhereParams); err != nil {
			o.err = err
			return o
		}
	}
	o.whereParams = tmpWhereParams

	return o
}

// buildWhereBracketsCondition 构建带括号的WHERE条件（第一个条件前不加操作符）
// 用于 AndOr 和 OrAnd 方法
// fns: 条件函数数组
// prefixOperator: 与已有条件的连接操作符
// innerOperator: 条件之间的连接操作符
func (o *orm) buildWhereBracketsCondition(fns []dialect.Condition, prefixOperator, innerOperator string) Builder {
	if len(fns) == 0 || o.err != nil {
		return o
	}

	if o.where.Len() == 0 {
		o.where.WriteString("(")
	} else {
		o.where.WriteString(prefixOperator + "(")
	}

	tmpWhereParams := make([]any, len(o.whereParams), len(o.whereParams)+len(fns))
	copy(tmpWhereParams, o.whereParams)

	for i, fn := range fns {
		cond, val := fn()

		// 空值检查：跳过空条件
		if cond == "" {
			continue
		}

		if i > 0 {
			// if strings.HasPrefix(cond, dialect.Operator_or) || strings.HasPrefix(cond, dialect.Operator_and) {
			// 	o.where.WriteString(" ")
			// } else {
			o.where.WriteString(innerOperator)
			// }
		}
		o.where.WriteString(cond)
		if err := parseWhereParams(val, &tmpWhereParams); err != nil {
			o.err = err
			return o
		}
	}
	o.whereParams = tmpWhereParams
	o.where.WriteString(")")

	return o
}

type Wherer interface {
	Where(fns ...dialect.Condition) Builder
	And(fns ...dialect.Condition) Builder
	Or(fns ...dialect.Condition) Builder
	AndOr(fns ...dialect.Condition) Builder
	OrAnd(fns ...dialect.Condition) Builder
}

// Where 添加查询条件，子条件之间为 and 关系。
func (o *orm) Where(fns ...dialect.Condition) Builder {
	return o.And(fns...)
}

// And 添加 AND 查询条件，子条件之间为 and 关系。
func (o *orm) And(fns ...dialect.Condition) Builder {
	return o.buildWhereSimpleCondition(fns, dialect.Operator_and)
}

// AndOr 添加 AND 查询条件，所有子条件之间为 or 关系。
func (o *orm) AndOr(fns ...dialect.Condition) Builder {
	return o.buildWhereBracketsCondition(fns, dialect.Operator_and, dialect.Operator_or)
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

// Or 添加 OR 查询条件，子条件之间为 or 关系。
func (o *orm) Or(fns ...dialect.Condition) Builder {
	return o.buildWhereSimpleCondition(fns, dialect.Operator_or)
}

// OrAnd 添加 OR 查询条件，子条件之间为 and 关系。
func (o *orm) OrAnd(fns ...dialect.Condition) Builder {
	return o.buildWhereBracketsCondition(fns, dialect.Operator_or, dialect.Operator_and)
}
