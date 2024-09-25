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
	Creator struct {
		db            Executer
		table         string
		affect        []dialect.Field
		cols          []dialect.Field
		params        []any
		command       strings.Builder
		commandString strings.Builder
		err           error
		inPool        bool //是否在池中
	}
)

var (
	createPool = sync.Pool{
		New: func() any {
			obj := &Creator{}
			return obj
		},
	}
)

// Creator
func newCreate(db Executer, tableName string) *Creator {
	obj := createPool.Get().(*Creator)
	if db == nil || tableName == "" {
		obj.err = errors.New("db or table is nil")
		return obj
	}
	obj.inPool = false
	obj.db = db
	obj.table = tableName
	obj.err = nil
	obj.commandString.Reset()

	return obj
}

func (c *Creator) Free() {
	if c == nil || c.inPool {
		return
	}

	_ = c.String()
	if c.db.Debug() {
		log.Info(c.String())
	}

	c.inPool = true
	c.table = ""
	c.affect = c.affect[:0]
	c.cols = c.cols[:0]
	c.command.Reset()
	c.params = c.params[:0]

	createPool.Put(c)
}

func (c *Creator) String() string {
	if c.commandString.Len() == 0 {
		c.commandString.WriteString(fmt.Sprintf("%s  %v \n", c.command.String(), c.params))
	}
	return c.commandString.String()
}

// Sets
func (c *Creator) Set(fns ...dialect.Setter) *Creator {
	if len(fns) == 0 || c.err != nil {
		return c
	}

	for _, fn := range fns {
		if fn == nil {
			continue
		}
		s, val := fn()
		// if v, ok := val.(error); ok {
		//	c.err = v
		//	return c
		// }
		c.cols = append(c.cols, s)
		c.params = append(c.params, val)
	}
	return c
}

func (c *Creator) Cols(cols ...dialect.Field) *Creator {
	for _, col := range cols {
		c.affect = append(c.affect, col)
	}
	return c
}

// Exec
func (c *Creator) Exec(ctx context.Context) (sql.Result, error) {
	defer c.Free()

	if c.err != nil {
		return nil, c.err
	}

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
	defer stmt.Close()

	// fmt.Println(c.command.String(), c.params)
	// return c.db.ExecContext(ctx, c.command.String(), c.params...)
	return stmt.ExecContext(ctx, c.params...)
}

// Struct
func (c *Creator) Struct(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error) {
	defer c.Free()

	if c.err != nil {
		return nil, c.err
	}

	lens := len(beans)
	if lens == 0 {
		return nil, dialect.ErrBeanEmpty
	}

	c.command.WriteString("INSERT INTO " + dialect.Quote_Char + c.table + dialect.Quote_Char + " (")

	_cols, _vals := beans[0].AssignValues(c.affect...)
	_colLens := len(_cols)
	c.command.WriteString(strings.Join(_cols, ","))
	c.command.WriteString(") VALUES ")
	// c.params = append(c.params, _vals...)
	c.command.WriteString("(" + strings.Repeat("?,", _colLens)[:_colLens*2-1] + ")")

	//
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.PrepareContext(ctx, c.command.String())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	c.params = _vals
	result, err := stmt.ExecContext(ctx, _vals...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := 1; i < lens; i++ {
		bean := beans[i]
		if bean == nil {
			return nil, dialect.ErrBeanEmpty
		}
		// c.command.WriteString(",")
		_, _vals = bean.AssignValues(c.affect...)
		// c.params = append(c.params, _vals...)
		// c.command.WriteString("(" + strings.Repeat("?,", _colLens)[:_colLens*2-1] + ")")
		c.params = _vals
		result, err = stmt.ExecContext(ctx, _vals...)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		bean.AssignPrimaryKeyValues(result)
	}
	err = tx.Commit()
	// fmt.Println(c.command.String(), c.params)
	// return c.db.ExecContext(ctx, c.command.String(), c.params...)
	return result, err
}
