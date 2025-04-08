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
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/log"
	"strings"
)

type (
	Creater interface {
		Sets(fns ...dialect.Setter) Creater
		Cols(cols ...dialect.Field) Creater
		Exec(ctx context.Context) (sql.Result, error)
		Struct(ctx context.Context, bean dialect.Modeler) (sql.Result, error)
		BatchStruct(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error)
	}
	CreateBuilder interface {
		Creater
		Table(name any) CreateBuilder
		Free()
		Reset()
		String() string
	}
	create struct {
		pool.Model
		db            Executer
		table         string
		affect        []dialect.Field
		cols          []dialect.Field
		params        []any
		command       strings.Builder
		commandString strings.Builder
	}
)

var (
	createPool = pool.New(app.Context, func() any {
		obj := &create{}
		obj.UUID()
		return obj
	})
)

// create
func newCreate(dbs ...Executer) CreateBuilder {
	obj := createPool.Get().(*create)
	obj.db = GetExec(dbs...)
	obj.commandString.Reset()

	return obj
}

func (c *create) Free() {
	if c == nil || c.table == "" {
		return
	}

	_ = c.String()
	if c.db.Debug() {
		log.Info(c.String())
	}

	createPool.Put(c)
}

func (c *create) Reset() {
	c.table = ""
	c.affect = c.affect[:0] // []dialect.Field{} // c.affect[:0]
	c.cols = c.cols[:0]     // []dialect.Field{}   // c.cols[:0]
	c.command.Reset()
	c.params = c.params[:0] // []any{} // c.params[:0]
}

func (c *create) String() string {
	if c.commandString.Len() == 0 {
		c.commandString.WriteString(fmt.Sprintf("%s  %v \n", c.command.String(), c.params))
	}
	return c.commandString.String()
}

// setTableName 设置表名
func (c *create) Table(n any) CreateBuilder {
	setTableName(&c.table, n)
	return c
}

// Sets 设置列名和值
func (c *create) Sets(fns ...dialect.Setter) Creater {
	if len(fns) == 0 {
		return c
	}

	for _, fn := range fns {
		if fn == nil {
			continue
		}
		s, val := fn()
		c.cols = append(c.cols, s)
		c.params = append(c.params, val)
	}
	return c
}

// Cols 设置列名
func (c *create) Cols(cols ...dialect.Field) Creater {
	for _, col := range cols {
		c.affect = append(c.affect, col)
	}
	return c
}

// Exec
func (c *create) Exec(ctx context.Context) (sql.Result, error) {
	defer c.Free()

	lens := len(c.cols)
	if lens == 0 {
		return nil, dialect.ErrCreateEmpty
	}

	c.command.WriteString("INSERT INTO " + dialect.Quote_Char + c.table + dialect.Quote_Char + " (")
	for i, col := range c.cols {
		if i > 0 {
			c.command.WriteString(",")
		}
		c.command.WriteString(col.Quote())
	}
	c.command.WriteString(") VALUES ")
	c.command.WriteString("(" + strings.Repeat("?,", lens)[:lens*2-1] + ")")

	stmt, err := c.db.PrepareContext(ctx, c.command.String())
	if err != nil {
		return nil, err
	}
	if c.db.IsDB() {
		defer stmt.Close()
	}

	return stmt.ExecContext(ctx, c.params...)
}

// Struct 执行插入，请不要在事务中使用; bean 必须是指针类型
func (c *create) Struct(ctx context.Context, bean dialect.Modeler) (sql.Result, error) {
	defer c.Free()

	c.command.WriteString("INSERT INTO " + dialect.Quote_Char + c.table + dialect.Quote_Char + " (")

	var _cols []string
	_cols, c.params = bean.AssignValues(c.affect...)
	_colLens := len(_cols)
	c.command.WriteString(strings.Join(_cols, ","))
	c.command.WriteString(") VALUES ")
	c.command.WriteString("(" + strings.Repeat("?,", _colLens)[:_colLens*2-1] + ")")

	stmt, err := c.db.PrepareContext(ctx, c.command.String())
	if err != nil {
		return nil, err
	}

	result, err := stmt.ExecContext(ctx, c.params...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// BatchStruct 执行批量插入，请不要在事务中使用
func (c *create) BatchStruct(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error) {
	defer c.Free()

	lens := len(beans)
	if lens == 0 {
		return nil, dialect.ErrBeanEmpty
	}

	c.command.WriteString("INSERT INTO " + dialect.Quote_Char + c.table + dialect.Quote_Char + " (")

	var _cols []string
	_cols, c.params = beans[0].RawAssignValues(c.affect...)
	_colLens := len(_cols)
	c.command.WriteString(strings.Join(_cols, ","))
	c.command.WriteString(") VALUES ")
	c.command.WriteString("(" + strings.Repeat("?,", _colLens)[:_colLens*2-1] + ")")

	// 启动事务批量执行Create
	ret, err := c.db.Transaction(ctx, func(tx *Tx) (any, error) {
		stmt, err := tx.PrepareContext(ctx, c.command.String())
		if err != nil {
			return nil, err
		}
		if c.db.IsDB() {
			defer stmt.Close()
		}

		result, err := stmt.ExecContext(ctx, c.params...)
		if err != nil {
			return nil, err
		}

		for i := 1; i < lens; i++ {
			bean := beans[i]
			if bean == nil {
				return nil, dialect.ErrBeanEmpty
			}
			_, c.params = bean.RawAssignValues(c.affect...)
			result, err = stmt.ExecContext(ctx, c.params...)
			if err != nil {
				return nil, err
			}
			bean.AssignPrimaryKeyValues(result)
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
