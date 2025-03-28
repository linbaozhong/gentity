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
	"github.com/linbaozhong/gentity/pkg/util"
	"reflect"
	"strconv"
	"strings"
)

type (
	Builder interface {
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

func New(x ...Executer) *orm {
	obj := ormPool.Get().(*orm)
	if len(x) > 0 {
		obj.db = x[0]
	} else {
		obj.db = GetDB()
	}

	obj.commandString.Reset()
	return obj
}

// Free 释放 orm 对象，将其重置并放回对象池。
func (s *orm) Free() {
	if s == nil || s.table == "" {
		return
	}

	if s.db.Debug() {
		log.Info(s.String())
	}

	selectPool.Put(s)
}

func (s *orm) Reset() {
	s.table = ""
	s.cols = s.cols[:0]   // []dialect.Field{} // s.cols[:0]
	s.funcs = s.funcs[:0] // []string{}       // s.funcs[:0]
	s.distinct = false
	s.join = s.join[:0]             // [][3]string{}      // s.join[:0]
	s.joinParams = s.joinParams[:0] // []any{}      // s.joinParams[:0]
	s.omits = s.omits[:0]           // []dialect.Field{} // s.omits[:0]
	s.where.Reset()
	s.whereParams = s.whereParams[:0] // []any{} // s.whereParams[:0]
	s.groupBy.Reset()
	s.having.Reset()
	s.havingParams = s.havingParams[:0] // []any{} // s.havingParams[:0]
	s.orderBy.Reset()
	s.limit = ""
	s.exprCols = s.exprCols[:0] // []expr{} // s.exprCols[:0]
	s.params = s.params[:0]     // []any{} // s.params[:0]
	s.command.Reset()
}

// String 返回 orm 对象的 SQL 语句和参数的字符串表示。
func (s *orm) String() string {
	if s.commandString.Len() == 0 {
		s.commandString.WriteString(fmt.Sprintf("%s  %v \n", s.command.String(), s.mergeParams()))
	}
	return s.commandString.String()
}

// Table 设置 orm 对象的表名。
func (s *orm) Table(a any) *orm {
	switch v := a.(type) {
	case string:
		s.table = v
	case dialect.TableNamer:
		s.table = v.TableName()
	default:
		// 避免多次调用 reflect.ValueOf 和 reflect.Indirect
		value := reflect.ValueOf(a)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		s.table = value.Type().Name()
	}
	return s
}

// GetTableName 获取 orm 对象的表名。
func (s *orm) GetTableName() string {
	return s.table
}

// GetCols 获取 orm 对象要查询的列。
func (s *orm) GetCols() []dialect.Field {
	return s.cols
}

// Distinct 设置查询结果去重，并指定去重的列。
func (s *orm) Distinct(cols ...dialect.Field) *orm {
	s.distinct = true
	for _, col := range cols {
		s.cols = append(s.cols, col)
	}

	return s
}

// Cols 指定要查询的列
func (s *orm) Cols(cols ...dialect.Field) *orm {
	for _, col := range cols {
		s.cols = append(s.cols, col)
	}
	return s
}

// Omits 忽略指定的列
func (s *orm) Omits(cols ...dialect.Field) *orm {
	for _, col := range cols {
		s.omits = append(s.omits, col)
	}
	return s
}

// Omit Deprecated: 此方法已弃用，请使用Omits
// 忽略指定的列
func (s *orm) Omit(cols ...dialect.Field) *orm {
	return s.Omits(cols...)
}

// Funcs 添加聚合函数到查询中
func (s *orm) Funcs(fns ...dialect.Function) *orm {
	for _, fn := range fns {
		s.funcs = append(s.funcs, fn())
	}
	return s
}

// Join 添加连接查询条件
func (s *orm) Join(joinType dialect.JoinType, left, right dialect.Field, fns ...dialect.Condition) *orm {
	var on strings.Builder
	for _, fn := range fns {
		on.WriteString(dialect.Operator_and)
		cond, val := fn()
		// if v, ok := val.(error); ok {
		//	s.err = v
		//	return s
		// }
		on.WriteString(cond)
		if vals, ok := val.([]any); ok {
			s.joinParams = append(s.joinParams, vals...)
		} else {
			s.joinParams = append(s.joinParams, val)
		}
	}
	s.join = append(s.join, [3]string{
		string(joinType),
		right.TableName(),
		left.Quote() + "=" + right.Quote() + on.String(),
	})
	return s
}

// LeftJoin 添加左连接查询条件。
func (s *orm) LeftJoin(left, right dialect.Field, fns ...dialect.Condition) *orm {
	return s.Join(dialect.Left_Join, left, right, fns...)
}

// RightJoin 添加右连接查询条件。
func (s *orm) RightJoin(left, right dialect.Field, fns ...dialect.Condition) *orm {
	return s.Join(dialect.Right_Join, left, right, fns...)
}

// Where 添加查询条件。
func (s *orm) Where(fns ...dialect.Condition) *orm {
	if len(fns) == 0 {
		return s
	}

	if s.where.Len() == 0 {
		s.where.WriteString("(")
	} else {
		s.where.WriteString(dialect.Operator_and + "(")
	}
	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if cond[:len(dialect.Operator_or)] == dialect.Operator_or || cond[:len(dialect.Operator_and)] == dialect.Operator_and {
				s.where.WriteString(" ")
			} else {
				s.where.WriteString(dialect.Operator_and)
			}
		}
		s.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			s.whereParams = append(s.whereParams, vals...)
		} else {
			s.whereParams = append(s.whereParams, val)
		}
	}
	s.where.WriteString(")")

	return s
}

// And 添加 AND 查询条件。
func (s *orm) And(fns ...dialect.Condition) *orm {
	if len(fns) == 0 {
		return s
	}

	if s.where.Len() == 0 {
		s.where.WriteString("(")
	} else {
		s.where.WriteString(dialect.Operator_and + "(")
	}
	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if cond[:len(dialect.Operator_or)] == dialect.Operator_or || cond[:len(dialect.Operator_and)] == dialect.Operator_and {
				s.where.WriteString(" ")
			} else {
				s.where.WriteString(dialect.Operator_or)
			}
		}
		s.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			s.whereParams = append(s.whereParams, vals...)
		} else {
			s.whereParams = append(s.whereParams, val)
		}
	}
	s.where.WriteString(")")
	return s
}

// Or 添加 OR 查询条件。
func (s *orm) Or(fns ...dialect.Condition) *orm {
	if len(fns) == 0 {
		return s
	}

	if s.where.Len() == 0 {
		s.where.WriteString("(")
	} else {
		s.where.WriteString(dialect.Operator_or + "(")
	}

	for i, fn := range fns {
		cond, val := fn()
		if i > 0 {
			if cond[:len(dialect.Operator_or)] == dialect.Operator_or || cond[:len(dialect.Operator_and)] == dialect.Operator_and {
				s.where.WriteString(" ")
			} else {
				s.where.WriteString(dialect.Operator_and)
			}
		}
		s.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			s.whereParams = append(s.whereParams, vals...)
		} else {
			s.whereParams = append(s.whereParams, val)
		}
	}
	s.where.WriteString(")")
	return s
}

// OrderFunc 方法用于根据传入的排序规则函数设置排序规则
// 它会遍历传入的排序规则函数，根据规则函数的返回值调用 Asc 或 Desc 方法
func (s *orm) OrderFunc(ords ...dialect.Order) *orm {
	for _, ord := range ords {
		sord, fs := ord()
		if sord == dialect.Operator_Desc {
			s.Desc(fs...)
		} else {
			s.Asc(fs...)
		}
	}
	return s
}

// OrderField
// Deprecated: 此方法后续版本可能会被移除，建议使用 OrderFunc 方法
func (s *orm) OrderField(ords ...dialect.Order) *orm {
	return s.OrderFunc(ords...)
}

// Order 指定查询结果的排序字段，默认升序。
func (s *orm) Order(cols ...dialect.Field) *orm {
	return s.Asc(cols...)
}

// Asc 指定查询结果按指定列升序排序。
func (s *orm) Asc(cols ...dialect.Field) *orm {
	if len(cols) == 0 {
		return s
	}
	for _, col := range cols {
		if s.orderBy.Len() > 0 {
			s.orderBy.WriteByte(',')
		}
		s.orderBy.WriteString(col.Quote())
	}
	return s
}

// Desc
func (s *orm) Desc(cols ...dialect.Field) *orm {
	if len(cols) == 0 {
		return s
	}
	for _, col := range cols {
		if s.orderBy.Len() > 0 {
			s.orderBy.WriteByte(',')
		}
		s.orderBy.WriteString(col.Quote() + dialect.Operator_Desc)
	}
	return s
}

// Group
func (s *orm) Group(cols ...dialect.Field) *orm {
	if len(cols) == 0 {
		return s
	}
	for _, col := range cols {
		if s.groupBy.Len() > 0 {
			s.groupBy.WriteByte(',')
		}
		s.groupBy.WriteString(col.Quote())
	}
	return s
}

// Group Having
func (s *orm) Having(fns ...dialect.Condition) *orm {
	if len(fns) == 0 {
		return s
	}

	s.having.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			s.having.WriteString(dialect.Operator_and)
		}
		cond, val := fn()
		// if v, ok := val.(error); ok {
		//	s.err = v
		//	return s
		// }
		s.having.WriteString(cond)
		if vals, ok := val.([]any); ok {
			s.havingParams = append(s.havingParams, vals...)
		}
	}
	return s
}

// Limit
// size 大小
// start 开始位置
func (s *orm) Limit(size uint, start ...uint) *orm {
	if size == 0 {
		s.limit = ""
		return s
	}
	if len(start) > 0 {
		s.limit = " LIMIT " + strconv.Itoa(int(size)) + " OFFSET " + strconv.Itoa(int(start[0]))
	} else {
		s.limit = " LIMIT " + strconv.Itoa(int(size))
	}

	return s
}

// Page
// pageIndex 页码
// pageSize 页大小
func (s *orm) Page(pageIndex, pageSize uint) *orm {
	if pageSize == 0 {
		return s.Limit(0)
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return s.Limit(pageSize, (pageIndex-1)*pageSize)
}

// Set
// 用于设置更新语句中的字段和值
// 例如：Set(dialect.F("name", "linbaozhong"))
func (s *orm) Set(fns ...dialect.Setter) *orm {
	if len(fns) == 0 {
		return s
	}

	for _, fn := range fns {
		c, val := fn()
		s.cols = append(s.cols, c)
		s.params = append(s.params, val)
	}
	return s
}

// SetExpr
// 用于设置更新语句中的表达式
// 例如：SetExpr(dialect.Expr("age", "age + 1"))
func (s *orm) SetExpr(fns ...dialect.ExprSetter) *orm {
	if len(fns) == 0 {
		return s
	}

	for _, fn := range fns {
		ex, val := fn()
		s.exprCols = append(s.exprCols, expr{colName: ex, arg: val})
	}
	return s
}

// Clone 克隆 orm
func (s *orm) Clone() *orm {
	_s := *s
	_s.cols = append([]dialect.Field(nil), s.cols...)
	_s.funcs = append([]string(nil), s.funcs...)
	_s.join = append([][3]string(nil), s.join...)
	_s.joinParams = append([]any(nil), s.joinParams...)
	_s.omits = append([]dialect.Field(nil), s.omits...)
	_s.whereParams = append([]any(nil), s.whereParams...)
	_s.havingParams = append([]any(nil), s.havingParams...)
	return &_s
}

// 合并参数
func (s *orm) mergeParams() []any {
	if len(s.joinParams) > 0 {
		if len(s.whereParams) > 0 {
			var params = make([]any, len(s.joinParams)+len(s.whereParams))
			copy(params, s.joinParams)
			copy(params[len(s.joinParams):], s.whereParams)
			return params
		}
		return s.joinParams
	}
	return s.whereParams
}

// parse
func (s *orm) parse() []dialect.Field {
	s.command.WriteString("SELECT ")

	var cols = util.SliceDiff(s.cols, s.omits)
	colens := len(cols)
	funlens := len(s.funcs)
	if colens+funlens == 0 {
		s.command.WriteString("*")
	} else {
		if s.distinct {
			s.command.WriteString("DISTINCT ")
		}
		for i, col := range cols {
			if i > 0 {
				s.command.WriteString(",")
			}
			s.command.WriteString(col.Quote())
		}
		if colens > 0 && funlens > 0 {
			s.command.WriteString(",")
		}
		s.command.WriteString(strings.Join(s.funcs, ","))
	}

	// FROM TABLE
	s.command.WriteString(" FROM " + dialect.Quote_Char + s.table + dialect.Quote_Char)
	for _, j := range s.join {
		s.command.WriteString(j[0] + " JOIN " + j[1] + " ON " + j[2] + " ")
	}

	// WHERE
	if s.where.Len() > 0 {
		s.command.WriteString(" WHERE " + s.where.String())
	}
	// GROUP BY
	if s.groupBy.Len() > 0 {
		s.command.WriteString(" GROUP BY " + s.groupBy.String())
		// HAVING
		if s.having.Len() > 0 {
			s.command.WriteString(" HAVING " + s.having.String())
			s.whereParams = append(s.whereParams, s.havingParams...)
		}
	}
	// ORDER BY
	if s.orderBy.Len() > 0 {
		s.command.WriteString(" ORDER BY " + s.orderBy.String())
	}

	// LIMIT
	if s.limit != "" {
		s.command.WriteString(s.limit)
	}

	return cols
}

// query
func (s *orm) query(ctx context.Context) (*sql.Rows, error) {
	_ = s.parse()
	return s.rows(ctx, s.command.String(), s.mergeParams()...)
}

func (s *orm) rows(ctx context.Context, sql string, params ...any) (*sql.Rows, error) {
	stmt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	if s.db.IsDB() {
		defer stmt.Close()
	}

	rows, err := stmt.QueryContext(ctx, params...)
	return rows, err
}

func (s *orm) row(ctx context.Context, sql string, params ...any) (*sql.Row, error) {
	stmt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	if s.db.IsDB() {
		defer stmt.Close()
	}

	return stmt.QueryRowContext(ctx, s.mergeParams()...), nil
}
