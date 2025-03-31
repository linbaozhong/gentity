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

package builder

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/util"
	"reflect"
	"strconv"
	"strings"
)

type (
	Wherer interface {
		Where(fns ...dialect.Condition) Builder
		And(fns ...dialect.Condition) Builder
		Or(fns ...dialect.Condition) Builder
	}
	Orderer interface {
		OrderFunc(ords ...dialect.Order) Builder
		Order(cols ...dialect.Field) Builder
		Asc(cols ...dialect.Field) Builder
		Desc(cols ...dialect.Field) Builder
	}

	Grouper interface {
		Group(cols ...dialect.Field) Builder
		Having(fns ...dialect.Condition) Builder
	}

	Selecter interface {
		Join(joinType dialect.JoinType, left, right dialect.Field, fns ...dialect.Condition) Builder
		LeftJoin(left, right dialect.Field, fns ...dialect.Condition) Builder
		RightJoin(left, right dialect.Field, fns ...dialect.Condition) Builder
		Page(pageIndex, pageSize uint) Builder
		Limit(size uint, start ...uint) Builder
		Distinct(cols ...dialect.Field) Builder
	}

	Columner interface {
		Cols(cols ...dialect.Field) Builder
		Funcs(fns ...dialect.Function) Builder
		Omits(cols ...dialect.Field) Builder
	}

	Builder interface {
		Free()
		String() string
		Clone() Builder
		Wherer
		Orderer
		Grouper
		Selecter
		Columner
		Create(a any) Creater
		Delete(a any) Deleter
		Read(a any) Reader
		Update(a any) Updater
	}

	orm struct {
		pool.Model
		db            ace.Executer
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

func New(x ...ace.Executer) Builder {
	obj := ormPool.Get().(*orm)
	if len(x) > 0 {
		obj.db = x[0]
	} else {
		obj.db = ace.GetDB()
	}

	obj.commandString.Reset()
	return obj
}

// Free 释放 orm 对象，将其重置并放回对象池。
func (o *orm) Free() {
	if o == nil || o.table == "" {
		return
	}

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
func (o *orm) setTable(a any) Builder {
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

// // GetTableName 获取 orm 对象的表名。
// func (s *orm) GetTableName() string {
//	return s.table
// }
//
// // GetCols 获取 orm 对象要查询的列。
// func (s *orm) GetCols() []dialect.Field {
//	return s.cols
// }

// Distinct 设置查询结果去重，并指定去重的列。
func (o *orm) Distinct(cols ...dialect.Field) Builder {
	o.distinct = true
	for _, col := range cols {
		o.cols = append(o.cols, col)
	}

	return o
}

// Cols 指定要查询的列
func (o *orm) Cols(cols ...dialect.Field) Builder {
	for _, col := range cols {
		o.cols = append(o.cols, col)
	}
	return o
}

// Omits 忽略指定的列
func (o *orm) Omits(cols ...dialect.Field) Builder {
	for _, col := range cols {
		o.omits = append(o.omits, col)
	}
	return o
}

// Omit Deprecated: 此方法已弃用，请使用Omits
// 忽略指定的列
func (o *orm) Omit(cols ...dialect.Field) Builder {
	return o.Omits(cols...)
}

// Funcs 添加聚合函数到查询中
func (o *orm) Funcs(fns ...dialect.Function) Builder {
	for _, fn := range fns {
		o.funcs = append(o.funcs, fn())
	}
	return o
}

// Join 添加连接查询条件
func (o *orm) Join(joinType dialect.JoinType, left, right dialect.Field, fns ...dialect.Condition) Builder {
	var on strings.Builder
	for _, fn := range fns {
		on.WriteString(dialect.Operator_and)
		cond, val := fn()
		// if v, ok := val.(error); ok {
		//	o.err = v
		//	return o
		// }
		on.WriteString(cond)
		if vals, ok := val.([]any); ok {
			o.joinParams = append(o.joinParams, vals...)
		} else {
			o.joinParams = append(o.joinParams, val)
		}
	}
	o.join = append(o.join, [3]string{
		string(joinType),
		right.TableName(),
		left.Quote() + "=" + right.Quote() + on.String(),
	})
	return o
}

// LeftJoin 添加左连接查询条件。
func (o *orm) LeftJoin(left, right dialect.Field, fns ...dialect.Condition) Builder {
	return o.Join(dialect.Left_Join, left, right, fns...)
}

// RightJoin 添加右连接查询条件。
func (o *orm) RightJoin(left, right dialect.Field, fns ...dialect.Condition) Builder {
	return o.Join(dialect.Right_Join, left, right, fns...)
}

// Where 添加查询条件。
func (o *orm) Where(fns ...dialect.Condition) Builder {
	if len(fns) == 0 {
		return o
	}

	if o.where.Len() == 0 {
		o.where.WriteString("(")
	} else {
		o.where.WriteString(dialect.Operator_and + "(")
	}
	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if cond[:len(dialect.Operator_or)] == dialect.Operator_or || cond[:len(dialect.Operator_and)] == dialect.Operator_and {
				o.where.WriteString(" ")
			} else {
				o.where.WriteString(dialect.Operator_and)
			}
		}
		o.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			o.whereParams = append(o.whereParams, vals...)
		} else {
			o.whereParams = append(o.whereParams, val)
		}
	}
	o.where.WriteString(")")

	return o
}

// And 添加 AND 查询条件。
func (o *orm) And(fns ...dialect.Condition) Builder {
	if len(fns) == 0 {
		return o
	}

	if o.where.Len() == 0 {
		o.where.WriteString("(")
	} else {
		o.where.WriteString(dialect.Operator_and + "(")
	}
	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if cond[:len(dialect.Operator_or)] == dialect.Operator_or || cond[:len(dialect.Operator_and)] == dialect.Operator_and {
				o.where.WriteString(" ")
			} else {
				o.where.WriteString(dialect.Operator_or)
			}
		}
		o.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			o.whereParams = append(o.whereParams, vals...)
		} else {
			o.whereParams = append(o.whereParams, val)
		}
	}
	o.where.WriteString(")")
	return o
}

// Or 添加 OR 查询条件。
func (o *orm) Or(fns ...dialect.Condition) Builder {
	if len(fns) == 0 {
		return o
	}

	if o.where.Len() == 0 {
		o.where.WriteString("(")
	} else {
		o.where.WriteString(dialect.Operator_or + "(")
	}

	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if cond[:len(dialect.Operator_or)] == dialect.Operator_or || cond[:len(dialect.Operator_and)] == dialect.Operator_and {
				o.where.WriteString(" ")
			} else {
				o.where.WriteString(dialect.Operator_and)
			}
		}
		o.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			o.whereParams = append(o.whereParams, vals...)
		} else {
			o.whereParams = append(o.whereParams, val)
		}
	}
	o.where.WriteString(")")
	return o
}

// OrderFunc 方法用于根据传入的排序规则函数设置排序规则
// 它会遍历传入的排序规则函数，根据规则函数的返回值调用 Asc 或 Desc 方法
func (o *orm) OrderFunc(ords ...dialect.Order) Builder {
	for _, ord := range ords {
		sord, fs := ord()
		if sord == dialect.Operator_Desc {
			o.Desc(fs...)
		} else {
			o.Asc(fs...)
		}
	}
	return o
}

// OrderField
// Deprecated: 此方法后续版本可能会被移除，建议使用 OrderFunc 方法
func (o *orm) OrderField(ords ...dialect.Order) Builder {
	return o.OrderFunc(ords...)
}

// Order 指定查询结果的排序字段，默认升序。
func (o *orm) Order(cols ...dialect.Field) Builder {
	return o.Asc(cols...)
}

// Asc 指定查询结果按指定列升序排序。
func (o *orm) Asc(cols ...dialect.Field) Builder {
	if len(cols) == 0 {
		return o
	}
	for _, col := range cols {
		if o.orderBy.Len() > 0 {
			o.orderBy.WriteByte(',')
		}
		o.orderBy.WriteString(col.Quote())
	}
	return o
}

// Desc
func (o *orm) Desc(cols ...dialect.Field) Builder {
	if len(cols) == 0 {
		return o
	}
	for _, col := range cols {
		if o.orderBy.Len() > 0 {
			o.orderBy.WriteByte(',')
		}
		o.orderBy.WriteString(col.Quote() + dialect.Operator_Desc)
	}
	return o
}

// Group
func (o *orm) Group(cols ...dialect.Field) Builder {
	if len(cols) == 0 {
		return o
	}
	for _, col := range cols {
		if o.groupBy.Len() > 0 {
			o.groupBy.WriteByte(',')
		}
		o.groupBy.WriteString(col.Quote())
	}
	return o
}

// Group Having
func (o *orm) Having(fns ...dialect.Condition) Builder {
	if len(fns) == 0 {
		return o
	}

	o.having.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			o.having.WriteString(dialect.Operator_and)
		}
		cond, val := fn()
		// if v, ok := val.(error); ok {
		//	o.err = v
		//	return o
		// }
		o.having.WriteString(cond)
		if vals, ok := val.([]any); ok {
			o.havingParams = append(o.havingParams, vals...)
		}
	}
	return o
}

// Limit
// size 大小
// start 开始位置
func (o *orm) Limit(size uint, start ...uint) Builder {
	if size == 0 {
		o.limit = ""
		return o
	}
	if len(start) > 0 {
		o.limit = " LIMIT " + strconv.Itoa(int(size)) + " OFFSET " + strconv.Itoa(int(start[0]))
	} else {
		o.limit = " LIMIT " + strconv.Itoa(int(size))
	}

	return o
}

// Page
// pageIndex 页码
// pageSize 页大小
func (o *orm) Page(pageIndex, pageSize uint) Builder {
	if pageSize == 0 {
		return o.Limit(0)
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return o.Limit(pageSize, (pageIndex-1)*pageSize)
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
