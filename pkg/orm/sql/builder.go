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

import "strings"

type JoinType string
type command string

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

// Expr represents an SQL express
type Expr struct {
	ColName string
	Arg     interface{}
}

type common struct {
	Command command
	Table   string
	Cols    []string
}

type Select struct {
	common
	Join        [][3]string
	Distinct    bool
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

type Update struct {
	common

	Params   []interface{}
	IncrCols []Expr
	DecrCols []Expr
	ExprCols []Expr
}

type Delete struct {
	common
}

type Create struct {
	common
	Params []interface{}
}
