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
	"errors"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/util"
	"strconv"
	"strings"
)

type (
	Select struct {
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
		command       strings.Builder
		commandString strings.Builder
		err           error
	}
)

var (
	selectPool = pool.New(app.Context, func() any {
		obj := &Select{}
		obj.UUID()
		return obj
	})
)

// Select
func newSelect(db Executer, tableName string) *Select {
	obj := selectPool.Get().(*Select)
	if db == nil || tableName == "" {
		obj.err = errors.New("db or table is nil")
		return obj
	}

	obj.db = db
	obj.table = tableName
	obj.err = nil
	obj.commandString.Reset()

	return obj
}

// Free 释放 Select 对象，将其重置并放回对象池。
func (s *Select) Free() {
	if s == nil || s.table == "" {
		return
	}

	_ = s.String()
	if s.db.Debug() {
		log.Info(s.String())
	}

	selectPool.Put(s)
}

func (s *Select) Reset() {
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
	s.command.Reset()
}

// String 返回 Select 对象的 SQL 语句和参数的字符串表示。
func (s *Select) String() string {
	if s.commandString.Len() == 0 {
		s.commandString.WriteString(fmt.Sprintf("%s  %v \n", s.command.String(), s.mergeParams()))
	}
	return s.commandString.String()
}

// SetTableName 设置 Select 对象的表名。
func (s *Select) SetTableName(n string) {
	s.table = n
}

// GetTableName 获取 Select 对象的表名。
func (s *Select) GetTableName() string {
	return s.table
}

// GetCols 获取 Select 对象要查询的列。
func (s *Select) GetCols() []dialect.Field {
	return s.cols
}

// Distinct 设置查询结果去重，并指定去重的列。
func (s *Select) Distinct(cols ...dialect.Field) *Select {
	s.distinct = true
	for _, col := range cols {
		s.cols = append(s.cols, col)
	}

	return s
}

// Cols 指定要查询的列
func (s *Select) Cols(cols ...dialect.Field) *Select {
	for _, col := range cols {
		s.cols = append(s.cols, col)
	}
	return s
}

// Omits 忽略指定的列
func (s *Select) Omits(cols ...dialect.Field) *Select {
	for _, col := range cols {
		s.omits = append(s.omits, col)
	}
	return s
}

// Omit Deprecated: 此方法已弃用，请使用Omits
// 忽略指定的列
func (s *Select) Omit(cols ...dialect.Field) *Select {
	return s.Omits(cols...)
}

// Funcs 添加聚合函数到查询中
func (s *Select) Funcs(fns ...dialect.Function) *Select {
	for _, fn := range fns {
		s.funcs = append(s.funcs, fn())
	}
	return s
}

// Join 添加连接查询条件
func (s *Select) Join(joinType dialect.JoinType, left, right dialect.Field, fns ...dialect.Condition) *Select {
	if s.err != nil {
		return s
	}

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
func (s *Select) LeftJoin(left, right dialect.Field, fns ...dialect.Condition) *Select {
	return s.Join(dialect.Left_Join, left, right, fns...)
}

// RightJoin 添加右连接查询条件。
func (s *Select) RightJoin(left, right dialect.Field, fns ...dialect.Condition) *Select {
	return s.Join(dialect.Right_Join, left, right, fns...)
}

// Where 添加查询条件。
func (s *Select) Where(fns ...dialect.Condition) *Select {
	if len(fns) == 0 || s.err != nil {
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
func (s *Select) And(fns ...dialect.Condition) *Select {
	if len(fns) == 0 || s.err != nil {
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
func (s *Select) Or(fns ...dialect.Condition) *Select {
	if len(fns) == 0 || s.err != nil {
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
func (s *Select) OrderFunc(ords ...dialect.Order) *Select {
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
func (s *Select) OrderField(ords ...dialect.Order) *Select {
	return s.OrderFunc(ords...)
}

// Order 指定查询结果的排序字段，默认升序。
func (s *Select) Order(cols ...dialect.Field) *Select {
	return s.Asc(cols...)
}

// Asc 指定查询结果按指定列升序排序。
func (s *Select) Asc(cols ...dialect.Field) *Select {
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
func (s *Select) Desc(cols ...dialect.Field) *Select {
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
func (s *Select) Group(cols ...dialect.Field) *Select {
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
func (s *Select) Having(fns ...dialect.Condition) *Select {
	if len(fns) == 0 || s.err != nil {
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
func (s *Select) Limit(size uint, start ...uint) *Select {
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
func (s *Select) Page(pageIndex, pageSize uint) *Select {
	if pageSize == 0 {
		return s.Limit(0)
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return s.Limit(pageSize, (pageIndex-1)*pageSize)
}

// Clone 克隆 Select
func (s *Select) Clone() *Select {
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

// Query
func (s *Select) Query(ctx context.Context) (*sql.Rows, error) {
	defer s.Free()

	if s.err != nil {
		return nil, s.err
	}

	return s.query(ctx)
}

// QueryRow
func (s *Select) QueryRow(ctx context.Context) (*sql.Row, error) {
	defer s.Free()

	if s.err != nil {
		return nil, s.err
	}

	_ = s.parse()

	// stmt, err := s.db.PrepareContext(ctx, s.command.String())
	// if err != nil {
	// 	return nil, err
	// }
	// if s.db.IsDB() {
	// 	defer stmt.Close()
	// }
	//
	// return stmt.QueryRowContext(ctx, s.mergeParams()...), nil
	return s.row(ctx, s.command.String(), s.mergeParams()...)
}

// Get 返回单个数据，dest 必须是指针
func (s *Select) Get(ctx context.Context, dest any) error {
	defer s.Free()

	if s.err != nil {
		return s.err
	}

	s.Limit(1)

	rows, err := s.query(ctx)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 如果 dest 实现了 Modeler 接口，直接调用 AssignPtr 方法，并 scan 数据
	// 否则，调用 scanAny 方法
	if d, ok := dest.(dialect.Modeler); ok && rows.Next() {
		vals := d.AssignPtr(s.cols...)
		return rows.Scan(vals...)
	}
	r := &Row{rows: rows, err: err, Mapper: s.db.Mapper()}
	return r.scanAny(dest, false)
}

// Gets 返回数据切片，dest 必须是slice指针
func (s *Select) Gets(ctx context.Context, dest any) error {
	defer s.Free()

	if s.err != nil {
		return s.err
	}

	rows, err := s.query(ctx)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scanAll(rows, dest, false)
}

// Map 返回 map[string]any，用于列数未知的情况
func (s *Select) Map(ctx context.Context) (map[string]any, error) {
	defer s.Free()

	if s.err != nil {
		return nil, s.err
	}

	s.Limit(1)

	rows, err := s.query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := &Row{rows: rows, err: err, Mapper: s.db.Mapper()}
	dest := make(map[string]any)
	return dest, r.MapScan(dest)
}

// Maps 返回 map[string]any 的切片 []map[string]any，用于列数未知的情况
func (s *Select) Maps(ctx context.Context) ([]map[string]any, error) {
	defer s.Free()

	if s.err != nil {
		return nil, s.err
	}

	rows, err := s.query(ctx)
	if err != nil {
		return nil, err
	}

	rs := &Rows{Rows: rows, Mapper: s.db.Mapper()}
	defer rs.Close()

	dests := make([]map[string]any, 0)
	for rs.Next() {
		dest := make(map[string]any)
		err = rs.MapScan(dest)
		if err != nil {
			break
		}
		dests = append(dests, dest)
	}

	return dests, rs.Err()
}

// Slice 返回切片 []any，用于列数未知的情况
func (s *Select) Slice(ctx context.Context) ([]any, error) {
	defer s.Free()

	if s.err != nil {
		return nil, s.err
	}

	s.Limit(1)

	rows, err := s.query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := &Row{rows: rows, err: err, Mapper: s.db.Mapper()}
	return r.SliceScan()
}

// Slices 返回 []any 的切片 [][]any，用于列数未知的情况
func (s *Select) Slices(ctx context.Context) ([][]any, error) {
	defer s.Free()

	if s.err != nil {
		return nil, s.err
	}

	rows, err := s.query(ctx)
	if err != nil {
		return nil, err
	}

	rs := &Rows{Rows: rows, Mapper: s.db.Mapper()}
	defer rs.Close()

	dests := make([][]any, 0)
	for rs.Next() {
		dest, err := rs.SliceScan()
		if err != nil {
			break
		}
		dests = append(dests, dest)
	}

	return dests, rs.Err()
}

// Count
func (s *Select) Count(ctx context.Context, cond ...dialect.Condition) (int64, error) {
	defer s.Free()

	if s.err != nil {
		return 0, s.err
	}

	s.Where(cond...)
	s.command.WriteString("SELECT COUNT(*)")

	// FROM TABLE
	s.command.WriteString(" FROM " + dialect.Quote_Char + s.table + dialect.Quote_Char)
	for _, j := range s.join {
		s.command.WriteString(j[0] + " JOIN " + j[1] + " ON " + j[2] + " ")
	}

	// WHERE
	if s.where.Len() > 0 {
		s.command.WriteString(" WHERE " + s.where.String())
	}

	// LIMIT
	if s.limit != "" {
		s.command.WriteString(s.limit)
	}

	// stmt, err := s.db.PrepareContext(ctx, s.command.String())
	// if err != nil {
	// 	return 0, err
	// }
	// if s.db.IsDB() {
	// 	defer stmt.Close()
	// }
	//
	// row := stmt.QueryRowContext(ctx, s.mergeParams()...)
	row, err := s.row(ctx, s.command.String(), s.mergeParams()...)
	if err != nil {
		return 0, err
	}
	var count int64
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Sum
func (s *Select) Sum(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error) {
	defer s.Free()

	if s.err != nil {
		return nil, s.err
	}

	for _, col := range cols {
		s.Funcs(col.Sum())
	}
	s.Where(cond...)
	s.command.WriteString("SELECT ")
	s.command.WriteString(strings.Join(s.funcs, ","))

	// FROM TABLE
	s.command.WriteString(" FROM " + dialect.Quote_Char + s.table + dialect.Quote_Char)
	for _, j := range s.join {
		s.command.WriteString(j[0] + " JOIN " + j[1] + " ON " + j[2] + " ")
	}

	// WHERE
	if s.where.Len() > 0 {
		s.command.WriteString(" WHERE " + s.where.String())
	}

	// LIMIT
	if s.limit != "" {
		s.command.WriteString(s.limit)
	}

	// stmt, err := s.db.PrepareContext(ctx, s.command.String())
	// if err != nil {
	// 	return nil, err
	// }
	// if s.db.IsDB() {
	// 	defer stmt.Close()
	// }
	//
	// row := stmt.QueryRowContext(ctx, s.mergeParams()...)
	row, err := s.row(ctx, s.command.String(), s.mergeParams()...)
	if err != nil {
		return nil, err
	}

	var sum = make([]any, len(cols))
	err = row.Scan(sum...)
	if err != nil {
		return nil, err
	}

	sums := make(map[string]any, len(cols))
	for i := range sum {
		sums[cols[i].Name] = sum[i]
	}
	return sums, nil
}

// Select 执行原生的 SQL 查询
// 此方法接受一个上下文、原生 SQL 语句和对应的参数，返回查询结果和可能的错误
func (s *Select) Select(ctx context.Context, sqlStr string, args ...any) (*sql.Rows, error) {
	defer s.Free()

	if s.err != nil {
		return nil, s.err
	}

	return s.rows(ctx, sqlStr, args...)
}

// SelectMap 执行原生 SQL 查询并返回 map[string]any
func (s *Select) SelectMap(ctx context.Context, sqlStr string, args ...any) (map[string]any, error) {
	rows, err := s.Select(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := &Row{rows: rows, err: err, Mapper: s.db.Mapper()}
	dest := make(map[string]any)
	return dest, r.MapScan(dest)
}

// SelectSlice 执行原生 SQL 查询并返回 []any
func (s *Select) SelectSlice(ctx context.Context, sqlStr string, args ...any) ([]any, error) {
	rows, err := s.Select(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := &Row{rows: rows, err: err, Mapper: s.db.Mapper()}
	return r.SliceScan()
}

// SelectModel 执行原生 SQL 查询并返回实现 dialect.Modeler 接口的结构体
func (s *Select) SelectModel(ctx context.Context, dest any, sqlStr string, args ...any) error {
	rows, err := s.Select(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 如果 dest 实现了 Modeler 接口，直接调用 AssignPtr 方法，并 scan 数据
	// 否则，调用 scanAny 方法
	if d, ok := dest.(dialect.Modeler); ok && rows.Next() {
		vals := d.AssignPtr(s.cols...)
		return rows.Scan(vals...)
	}
	r := &Row{rows: rows, err: err, Mapper: s.db.Mapper()}
	return r.scanAny(dest, false)
}

// 合并参数
func (s *Select) mergeParams() []any {
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
func (s *Select) parse() []dialect.Field {
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
		if colens > 0 {
			for i, col := range cols {
				if i > 0 {
					s.command.WriteString(",")
				}
				s.command.WriteString(col.Quote())
			}
			if funlens > 0 {
				s.command.WriteString(",")
			}
		}
		if funlens > 0 {
			s.command.WriteString(strings.Join(s.funcs, ","))
		} else {
			s.cols = cols
		}
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
func (s *Select) query(ctx context.Context) (*sql.Rows, error) {
	_ = s.parse()
	return s.rows(ctx, s.command.String(), s.mergeParams()...)
}

func (s *Select) rows(ctx context.Context, sql string, params ...any) (*sql.Rows, error) {
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

func (s *Select) row(ctx context.Context, sql string, params ...any) (*sql.Row, error) {
	stmt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	if s.db.IsDB() {
		defer stmt.Close()
	}

	return stmt.QueryRowContext(ctx, s.mergeParams()...), nil
}
