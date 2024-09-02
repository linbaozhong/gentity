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
	"strconv"
	"strings"
	"sync"
)

type (
	Selector struct {
		db           types.Executer
		table        string
		join         [][3]string
		distinct     bool
		cols         []types.Field
		funcs        []string
		omit         []any
		groupBy      strings.Builder
		having       strings.Builder
		havingParams []any
		orderBy      strings.Builder
		limit        string
		where        strings.Builder
		whereParams  []any
		command      strings.Builder
	}
)

var (
	selectPool = sync.Pool{
		New: func() any {
			obj := &Selector{}
			return obj
		},
	}
)

// Selector
func NewSelect(db types.Executer, tableName string) *Selector {
	if db == nil || tableName == "" {
		panic("db or table is nil")
		return nil
	}
	obj := selectPool.Get().(*Selector)
	obj.db = db
	obj.table = tableName
	obj.command.Reset()
	return obj
}

func (s *Selector) Free() {
	s.table = ""
	s.cols = s.cols[:]
	s.funcs = s.funcs[:]
	s.distinct = false
	s.join = s.join[:]
	s.omit = s.omit[:]
	s.where.Reset()
	s.whereParams = s.whereParams[:]
	s.groupBy.Reset()
	s.having.Reset()
	s.havingParams = s.havingParams[:]
	s.orderBy.Reset()
	s.limit = ""
	selectPool.Put(s)
}

func (s *Selector) String() string {
	return fmt.Sprintf("%s  %v", s.command.String(), s.whereParams)
}

// distinct
func (c *Selector) Distinct(cols ...types.Field) *Selector {
	c.distinct = true
	for _, col := range cols {
		c.cols = append(c.cols, col)
	}

	return c
}

// cols
func (c *Selector) Cols(cols ...types.Field) *Selector {
	for _, col := range cols {
		c.cols = append(c.cols, col)
	}
	return c
}

// funcs
func (c *Selector) Funcs(fns ...types.Function) *Selector {
	for _, fn := range fns {
		c.funcs = append(c.funcs, fn())
	}
	return c
}

// join
func (c *Selector) Join(joinType types.JoinType, left, right types.Field, fns ...types.Condition) *Selector {
	var on strings.Builder
	for _, fn := range fns {
		on.WriteString(types.Operator_and)
		cond, val := fn()
		on.WriteString(cond)
		if vals, ok := val.([]any); ok {
			c.whereParams = append(c.whereParams, vals...)
		} else {
			c.whereParams = append(c.whereParams, val)
		}
	}
	c.join = append(c.join, [3]string{
		string(joinType),
		right.TableName(),
		left.Quote() + "=" + right.Quote() + on.String(),
	})
	return c
}

func (c *Selector) LeftJoin(left, right types.Field, fns ...types.Condition) *Selector {
	return c.Join(types.Left_Join, left, right, fns...)
}
func (c *Selector) RightJoin(left, right types.Field, fns ...types.Condition) *Selector {
	return c.Join(types.Right_Join, left, right, fns...)
}

// Where
func (c *Selector) Where(fns ...types.Condition) *Selector {
	if len(fns) == 0 {
		return c
	}
	if c.where.Len() == 0 {
		c.where.WriteString("(")
	} else {
		c.where.WriteString(types.Operator_and + "(")
	}
	for i, fn := range fns {
		if i > 0 {
			c.where.WriteString(types.Operator_and)
		}
		cond, val := fn()
		c.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			c.whereParams = append(c.whereParams, vals...)
		} else {
			c.whereParams = append(c.whereParams, val)
		}
	}
	c.where.WriteString(")")

	return c
}

// And
func (c *Selector) And(fns ...types.Condition) *Selector {
	if len(fns) == 0 {
		return c
	}

	if c.where.Len() == 0 {
		c.where.WriteString("(")
	} else {
		c.where.WriteString(types.Operator_and + "(")
	}

	for i, fn := range fns {
		if i > 0 {
			c.where.WriteString(types.Operator_or)
		}
		cond, val := fn()
		c.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			c.whereParams = append(c.whereParams, vals...)
		} else {
			c.whereParams = append(c.whereParams, val)
		}
	}
	c.where.WriteString(")")
	return c
}

// Or
func (c *Selector) Or(fns ...types.Condition) *Selector {
	if len(fns) == 0 {
		return c
	}

	if c.where.Len() == 0 {
		c.where.WriteString("(")
	} else {
		c.where.WriteString(types.Operator_or + "(")
	}

	for i, fn := range fns {
		if i > 0 {
			c.where.WriteString(types.Operator_and)
		}
		cond, val := fn()
		c.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			c.whereParams = append(c.whereParams, vals...)
		} else {
			c.whereParams = append(c.whereParams, val)
		}
	}
	c.where.WriteString(")")
	return c
}

// Order
func (c *Selector) Order(cols ...types.Field) *Selector {
	return c.Asc(cols...)
}

// Order Asc
func (c *Selector) Asc(cols ...types.Field) *Selector {
	if len(cols) == 0 {
		return c
	}
	for _, col := range cols {
		if c.orderBy.Len() > 0 {
			c.orderBy.WriteByte(',')
		}
		c.orderBy.WriteString(col.Quote())
	}
	return c
}

// Order Desc
func (c *Selector) Desc(cols ...types.Field) *Selector {
	if len(cols) == 0 {
		return c
	}
	for _, col := range cols {
		if c.orderBy.Len() > 0 {
			c.orderBy.WriteByte(',')
		}
		c.orderBy.WriteString(col.Quote() + " DESC")
	}
	return c
}

// Group
func (c *Selector) Group(cols ...types.Field) *Selector {
	if len(cols) == 0 {
		return c
	}
	for _, col := range cols {
		if c.groupBy.Len() > 0 {
			c.groupBy.WriteByte(',')
		}
		c.groupBy.WriteString(col.Quote())
	}
	return c
}

// Group Having
func (c *Selector) Having(fns ...types.Condition) *Selector {
	if len(fns) == 0 {
		return c
	}
	c.having.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			c.having.WriteString(types.Operator_and)
		}
		cond, val := fn()
		c.having.WriteString(cond)
		if vals, ok := val.([]any); ok {
			c.havingParams = append(c.havingParams, vals...)
		}
	}
	return c
}
func (c *Selector) Limit(size uint, start ...uint) *Selector {
	if len(start) > 0 {
		c.limit = " LIMIT " + strconv.Itoa(int(size)) + " OFFSET " + strconv.Itoa(int(start[0]))
	} else {
		c.limit = " LIMIT " + strconv.Itoa(int(size))
	}

	return c
}

func (c *Selector) stmt() {
	c.command.WriteString("SELECT ")

	colens := len(c.cols)
	funlens := len(c.funcs)
	if colens+funlens == 0 {
		c.command.WriteString("*")
	} else {
		if c.distinct {
			c.command.WriteString("DISTINCT ")
		}
		for i, col := range c.cols {
			if i > 0 {
				c.command.WriteString(",")
			}
			c.command.WriteString(col.Quote())
		}
		if colens > 0 && funlens > 0 {
			c.command.WriteString(",")
		}
		c.command.WriteString(strings.Join(c.funcs, ","))
	}
	// FROM TABLE
	c.command.WriteString(" FROM " + types.Quote_Char + c.table + types.Quote_Char)
	for _, j := range c.join {
		c.command.WriteString(j[0] + " JOIN " + j[1] + " ON " + j[2] + " ")
	}
	// WHERE
	if c.where.Len() > 0 {
		c.command.WriteString(" WHERE " + c.where.String())
	}
	// GROUP BY
	if c.groupBy.Len() > 0 {
		c.command.WriteString(" GROUP BY " + c.groupBy.String())
		// HAVING
		if c.having.Len() > 0 {
			c.command.WriteString(" HAVING " + c.having.String())
		}
	}
	// ORDER BY
	if c.orderBy.Len() > 0 {
		c.command.WriteString(" ORDER BY " + c.orderBy.String())
	}

	// LIMIT
	if c.limit != "" {
		c.command.WriteString(c.limit)
	}
}

// Query
func (c *Selector) Query(ctx context.Context) (*sql.Rows, error) {
	c.stmt()
	return c.db.QueryContext(ctx, c.command.String(), c.whereParams...)
}

// Count
func (c *Selector) Count(ctx context.Context, cond ...types.Condition) (int64, error) {
	c.Where(cond...)
	c.command.WriteString("SELECT COUNT(*)")
	// FROM TABLE
	c.command.WriteString(" FROM " + types.Quote_Char + c.table + types.Quote_Char)
	// WHERE
	if c.where.Len() > 0 {
		c.command.WriteString(" WHERE " + c.where.String())
	}

	rows, err := c.db.QueryContext(ctx, c.command.String(), c.whereParams...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Next() {
		var count int64
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
	return 0, nil
}

// Sum
func (c *Selector) Sum(ctx context.Context, col types.Field, cond ...types.Condition) (int64, error) {
	c.Funcs(col.Sum()).Where(cond...)
	c.command.WriteString("SELECT ")
	c.command.WriteString(c.funcs[0])
	// FROM TABLE
	c.command.WriteString(" FROM " + types.Quote_Char + c.table + types.Quote_Char)
	// WHERE
	if c.where.Len() > 0 {
		c.command.WriteString(" WHERE " + c.where.String())
	}

	rows, err := c.db.QueryContext(ctx, c.command.String(), c.whereParams...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Next() {
		var n int64
		err := rows.Scan(&n)
		if err != nil {
			return 0, err
		}
		return n, nil
	}
	return 0, nil
}
