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
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/log"
	"strings"
)

type (
	Update struct {
		pool.Model
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
	updatePool = pool.New(app.Context, func() any {
		obj := &Update{}
		obj.UUID()
		return obj
	})
)

// Update
func NewUpdate(db Executer, tableName string) *Update {
	obj := updatePool.Get().(*Update)
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

func (u *Update) Free() {
	if u == nil || u.table == "" {
		return
	}

	if u.db.Debug() {
		log.Info(u.String())
	}

	updatePool.Put(u)
}

func (u *Update) Reset() {
	u.table = ""
	u.affect = u.affect[:0]     // []dialect.Field{} // u.affect[:0]
	u.cols = u.cols[:0]         // []dialect.Field{}   // u.cols[:0]
	u.exprCols = u.exprCols[:0] // []expr{}        // u.exprCols[:0]
	u.where.Reset()
	u.whereParams = u.whereParams[:0] // []any{} // u.whereParams[:0]
	u.command.Reset()
	u.params = u.params[:0] // []any{} // u.params[:0]
}

func (u *Update) String() string {
	if u.table == "" {
		u.commandString.WriteString(fmt.Sprintf("%s  %v \n", u.command.String(), u.params))
	}
	return u.commandString.String()
}

// Set
func (u *Update) Set(fns ...dialect.Setter) *Update {
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

func (u *Update) SetExpr(fns ...dialect.ExprSetter) *Update {
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
func (u *Update) Where(fns ...dialect.Condition) *Update {
	if len(fns) == 0 || u.err != nil {
		return u
	}

	if u.where.Len() == 0 {
		u.where.WriteString("(")
	} else {
		u.where.WriteString(dialect.Operator_and + "(")
	}
	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if cond[:len(dialect.Operator_or)] == dialect.Operator_or || cond[:len(dialect.Operator_and)] == dialect.Operator_and {
				u.where.WriteString(" ")
			} else {
				u.where.WriteString(dialect.Operator_and)
			}
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
func (u *Update) And(fns ...dialect.Condition) *Update {
	if len(fns) == 0 || u.err != nil {
		return u
	}

	if u.where.Len() == 0 {
		u.where.WriteString("(")
	} else {
		u.where.WriteString(dialect.Operator_and + "(")
	}

	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if cond[:len(dialect.Operator_or)] == dialect.Operator_or || cond[:len(dialect.Operator_and)] == dialect.Operator_and {
				u.where.WriteString(" ")
			} else {
				u.where.WriteString(dialect.Operator_or)
			}
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
func (u *Update) Or(fns ...dialect.Condition) *Update {
	if len(fns) == 0 || u.err != nil {
		return u
	}

	if u.where.Len() == 0 {
		u.where.WriteString("(")
	} else {
		u.where.WriteString(dialect.Operator_or + "(")
	}

	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if cond[:len(dialect.Operator_or)] == dialect.Operator_or || cond[:len(dialect.Operator_and)] == dialect.Operator_and {
				u.where.WriteString(" ")
			} else {
				u.where.WriteString(dialect.Operator_and)
			}
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

func (u *Update) Cols(cols ...dialect.Field) *Update {
	for _, col := range cols {
		u.affect = append(u.affect, col)
	}
	return u
}

// Exec
func (u *Update) Exec(ctx context.Context) (sql.Result, error) {
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
func (u *Update) Struct(ctx context.Context, bean dialect.Modeler) (sql.Result, error) {
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
func (u *Update) StructBatch(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error) {
	defer u.Free()

	if u.err != nil {
		return nil, u.err
	}

	lens := len(beans)
	if lens == 0 {
		return nil, dialect.ErrCreateEmpty
	}

	u.command.WriteString("UPDATE " + dialect.Quote_Char + u.table + dialect.Quote_Char + " SET ")
	cols, vals := beans[0].RawAssignValues(u.affect...)
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
			_, vals = bean.RawAssignValues(u.affect...)
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
