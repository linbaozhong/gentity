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
	"context"
	"database/sql"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/ace/types"
	"strings"
	"sync"
)

type (
	Updater struct {
		db          Executer
		object      types.Modeler
		affect      []types.Field
		cols        []types.Field
		params      []any
		exprCols    []expr
		where       strings.Builder
		whereParams []interface{}
		command     strings.Builder
	}
	expr struct {
		colName string
		arg     interface{}
	}
)

var (
	updatePool = sync.Pool{
		New: func() any {
			obj := &Updater{}
			return obj
		},
	}
)

// Updater
func NewUpdate(db Executer, mod types.Modeler) *Updater {
	if db == nil || mod == nil {
		panic("db or table is nil")
		return nil
	}
	obj := updatePool.Get().(*Updater)
	obj.db = db
	obj.object = mod
	obj.command.Reset()

	return obj

}

func (u *Updater) Free() {
	u.object = nil
	u.affect = u.affect[:]
	u.cols = u.cols[:]
	u.params = u.params[:]
	u.exprCols = u.exprCols[:]
	u.where.Reset()
	u.whereParams = u.whereParams[:]
	updatePool.Put(u)
}

func (u *Updater) String() string {
	return fmt.Sprintf("%s  %v", u.command.String(), u.params)
}

// Set
func (u *Updater) Set(fns ...types.Setter) *Updater {
	for _, fn := range fns {
		s, val := fn()
		u.cols = append(u.cols, s)
		u.params = append(u.params, val)
	}
	return u
}

func (u *Updater) SetExpr(fns ...types.ExprSetter) *Updater {
	for _, fn := range fns {
		s, val := fn()
		u.exprCols = append(u.exprCols, expr{colName: s, arg: val})
	}
	return u
}

// Where
func (u *Updater) Where(fns ...types.Condition) *Updater {
	if len(fns) == 0 {
		return u
	}
	u.where.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			u.where.WriteString(types.Operator_and)
		}
		cond, val := fn()
		u.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			u.whereParams = append(u.whereParams, vals...)
		} else {
			u.whereParams = append(u.whereParams, val)
		}
	}
	u.where.WriteString(")")

	return u
}

// And
func (u *Updater) And(fns ...types.Condition) *Updater {
	if len(fns) == 0 {
		return u
	}

	if u.where.Len() == 0 {
		u.where.WriteString("(")
	} else {
		u.where.WriteString(types.Operator_and + "(")
	}

	for i, fn := range fns {
		if i > 0 {
			u.where.WriteString(types.Operator_or)
		}
		cond, val := fn()
		u.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			u.whereParams = append(u.whereParams, vals...)
		} else {
			u.whereParams = append(u.whereParams, val)
		}
	}
	u.where.WriteString(")")
	return u
}

// Or
func (u *Updater) Or(fns ...types.Condition) *Updater {
	if len(fns) == 0 {
		return u
	}

	if u.where.Len() == 0 {
		u.where.WriteString("(")
	} else {
		u.where.WriteString(types.Operator_or + "(")
	}

	for i, fn := range fns {
		if i > 0 {
			u.where.WriteString(types.Operator_and)
		}
		cond, val := fn()
		u.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			u.whereParams = append(u.whereParams, vals...)
		} else {
			u.whereParams = append(u.whereParams, val)
		}
	}
	u.where.WriteString(")")
	return u
}

func (u *Updater) Cols(cols ...types.Field) *Updater {
	for _, col := range cols {
		u.affect = append(u.affect, col)
	}
	return u
}

// Do
func (u *Updater) Do(ctx context.Context, beans ...types.Modeler) (sql.Result, error) {
	defer u.Free()

	u.command.WriteString("UPDATE " + types.Quote_Char + u.object.TableName() + types.Quote_Char + " SET ")

	lens := len(u.cols) + len(u.exprCols)
	if lens > 0 {
		_cols := make([]string, 0, lens)
		for _, col := range u.cols {
			_cols = append(_cols, col.Quote()+" = ?")
		}
		for _, col := range u.exprCols {
			_cols = append(_cols, col.colName)
			if col.arg != nil {
				u.params = append(u.params, col.arg)
			}
		}
		u.command.WriteString(strings.Join(_cols, ","))
	} else {
		if len(beans) == 0 {
			return nil, types.ErrCreateEmpty
		}
		_cols := beans[0].AssignColumns(u.affect...)
		u.params = beans[0].AssignValues(u.affect...)
		for i, col := range _cols {
			if i > 0 {
				u.command.WriteString(",")
			}
			u.command.WriteString(col + " = ?")
		}
	}
	// WHERE
	if u.where.Len() > 0 {
		u.command.WriteString(" WHERE " + u.where.String())
	}

	u.params = append(u.params, u.whereParams...)

	return u.db.ExecContext(ctx, u.command.String(), u.params...)
}
