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
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/log"
	"strings"
)

type (
	Deleter interface {
		Where(fns ...dialect.Condition) Deleter
		And(fns ...dialect.Condition) Deleter
		Or(fns ...dialect.Condition) Deleter
		Exec(ctx context.Context) (sql.Result, error)
	}
	DeleteBuilder interface {
		Deleter
		Table(name any) DeleteBuilder
		Free()
		Reset()
		String() string
	}
	delete struct {
		pool.Model
		db            Executer
		table         string
		where         strings.Builder
		whereParams   []any
		command       strings.Builder
		commandString strings.Builder
	}
)

var (
	deletePool = pool.New(app.Context, func() any {
		obj := &delete{}
		obj.UUID()
		return obj
	})
)

// delete
func newDelete(dbs ...Executer) DeleteBuilder {
	obj := deletePool.Get().(*delete)
	obj.db = GetExec(dbs...)
	obj.commandString.Reset()

	return obj

}

func (d *delete) Free() {
	if d == nil || d.table == "" {
		return
	}

	_ = d.String()
	if d.db.Debug() {
		log.Info(d.String())
	}

	deletePool.Put(d)
}

func (d *delete) Reset() {
	d.table = ""
	d.where.Reset()
	d.whereParams = d.whereParams[:0] // []any{} // d.whereParams[:0]
	d.command.Reset()
}

// setTableName 设置表名
func (d *delete) Table(n any) DeleteBuilder {
	setTableName(&d.table, n)
	return d
}

func (d *delete) String() string {
	if d.commandString.Len() == 0 {
		d.commandString.WriteString(fmt.Sprintf("%s  %v \n", d.command.String(), d.whereParams))
	}
	return d.commandString.String()
}

// Where
func (d *delete) Where(fns ...dialect.Condition) Deleter {
	if len(fns) == 0 {
		return d
	}

	if d.where.Len() == 0 {
		d.where.WriteString("(")
	} else {
		d.where.WriteString(dialect.Operator_and + "(")
	}
	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if cond[:len(dialect.Operator_or)] == dialect.Operator_or || cond[:len(dialect.Operator_and)] == dialect.Operator_and {
				d.where.WriteString(" ")
			} else {
				d.where.WriteString(dialect.Operator_and)
			}
		}
		d.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			d.whereParams = append(d.whereParams, vals...)
		} else {
			d.whereParams = append(d.whereParams, val)
		}
	}
	d.where.WriteString(")")

	return d
}

// And
func (d *delete) And(fns ...dialect.Condition) Deleter {
	if len(fns) == 0 {
		return d
	}

	if d.where.Len() == 0 {
		d.where.WriteString("(")
	} else {
		d.where.WriteString(dialect.Operator_and + "(")
	}

	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if cond[:len(dialect.Operator_or)] == dialect.Operator_or || cond[:len(dialect.Operator_and)] == dialect.Operator_and {
				d.where.WriteString(" ")
			} else {
				d.where.WriteString(dialect.Operator_or)
			}
		}
		d.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			d.whereParams = append(d.whereParams, vals...)
		} else {
			d.whereParams = append(d.whereParams, val)
		}
	}
	d.where.WriteString(")")
	return d
}

// Or
func (d *delete) Or(fns ...dialect.Condition) Deleter {
	if len(fns) == 0 {
		return d
	}

	if d.where.Len() == 0 {
		d.where.WriteString("(")
	} else {
		d.where.WriteString(dialect.Operator_or + "(")
	}

	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if cond[:len(dialect.Operator_or)] == dialect.Operator_or || cond[:len(dialect.Operator_and)] == dialect.Operator_and {
				d.where.WriteString(" ")
			} else {
				d.where.WriteString(dialect.Operator_and)
			}
		}
		d.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			d.whereParams = append(d.whereParams, vals...)
		} else {
			d.whereParams = append(d.whereParams, val)
		}
	}
	d.where.WriteString(")")
	return d
}

// Exec
func (d *delete) Exec(ctx context.Context) (sql.Result, error) {
	defer d.Free()

	d.command.WriteString("DELETE FROM " + dialect.Quote_Char + d.table + dialect.Quote_Char)
	// WHERE
	if d.where.Len() > 0 {
		d.command.WriteString(" WHERE " + d.where.String())
	}

	stmt, err := d.db.PrepareContext(ctx, d.command.String())
	if err != nil {
		return nil, err
	}
	if d.db.IsDB() {
		defer stmt.Close()
	}

	return stmt.ExecContext(ctx, d.whereParams...)
}
