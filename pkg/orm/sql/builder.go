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

	command_insert command = "INSERT INTO "
	command_select command = "SELECT "
	command_update command = "UPDATE "
	command_delete command = "DELETE FROM "
)

type (
	JoinType string
	command  string
	// expr represents an SQL express
	expr struct {
		ColName string
		Arg     interface{}
	}
	condition struct {
		Where       strings.Builder
		WhereParams []interface{}
	}

	Creator struct {
		Command command
		Table   string
		Cols    []string
		Params  []interface{}
	}
	Selector struct {
		Command  command
		Table    string
		Join     [][3]string
		Distinct bool
		Cols     []string
		Omit     []interface{}
		// Where       strings.Builder
		// WhereParams []interface{}
		GroupBy    strings.Builder
		Having     strings.Builder
		OrderBy    strings.Builder
		Limit      string
		LimitSize  int
		LimitStart int

		condition
	}
	Updater struct {
		Command  command
		Table    string
		Cols     []string
		Params   []interface{}
		IncrCols []expr
		DecrCols []expr
		ExprCols []expr

		condition
	}
	Deleter struct {
		Command command
		Table   string

		condition
	}
)

var (
	createPool = sync.Pool{
		New: func() interface{} {
			return &Creator{
				Command: command_insert,
			}
		},
	}
	selectPool = sync.Pool{
		New: func() interface{} {
			return &Selector{
				Command: command_select,
			}
		},
	}
	updatePool = sync.Pool{
		New: func() interface{} {
			return &Updater{
				Command: command_update,
			}
		},
	}
	deletePool = sync.Pool{
		New: func() interface{} {
			return &Deleter{
				Command: command_delete,
			}
		},
	}
)

// ////////////////////////////////////////
// Creator
func NewCreate() *Creator {
	return createPool.Get().(*Creator)
}

func (c *Creator) Free() {
	c.Table = ""
	c.Cols = c.Cols[:]
	c.Params = c.Params[:]
	createPool.Put(c)
}

// Set
func (c *Creator) Set(fn Setter) *Creator {
	s, val := fn()
	c.Cols = append(c.Cols, s)
	c.Params = append(c.Params, val)
	return c
}

// Sets
func (c *Creator) Sets(fns ...Setter) *Creator {
	for _, fn := range fns {
		s, val := fn()
		c.Cols = append(c.Cols, s)
		c.Params = append(c.Params, val)
	}
	return c
}

// Do
func (c *Creator) Do(ctx context.Context) (sql.Result, error) {
	return nil, nil
}

// ///////////////////////////////////////////////
// Selector
func NewSelect() *Selector {
	return selectPool.Get().(*Selector)
}

func (s *Selector) Free() {
	s.Table = ""
	s.Cols = s.Cols[:]
	s.Distinct = false
	s.Join = s.Join[:]
	s.Omit = s.Omit[:]
	s.Where.Reset()
	s.WhereParams = s.WhereParams[:]
	// s.AndOr = false
	s.GroupBy.Reset()
	s.Having.Reset()
	s.OrderBy.Reset()
	s.Limit = ""
	s.LimitSize = 0
	s.LimitStart = 0
	selectPool.Put(s)
}

// Do
func (s *Selector) Do(ctx context.Context) (sql.Result, error) {
	return nil, nil
}

// /////////////////////////////////////////////////
// Updater
func NewUpdate() *Updater {
	return updatePool.Get().(*Updater)
}

func (u *Updater) Free() {
	u.Table = ""
	u.Cols = u.Cols[:]
	u.Params = u.Params[:]
	u.IncrCols = u.IncrCols[:]
	u.DecrCols = u.DecrCols[:]
	u.ExprCols = u.ExprCols[:]
	u.condition.Where.Reset()
	u.condition.WhereParams = u.condition.WhereParams[:]
	// u.AndOr = false
	updatePool.Put(u)
}

// Set
func (c *Updater) Set(fns ...Setter) *Updater {
	for _, fn := range fns {
		s, val := fn()
		c.Cols = append(c.Cols, s)
		c.Params = append(c.Params, val)
	}
	return c
}

// Where
func (c *Updater) Where(fns ...Condition) *Updater {
	if len(fns) == 0 {
		return c
	}
	c.condition.Where.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			c.condition.Where.WriteString(operator_and)
		}
		cond, val := fn()
		c.condition.Where.WriteString(cond)
		c.condition.WhereParams = append(c.condition.WhereParams, val)
	}
	c.condition.Where.WriteString(")")

	return c
}

// And
func (c *Updater) AndOr(fns ...Condition) *Updater {
	if len(fns) == 0 {
		return c
	}
	c.condition.Where.WriteString(operator_and + "(")
	for i, fn := range fns {
		if i > 0 {
			c.condition.Where.WriteString(operator_or)
		}
		cond, val := fn()
		c.condition.Where.WriteString(cond)
		c.condition.WhereParams = append(c.condition.WhereParams, val)
	}
	c.condition.Where.WriteString(")")
	return c
}

// Or
func (c *Updater) OrAnd(fns ...Condition) *Updater {
	if len(fns) == 0 {
		return c
	}
	c.condition.Where.WriteString(operator_or + "(")
	for i, fn := range fns {
		if i > 0 {
			c.condition.Where.WriteString(operator_and)
		}
		cond, val := fn()
		c.condition.Where.WriteString(cond)
		c.condition.WhereParams = append(c.condition.WhereParams, val)
	}
	c.condition.Where.WriteString(")")
	return c
}

// Do
func (c *Updater) Do(ctx context.Context) (sql.Result, error) {
	return nil, nil
}

// //////////////////////////////////////////////////
// Deleter
func NewDelete() *Deleter {
	return deletePool.Get().(*Deleter)
}

func (d *Deleter) Free() {
	deletePool.Put(d)
}

// Do
func (c *Deleter) Do(ctx context.Context) (sql.Result, error) {
	return nil, nil
}
