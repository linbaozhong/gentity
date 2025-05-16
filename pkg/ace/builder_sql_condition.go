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

type Wherer interface {
	Where(fns ...dialect.Condition) Builder
	And(fns ...dialect.Condition) Builder
	Or(fns ...dialect.Condition) Builder
}

// Where 添加查询条件。
func (o *orm) Where(fns ...dialect.Condition) Builder {
	if len(fns) == 0 || o.err != nil {
		return o
	}

	if o.where.Len() == 0 {
		o.where.WriteString("(")
	} else {
		o.where.WriteString(dialect.Operator_and + "(")
	}

	tmpWhereParams := make([]any, len(o.whereParams), len(o.whereParams)+len(fns))
	copy(tmpWhereParams, o.whereParams)

	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if strings.HasPrefix(cond, dialect.Operator_or) || strings.HasPrefix(cond, dialect.Operator_and) {
				o.where.WriteString(" ")
			} else {
				o.where.WriteString(dialect.Operator_and)
			}
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

// And 添加 AND 查询条件。
func (o *orm) And(fns ...dialect.Condition) Builder {
	if len(fns) == 0 || o.err != nil {
		return o
	}

	if o.where.Len() == 0 {
		o.where.WriteString("(")
	} else {
		o.where.WriteString(dialect.Operator_and + "(")
	}

	tmpWhereParams := make([]any, len(o.whereParams), len(o.whereParams)+len(fns))
	copy(tmpWhereParams, o.whereParams)

	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if strings.HasPrefix(cond, dialect.Operator_or) || strings.HasPrefix(cond, dialect.Operator_and) {
				o.where.WriteString(" ")
			} else {
				o.where.WriteString(dialect.Operator_or)
			}
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

func parseWhereParams(val any, params *[]any) error {
	switch v := val.(type) {
	case error:
		return v
	case []any:
		*params = append(*params, v...)
	default:
		*params = append(*params, val)
	}
	return nil
}

// Or 添加 OR 查询条件。
func (o *orm) Or(fns ...dialect.Condition) Builder {
	if len(fns) == 0 || o.err != nil {
		return o
	}

	if o.where.Len() == 0 {
		o.where.WriteString("(")
	} else {
		o.where.WriteString(dialect.Operator_or + "(")
	}

	tmpWhereParams := make([]any, len(o.whereParams), len(o.whereParams)+len(fns))
	copy(tmpWhereParams, o.whereParams)

	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if strings.HasPrefix(cond, dialect.Operator_or) || strings.HasPrefix(cond, dialect.Operator_and) {
				o.where.WriteString(" ")
			} else {
				o.where.WriteString(dialect.Operator_and)
			}
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
