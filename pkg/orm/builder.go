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

	Creator struct {
		db      ExtContext
		table   string
		cols    []string
		params  []interface{}
		command strings.Builder
	}
	Selector struct {
		db          ExtContext
		table       string
		join        [][3]string
		distinct    bool
		cols        []string
		omit        []interface{}
		groupBy     strings.Builder
		having      strings.Builder
		orderBy     strings.Builder
		limit       string
		limitSize   int
		limitStart  int
		where       strings.Builder
		whereParams []interface{}
		command     strings.Builder
	}
	Updater struct {
		db          ExtContext
		table       string
		cols        []string
		params      []interface{}
		incrCols    []expr
		decrCols    []expr
		exprCols    []expr
		where       strings.Builder
		whereParams []interface{}
		command     strings.Builder
	}
	Deleter struct {
		db          ExtContext
		table       string
		where       strings.Builder
		whereParams []interface{}
		command     strings.Builder
	}
)

var (
	createPool = sync.Pool{
		New: func() any {
			obj := &Creator{}
			return obj
		},
	}
	selectPool = sync.Pool{
		New: func() any {
			obj := &Selector{}
			return obj
		},
	}
	updatePool = sync.Pool{
		New: func() any {
			obj := &Updater{}
			return obj
		},
	}
	deletePool = sync.Pool{
		New: func() interface{} {
			obj := &Deleter{}
			return obj
		},
	}
)

// ////////////////////////////////////////
// Creator
func NewCreate(db ExtContext, table string) *Creator {
	if db == nil || table == "" {
		panic("db or table is nil")
		return nil
	}
	obj := createPool.Get().(*Creator)
	obj.db = db
	obj.table = table
	obj.command.Reset()
	return obj
}

func (c *Creator) Free() {
	c.table = ""
	c.cols = c.cols[:]
	c.params = c.params[:]
	createPool.Put(c)
}

// Sets
func (c *Creator) Set(fns ...Setter) *Creator {
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
	c.command.WriteString("INSERT INTO " + c.table + " (")
	c.command.WriteString(strings.Join(c.cols, ",") + ") VALUES (")
	c.command.WriteString(strings.Repeat("?,", len(c.cols))[:len(c.cols)*2-1])
	c.command.WriteString(")")
	fmt.Println(c.command.String(), c.params)
	return nil, nil

	return c.db.ExecContext(ctx, c.command.String(), c.params...)
}

// /////////////////////////////////////////////////
// Updater
func NewUpdate(db ExtContext, table string) *Updater {
	if db == nil || table == "" {
		panic("db or table is nil")
		return nil
	}
	obj := updatePool.Get().(*Updater)
	obj.db = db
	obj.table = table
	obj.command.Reset()

	return obj

}

func (u *Updater) Free() {
	u.table = ""
	u.cols = u.cols[:]
	u.params = u.params[:]
	u.incrCols = u.incrCols[:]
	u.decrCols = u.decrCols[:]
	u.exprCols = u.exprCols[:]
	u.where.Reset()
	u.whereParams = u.whereParams[:]
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
func (c *Updater) Incr(fns ...ExprSetter) *Updater {
	for _, fn := range fns {
		s, val := fn()
		c.incrCols = append(c.incrCols, expr{colName: s, arg: val})
	}
	return c
}
func (c *Updater) Decr(fns ...ExprSetter) *Updater {
	for _, fn := range fns {
		s, val := fn()
		c.decrCols = append(c.decrCols, expr{colName: s, arg: val})
	}
	return c
}
func (c *Updater) SetExpr(fns ...ExprSetter) *Updater {
	for _, fn := range fns {
		s, val := fn()
		c.exprCols = append(c.exprCols, expr{colName: s, arg: val})
	}
	return c
}

// Where
func (c *Updater) Where(fns ...Condition) *Updater {
	if len(fns) == 0 {
		return c
	}
	c.where.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			c.where.WriteString(operator_and)
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
func (c *Updater) AndOr(fns ...Condition) *Updater {
	if len(fns) == 0 {
		return c
	}
	c.where.WriteString(operator_and + "(")
	for i, fn := range fns {
		if i > 0 {
			c.where.WriteString(operator_or)
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
func (c *Updater) OrAnd(fns ...Condition) *Updater {
	if len(fns) == 0 {
		return c
	}
	c.where.WriteString(operator_or + "(")
	for i, fn := range fns {
		if i > 0 {
			c.where.WriteString(operator_and)
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

// Do
func (c *Updater) Do(ctx context.Context) (sql.Result, error) {
	if len(c.cols) == 0 &&
		len(c.incrCols) == 0 &&
		len(c.decrCols) == 0 &&
		len(c.exprCols) == 0 {
		return nil, ErrCreateEmpty
	}
	_cols := make([]string, 0, 5)
	for _, col := range c.cols {
		_cols = append(_cols, col+" = ?")
	}
	for _, col := range c.incrCols {
		_cols = append(_cols, col.colName)
		c.params = append(c.params, col.arg)
	}
	for _, col := range c.decrCols {
		_cols = append(_cols, col.colName)
		c.params = append(c.params, col.arg)
	}

	for _, col := range c.exprCols {
		_cols = append(_cols, col.colName)
	}
	c.command.WriteString("UPDATE " + c.table + " SET ")
	c.command.WriteString(strings.Join(_cols, ","))
	// WHERE
	if c.where.Len() > 0 {
		c.command.WriteString(" WHERE " + c.where.String())
	}

	c.params = append(c.params, c.whereParams...)
	fmt.Println(c.command.String(), c.params)
	return nil, nil
	return c.db.ExecContext(ctx, c.command.String(), c.params...)
}

// //////////////////////////////////////////////////
// Deleter
func NewDelete(db ExtContext, table string) *Deleter {
	if db == nil || table == "" {
		panic("db or table is nil")
		return nil
	}
	obj := deletePool.Get().(*Deleter)
	obj.db = db
	obj.table = table
	obj.command.Reset()
	return obj

}

func (d *Deleter) Free() {
	deletePool.Put(d)
}

// Where
func (c *Deleter) Where(fns ...Condition) *Deleter {
	if len(fns) == 0 {
		return c
	}
	c.where.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			c.where.WriteString(operator_and)
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
func (c *Deleter) AndOr(fns ...Condition) *Deleter {
	if len(fns) == 0 {
		return c
	}
	c.where.WriteString(operator_and + "(")
	for i, fn := range fns {
		if i > 0 {
			c.where.WriteString(operator_or)
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
func (c *Deleter) OrAnd(fns ...Condition) *Deleter {
	if len(fns) == 0 {
		return c
	}
	c.where.WriteString(operator_or + "(")
	for i, fn := range fns {
		if i > 0 {
			c.where.WriteString(operator_and)
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

// Do
func (c *Deleter) Do(ctx context.Context) (sql.Result, error) {
	c.command.WriteString("DELETE FROM " + c.table)
	return c.db.ExecContext(ctx, c.command.String(), c.whereParams...)
}

// ///////////////////////////////////////////////
// Selector
func NewSelect(db ExtContext, table string) *Selector {
	if db == nil || table == "" {
		panic("db or table is nil")
		return nil
	}
	obj := selectPool.Get().(*Selector)
	obj.db = db
	obj.table = table
	obj.command.Reset()
	return obj
}

func (s *Selector) Free() {
	s.table = ""
	s.cols = s.cols[:]
	s.distinct = false
	s.join = s.join[:]
	s.omit = s.omit[:]
	s.where.Reset()
	s.whereParams = s.whereParams[:]
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
	c.where.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			c.where.WriteString(operator_and)
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
func (c *Selector) AndOr(fns ...Condition) *Selector {
	if len(fns) == 0 {
		return c
	}
	c.where.WriteString(operator_and + "(")
	for i, fn := range fns {
		if i > 0 {
			c.where.WriteString(operator_or)
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
func (c *Selector) OrAnd(fns ...Condition) *Selector {
	if len(fns) == 0 {
		return c
	}
	c.where.WriteString(operator_or + "(")
	for i, fn := range fns {
		if i > 0 {
			c.where.WriteString(operator_and)
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

// Do
func (c *Selector) Do(ctx context.Context) (sql.Result, error) {
	if len(c.cols) == 0 {
		return nil, ErrCreateEmpty
	}
	c.command.WriteString("SELECT ")
	fmt.Println(c.command, c.table, c.cols, c.whereParams)
	return c.db.ExecContext(ctx, c.command.String(), c.whereParams...)
}
