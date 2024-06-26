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

package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

const (
	Inner_Join JoinType = " INNER"
	Left_Join  JoinType = " LEFT"
	Right_Join JoinType = " RIGHT"

	operator_and = " AND "
	operator_or  = " OR "
	placeholder  = "?"
	Quote_Char   = "`"

	command_insert command = "INSERT INTO "
	command_select command = "SELECT "
	command_update command = "UPDATE "
	command_delete command = "DELETE FROM "
)

var (
	ErrCreateEmpty = fmt.Errorf("")
)

type (
	JoinType string
	command  string
	// expr represents an SQL express
	expr struct {
		colName string
		arg     interface{}
	}
	condition struct {
		where       strings.Builder
		whereParams []interface{}
	}

	Creator struct {
		command command
		table   string
		cols    []string
		params  []interface{}
	}
	Selector struct {
		command    command
		table      string
		join       [][3]string
		distinct   bool
		cols       []string
		omit       []interface{}
		groupBy    strings.Builder
		having     strings.Builder
		orderBy    strings.Builder
		limit      string
		limitSize  int
		limitStart int

		condition
	}
	Updater struct {
		command  command
		table    string
		cols     []string
		params   []interface{}
		incrCols []expr
		decrCols []expr
		exprCols []expr

		condition
	}
	Deleter struct {
		command command
		table   string

		condition
	}
)

var (
	createPool = sync.Pool{
		New: func() interface{} {
			return &Creator{
				command: command_insert,
			}
		},
	}
	selectPool = sync.Pool{
		New: func() interface{} {
			return &Selector{
				command: command_select,
			}
		},
	}
	updatePool = sync.Pool{
		New: func() interface{} {
			return &Updater{
				command: command_update,
			}
		},
	}
	deletePool = sync.Pool{
		New: func() interface{} {
			return &Deleter{
				command: command_delete,
			}
		},
	}
)

// ////////////////////////////////////////
// Creator
func NewCreate(table string) *Creator {
	obj := createPool.Get().(*Creator)
	obj.table = table
	return obj
}

func (c *Creator) Free() {
	c.table = ""
	c.cols = c.cols[:]
	c.params = c.params[:]
	createPool.Put(c)
}

// Set
func (c *Creator) Set(fn Setter) *Creator {
	s, val := fn()
	c.cols = append(c.cols, s)
	c.params = append(c.params, val)
	return c
}

// Sets
func (c *Creator) Sets(fns ...Setter) *Creator {
	for _, fn := range fns {
		s, val := fn()
		c.cols = append(c.cols, s)
		c.params = append(c.params, val)
	}
	return c
}

// Do
func (c *Creator) Do(ctx context.Context) (sql.Result, error) {
	if len(c.cols) == 0 {
		return nil, ErrCreateEmpty
	}
	fmt.Println(c.command, c.table, c.cols, c.params)
	return nil, nil
}

// /////////////////////////////////////////////////
// Updater
func NewUpdate(table string) *Updater {
	obj := updatePool.Get().(*Updater)
	obj.table = table
	return obj

}

func (u *Updater) Free() {
	u.table = ""
	u.cols = u.cols[:]
	u.params = u.params[:]
	u.incrCols = u.incrCols[:]
	u.decrCols = u.decrCols[:]
	u.exprCols = u.exprCols[:]
	u.condition.where.Reset()
	u.condition.whereParams = u.condition.whereParams[:]
	// u.AndOr = false
	updatePool.Put(u)
}

// Set
func (c *Updater) Set(fns ...Setter) *Updater {
	for _, fn := range fns {
		s, val := fn()
		c.cols = append(c.cols, s)
		c.params = append(c.params, val)
	}
	return c
}

// Where
func (c *Updater) Where(fns ...Condition) *Updater {
	if len(fns) == 0 {
		return c
	}
	c.condition.where.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			c.condition.where.WriteString(operator_and)
		}
		cond, val := fn()
		c.condition.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			c.condition.whereParams = append(c.condition.whereParams, vals...)
		} else {
			c.condition.whereParams = append(c.condition.whereParams, val)
		}
	}
	c.condition.where.WriteString(")")

	return c
}

// And
func (c *Updater) AndOr(fns ...Condition) *Updater {
	if len(fns) == 0 {
		return c
	}
	c.condition.where.WriteString(operator_and + "(")
	for i, fn := range fns {
		if i > 0 {
			c.condition.where.WriteString(operator_or)
		}
		cond, val := fn()
		c.condition.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			c.condition.whereParams = append(c.condition.whereParams, vals...)
		} else {
			c.condition.whereParams = append(c.condition.whereParams, val)
		}
	}
	c.condition.where.WriteString(")")
	return c
}

// Or
func (c *Updater) OrAnd(fns ...Condition) *Updater {
	if len(fns) == 0 {
		return c
	}
	c.condition.where.WriteString(operator_or + "(")
	for i, fn := range fns {
		if i > 0 {
			c.condition.where.WriteString(operator_and)
		}
		cond, val := fn()
		c.condition.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			c.condition.whereParams = append(c.condition.whereParams, vals...)
		} else {
			c.condition.whereParams = append(c.condition.whereParams, val)
		}
	}
	c.condition.where.WriteString(")")
	return c
}

// Do
func (c *Updater) Do(ctx context.Context) (sql.Result, error) {
	return nil, nil
}

// //////////////////////////////////////////////////
// Deleter
func NewDelete(table string) *Deleter {
	obj := deletePool.Get().(*Deleter)
	obj.table = table
	return obj

}

func (d *Deleter) Free() {
	deletePool.Put(d)
}

// Do
func (c *Deleter) Do(ctx context.Context) (sql.Result, error) {
	return nil, nil
}

// ///////////////////////////////////////////////
// Selector
func NewSelect(table string) *Selector {
	obj := selectPool.Get().(*Selector)
	obj.table = table
	return obj

}

func (s *Selector) Free() {
	s.table = ""
	s.cols = s.cols[:]
	s.distinct = false
	s.join = s.join[:]
	s.omit = s.omit[:]
	s.condition.where.Reset()
	s.condition.whereParams = s.condition.whereParams[:]
	// s.AndOr = false
	s.groupBy.Reset()
	s.having.Reset()
	s.orderBy.Reset()
	s.limit = ""
	s.limitSize = 0
	s.limitStart = 0
	selectPool.Put(s)
}

// distinct
func (c *Selector) Distinct(cols ...string) *Selector {
	c.distinct = true
	return c
}

// cols
func (c *Selector) Cols(cols ...string) *Selector {
	c.cols = append(c.cols, cols...)
	return c
}

// join
func (c *Selector) Join(joinType JoinType, left, right Field, fns ...Condition) *Selector {
	var on strings.Builder
	for _, fn := range fns {
		on.WriteString(operator_and)
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
		left.FieldName() + "=" + right.FieldName() + on.String(),
	})
	return c
}

func (c *Selector) LeftJoin(left, right Field, fns ...Condition) *Selector {
	return c.Join(Left_Join, left, right, fns...)
}
func (c *Selector) RightJoin(left, right Field, fns ...Condition) *Selector {
	return c.Join(Right_Join, left, right, fns...)
}

// Where
func (c *Selector) Where(fns ...Condition) *Selector {
	if len(fns) == 0 {
		return c
	}
	c.condition.where.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			c.condition.where.WriteString(operator_and)
		}
		cond, val := fn()
		c.condition.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			c.condition.whereParams = append(c.condition.whereParams, vals...)
		} else {
			c.condition.whereParams = append(c.condition.whereParams, val)
		}
	}
	c.condition.where.WriteString(")")

	return c
}

// And
func (c *Selector) AndOr(fns ...Condition) *Selector {
	if len(fns) == 0 {
		return c
	}
	c.condition.where.WriteString(operator_and + "(")
	for i, fn := range fns {
		if i > 0 {
			c.condition.where.WriteString(operator_or)
		}
		cond, val := fn()
		c.condition.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			c.condition.whereParams = append(c.condition.whereParams, vals...)
		} else {
			c.condition.whereParams = append(c.condition.whereParams, val)
		}
	}
	c.condition.where.WriteString(")")
	return c
}

// Or
func (c *Selector) OrAnd(fns ...Condition) *Selector {
	if len(fns) == 0 {
		return c
	}
	c.condition.where.WriteString(operator_or + "(")
	for i, fn := range fns {
		if i > 0 {
			c.condition.where.WriteString(operator_and)
		}
		cond, val := fn()
		c.condition.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			c.condition.whereParams = append(c.condition.whereParams, vals...)
		} else {
			c.condition.whereParams = append(c.condition.whereParams, val)
		}
	}
	c.condition.where.WriteString(")")
	return c
}

// Do
func (s *Selector) Do(ctx context.Context) (sql.Result, error) {
	return nil, nil
}
