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
	// Expr represents an SQL express
	Expr struct {
		ColName string
		Arg     interface{}
	}
	common struct {
		Command command
		Table   string
	}
	Create struct {
		common
		Cols   []string
		Params []interface{}
	}
	Select struct {
		common
		Join        [][3]string
		Distinct    bool
		Cols        []string
		Omit        []interface{}
		Where       strings.Builder
		WhereParams []interface{}
		GroupBy     strings.Builder
		Having      strings.Builder
		OrderBy     strings.Builder
		Limit       string
		LimitSize   int
		LimitStart  int

		AndOr bool
	}
	Update struct {
		common
		Cols []string

		Params   []interface{}
		IncrCols []Expr
		DecrCols []Expr
		ExprCols []Expr
	}
	Delete struct {
		common
	}
)

var (
	createPool = sync.Pool{
		New: func() interface{} {
			return &Create{
				common: common{
					Command: command_insert,
				},
			}
		},
	}
	selectPool = sync.Pool{
		New: func() interface{} {
			return &Select{
				common: common{
					Command: command_select,
				},
			}
		},
	}
	updatePool = sync.Pool{
		New: func() interface{} {
			return &Update{
				common: common{
					Command: command_update,
				},
			}
		},
	}
	deletePool = sync.Pool{
		New: func() interface{} {
			return &Delete{
				common: common{
					Command: command_delete,
				},
			}
		},
	}
)

func NewCreate() *Create {
	return createPool.Get().(*Create)
}

func NewSelect() *Select {
	return selectPool.Get().(*Select)
}

func NewUpdate() *Update {
	return updatePool.Get().(*Update)
}

func NewDelete() *Delete {
	return deletePool.Get().(*Delete)
}

func (c *Create) Free() {
	c.Cols = c.Cols[:]
	c.Params = c.Params[:]
	createPool.Put(c)
}

func (c *Select) Free() {
	c.Cols = c.Cols[:]
	c.Distinct = false
	c.Join = c.Join[:]
	c.Omit = c.Omit[:]
	c.Where.Reset()
	c.WhereParams = c.WhereParams[:]
	c.GroupBy.Reset()
	c.Having.Reset()
	c.OrderBy.Reset()
	c.Limit = ""
	c.LimitSize = 0
	c.LimitStart = 0
	selectPool.Put(c)
}

func (u *Update) Free() {
	u.Cols = u.Cols[:]
	u.Params = u.Params[:]
	u.IncrCols = u.IncrCols[:]
	u.DecrCols = u.DecrCols[:]
	u.ExprCols = u.ExprCols[:]
	updatePool.Put(u)
}

func (d *Delete) Free() {
	deletePool.Put(d)
}
