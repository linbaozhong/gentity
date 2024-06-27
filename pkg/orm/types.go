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
	"database/sql"
	"strings"
)

type (
	Field struct {
		Name  string
		Table string
	}
	Condition func() (string, any)
	Setter    func() (string, any)

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
)

func (f Field) quote() string {
	return f.TableName() + "." + f.FieldName()
}

func (f Field) TableName() string {
	return Quote_Char + f.Table + Quote_Char
}
func (f Field) FieldName() string {
	return Quote_Char + f.Name + Quote_Char
}

func (f Field) Set(val any) Setter {
	return func() (string, any) {
		return f.quote(), val
	}
}

func (f Field) Eq(val any) Condition {
	return func() (string, any) {
		return f.quote() + " = " + placeholder, val
	}
}

func (f Field) NotEq(val any) Condition {
	return func() (string, any) {
		return f.quote() + " != " + placeholder, val
	}
}

func (f Field) Gt(val any) Condition {
	return func() (string, any) {
		return f.quote() + " > " + placeholder, val
	}
}

func (f Field) Gte(val any) Condition {
	return func() (string, any) {
		return f.quote() + " >= " + placeholder, val
	}
}

func (f Field) Lt(val any) Condition {
	return func() (string, any) {
		return f.quote() + " < " + placeholder, val
	}
}

func (f Field) Lte(val any) Condition {
	return func() (string, any) {
		return f.quote() + " <= " + placeholder, val
	}
}

func (f Field) In(vals ...any) Condition {
	return func() (string, any) {
		l := len(vals)
		return f.quote() + " In (" + strings.Repeat(placeholder+",", l)[:2*l-1] + ") ", vals
	}
}

func (f Field) NotIn(vals ...any) Condition {
	return func() (string, any) {
		l := len(vals)
		return f.quote() + " Not In (" + strings.Repeat(placeholder+",", l)[:2*l-1] + ") ", vals
	}
}

func (f Field) Between(vals ...any) Condition {
	return func() (string, any) {
		return f.quote() + " BETWEEN " + placeholder + " AND " + placeholder, vals
	}
}

func (f Field) Like(val any) Condition {
	return func() (string, any) {
		return "(" + f.quote() + " LIKE CONCAT('%'," + placeholder + ",'%'))", val
	}
}

func (f Field) Llike(val any) Condition {
	return func() (string, any) {
		return "(" + f.quote() + " LIKE CONCAT('%'," + placeholder + "))", val
	}
}

func (f Field) Rlike(val any) Condition {
	return func() (string, any) {
		return "(" + f.quote() + " LIKE CONCAT(" + placeholder + ",'%'))", val
	}
}

func (f Field) Null(val any) Condition {
	return func() (string, any) {
		return " ISNULL(" + placeholder + ")", val
	}
}

func (f Field) NotNull(val any) Condition {
	return func() (string, any) {
		return " NOT ISNULL(" + placeholder + ")", val
	}
}
