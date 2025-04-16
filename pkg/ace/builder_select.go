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
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"strings"
)

type SelectBuilder interface {
	Columner
	Wherer
	Orderer
	Grouper
	Join(joinType dialect.JoinType, left, right dialect.Field, fns ...dialect.Condition) Builder
	LeftJoin(left, right dialect.Field, fns ...dialect.Condition) Builder
	RightJoin(left, right dialect.Field, fns ...dialect.Condition) Builder
	Page(pageIndex, pageSize uint) Builder
	PageByBookmark(size uint, bm dialect.Condition) Builder
	Limit(size uint, start ...uint) Builder
	Distinct(cols ...dialect.Field) Builder
	Select(x ...Executer) Selecter
}

type Selecter interface {
	// Query
	Query(ctx context.Context) (*sql.Rows, error)
	// QueryRow
	QueryRow(ctx context.Context) (*sql.Row, error)
	// Get
	Get(ctx context.Context, dest any) error
	// Gets
	Gets(ctx context.Context, dest any) error
	// Map
	Map(ctx context.Context) (map[string]any, error)
	// Maps
	Maps(ctx context.Context) ([]map[string]any, error)
	// 	Slice
	Slice(ctx context.Context) ([]any, error)
	// 	Slices
	Slices(ctx context.Context) ([][]any, error)
	// Count 返回数量
	Count(ctx context.Context, cond ...dialect.Condition) (int64, error)
	// Sum 返回总和
	Sum(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error)
	// Select 执行原生查询，返回指定列的数据
	RawQuery(ctx context.Context, sqlStr string, args ...any) (*sql.Rows, error)
	// SelectMap 执行原生查询，返回 map[string]any
	RawQueryMap(ctx context.Context, sqlStr string, args ...any) (map[string]any, error)
	// SelectSlice 执行原生查询，返回 []any
	RawQuerySlice(ctx context.Context, sqlStr string, args ...any) ([]any, error)
	// SelectStruct 执行原生查询，返回结构体对象
	RawQueryStruct(ctx context.Context, dest any, sqlStr string, args ...any) error
}

type read struct {
	*orm
}

// Read 创建查询器
func (o *orm) Select(x ...Executer) Selecter {
	o.connect(x...)
	return &read{
		orm: o,
	}
}

// Query
func (s *read) Query(ctx context.Context) (*sql.Rows, error) {
	defer s.Free()

	return s.query(ctx)
}

// QueryRow
func (s *read) QueryRow(ctx context.Context) (*sql.Row, error) {
	defer s.Free()

	_ = s.parse()

	return s.row(ctx, s.command.String(), s.mergeParams()...)
}

// Get 返回单个数据，dest 必须是指针
func (s *read) Get(ctx context.Context, dest any) error {
	defer s.Free()

	s.Limit(1)

	rows, err := s.query(ctx)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 如果 dest 实现了 Modeler 接口，直接调用 AssignPtr 方法，并 scan 数据
	// 否则，调用 scanAny 方法
	if d, ok := dest.(dialect.Modeler); ok {
		if !rows.Next() {
			return sql.ErrNoRows
		}
		vals := d.AssignPtr()
		return rows.Scan(vals...)
	}
	r := &Row{rows: rows, err: err, Mapper: s.db.Mapper()}
	return r.scanAny(dest, false)
}

// Gets 返回数据切片，dest 必须是slice指针
func (s *read) Gets(ctx context.Context, dest any) error {
	defer s.Free()

	rows, err := s.query(ctx)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scanAll(rows, dest, false)
}

// Map 返回 map[string]any，用于列数未知的情况
func (s *read) Map(ctx context.Context) (map[string]any, error) {
	defer s.Free()

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
func (s *read) Maps(ctx context.Context) ([]map[string]any, error) {
	defer s.Free()

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
func (s *read) Slice(ctx context.Context) ([]any, error) {
	defer s.Free()

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
func (s *read) Slices(ctx context.Context) ([][]any, error) {
	defer s.Free()

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
func (s *read) Count(ctx context.Context, cond ...dialect.Condition) (int64, error) {
	defer s.Free()

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
func (s *read) Sum(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error) {
	defer s.Free()

	for _, col := range cols {
		s.Func(col.Sum())
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
func (s *read) RawQuery(ctx context.Context, sqlStr string, args ...any) (*sql.Rows, error) {
	defer s.Free()

	return s.rows(ctx, sqlStr, args...)
}

// SelectMap 执行原生 SQL 查询并返回 map[string]any
func (se *read) RawQueryMap(ctx context.Context, sqlStr string, args ...any) (map[string]any, error) {
	rows, err := se.RawQuery(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := &Row{rows: rows, err: err, Mapper: se.db.Mapper()}
	dest := make(map[string]any)
	return dest, r.MapScan(dest)
}

// SelectSlice 执行原生 SQL 查询并返回 []any
func (se *read) RawQuerySlice(ctx context.Context, sqlStr string, args ...any) ([]any, error) {
	rows, err := se.RawQuery(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := &Row{rows: rows, err: err, Mapper: se.db.Mapper()}
	return r.SliceScan()
}

// SelectStruct 执行原生 SQL 查询并返回实现 dialect.Modeler 接口的结构体
func (se *read) RawQueryStruct(ctx context.Context, dest any, sqlStr string, args ...any) error {
	rows, err := se.RawQuery(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 如果 dest 实现了 Modeler 接口，直接调用 AssignPtr 方法，并 scan 数据
	// 否则，调用 scanAny 方法
	if d, ok := dest.(dialect.Modeler); ok {
		if !rows.Next() {
			return sql.ErrNoRows
		}
		vals := d.AssignPtr()
		return rows.Scan(vals...)
	}
	r := &Row{rows: rows, err: err, Mapper: se.db.Mapper()}
	return r.scanAny(dest, false)
}
