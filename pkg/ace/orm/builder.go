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
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/util"
	"reflect"
	"strings"
)

type (
	Builder interface {
		Free()
		String() string
		Clone() Builder

		Columner
		Wherer

		SelectBuilder
		CreateBuilder
		UpdateBuilder
		DeleteBuilder
	}

	orm struct {
		pool.Model
		db            Executer
		table         string
		join          [][3]string
		joinParams    []any
		distinct      bool
		cols          []dialect.Field
		funcs         []string
		omits         []dialect.Field
		groupBy       strings.Builder
		having        strings.Builder
		havingParams  []any
		orderBy       strings.Builder
		limit         string
		where         strings.Builder
		whereParams   []any
		exprCols      []expr
		params        []any
		command       strings.Builder
		commandString strings.Builder
		// err           error
	}
)

var (
	ormPool = pool.New(app.Context, func() any {
		obj := &orm{}
		obj.UUID()
		return obj
	})
)

func new() *orm {
	obj := ormPool.Get().(*orm)
	obj.commandString.Reset()
	return obj
}

func New(opts ...Option) Builder {
	obj := new()

	for _, opt := range opts {
		opt(obj)
	}
	return obj
}

// Free 释放 orm 对象，将其重置并放回对象池。
func (o *orm) Free() {
	if o == nil || o.table == "" {
		return
	}

	_ = o.String()
	if o.db.Debug() {
		log.Info(o.String())
	}

	ormPool.Put(o)
}

func (o *orm) Reset() {
	o.table = ""
	o.cols = o.cols[:0]   // []dialect.Field{} // o.cols[:0]
	o.funcs = o.funcs[:0] // []string{}       // o.funcs[:0]
	o.distinct = false
	o.join = o.join[:0]             // [][3]string{}      // o.join[:0]
	o.joinParams = o.joinParams[:0] // []any{}      // o.joinParams[:0]
	o.omits = o.omits[:0]           // []dialect.Field{} // o.omits[:0]
	o.where.Reset()
	o.whereParams = o.whereParams[:0] // []any{} // o.whereParams[:0]
	o.groupBy.Reset()
	o.having.Reset()
	o.havingParams = o.havingParams[:0] // []any{} // o.havingParams[:0]
	o.orderBy.Reset()
	o.limit = ""
	o.exprCols = o.exprCols[:0] // []expr{} // o.exprCols[:0]
	o.params = o.params[:0]     // []any{} // o.params[:0]
	o.command.Reset()
}

// String 返回 orm 对象的 SQL 语句和参数的字符串表示。
func (o *orm) String() string {
	if o.commandString.Len() == 0 {
		o.commandString.WriteString(fmt.Sprintf("%s  %v \n", o.command.String(), o.mergeParams()))
	}
	return o.commandString.String()
}

// Table 设置 orm 对象的表名。
func (o *orm) Table(a any) Builder {
	switch v := a.(type) {
	case string:
		o.table = v
	case dialect.TableNamer:
		o.table = v.TableName()
	default:
		// 避免多次调用 reflect.ValueOf 和 reflect.Indirect
		value := reflect.ValueOf(a)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		o.table = value.Type().Name()
	}
	return o
}

// connect 连接数据库
func (o *orm) connect(x ...Executer) Builder {
	if len(x) > 0 {
		o.db = x[0]
	} else {
		o.db = GetDB()
	}
	return o
}

// GetTableName 获取 orm 对象的表名。
func (s *orm) GetTableName() string {
	return s.table
}

// Set
// 用于设置更新语句中的字段和值
// 例如：Set(dialect.F("name", "linbaozhong"))
func (o *orm) Set(fns ...dialect.Setter) Builder {
	if len(fns) == 0 {
		return o
	}

	for _, fn := range fns {
		c, val := fn()
		o.cols = append(o.cols, c)
		o.params = append(o.params, val)
	}
	return o
}

// SetExpr
// 用于设置更新语句中的表达式
// 例如：SetExpr(dialect.Expr("age", "age + 1"))
func (o *orm) SetExpr(fns ...dialect.ExprSetter) Builder {
	if len(fns) == 0 {
		return o
	}

	for _, fn := range fns {
		ex, val := fn()
		o.exprCols = append(o.exprCols, expr{colName: ex, arg: val})
	}
	return o
}

// Clone 克隆 orm
func (o *orm) Clone() Builder {
	_s := *o
	_s.cols = append([]dialect.Field(nil), o.cols...)
	_s.funcs = append([]string(nil), o.funcs...)
	_s.join = append([][3]string(nil), o.join...)
	_s.joinParams = append([]any(nil), o.joinParams...)
	_s.omits = append([]dialect.Field(nil), o.omits...)
	_s.whereParams = append([]any(nil), o.whereParams...)
	_s.havingParams = append([]any(nil), o.havingParams...)
	return &_s
}

// 合并参数
func (o *orm) mergeParams() []any {
	if len(o.joinParams) > 0 {
		if len(o.whereParams) > 0 {
			var params = make([]any, len(o.joinParams)+len(o.whereParams))
			copy(params, o.joinParams)
			copy(params[len(o.joinParams):], o.whereParams)
			return params
		}
		return o.joinParams
	}
	return o.whereParams
}

// parse
func (o *orm) parse() []dialect.Field {
	o.command.WriteString("SELECT ")

	var cols = util.SliceDiff(o.cols, o.omits)
	colens := len(cols)
	funlens := len(o.funcs)
	if colens+funlens == 0 {
		o.command.WriteString("*")
	} else {
		if o.distinct {
			o.command.WriteString("DISTINCT ")
		}
		for i, col := range cols {
			if i > 0 {
				o.command.WriteString(",")
			}
			o.command.WriteString(col.Quote())
		}
		if colens > 0 && funlens > 0 {
			o.command.WriteString(",")
		}
		o.command.WriteString(strings.Join(o.funcs, ","))
	}

	// FROM TABLE
	o.command.WriteString(" FROM " + dialect.Quote_Char + o.table + dialect.Quote_Char)
	for _, j := range o.join {
		o.command.WriteString(j[0] + " JOIN " + j[1] + " ON " + j[2] + " ")
	}

	// WHERE
	if o.where.Len() > 0 {
		o.command.WriteString(" WHERE " + o.where.String())
	}
	// GROUP BY
	if o.groupBy.Len() > 0 {
		o.command.WriteString(" GROUP BY " + o.groupBy.String())
		// HAVING
		if o.having.Len() > 0 {
			o.command.WriteString(" HAVING " + o.having.String())
			o.whereParams = append(o.whereParams, o.havingParams...)
		}
	}
	// ORDER BY
	if o.orderBy.Len() > 0 {
		o.command.WriteString(" ORDER BY " + o.orderBy.String())
	}

	// LIMIT
	if o.limit != "" {
		o.command.WriteString(o.limit)
	}

	return cols
}

// query
func (o *orm) query(ctx context.Context) (*sql.Rows, error) {
	_ = o.parse()
	return o.rows(ctx, o.command.String(), o.mergeParams()...)
}

func (o *orm) rows(ctx context.Context, sql string, params ...any) (*sql.Rows, error) {
	stmt, err := o.db.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	if o.db.IsDB() {
		defer stmt.Close()
	}

	rows, err := stmt.QueryContext(ctx, params...)
	return rows, err
}

func (o *orm) row(ctx context.Context, sql string, params ...any) (*sql.Row, error) {
	stmt, err := o.db.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	if o.db.IsDB() {
		defer stmt.Close()
	}

	return stmt.QueryRowContext(ctx, o.mergeParams()...), nil
}
