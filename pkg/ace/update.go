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
	"github.com/linbaozhong/gentity/pkg/ace/types"
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
	if u == nil {
		return
	}

	u.commandString.WriteString(fmt.Sprintf("%s  %v", u.command.String(), u.params))

	if u.db.Debug() {
		log.Info(u.String())
	}
	u.table = ""
	u.affect = u.affect[:]
	u.cols = u.cols[:]
	u.exprCols = u.exprCols[:]
	u.where.Reset()
	u.whereParams = u.whereParams[:]
	u.command.Reset()
	u.params = u.params[:]

	updatePool.Put(u)
}

func (u *Updater) String() string {
	if u.table == "" {
		return u.commandString.String()
	}
	return fmt.Sprintf("%s  %v", u.command.String(), u.params)
}

// Set
func (u *Updater) Set(fns ...dialect.Setter) *Updater {
	if len(fns) == 0 || u.err != nil {
		return u
	}

	for _, fn := range fns {
		s, val := fn()
		if v, ok := val.(error); ok {
			u.err = v
			return u
		}
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
		if v, ok := val.(error); ok {
			u.err = v
			return u
		}
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
		u.where.WriteString(types.Operator_and + "(")
	}
	for i, fn := range fns {
		if i > 0 {
			u.where.WriteString(types.Operator_and)
		}
		cond, val := fn()
		if v, ok := val.(error); ok {
			u.err = v
			return u
		}
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
		u.where.WriteString(types.Operator_and + "(")
	}

	for i, fn := range fns {
		if i > 0 {
			u.where.WriteString(types.Operator_or)
		}
		cond, val := fn()
		if v, ok := val.(error); ok {
			u.err = v
			return u
		}
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
		u.where.WriteString(types.Operator_or + "(")
	}

	for i, fn := range fns {
		if i > 0 {
			u.where.WriteString(types.Operator_and)
		}
		cond, val := fn()
		if v, ok := val.(error); ok {
			u.err = v
			return u
		}
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

// Do
func (u *Updater) Do(ctx context.Context) (sql.Result, error) {
	defer u.Free()

	lens := len(u.cols) + len(u.exprCols)
	if lens == 0 {
		return nil, types.ErrCreateEmpty
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
	defer stmt.Close()

	u.params = append(u.params, u.whereParams...)
	return stmt.ExecContext(ctx, u.params...)
}

// Struct
func (u *Updater) Struct(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error) {
	defer u.Free()

	lens := len(beans)
	if lens == 0 {
		return nil, types.ErrCreateEmpty
	}

	u.command.WriteString("UPDATE " + dialect.Quote_Char + u.table + dialect.Quote_Char + " SET ")
	cols, params := beans[0].AssignValues(u.affect...)
	for i, col := range cols {
		if i > 0 {
			u.command.WriteString(",")
		}
		u.command.WriteString(col + " = ?")
	}
	u.params = append(u.params, params...)
	//
	keys, values := beans[0].AssignKeys()
	for i := 0; i < len(keys); i++ {
		u.Where(keys[i].Eq(values[i]))
	}
	// WHERE
	if u.where.Len() > 0 {
		u.command.WriteString(" WHERE " + u.where.String())
	}
	u.params = append(u.params, u.whereParams...)
	//
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.PrepareContext(ctx, u.command.String())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, u.params...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := 1; i < lens; i++ {
		bean := beans[i]
		if bean == nil {
			return nil, types.ErrBeanEmpty
		}
		_, params = bean.AssignValues(u.affect...)
		u.params = u.params[:]
		u.params = append(u.params, params...)
		//
		_, values = bean.AssignKeys()
		u.params = append(u.params, values...)

		result, err = stmt.ExecContext(ctx, u.params...)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

	}
	err = tx.Commit()
	return result, err

	//for n, bean := range beans {
	//	if n > 0 {
	//		u.command.WriteString(";")
	//	}
	//	cols, params := bean.AssignValues(u.affect...)
	//	u.params = append(u.params, params...)
	//	//
	//	keys, values := bean.AssignKeys()
	//	for i := 0; i < len(keys); i++ {
	//		u.Where(keys[i].Eq(values[i]))
	//	}
	//	for i, col := range cols {
	//		if i > 0 {
	//			u.command.WriteString(",")
	//		}
	//		u.command.WriteString(col + " = ?")
	//	}
	//	// WHERE
	//	if u.where.Len() > 0 {
	//		u.command.WriteString(" WHERE " + u.where.String())
	//	}
	//	u.params = append(u.params, u.whereParams...)
	//}
	//return u.db.ExecContext(ctx, u.command.String(), u.params...)
}
