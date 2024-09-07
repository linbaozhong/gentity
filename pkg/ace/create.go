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
	"fmt"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/ace/types"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/sql"
	"strings"
	"sync"
)

type (
	Creator struct {
		db            types.Executer
		table         string
		affect        []types.Field
		cols          []types.Field
		params        []any
		command       strings.Builder
		commandString strings.Builder
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
func NewCreate(db types.Executer, tableName string) *Creator {
	if db == nil || tableName == "" {
		panic("db or table is nil")
		return nil
	}
	obj := createPool.Get().(*Creator)
	obj.db = db
	obj.table = tableName

	return obj
}

func (c *Creator) Free() {
	if c == nil {
		return
	}
	c.commandString.Reset()
	c.commandString.WriteString(fmt.Sprintf("%s  %v", c.command.String(), c.params))

	if c.db.Debug() {
		log.Info(c.String())
	}
	c.table = ""
	c.affect = c.affect[:]
	c.cols = c.cols[:]
	c.command.Reset()
	c.params = c.params[:]
	createPool.Put(c)
}

func (c *Creator) String() string {
	return c.commandString.String()
}

// Sets
func (c *Creator) Set(fns ...types.Setter) *Creator {
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

func (c *Creator) Cols(cols ...types.Field) *Creator {
	for _, col := range cols {
		c.affect = append(c.affect, col)
	}
	return c
}

// Do
func (c *Creator) Do(ctx context.Context) (sql.Result, error) {
	defer c.Free()
	lens := len(c.cols)
	if lens == 0 {
		return nil, types.ErrCreateEmpty
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

	// fmt.Println(c.command.String(), c.params)
	return c.db.ExecContext(ctx, c.command.String(), c.params...)
}

// Struct
func (c *Creator) Struct(ctx context.Context, beans ...types.Modeler) (sql.Result, error) {
	defer c.Free()

	lens := len(beans)
	if lens == 0 || lens > 100 || beans[0] == nil {
		return nil, types.ErrBeanEmpty
	}

	c.command.WriteString("INSERT INTO " + dialect.Quote_Char + c.table + dialect.Quote_Char + " (")

	_cols, _vals := beans[0].AssignValues(c.affect...)
	_colLens := len(_cols)
	c.command.WriteString(strings.Join(_cols, ","))
	c.command.WriteString(") VALUES ")
	c.params = append(c.params, _vals...)
	c.command.WriteString("(" + strings.Repeat("?,", _colLens)[:_colLens*2-1] + ")")

	for i := 1; i < lens; i++ {
		bean := beans[i]
		if bean == nil {
			return nil, types.ErrBeanEmpty
		}
		c.command.WriteString(",")
		_, _vals = bean.AssignValues(c.affect...)
		c.params = append(c.params, _vals...)
		c.command.WriteString("(" + strings.Repeat("?,", _colLens)[:_colLens*2-1] + ")")
	}
	// fmt.Println(c.command.String(), c.params)
	return c.db.ExecContext(ctx, c.command.String(), c.params...)
}
