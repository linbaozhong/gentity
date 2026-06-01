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

		// Debug 设置调试模式
		//
		// 启用调试模式后，SQL 语句将被打印到日志中，但不会实际执行数据库操作。
		// 这对于调试和查看生成的 SQL 语句非常有用。
		//
		// 参数说明:
		//   - b: 可变参数，用于控制调试模式的开关
		//        * 不传参数：启用调试模式（等同于传 true）
		//        * 传 true：启用调试模式
		//        * 传 false：禁用调试模式
		//
		// 返回值说明:
		//   - Builder: 返回构建器实例，支持链式调用
		//
		// 使用示例:
		//   // 示例1: 启用调试模式（不传参数）
		//   db.Table("users").
		//     Where(tblUsers.Id.Eq(1)).
		//     Set(tblUsers.Name.Set("张三")).
		//     Update().
		//     Debug().  // 仅打印 SQL，不执行
		//     Exec(ctx)
		//
		//   // 示例2: 明确启用调试模式
		//   db.Table("users").
		//     Debug(true).
		//     Where(tblUsers.Id.Eq(1)).
		//     Update().
		//     Exec(ctx)
		//
		// 注意:
		//   - 调试模式下，Exec/Struct/BatchStruct 等方法会返回 Err_ToSql 错误
		//   - 可以通过 DB.Debug() 方法全局设置调试模式
		//   - 调试模式下的 SQL 语句会记录到日志系统（log.Info）
		Debug(...bool) Builder

		SelectBuilder
		CreateBuilder
		UpdateBuilder
		DeleteBuilder

		parse() (strings.Builder, []any, error)
	}
	join struct {
		joinType   dialect.JoinType
		table      dialect.Field
		left       dialect.Field
		right      dialect.Field
		conditions []dialect.Condition
	}
	order struct {
		col   dialect.Field
		order dialect.OrderType
	}
	orm struct {
		pool.Model
		db         Executer
		paramIndex uint16 // 参数索引计数器
		table      string
		// join         [][3]string
		join       []join
		joinParams []any
		distinct   bool
		cols       []dialect.Field
		funcs      []dialect.Function
		omits      []dialect.Field
		groupBy    []dialect.Field
		having     []dialect.Condition
		// havingParams []any
		orderBy []order
		limit   string
		// where        strings.Builder
		// 条件
		cond           []dialect.Condition
		whereParams    []any
		subQueryParams []any // 子查询参数
		exprCols       []expr
		params         []any
		command        strings.Builder
		// debug 为true时，仅打印SQL语句，不执行
		debug bool
		err   error
	}
)

var (
	ormPool = pool.New[*orm](func() any {
		obj := &orm{}
		return obj
	})
)

func newOrm(dbs ...*DB) *orm {
	obj := ormPool.Get()
	if len(dbs) > 0 {
		obj.db = dbs[0]
	}
	return obj
}

// func New(db *DB, opts ...Option) Builder {
// 	obj := newOrm()
// 	obj.db = db
//
// 	for _, opt := range opts {
// 		opt(obj)
// 	}
// 	return obj
// }

// Free 释放 orm 对象，将其重置并放回对象池。
func (o *orm) Free() {
	if o == nil {
		return
	}

	if o.debug || (o.db != nil && o.db.Debug()) {
		log.Info(o.String())
	}

	ormPool.Put(o)
}

func (o *orm) Reset() {
	o.db = nil
	o.table = ""
	o.paramIndex = 0
	o.cols = o.cols[:0]   // []dialect.Field{} // o.cols[:0]
	o.funcs = o.funcs[:0] // []string{}       // o.funcs[:0]
	o.distinct = false
	o.join = o.join[:0]             // [][3]string{}      // o.join[:0]
	o.joinParams = o.joinParams[:0] // []any{}      // o.joinParams[:0]
	o.omits = o.omits[:0]           // []dialect.Field{} // o.omits[:0]
	o.cond = o.cond[:0]
	o.whereParams = o.whereParams[:0] // []any{} // o.whereParams[:0]
	o.subQueryParams = o.subQueryParams[:0]
	o.groupBy = o.groupBy[:0]
	o.having = o.having[:0]
	// o.havingParams = o.havingParams[:0] // []any{} // o.havingParams[:0]
	o.orderBy = o.orderBy[:0]
	o.limit = ""
	o.exprCols = o.exprCols[:0] // []expr{} // o.exprCols[:0]
	o.params = o.params[:0]     // []any{} // o.params[:0]
	o.command.Reset()
	o.debug = false
	o.err = nil
}

// String 返回 orm 对象的 SQL 语句和参数的字符串表示。
func (o *orm) String() string {
	return fmt.Sprintf("%s  %v \n", o.command.String(), o.mergeParams())
}

// SetDB 设置 orm 对象的数据库连接。
func (o *orm) SetDB(d *DB) Builder {
	o.db = d
	return o
}

// Table 设置 orm 对象的表名。
// 如果 a 是字符串，代表数据库表名
// 如果 a 是实现了 TableNamer 接口的结构体，可以通过 TableName()方法提取数据库表名
// 如果 a 是 Builder 接口的对象，表示该查询使用了子查询。
func (o *orm) Table(a any, as ...string) Builder {
	switch v := a.(type) {
	case string:
		o.table = v
	case dialect.TableNamer:
		o.table = v.TableName()
	case Builder:
		cmd, params, e := v.parse()
		if e != nil {
			o.err = e
			return o
		}
		o.table = "(" + cmd.String() + ")"
		o.subQueryParams = append(o.subQueryParams, params...)
		if len(as) > 0 {
			o.table = fmt.Sprintf("%s AS %s", o.table, as[0])
		}
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

// GetTableName 获取 orm 对象的表名。
func (s *orm) GetTableName() string {
	return s.table
}

// Set
// 用于设置更新语句中的字段和值
// 例如：Set(dialect.F("name", "哈利波特"))
func (o *orm) Set(fns ...dialect.Setter) Builder {
	l := len(fns)
	if l == 0 || o.err != nil {
		return o
	}

	tmpCols := make([]dialect.Field, len(o.cols), len(o.cols)+l)
	copy(tmpCols, o.cols)
	tmpParams := make([]any, len(o.params), len(o.params)+l)
	copy(tmpParams, o.params)
	tmpExprCols := make([]expr, len(o.exprCols), len(o.exprCols)+l)
	copy(tmpExprCols, o.exprCols)

	d := o.db.Dialect()
	for _, fn := range fns {
		c, val, op := fn()
		if e, ok := val.(error); ok {
			o.err = e
			return o
		}
		if op == dialect.Op_Normal {
			tmpCols = append(tmpCols, c)
			tmpParams = append(tmpParams, val)
		} else {
			ex, val, e := dialect.ParseSetter(fn, &o.paramIndex, d)
			if e != nil {
				o.err = e
				return o
			}
			tmpExprCols = append(tmpExprCols, expr{colName: ex, arg: val})
		}
	}
	o.cols = tmpCols
	o.params = tmpParams
	o.exprCols = tmpExprCols
	return o
}

// Deprecated: 请使用 Set
//
// SetExpr
// 用于设置更新语句中的表达式
// 例如：SetExpr(dialect.Expr("age", "age + 1"))
func (o *orm) SetExpr(fns ...dialect.Setter) Builder {
	l := len(fns)
	if l == 0 || o.err != nil {
		return o
	}

	tmpExprCols := make([]expr, len(o.exprCols), len(o.exprCols)+l)
	copy(tmpExprCols, o.exprCols)

	d := o.db.Dialect()
	for _, fn := range fns {
		ex, val, e := dialect.ParseSetter(fn, &o.paramIndex, d)
		if e != nil {
			o.err = e
			return o
		}
		tmpExprCols = append(tmpExprCols, expr{colName: ex, arg: val})
	}
	o.exprCols = tmpExprCols
	return o
}

// Debug 不传参数或者参数为true时，仅打印SQL语句，不执行
func (o *orm) Debug(b ...bool) Builder {
	if len(b) > 0 {
		o.debug = b[0]
	} else {
		o.debug = true
	}
	return o
}

// Clone 克隆 orm
func (o *orm) Clone() Builder {
	// 创建一个新的 orm 实例，并复制非引用类型字段
	newOrm := *o

	// 复制切片类型的字段，创建新的底层数组
	newOrm.cols = make([]dialect.Field, len(o.cols))
	copy(newOrm.cols, o.cols)

	newOrm.funcs = make([]dialect.Function, len(o.funcs))
	copy(newOrm.funcs, o.funcs)

	// newOrm.join = make([][3]string, len(o.join))
	newOrm.join = make([]join, len(o.join))
	copy(newOrm.join, o.join)

	newOrm.joinParams = make([]any, len(o.joinParams))
	copy(newOrm.joinParams, o.joinParams)

	newOrm.omits = make([]dialect.Field, len(o.omits))
	copy(newOrm.omits, o.omits)

	newOrm.whereParams = make([]any, len(o.whereParams))
	copy(newOrm.whereParams, o.whereParams)

	newOrm.subQueryParams = make([]any, len(o.subQueryParams))
	copy(newOrm.subQueryParams, o.subQueryParams)

	newOrm.exprCols = make([]expr, len(o.exprCols))
	copy(newOrm.exprCols, o.exprCols)

	newOrm.params = make([]any, len(o.params))
	copy(newOrm.params, o.params)

	// 复制 strings.Builder 类型的字段
	newOrm.groupBy = make([]dialect.Field, len(o.groupBy))
	copy(newOrm.groupBy, o.groupBy)

	newOrm.having = make([]dialect.Condition, len(o.having))
	copy(newOrm.having, o.having)

	newOrm.orderBy = make([]order, len(o.orderBy))
	copy(newOrm.orderBy, o.orderBy)

	newOrm.cond = make([]dialect.Condition, len(o.cond))
	copy(newOrm.cond, o.cond)

	newOrm.command.Reset()
	newOrm.command.WriteString(o.command.String())

	return &newOrm
}

// 合并参数
func (o *orm) mergeParams() []any {
	params := make([]any, 0, len(o.joinParams)+len(o.subQueryParams)+len(o.params)+len(o.whereParams))

	params = append(params, o.joinParams...)
	params = append(params, o.subQueryParams...)
	params = append(params, o.params...)
	params = append(params, o.whereParams...)

	return params
}

// parse
func (o *orm) parse() (strings.Builder, []any, error) {
	d := o.db.Dialect()
	o.command.Reset()
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
			o.command.WriteString(col.Quote(d))
		}
		if colens > 0 && funlens > 0 {
			o.command.WriteString(",")
		}
		o.command.WriteString(strings.Join(o.parseFunc(o.funcs), ","))
	}

	// FROM TABLE
	if o.table == "" {
		if colens > 0 {
			o.table = cols[0].Table
		} else {
			o.err = Err_TableName
			return o.command, nil, o.err
		}
	}
	if strings.HasPrefix(o.table, "(") {
		// 如果表名以(开头，则不加引号
		o.command.WriteString(" FROM " + o.table)
	} else {
		o.command.WriteString(" FROM " + d.Quote(o.table))
	}

	if len(o.join) > 0 {
		joinStr, params, e := o.parseJoin(o.join)
		if e != nil {
			o.err = e
			return o.command, nil, o.err
		}
		o.joinParams = params
		if joinStr.Len() > 0 {
			o.command.WriteString(joinStr.String())
		}
	}

	// WHERE
	if len(o.cond) > 0 {
		where, params, e := o.parseCond(o.cond)
		if e != nil {
			o.err = e
			return o.command, nil, o.err
		}
		if where.Len() > 0 {
			o.command.WriteString(" WHERE " + where.String())
			o.whereParams = params
		}
	}
	// GROUP BY
	if len(o.groupBy) > 0 {
		o.command.WriteString(" GROUP BY ")
		for i, col := range o.groupBy {
			if i > 0 {
				o.command.WriteByte(',')
			}
			o.command.WriteString(col.Quote(d))
		}

		// HAVING
		if len(o.having) > 0 {
			where, havingParams, e := o.parseCond(o.having)
			if e != nil {
				o.err = e
				return o.command, nil, o.err
			}
			if where.Len() > 0 {
				o.command.WriteString(" HAVING " + where.String())
				o.subQueryParams = append(o.subQueryParams, havingParams...)
			}
		}
	}
	// ORDER BY
	if len(o.orderBy) > 0 {
		o.command.WriteString(" ORDER BY ")
		for i, ord := range o.orderBy {
			if i > 0 {
				o.command.WriteByte(',')
			}
			o.command.WriteString(ord.col.Quote(d) + " " + ord.order.String())
		}
	}

	// LIMIT
	if o.limit != "" {
		o.command.WriteString(o.limit)
	}

	return o.command, o.mergeParams(), nil
}

// query
func (o *orm) query(ctx context.Context) (*sql.Rows, error) {
	cmd, params, e := o.parse()
	if e != nil {
		return nil, e
	}
	return o.rows(ctx, cmd.String(), params...)
}

func (o *orm) rows(ctx context.Context, sqlStr string, params ...any) (*sql.Rows, error) {
	if o.debug || (o.db != nil && o.db.Debug()) {
		log.Info(o.String())
		return &sql.Rows{}, Err_ToSql
	}
	stmt, err := o.db.PrepareContext(ctx, sqlStr)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	return stmt.QueryContext(ctx, params...)
}

func (o *orm) row(ctx context.Context, sqlStr string, params ...any) (*sql.Row, error) {
	if o.debug || (o.db != nil && o.db.Debug()) {
		log.Info(o.String())
		return &sql.Row{}, Err_ToSql
	}
	stmt, err := o.db.PrepareContext(ctx, sqlStr)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	return stmt.QueryRowContext(ctx, params...), nil
}
