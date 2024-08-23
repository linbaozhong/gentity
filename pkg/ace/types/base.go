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

package types

import (
	"database/sql"
	"fmt"
)

const (
	Inner_Join JoinType = " INNER"
	Left_Join  JoinType = " LEFT"
	Right_Join JoinType = " RIGHT"

	Operator_and = " AND "
	Operator_or  = " OR "
	placeholder  = "?"
	Quote_Char   = "`"
)

var (
	ErrCreateEmpty  = fmt.Errorf("No data is created")
	ErrBeanEmpty    = fmt.Errorf("bean=nil 或者 len(beans)=0 或者 len(beans)>100")
	ErrNotFound     = fmt.Errorf("not found")
	ErrSetterEmpty  = fmt.Errorf("setter=nil 或者 len(setter)=0")
	ErrArgsNotMatch = fmt.Errorf("args not match")
)

type (
	JoinType string
	// NullBool is an alias to sql.NullBool.
	NullBool = sql.NullBool
	// NullInt64 is an alias to sql.NullInt64.
	NullInt64 = sql.NullInt64
	// NullInt32 is an alias to sql.NullInt32.
	NullInt32 = sql.NullInt32
	// NullInt16 is an alias to sql.NullInt16.
	NullInt16 = sql.NullInt16
	// NullByte is an alias to sql.NullByte.
	NullByte = sql.NullByte
	// NullString is an alias to sql.NullString.
	NullString = sql.NullString
	// NullFloat64 is an alias to sql.NullFloat64.
	NullFloat64 = sql.NullFloat64
	// NullTime represents a time.Time that may be null.
	NullTime = sql.NullTime
	// UnknownType is a named type to any indicates the info
	// needs to be extracted from the underlying rows.
	UnknownType any

	Modeler interface {
		// New() Modeler
		TableName() string
		// Set(args ...Field) ([]string, []any)
		// AssignColumns(args ...Field) []string
		AssignValues(args ...Field) ([]string, []any)
		Scan(rows *sql.Rows, args ...Field) ([]Modeler, error)
		AssignKeys() ([]Field, []any)
		// AssignValues(columns []string, values []any) error
		// ScanValues(columns []string) ([]any, error)
	}
)
