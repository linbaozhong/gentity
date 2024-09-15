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
	"github.com/linbaozhong/gentity/pkg/ace/types"
	"github.com/linbaozhong/gentity/pkg/log"
	"strconv"
	"strings"
	"sync"
)

type (
	Selector struct {
		db            Executer
		table         string
		join          [][3]string
		distinct      bool
		cols          []dialect.Field
		funcs         []string
		omit          []any
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
	selectPool = sync.Pool{
		New: func() any {
			obj := &Selector{}
			return obj
		},
	}
)

// Selector
func newSelect(db Executer, tableName string) *Selector {
	obj := selectPool.Get().(*Selector)
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

func (s *Selector) Free() {
	if s == nil {
		return
	}

	s.commandString.WriteString(fmt.Sprintf("%s  %v \n", s.command.String(), s.whereParams))

	if s.db.Debug() {
		log.Info(s.String())
	}
	s.table = ""
	s.cols = s.cols[:]
	s.funcs = s.funcs[:]
	s.distinct = false
	s.join = s.join[:]
	s.omit = s.omit[:]
	s.where.Reset()
	s.whereParams = s.whereParams[:]
	s.groupBy.Reset()
	s.having.Reset()
	s.havingParams = s.havingParams[:]
	s.orderBy.Reset()
	s.limit = ""
	s.command.Reset()

	selectPool.Put(s)
}

func (s *Selector) String() string {
	if s.table == "" {
		return s.commandString.String()
	}
	return fmt.Sprintf("%s  %v", s.command.String(), s.whereParams)
}

// distinct
func (s *Selector) Distinct(cols ...dialect.Field) *Selector {
	s.distinct = true
	for _, col := range cols {
		s.cols = append(s.cols, col)
	}

	return s
}

// cols
func (s *Selector) Cols(cols ...dialect.Field) *Selector {
	for _, col := range cols {
		s.cols = append(s.cols, col)
	}
	return s
}

// funcs
func (s *Selector) Funcs(fns ...dialect.Function) *Selector {
	for _, fn := range fns {
		s.funcs = append(s.funcs, fn())
	}
	return s
}

// join
func (s *Selector) Join(joinType types.JoinType, left, right dialect.Field, fns ...dialect.Condition) *Selector {
	if s.err != nil {
		return s
	}

	var on strings.Builder
	for _, fn := range fns {
		on.WriteString(types.Operator_and)
		cond, val := fn()
		//if v, ok := val.(error); ok {
		//	s.err = v
		//	return s
		//}
		on.WriteString(cond)
		if vals, ok := val.([]any); ok {
			s.whereParams = append(s.whereParams, vals...)
		} else {
			s.whereParams = append(s.whereParams, val)
		}
	}
	s.join = append(s.join, [3]string{
		string(joinType),
		right.TableName(),
		left.Quote() + "=" + right.Quote() + on.String(),
	})
	return s
}

func (s *Selector) LeftJoin(left, right dialect.Field, fns ...dialect.Condition) *Selector {
	return s.Join(types.Left_Join, left, right, fns...)
}
func (s *Selector) RightJoin(left, right dialect.Field, fns ...dialect.Condition) *Selector {
	return s.Join(types.Right_Join, left, right, fns...)
}

// Where
func (s *Selector) Where(fns ...dialect.Condition) *Selector {
	if len(fns) == 0 || s.err != nil {
		return s
	}

	if s.where.Len() == 0 {
		s.where.WriteString("(")
	} else {
		s.where.WriteString(types.Operator_and + "(")
	}
	for i, fn := range fns {
		if i > 0 {
			s.where.WriteString(types.Operator_and)
		}
		cond, val := fn()
		//if v, ok := val.(error); ok {
		//	s.err = v
		//	return s
		//}
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

// And
func (s *Selector) And(fns ...dialect.Condition) *Selector {
	if len(fns) == 0 || s.err != nil {
		return s
	}

	if s.where.Len() == 0 {
		s.where.WriteString("(")
	} else {
		s.where.WriteString(types.Operator_and + "(")
	}

	for i, fn := range fns {
		if i > 0 {
			s.where.WriteString(types.Operator_or)
		}
		cond, val := fn()
		//if v, ok := val.(error); ok {
		//	s.err = v
		//	return s
		//}
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

// Or
func (s *Selector) Or(fns ...dialect.Condition) *Selector {
	if len(fns) == 0 || s.err != nil {
		return s
	}

	if s.where.Len() == 0 {
		s.where.WriteString("(")
	} else {
		s.where.WriteString(types.Operator_or + "(")
	}

	for i, fn := range fns {
		if i > 0 {
			s.where.WriteString(types.Operator_and)
		}
		cond, val := fn()
		//if v, ok := val.(error); ok {
		//	s.err = v
		//	return s
		//}
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

// Order
func (s *Selector) Order(cols ...dialect.Field) *Selector {
	return s.Asc(cols...)
}

// Order Asc
func (s *Selector) Asc(cols ...dialect.Field) *Selector {
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

// Order Desc
func (s *Selector) Desc(cols ...dialect.Field) *Selector {
	if len(cols) == 0 {
		return s
	}
	for _, col := range cols {
		if s.orderBy.Len() > 0 {
			s.orderBy.WriteByte(',')
		}
		s.orderBy.WriteString(col.Quote() + " DESC")
	}
	return s
}

// Group
func (s *Selector) Group(cols ...dialect.Field) *Selector {
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
func (s *Selector) Having(fns ...dialect.Condition) *Selector {
	if len(fns) == 0 || s.err != nil {
		return s
	}

	s.having.WriteString("(")
	for i, fn := range fns {
		if i > 0 {
			s.having.WriteString(types.Operator_and)
		}
		cond, val := fn()
		//if v, ok := val.(error); ok {
		//	s.err = v
		//	return s
		//}
		s.having.WriteString(cond)
		if vals, ok := val.([]any); ok {
			s.havingParams = append(s.havingParams, vals...)
		}
	}
	return s
}
func (s *Selector) Limit(size uint, start ...uint) *Selector {
	if len(start) > 0 {
		s.limit = " LIMIT " + strconv.Itoa(int(size)) + " OFFSET " + strconv.Itoa(int(start[0]))
	} else {
		s.limit = " LIMIT " + strconv.Itoa(int(size))
	}

	return s
}

func (s *Selector) parse() {
	s.command.WriteString("SELECT ")

	colens := len(s.cols)
	funlens := len(s.funcs)
	if colens+funlens == 0 {
		s.command.WriteString("*")
	} else {
		if s.distinct {
			s.command.WriteString("DISTINCT ")
		}
		for i, col := range s.cols {
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
}

// Query
func (s *Selector) Query(ctx context.Context) (*sql.Rows, error) {
	defer s.Free()
	if s.err != nil {
		return nil, s.err
	}

	s.parse()
	s.whereParams = append(s.whereParams, s.havingParams...)

	stmt, err := s.db.PrepareContext(ctx, s.command.String())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryContext(ctx, s.whereParams...)
}

// QueryRow
func (s *Selector) QueryRow(ctx context.Context) (*sql.Row, error) {
	defer s.Free()
	if s.err != nil {
		return nil, s.err
	}

	s.parse()
	s.whereParams = append(s.whereParams, s.havingParams...)

	stmt, err := s.db.PrepareContext(ctx, s.command.String())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRowContext(ctx, s.whereParams...), nil
}

// Count
func (s *Selector) Count(ctx context.Context, cond ...dialect.Condition) (int64, error) {
	defer s.Free()
	if s.err != nil {
		return 0, s.err
	}

	s.Where(cond...)
	s.command.WriteString("SELECT COUNT(*)")
	// FROM TABLE
	s.command.WriteString(" FROM " + dialect.Quote_Char + s.table + dialect.Quote_Char)
	// WHERE
	if s.where.Len() > 0 {
		s.command.WriteString(" WHERE " + s.where.String())
	}

	stmt, err := s.db.PrepareContext(ctx, s.command.String())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, s.whereParams...)
	var count int64
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Sum
func (s *Selector) Sum(ctx context.Context, col dialect.Field, cond ...dialect.Condition) (int64, error) {
	defer s.Free()
	if s.err != nil {
		return 0, s.err
	}

	s.Funcs(col.Sum()).Where(cond...)
	s.command.WriteString("SELECT ")
	s.command.WriteString(s.funcs[0])
	// FROM TABLE
	s.command.WriteString(" FROM " + dialect.Quote_Char + s.table + dialect.Quote_Char)
	// WHERE
	if s.where.Len() > 0 {
		s.command.WriteString(" WHERE " + s.where.String())
	}

	stmt, err := s.db.PrepareContext(ctx, s.command.String())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, s.whereParams...)
	var sum int64
	err = row.Scan(&sum)
	if err != nil {
		return 0, err
	}
	return sum, nil
}
