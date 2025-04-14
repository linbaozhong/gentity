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

package orm

import (
	"context"
	"database/sql"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"strings"
)

type UpdateBuilder interface {
	Columner
	Set(fns ...dialect.Setter) Builder
	SetExpr(fns ...dialect.ExprSetter) Builder
	Wherer
	Update(x ...Executer) Updater
}

type Updater interface {
	Exec(ctx context.Context) (sql.Result, error)
	Struct(ctx context.Context, bean dialect.Modeler) (sql.Result, error)
	BatchStruct(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error)
}

type update struct {
	*orm
}

type expr struct {
	colName string
	arg     any
}

// Update 更新器
func (o *orm) Update(x ...Executer) Updater {
	o.connect(x...)
	return &update{
		orm: o,
	}
}

// Exec 执行更新
func (u *update) Exec(ctx context.Context) (sql.Result, error) {
	defer u.Free()

	// if u.err != nil {
	// 	return nil, u.err
	// }

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
	return nil, nil
	return stmt.ExecContext(ctx, u.params...)
}

// Struct 更新一个结构体
func (u *update) Struct(ctx context.Context, bean dialect.Modeler) (sql.Result, error) {
	defer u.Free()

	// if u.err != nil {
	// 	return nil, u.err
	// }

	u.command.WriteString("UPDATE " + dialect.Quote_Char + u.table + dialect.Quote_Char + " SET ")
	cols, vals := bean.AssignValues(u.cols...)
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

// BatchStruct 执行批量更新,请不要在事务中使用
func (u *update) BatchStruct(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error) {
	defer u.Free()

	lens := len(beans)
	if lens == 0 {
		return nil, dialect.ErrCreateEmpty
	}

	u.command.WriteString("UPDATE " + dialect.Quote_Char + u.table + dialect.Quote_Char + " SET ")
	cols, vals := beans[0].RawAssignValues(u.cols...)
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
			_, vals = bean.RawAssignValues(u.cols...)
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
