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
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"strings"
)

type Creater interface {
	Insert(ctx context.Context) (sql.Result, error)
	InsertStruct(ctx context.Context, bean dialect.Modeler) (sql.Result, error)
	InsertBatchStruct(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error)
}

type cc struct {
	*orm
}

// C 创建插入器
func (c *orm) C(name string) Creater {
	c.table = name
	return &cc{
		orm: c,
	}
}

// Insert 执行插入
func (c *cc) Insert(ctx context.Context) (sql.Result, error) {
	defer c.Free()

	// if c.err != nil {
	// 	return nil, c.err
	// }

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

// InsertStruct 执行插入一个结构体
func (c *cc) InsertStruct(ctx context.Context, bean dialect.Modeler) (sql.Result, error) {
	defer c.Free()

	// if c.err != nil {
	// 	return nil, c.err
	// }

	c.command.WriteString("INSERT INTO " + dialect.Quote_Char + c.table + dialect.Quote_Char + " (")

	var _cols []string
	_cols, c.params = bean.AssignValues(c.cols...)
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

// InsertBatchStruct 执行批量插入，请不要在事务中使用
func (c *cc) InsertBatchStruct(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error) {
	defer c.Free()

	// if c.err != nil {
	// 	return nil, c.err
	// }

	lens := len(beans)
	if lens == 0 {
		return nil, dialect.ErrBeanEmpty
	}

	c.command.WriteString("INSERT INTO " + dialect.Quote_Char + c.table + dialect.Quote_Char + " (")

	var _cols []string
	_cols, c.params = beans[0].RawAssignValues(c.cols...)
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
			_, c.params = bean.RawAssignValues(c.cols...)
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
