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
	"errors"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/log"
	"strings"
	"sync"
)

type (
	Updater struct {
		db            Executer
		table         string
		affect        []dialect.Field
		cols          []dialect.Field
		params        []any
		exprCols      []expr
		where         strings.Builder
		whereParams   []any
		command       strings.Builder
		commandString strings.Builder
		err           error
	}
	expr struct {
		colName string
		arg     any
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
func NewUpdate(db Executer, tableName string) *Updater {
	obj := updatePool.Get().(*Updater)
	if db == nil || tableName == "" {
		obj.err = errors.New("db or table is nil")
		return obj
	}

	obj.db = db
	obj.table = tableName
	obj.err = nil
	obj.commandString.Reset()

	return obj

}

func (u *Updater) Free() {
	if u == nil || u.table == "" {
		return
	}

	_ = u.String()
	if u.db.Debug() {
		log.Info(u.String())
	}

	u.table = ""
	u.affect = u.affect[:0]
	u.cols = u.cols[:0]
	u.exprCols = u.exprCols[:0]
	u.where.Reset()
	u.whereParams = u.whereParams[:0]
	u.command.Reset()
	u.params = u.params[:0]

	updatePool.Put(u)
}

func (u *Updater) String() string {
	if u.table == "" {
		u.commandString.WriteString(fmt.Sprintf("%s  %v \n", u.command.String(), u.params))
	}
	return u.commandString.String()
}

// Set
func (u *Updater) Set(fns ...dialect.Setter) *Updater {
	if len(fns) == 0 || u.err != nil {
		return u
	}

	for _, fn := range fns {
		s, val := fn()
		// if v, ok := val.(error); ok {
		//	u.err = v
		//	return u
		// }
		u.cols = append(u.cols, s)
		u.params = append(u.params, val)
	}
	return u
}

func (u *Updater) SetExpr(fns ...dialect.ExprSetter) *Updater {
	if len(fns) == 0 || u.err != nil {
		return u
	}

	for _, fn := range fns {
		s, val := fn()
		// if v, ok := val.(error); ok {
		//	u.err = v
		//	return u
		// }
		u.exprCols = append(u.exprCols, expr{colName: s, arg: val})
	}
	return u
}

// Where
func (u *Updater) Where(fns ...dialect.Condition) *Updater {
	if len(fns) == 0 || u.err != nil {
		return u
	}

	if u.where.Len() == 0 {
		u.where.WriteString("(")
	} else {
		u.where.WriteString(dialect.Operator_and + "(")
	}
	for i, fn := range fns {
		if i > 0 {
			u.where.WriteString(dialect.Operator_and)
		}
		cond, val := fn()
		// if v, ok := val.(error); ok {
		//	u.err = v
		//	return u
		// }
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
func (u *Updater) And(fns ...dialect.Condition) *Updater {
	if len(fns) == 0 || u.err != nil {
		return u
	}

	if u.where.Len() == 0 {
		u.where.WriteString("(")
	} else {
		u.where.WriteString(dialect.Operator_and + "(")
	}

	for i, fn := range fns {
		if i > 0 {
			u.where.WriteString(dialect.Operator_or)
		}
		cond, val := fn()
		// if v, ok := val.(error); ok {
		//	u.err = v
		//	return u
		// }
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
func (u *Updater) Or(fns ...dialect.Condition) *Updater {
	if len(fns) == 0 || u.err != nil {
		return u
	}

	if u.where.Len() == 0 {
		u.where.WriteString("(")
	} else {
		u.where.WriteString(dialect.Operator_or + "(")
	}

	for i, fn := range fns {
		if i > 0 {
			u.where.WriteString(dialect.Operator_and)
		}
		cond, val := fn()
		// if v, ok := val.(error); ok {
		//	u.err = v
		//	return u
		// }
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

func (u *Updater) Cols(cols ...dialect.Field) *Updater {
	for _, col := range cols {
		u.affect = append(u.affect, col)
	}
	return u
}

// Exec
func (u *Updater) Exec(ctx context.Context) (sql.Result, error) {
	defer u.Free()

	if u.err != nil {
		return nil, u.err
	}

	lens := len(u.cols) + len(u.exprCols)
	if lens == 0 {
		return nil, dialect.ErrCreateEmpty
	}

	u.command.WriteString("UPDATE " + dialect.Quote_Char + u.table + dialect.Quote_Char + " SET ")
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
	// WHERE
	if u.where.Len() > 0 {
		u.command.WriteString(" WHERE " + u.where.String())
	}

	stmt, err := u.db.PrepareContext(ctx, u.command.String())
	if err != nil {
		return nil, err
	}
	if u.db.IsDB() {
		defer stmt.Close()
	}

	u.params = append(u.params, u.whereParams...)
	return stmt.ExecContext(ctx, u.params...)
}

// Struct
func (u *Updater) Struct(ctx context.Context, bean dialect.Modeler) (sql.Result, error) {
	defer u.Free()

	if u.err != nil {
		return nil, u.err
	}

	u.command.WriteString("UPDATE " + dialect.Quote_Char + u.table + dialect.Quote_Char + " SET ")
	cols, vals := bean.AssignValues(u.affect...)
	for i, col := range cols {
		if i > 0 {
			u.command.WriteString(",")
		}
		u.command.WriteString(col + " = ?")
	}
	u.params = append(u.params, vals...)
	//
	keys, values := bean.AssignKeys()
	u.Where(keys.Eq(values))

	// WHERE
	if u.where.Len() > 0 {
		u.command.WriteString(" WHERE " + u.where.String())
	}
	u.params = append(u.params, u.whereParams...)

	stmt, err := u.db.PrepareContext(ctx, u.command.String())
	if err != nil {
		return nil, err
	}
	if u.db.IsDB() {
		defer stmt.Close()
	}

	result, err := stmt.ExecContext(ctx, u.params...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Struct 执行更新,请不要在事务中使用
func (u *Updater) StructBatch(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error) {
	defer u.Free()

	if u.err != nil {
		return nil, u.err
	}

	lens := len(beans)
	if lens == 0 {
		return nil, dialect.ErrCreateEmpty
	}

	u.command.WriteString("UPDATE " + dialect.Quote_Char + u.table + dialect.Quote_Char + " SET ")
	cols, vals := beans[0].AssignValues(u.affect...)
	for i, col := range cols {
		if i > 0 {
			u.command.WriteString(",")
		}
		u.command.WriteString(col + " = ?")
	}
	u.params = append(u.params, vals...)
	//
	keys, values := beans[0].AssignKeys()
	u.Where(keys.Eq(values))

	// WHERE
	if u.where.Len() > 0 {
		u.command.WriteString(" WHERE " + u.where.String())
	}
	u.params = append(u.params, u.whereParams...)

	// 启动事务批量执行更新
	ret, err := u.db.Transaction(ctx, func(tx *Tx) (any, error) {
		stmt, err := tx.PrepareContext(ctx, u.command.String())
		if err != nil {
			return nil, err
		}
		if u.db.IsDB() {
			defer stmt.Close()
		}

		result, err := stmt.ExecContext(ctx, u.params...)
		if err != nil {
			return nil, err
		}

		for i := 1; i < lens; i++ {
			bean := beans[i]
			if bean == nil {
				return nil, dialect.ErrBeanEmpty
			}
			_, vals = bean.AssignValues(u.affect...)
			u.params = u.params[:0]
			u.params = append(u.params, vals...)
			//
			_, values = bean.AssignKeys()
			u.params = append(u.params, values)

			result, err = stmt.ExecContext(ctx, u.params...)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	})
	if err != nil {
		return nil, err
	}
	if result, ok := ret.(sql.Result); ok {
		return result, nil
	}
	return nil, err
}
