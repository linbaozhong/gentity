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
	"strconv"
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
	ErrNotFound    = fmt.Errorf("")
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
		object  Modeler
		cols    []string
		params  []interface{}
		command strings.Builder
	}
	Selector struct {
		db           ExtContext
		object       Modeler
		join         [][3]string
		distinct     bool
		cols         []string
		omit         []interface{}
		groupBy      strings.Builder
		having       strings.Builder
		havingParams []interface{}
		orderBy      strings.Builder
		limit        string
		limitSize    int
		limitStart   int
		where        strings.Builder
		whereParams  []interface{}
		command      strings.Builder
	}
	Updater struct {
		db     ExtContext
		object Modeler
		cols   []string
		params []interface{}
		// incrCols    []expr
		// decrCols    []expr
		exprCols    []expr
		where       strings.Builder
		whereParams []interface{}
		command     strings.Builder
	}
	Deleter struct {
		db          ExtContext
		object      Modeler
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
func NewCreate(db ExtContext, mod Modeler) *Creator {
	if db == nil || mod == nil {
		panic("db or table is nil")
		return nil
	}
	obj := createPool.Get().(*Creator)
	obj.db = db
	obj.object = mod
	obj.command.Reset()
	return obj
}

func (c *Creator) Free() {
	c.object = nil
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
	c.command.WriteString("INSERT INTO " + c.object.TableName() + " (")
	c.command.WriteString(strings.Join(c.cols, ",") + ") VALUES (")
	c.command.WriteString(strings.Repeat("?,", len(c.cols))[:len(c.cols)*2-1])
	c.command.WriteString(")")
	fmt.Println(c.command.String(), c.params)
	return nil, nil

	return c.db.ExecContext(ctx, c.command.String(), c.params...)
}

// /////////////////////////////////////////////////
// Updater
func NewUpdate(db ExtContext, mod Modeler) *Updater {
	if db == nil || mod == nil {
		panic("db or table is nil")
		return nil
	}
	obj := updatePool.Get().(*Updater)
	obj.db = db
	obj.object = mod
	obj.command.Reset()

	return obj

}

func (u *Updater) Free() {
	u.object = nil
	u.cols = u.cols[:]
	u.params = u.params[:]
	// u.incrCols = u.incrCols[:]
	// u.decrCols = u.decrCols[:]
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
func (c *Updater) And(fns ...Condition) *Updater {
	if len(fns) == 0 {
		return c
	}

	if c.where.Len() == 0 {
		c.where.WriteString("(")
	} else {
		c.where.WriteString(operator_and + "(")
	}

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
func (c *Updater) Or(fns ...Condition) *Updater {
	if len(fns) == 0 {
		return c
	}

	if c.where.Len() == 0 {
		c.where.WriteString("(")
	} else {
		c.where.WriteString(operator_or + "(")
	}

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
		// len(c.incrCols) == 0 &&
		// len(c.decrCols) == 0 &&
		len(c.exprCols) == 0 {
		return nil, ErrCreateEmpty
	}
	_cols := make([]string, 0, 5)
	for _, col := range c.cols {
		_cols = append(_cols, col+" = ?")
	}
	for _, col := range c.exprCols {
		_cols = append(_cols, col.colName)
		if col.arg != nil {
			c.params = append(c.params, col.arg)
		}
	}
	// for _, col := range c.decrCols {
	// 	_cols = append(_cols, col.colName)
	// 	c.params = append(c.params, col.arg)
	// }
	//
	// for _, col := range c.exprCols {
	// 	_cols = append(_cols, col.colName)
	// }
	c.command.WriteString("UPDATE " + c.object.TableName() + " SET ")
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
func NewDelete(db ExtContext, mod Modeler) *Deleter {
	if db == nil || mod == nil {
		panic("db or table is nil")
		return nil
	}
	obj := deletePool.Get().(*Deleter)
	obj.db = db
	obj.object = mod
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
func (c *Deleter) And(fns ...Condition) *Deleter {
	if len(fns) == 0 {
		return c
	}

	if c.where.Len() == 0 {
		c.where.WriteString("(")
	} else {
		c.where.WriteString(operator_and + "(")
	}

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
func (c *Deleter) Or(fns ...Condition) *Deleter {
	if len(fns) == 0 {
		return c
	}

	if c.where.Len() == 0 {
		c.where.WriteString("(")
	} else {
		c.where.WriteString(operator_or + "(")
	}

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
	c.command.WriteString("DELETE FROM " + c.object.TableName())
	return c.db.ExecContext(ctx, c.command.String(), c.whereParams...)
}

// ///////////////////////////////////////////////
// Selector
func NewSelect(db ExtContext, mod Modeler) *Selector {
	if db == nil || mod == nil {
		panic("db or table is nil")
		return nil
	}
	obj := selectPool.Get().(*Selector)
	obj.db = db
	obj.object = mod
	obj.command.Reset()
	return obj
}

func (s *Selector) Free() {
	s.object = nil
	s.cols = s.cols[:]
	s.distinct = false
	s.join = s.join[:]
	s.omit = s.omit[:]
	s.where.Reset()
	s.whereParams = s.whereParams[:]
	// s.AndOr = false
	s.groupBy.Reset()
	s.having.Reset()
	s.havingParams = s.havingParams[:]
	s.orderBy.Reset()
	s.limit = ""
	s.limitSize = 0
	s.limitStart = 0
	selectPool.Put(s)
}

// distinct
func (c *Selector) Distinct(cols ...Field) *Selector {
	c.distinct = true
	for _, col := range cols {
		c.cols = append(c.cols, col.quote())
	}

	return c
}

// cols
func (c *Selector) Cols(cols ...Field) *Selector {
	for _, col := range cols {
		c.cols = append(c.cols, col.quote())
	}
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
		left.quote() + "=" + right.quote() + on.String(),
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
func (c *Selector) And(fns ...Condition) *Selector {
	if len(fns) == 0 {
		return c
	}

	if c.where.Len() == 0 {
		c.where.WriteString("(")
	} else {
		c.where.WriteString(operator_and + "(")
	}

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
func (c *Selector) Or(fns ...Condition) *Selector {
	if len(fns) == 0 {
		return c
	}

	if c.where.Len() == 0 {
		c.where.WriteString("(")
	} else {
		c.where.WriteString(operator_or + "(")
	}

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

// Order
func (c *Selector) Order(cols ...Field) *Selector {
	return c.Asc(cols...)
}

// Order Asc
func (c *Selector) Asc(cols ...Field) *Selector {
	if len(cols) == 0 {
		return c
	}
	for _, col := range cols {
		if c.orderBy.Len() > 0 {
			c.orderBy.WriteByte(',')
		}
		c.orderBy.WriteString(col.quote())
	}
	return c
}

// Order Desc
func (c *Selector) Desc(cols ...Field) *Selector {
	if len(cols) == 0 {
		return c
	}
	for _, col := range cols {
		if c.orderBy.Len() > 0 {
			c.orderBy.WriteByte(',')
		}
		c.orderBy.WriteString(col.quote() + " DESC")
	}
	return c
}

// Group
func (c *Selector) Group(cols ...Field) *Selector {
	if len(cols) == 0 {
		return c
	}
	for _, col := range cols {
		if c.groupBy.Len() > 0 {
			c.groupBy.WriteByte(',')
		}
		c.groupBy.WriteString(col.quote())
	}
	return c
}

// Group Having
func (c *Selector) Having(fns ...Condition) *Selector {
	if len(fns) == 0 {
		return c
	}
	c.having.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			c.having.WriteString(operator_and)
		}
		cond, val := fn()
		c.having.WriteString(cond)
		if vals, ok := val.([]any); ok {
			c.havingParams = append(c.havingParams, vals...)
		}
	}
	return c
}
func (c *Selector) Limit(size int, start ...int) *Selector {
	c.limitSize = size
	if len(start) > 0 {
		c.limitStart = start[0]
	}
	c.limit = " LIMIT " + strconv.Itoa(c.limitSize) + " OFFSET " + strconv.Itoa(c.limitStart)

	return c
}

func (c *Selector) stmt() {
	c.command.WriteString("SELECT ")
	if len(c.cols) == 0 {
		c.command.WriteString("*")
	} else {
		if c.distinct {
			c.command.WriteString("DISTINCT ")
		}
		c.command.WriteString(strings.Join(c.cols, ","))
	}
	// FROM TABLE
	c.command.WriteString(" FROM " + Quote_Char + c.object.TableName() + Quote_Char)
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

// // Query
// func (c *Selector) QueryRow(ctx context.Context) *sql.Row {
// 	c.Limit(1)
// 	c.stmt()
// 	return c.db.QueryRowContext(ctx, c.command.String(), c.whereParams...)
// }

// Query
func (c *Selector) Query(ctx context.Context) (*sql.Rows, error) {
	c.stmt()
	return c.db.QueryContext(ctx, c.command.String(), c.whereParams...)
}

func (c *Selector) Get(ctx context.Context) (Modeler, error) {
	rows, err := c.Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrNotFound
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	vals, err := c.object.ScanValues(cols)
	if err != nil {
		return nil, err
	}
	err = rows.Scan(vals...)
	if err != nil {
		return nil, err
	}

	err = c.object.AssignValues(cols, vals)

	return c.object, err
}

func (c *Selector) Gets(ctx context.Context) ([]Modeler, error) {
	rows, err := c.Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	vals, err := c.object.ScanValues(cols)
	if err != nil {
		return nil, err
	}

	dests := make([]Modeler, 0, 1)
	i := 0
	for rows.Next() {
		_vals := make([]any, 0, len(vals))
		copy(_vals, vals)

		err = rows.Scan(_vals...)
		if err != nil {
			return nil, err
		}

		dest := c.object.New()
		err = dest.AssignValues(cols, _vals)
		if err != nil {
			return nil, err
		}
		dests = append(dests, dest)
		i++
	}

	return dests, nil
}
