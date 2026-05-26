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
	"github.com/linbaozhong/gentity/pkg/log"
	"strings"
)

type CreateBuilder interface {
	Columner

	// Set 设置要插入的字段和值
	//
	// 参数说明:
	//   - fns: 可变数量的 Setter 函数，用于定义字段及其对应的值
	//          例如: tblUsers.Name.Set("张三"), tblUsers.Age.Set(18)
	//
	// 返回值说明:
	//   - Builder: 返回构建器实例，支持链式调用
	//
	// 使用示例:
	//   db.Table("users").Set(tblUsers.Name.Set("张三"), tblUsers.Age.Set(18)).Create().Exec(ctx)
	Set(fns ...dialect.Setter) Builder

	// SetExpr 使用表达式设置要更新的字段和值
	//
	// 参数说明:
	//   - fns: 可变数量的 Setter 函数，用于定义字段及其对应的表达式
	//          支持自增(incr)、自减(decr)、替换(replace)等表达式操作
	//          例如: tblUsers.Age.Incr(), tblUsers.Name.Replace("旧", "新")
	//
	// 返回值说明:
	//   - Builder: 返回构建器实例，支持链式调用
	//
	// 注意:
	//   - 该方法已废弃，请使用 Set 方法替代
	//   - Set 方法现在同时支持普通赋值和表达式赋值
	//
	// 使用示例:
	//   db.Table("users").Where(condition).Set(tblUsers.Age.Incr()).Update().Exec(ctx)
	SetExpr(...dialect.Setter) Builder

	// Create 创建插入器
	//
	// 参数说明:
	//   - db: 可选的数据库连接对象，如果 Builder 实例化时未传入 DB，则此处必须传入
	//
	// 返回值说明:
	//   - Creater: 返回插入器接口，提供 Exec、Struct、BatchStruct 等插入方法
	//
	// 使用示例:
	//   // 方式1: 通过 Cols 指定列后执行
	//   db.Table("users").Cols(col1, col2).Create().Exec(ctx)
	//
	//   // 方式2: 通过结构体插入
	//   db.Table("users").Create().Struct(ctx, &user)
	//
	//   // 方式3: 批量插入
	//   db.Table("users").Create().BatchStruct(ctx, &user1, &user2)
	Create(...*DB) Creater

	// Insert 创建插入器（Create 的别名方法）
	//
	// 参数说明:
	//   - db: 可选的数据库连接对象，如果 Builder 实例化时未传入 DB，则此处必须传入
	//
	// 返回值说明:
	//   - Creater: 返回插入器接口，提供 Exec、Struct、BatchStruct 等插入方法
	//
	// 注意:
	//   - 该方法与 Create 功能完全相同，仅提供不同的命名风格
	//   - 推荐使用 Create 以保持代码一致性
	//
	// 使用示例:
	//   db.Table("users").Insert().Struct(ctx, &user)
	Insert(...*DB) Creater
}

type Creater interface {

	// Exec 执行插入操作
	//
	// 该方法根据之前通过 Set 等方法设置的列和值，生成并执行 INSERT SQL 语句。
	// 必须在调用 Exec 之前至少设置一个列（通过 Set 方法）。
	//
	// 参数说明:
	//   - ctx: 上下文对象，用于控制请求的生命周期、超时和取消操作
	//
	// 返回值说明:
	//   - sql.Result: SQL 执行结果对象，可通过 LastInsertId() 获取最后插入的 ID，
	//                 通过 RowsAffected() 获取受影响的行数
	//   - error: 错误信息，如果插入成功则为 nil；可能的错误包括：
	//              * dialect.ErrCreateEmpty: 未设置任何列时返回
	//              * Err_ToSql: 调试模式下仅返回 SQL 不执行时返回
	//              * 数据库相关错误：如连接失败、约束违反等
	//
	// 使用示例:
	//   // 通过 Set 设置值
	//   result, err := db.Table("users").
	//     Set(tblUsers.Name.Set("张三"), tblUsers.Age.Set(18)).
	//     Create().
	//     Exec(ctx)
	//
	// 注意:
	//   - 该方法执行后会自动释放 Builder 资源，不可再次使用
	//   - 如果在事务中使用，请确保在事务提交前完成所有操作
	//   - 调试模式（Debug）下不会实际执行插入，仅打印 SQL 语句
	Exec(context.Context) (sql.Result, error)

	// Insert(context.Context, dialect.Modeler) (sql.Result, error)

	// Struct 执行插入一个结构体
	//
	// 参数说明:
	//   - ctx: 上下文对象，用于控制请求的生命周期和取消操作
	//   - bean: 实现了 dialect.Modeler 接口的数据模型对象，包含要插入的数据
	//
	// 返回值说明:
	//   - sql.Result: SQL 执行结果，包含最后插入 ID 和受影响的行数
	//   - error: 错误信息，如果插入成功则为 nil
	Struct(context.Context, dialect.Modeler) (sql.Result, error)
	// InsertBatch(context.Context, ...dialect.Modeler) (sql.Result, error)

	// BatchStruct 执行批量插入，请不要在事务中使用
	//
	// 参数说明:
	//   - ctx: 上下文对象，用于控制请求的生命周期和取消操作
	//   - beans: 实现了 dialect.Modeler 接口的数据模型对象切片，包含要批量插入的多条数据
	//
	// 返回值说明:
	//   - sql.Result: SQL 执行结果，包含最后插入 ID 和受影响的行数
	//   - error: 错误信息，如果批量插入成功则为 nil
	BatchStruct(context.Context, ...dialect.Modeler) (sql.Result, error)
}

type create struct {
	*orm
}

func (o *orm) Insert(x ...*DB) Creater {
	if len(x) > 0 {
		o.db = x[0]
	}
	return &create{
		orm: o,
	}
}

// Create 创建插入器
func (o *orm) Create(x ...*DB) Creater {
	if len(x) > 0 {
		o.db = x[0]
	}
	return &create{
		orm: o,
	}
}

// Exec 执行插入
func (c *create) Exec(ctx context.Context) (sql.Result, error) {
	defer c.Free()
	if c.err != nil {
		return nil, c.err
	}

	lens := len(c.cols)
	if lens == 0 {
		return nil, dialect.ErrCreateEmpty
	}

	d := c.db.Dialect()
	c.command.WriteString("INSERT INTO " + d.Quote(c.table) + " (")
	for i, col := range c.cols {
		if i > 0 {
			c.command.WriteString(",")
		}
		c.command.WriteString(col.Quote(d))
	}
	c.command.WriteString(") VALUES ")

	values := make([]string, lens)
	for i := range values {
		values[i] = d.Placeholder(&c.paramIndex)
	}
	c.command.WriteString("(" + strings.Join(values, ",") + ")")
	// 只返回SQL语句，不执行
	if c.debug || c.db.Debug() {
		log.Info(c.String())
		return &noRows{}, Err_ToSql
	}

	// 执行SQL语句
	stmt, err := c.db.PrepareContext(ctx, c.command.String())
	if err != nil {
		return nil, err
	}
	if c.db.IsDB() {
		defer stmt.Close()
	}

	return stmt.ExecContext(ctx, c.params...)
}

// Insert 插入单条记录到数据库表中
//
// 参数说明:
//   - ctx: 上下文对象，用于控制请求的生命周期和取消操作
//   - bean: 实现了 dialect.Modeler 接口的数据模型对象，包含要插入的数据
//
// 返回值说明:
//   - sql.Result: SQL 执行结果，包含最后插入 ID 和受影响的行数
//   - error: 错误信息，如果插入成功则为 nil
//
// 注意: 该方法内部调用 Struct 方法执行实际的插入操作
// func (c *create) Insert(ctx context.Context, bean dialect.Modeler) (sql.Result, error) {
// 	return c.Struct(ctx, bean)
// }

// Struct 执行插入一个结构体
//
// 参数说明:
//   - ctx: 上下文对象，用于控制请求的生命周期和取消操作
//   - bean: 实现了 dialect.Modeler 接口的数据模型对象，包含要插入的数据
//
// 返回值说明:
//   - sql.Result: SQL 执行结果，包含最后插入 ID 和受影响的行数
//   - error: 错误信息，如果插入成功则为 nil
func (c *create) Struct(ctx context.Context, bean dialect.Modeler) (sql.Result, error) {
	defer c.Free()
	if c.err != nil {
		return nil, c.err
	}

	d := c.db.Dialect()
	c.command.WriteString("INSERT INTO " + d.Quote(c.table) + " (")

	var _cols []string
	_cols, c.params = bean.AssignValues(d, c.cols...)
	_colLens := len(_cols)
	c.command.WriteString(strings.Join(_cols, ","))
	c.command.WriteString(") VALUES ")
	//
	values := make([]string, _colLens)
	for i := range values {
		values[i] = d.Placeholder(&c.paramIndex)
	}
	c.command.WriteString("(" + strings.Join(values, ",") + ")")
	// 只返回SQL语句，不执行
	if c.debug || c.db.Debug() {
		log.Info(c.String())
		return &noRows{}, Err_ToSql
	}

	// 执行SQL语句
	stmt, err := c.db.PrepareContext(ctx, c.command.String())
	if err != nil {
		return nil, err
	}

	if c.db.IsDB() {
		defer stmt.Close()
	}

	return stmt.ExecContext(ctx, c.params...)
}

// InsertBatch 批量插入多条记录到数据库表中
//
// 参数说明:
//   - ctx: 上下文对象，用于控制请求的生命周期和取消操作
//   - beans: 实现了 dialect.Modeler 接口的数据模型对象切片，包含要批量插入的多条数据
//
// 返回值说明:
//   - sql.Result: SQL 执行结果，包含最后插入 ID 和受影响的行数
//   - error: 错误信息，如果批量插入成功则为 nil
//
// 注意: 该方法内部调用 BatchStruct 方法执行实际的批量插入操作，建议不要在事务中使用
// func (c *create) InsertBatch(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error) {
// 	return c.BatchStruct(ctx, beans...)
// }

// BatchStruct 执行批量插入，请不要在事务中使用
//
// 参数说明:
//   - ctx: 上下文对象，用于控制请求的生命周期和取消操作
//   - beans: 实现了 dialect.Modeler 接口的数据模型对象切片，包含要批量插入的多条数据
//
// 返回值说明:
//   - sql.Result: SQL 执行结果，包含最后插入 ID 和受影响的行数
//   - error: 错误信息，如果批量插入成功则为 nil
func (c *create) BatchStruct(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error) {
	defer c.Free()
	if c.err != nil {
		return nil, c.err
	}

	lens := len(beans)
	if lens == 0 {
		return nil, dialect.ErrBeanEmpty
	}

	d := c.db.Dialect()
	c.command.WriteString("INSERT INTO " + d.Quote(c.table) + " (")

	var _cols []string
	_cols, c.params = beans[0].RawAssignValues(d, c.cols...)
	_colLens := len(_cols)
	c.command.WriteString(strings.Join(_cols, ","))
	c.command.WriteString(") VALUES ")
	values := make([]string, _colLens)
	for i := range values {
		values[i] = d.Placeholder(&c.paramIndex)
	}
	c.command.WriteString("(" + strings.Join(values, ",") + ")")
	// 只返回SQL语句，不执行
	if c.debug || c.db.Debug() {
		log.Info(c.String())
		return &noRows{}, Err_ToSql
	}

	// 启动事务批量执行Create
	ret, err := c.db.Transaction(ctx, func(tx *Tx) (any, error) {
		stmt, err := tx.PrepareContext(ctx, c.command.String())
		if err != nil {
			return nil, err
		}
		if c.db.IsDB() {
			defer stmt.Close()
		}

		result, err := stmt.ExecContext(ctx, c.params...)
		if err != nil {
			return nil, err
		}
		beans[0].AssignPrimaryKeyValues(result)

		for i := 1; i < lens; i++ {
			bean := beans[i]
			if bean == nil {
				return nil, dialect.ErrBeanEmpty
			}
			_, c.params = bean.RawAssignValues(d, c.cols...)
			result, err = stmt.ExecContext(ctx, c.params...)
			if err != nil {
				return nil, err
			}
			bean.AssignPrimaryKeyValues(result)
		}
		return result, nil
	})
	if err != nil {
		return nil, err
	}
	if result, ok := ret.(sql.Result); ok {
		return result, nil
	}
	return nil, err
}
