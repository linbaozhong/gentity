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

type UpdateBuilder interface {
	Columner

	// Set 设置要更新的字段和值
	//
	// 该方法用于指定 UPDATE 语句中 SET 子句的字段及其对应的值。支持普通赋值和表达式赋值两种方式：
	// - 普通赋值：直接设置字段值为指定值，如 tblUsers.Name.Set("张三")
	// - 表达式赋值：使用 SQL 表达式，如自增(tblUsers.Age.Incr())、自减、字符串替换等
	//
	// 参数说明:
	//   - fns: 可变数量的 Setter 函数，每个函数定义一个字段及其值的设置方式
	//          例如: tblUsers.Name.Set("张三"), tblUsers.Age.Set(18), tblUsers.Score.Incr(10)
	//
	// 返回值说明:
	//   - Builder: 返回构建器实例，支持链式调用
	//
	// 使用示例:
	//   // 普通赋值
	//   db.Table("users").
	//     Where(tblUsers.Id.Eq(1)).
	//     Set(tblUsers.Name.Set("张三"), tblUsers.Age.Set(18)).
	//     Update().
	//     Exec(ctx)
	//
	//   // 表达式赋值（自增）
	//   db.Table("users").
	//     Where(tblUsers.Id.Eq(1)).
	//     Set(tblUsers.ViewCount.Incr(1)).
	//     Update().
	//     Exec(ctx)
	//
	//   // 混合使用
	//   db.Table("users").
	//     Where(tblUsers.Id.Eq(1)).
	//     Set(
	//       tblUsers.Name.Set("李四"),
	//       tblUsers.LoginCount.Incr(1),
	//       tblUsers.LastLoginTime.Set(time.Now()),
	//     ).
	//     Update().
	//     Exec(ctx)
	//
	// 注意:
	//   - 必须在调用 Update() 之前至少设置一个字段
	//   - 如果未设置任何字段，执行时会返回 dialect.ErrCreateEmpty 错误
	//   - Set 方法现在同时支持普通赋值和表达式赋值，推荐统一使用 Set
	Set(fns ...dialect.Setter) Builder

	// SetExpr 使用表达式设置要更新的字段
	//
	// 该方法专门用于设置基于 SQL 表达式的更新操作，如字段自增、自减、字符串替换等。
	// 与 Set 方法不同，SetExpr 只处理表达式类型的赋值，不支持普通值赋值。
	//
	// 参数说明:
	//   - fns: 可变数量的 Setter 函数，每个函数定义一个字段的表达式更新方式
	//          支持的表达式包括: Incr(自增), Decr(自减), Replace(替换), Concat(拼接) 等
	//          例如: tblUsers.Age.Incr(), tblUsers.Name.Replace("旧", "新")
	//
	// 返回值说明:
	//   - Builder: 返回构建器实例，支持链式调用
	//
	// 使用示例:
	//   // 字段自增
	//   db.Table("articles").
	//     Where(tblArticles.Id.Eq(1)).
	//     SetExpr(tblArticles.ViewCount.Incr(1)).
	//     Update().
	//     Exec(ctx)
	//
	//   // 字段自减
	//   db.Table("products").
	//     Where(tblProducts.Id.Eq(1)).
	//     SetExpr(tblProducts.Stock.Decr(5)).
	//     Update().
	//     Exec(ctx)
	//
	//   // 字符串替换
	//   db.Table("users").
	//     Where(tblUsers.Id.Eq(1)).
	//     SetExpr(tblUsers.Email.Replace("@old.com", "@new.com")).
	//     Update().
	//     Exec(ctx)
	//
	// 注意:
	//   - 该方法已废弃（Deprecated），请使用 Set 方法替代
	//   - Set 方法现在同时支持普通赋值和表达式赋值，功能更全面
	//   - 为了保持代码一致性和可维护性，建议统一使用 Set 方法
	//   - 该方法仍然可用，以确保向后兼容
	SetExpr(fns ...dialect.Setter) Builder

	Wherer

	// Update 创建更新器
	//
	// 该方法将当前的 Builder 转换为更新器（Updater），用于执行 UPDATE 操作。
	// 在调用 Update() 之前，应该已经通过 Table() 设置了表名，通过 Where() 设置了更新条件，
	// 并通过 Set() 或 SetExpr() 设置了要更新的字段。
	//
	// 参数说明:
	//   - db: 可选的数据库连接对象。如果 Builder 实例化时已经传入了 DB，则此处可以省略；
	//         如果实例化时未传入 DB，则此处必须传入，否则执行时会出错
	//
	// 返回值说明:
	//   - Updater: 返回更新器接口，提供以下执行方法:
	//                * Exec(ctx): 执行更新操作，适用于已通过 Set 设置字段的场景
	//                * Struct(ctx, bean): 根据结构体自动提取字段和主键进行更新
	//                * BatchStruct(ctx, beans...): 批量更新多个结构体记录
	//
	// 使用示例:
	//   // 示例1: 使用 Set + Exec 更新
	//   result, err := db.Table("users").
	//     Where(tblUsers.Id.Eq(1)).
	//     Set(tblUsers.Name.Set("张三")).
	//     Update().
	//     Exec(ctx)
	//
	//   // 示例2: 使用 Struct 自动更新（根据主键）
	//   user.Name = "李四"
	//   result, err := db.Table("users").
	//     Update().
	//     Struct(ctx, &user)
	//
	//   // 示例3: 批量更新
	//   result, err := db.Table("users").
	//     Update().
	//     BatchStruct(ctx, &user1, &user2, &user3)
	//
	//   // 示例4: 如果实例化时未传入 DB，则此处必须传入
	//   builder := ace.New(nil)
	//   result, err := builder.Table("users").
	//     Where(tblUsers.Id.Eq(1)).
	//     Set(tblUsers.Name.Set("张三")).
	//     Update(db).  // 这里必须传入 db
	//     Exec(ctx)
	//
	// 注意:
	//   - 建议在调用 Update() 之前设置 Where 条件，避免全表更新
	//   - 如果使用 Struct 方法，bean 必须实现 dialect.Modeler 接口
	//   - BatchStruct 内部会使用事务处理，不建议在外部事务中调用
	Update(...*DB) Updater
}

type Updater interface {
	// Exec 执行更新操作
	//
	// 该方法用于执行已构建的 UPDATE SQL 语句。在调用 Exec 之前，必须满足以下条件：
	// 1. 已通过 Table() 设置表名
	// 2. 已通过 Set() 或 SetExpr() 设置至少一个要更新的字段
	// 3. 建议通过 Where() 设置更新条件，避免全表更新
	//
	// 参数说明:
	//   - ctx: 上下文对象，用于控制请求的生命周期、超时和取消操作
	//
	// 返回值说明:
	//   - sql.Result: SQL 执行结果对象，可通过以下方法获取详细信息:
	//                   * RowsAffected() (int64, error): 获取受影响的行数
	//                   * LastInsertId() (int64, error): 获取最后插入的 ID（更新操作通常不适用）
	//   - error: 错误信息，可能的错误包括:
	//              * dialect.ErrCreateEmpty: 未设置任何要更新的字段
	//              * 数据库连接错误、SQL 语法错误等
	//              * Err_ToSql: 调试模式下仅打印 SQL，不执行
	//
	// 使用示例:
	//   // 示例1: 基本更新操作
	//   result, err := db.Table("users").
	//     Where(tblUsers.Id.Eq(1)).
	//     Set(tblUsers.Name.Set("张三"), tblUsers.Age.Set(18)).
	//     Update().
	//     Exec(ctx)
	//   if err != nil {
	//     log.Fatal(err)
	//   }
	//   rows, _ := result.RowsAffected()
	//   fmt.Printf("更新了 %d 条记录\n", rows)
	//
	//   // 示例2: 使用表达式更新（自增）
	//   result, err := db.Table("articles").
	//     Where(tblArticles.Id.Eq(1)).
	//     Set(tblArticles.ViewCount.Incr(1)).
	//     Update().
	//     Exec(ctx)
	//
	//   // 示例3: 带条件的批量更新
	//   result, err := db.Table("users").
	//     Where(tblUsers.Status.Eq(1)).
	//     Set(tblUsers.Status.Set(1), tblUsers.UpdateTime.Set(time.Now())).
	//     Update().
	//     Exec(ctx)
	//
	//   // 示例4: 多个字段混合更新
	//   result, err := db.Table("products").
	//     Where(tblProducts.Id.Eq(100)).
	//     Set(
	//       tblProducts.Price.Set(99.99),
	//       tblProducts.Stock.Decr(5),
	//       tblProducts.SalesCount.Incr(5),
	//     ).
	//     Update().
	//     Exec(ctx)
	//
	// 注意:
	//   - 必须在调用 Update() 之前至少设置一个字段，否则返回 dialect.ErrCreateEmpty 错误
	//   - 强烈建议在调用 Exec 前设置 Where 条件，避免意外全表更新
	//   - 如果未设置 Where 条件，将更新表中的所有记录
	//   - 该方法会自动释放 Builder 资源（调用 Free()），不可重复使用
	//   - 调试模式下（Debug(true)），仅打印 SQL 语句，返回 Err_ToSql 错误
	Exec(ctx context.Context) (sql.Result, error)

	// Update(ctx context.Context, bean dialect.Modeler) (sql.Result, error)

	// Struct 根据结构体自动更新单条记录
	//
	// 该方法从实现了 dialect.Modeler 接口的结构体中自动提取字段值和主键信息，
	// 生成并执行 UPDATE SQL 语句。系统会自动使用结构体的主键作为 WHERE 条件，
	// 确保只更新对应的单条记录。
	//
	// 参数说明:
	//   - ctx: 上下文对象，用于控制请求的生命周期、超时和取消操作
	//   - bean: 实现了 dialect.Modeler 接口的结构体指针，包含:
	//             * 要更新的字段值（非零值字段会被更新）
	//             * 主键字段值（用于定位要更新的记录）
	//             * 表名信息（如果未在 Builder 中设置表名，可从 bean 中提取）
	//
	// 返回值说明:
	//   - sql.Result: SQL 执行结果对象，可通过以下方法获取详细信息:
	//                   * RowsAffected() (int64, error): 获取受影响的行数
	//                     - 返回 0: 未找到匹配的记录
	//                     - 返回 1: 成功更新一条记录
	//                     - 返回 >1: 更新了多条记录（可能主键不唯一）
	//   - error: 错误信息，可能的错误包括:
	//              * sql.ErrNoRows: 未找到匹配主键的记录
	//              * 数据库连接错误、SQL 语法错误等
	//              * Err_ToSql: 调试模式下仅打印 SQL，不执行
	//
	// 使用示例:
	//   // 示例1: 基本更新操作
	//   user := &User{
	//     Id:   1,
	//     Name: "李四",
	//     Age:  25,
	//   }
	//   result, err := db.Table("users").
	//     Update().
	//     Struct(ctx, user)
	//   if err != nil {
	//     log.Fatal(err)
	//   }
	//
	//   // 示例2: 先查询再更新
	//   var user User
	//   db.Table("users").
	//     Where(tblUsers.Id.Eq(1)).
	//     Select().
	//     Get(ctx, &user)
	//
	//   user.Name = "新名字"
	//   user.Age = 30
	//   db.Table("users").
	//     Update().
	//     Struct(ctx, &user)
	//
	//   // 示例3: 指定要更新的列（只更新部分字段）
	//   user.Name = "新名字"
	//   user.Age = 30  // 这个字段不会被更新
	//   result, err := db.Table("users").
	//     Cols(tblUsers.Name).  // 只更新 Name 字段
	//     Update().
	//     Struct(ctx, &user)
	//
	//   // 示例4: 更新时添加额外条件
	//   user.Status = 1
	//   result, err := db.Table("users").
	//     Where(tblUsers.Version.Eq(user.Version)).  // 乐观锁
	//     Update().
	//     Struct(ctx, &user)
	//
	// 注意:
	//   - bean 必须是指针类型，且实现了 dialect.Modeler 接口
	//   - 系统会自动使用 bean 的主键字段作为 WHERE 条件
	//   - 默认情况下，bean 中所有非零值字段都会被更新
	//   - 可通过 Cols() 方法指定只更新特定字段
	//   - 可在调用 Struct 前通过 Where() 添加额外的更新条件
	//   - 该方法会自动释放 Builder 资源（调用 Free()），不可重复使用
	//   - 如果表中不存在对应主键的记录，RowsAffected() 返回 0
	Struct(ctx context.Context, bean dialect.Modeler) (sql.Result, error)

	// UpdateBatch(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error)

	// BatchStruct 批量更新多个结构体记录
	//
	// 该方法用于批量更新多条记录，内部使用事务机制确保数据一致性。
	// 所有更新操作要么全部成功，要么全部回滚。每个结构体的主键将作为
	// 对应记录的 WHERE 条件。
	//
	// 重要提示：该方法内部会启动事务，请不要在外部事务中调用此方法，
	// 否则可能导致嵌套事务问题或死锁。如需在事务中批量更新，请手动
	// 循环调用 Struct 方法。
	//
	// 参数说明:
	//   - ctx: 上下文对象，用于控制请求的生命周期、超时和取消操作
	//   - beans: 可变数量的模型对象参数，每个对象必须:
	//              * 是指针类型
	//              * 实现 dialect.Modeler 接口
	//              * 包含有效的主键值
	//              * 属于同一张表
	//            例如: &user1, &user2, &user3
	//
	// 返回值说明:
	//   - sql.Result: SQL 执行结果对象，通常是最后一个更新操作的结果:
	//                   * RowsAffected(): 最后一次更新的受影响行数
	//                   * 注意：不是所有更新的总和
	//   - error: 错误信息，可能的错误包括:
	//              * dialect.ErrBeanEmpty: beans 为空或包含 nil 元素
	//              * 事务错误：任何一个更新失败都会导致整个事务回滚
	//              * 数据库连接错误、SQL 语法错误等
	//              * Err_ToSql: 调试模式下仅打印 SQL，不执行
	//
	// 使用示例:
	//   // 示例1: 基本批量更新
	//   users := []*User{
	//     {Id: 1, Name: "用户1", Age: 20},
	//     {Id: 2, Name: "用户2", Age: 25},
	//     {Id: 3, Name: "用户3", Age: 30},
	//   }
	//   result, err := db.Table("users").
	//     Update().
	//     BatchStruct(ctx, users...)
	//   if err != nil {
	//     log.Fatal(err)
	//   }
	//
	//   // 示例2: 指定要更新的列
	//   result, err := db.Table("users").
	//     Cols(tblUsers.Name, tblUsers.Age).  // 只更新这两个字段
	//     Update().
	//     BatchStruct(ctx, &user1, &user2, &user3)
	//
	//   // 示例3: 从切片批量更新
	//   userList := []User{user1, user2, user3}
	//   ptrs := make([]*User, len(userList))
	//   for i := range userList {
	//     ptrs[i] = &userList[i]
	//   }
	//   result, err := db.Table("users").
	//     Update().
	//     BatchStruct(ctx, ptrs...)
	//
	//   // 示例4: 在事务外批量更新（推荐方式）
	//   // ✓ 正确：直接调用 BatchStruct
	//   result, err := db.Table("users").
	//     Update().
	//     BatchStruct(ctx, &user1, &user2, &user3)
	//
	//   // ✗ 错误：不要在外部事务中调用 BatchStruct
	//   // db.Transaction(ctx, func(tx *Tx) (any, error) {
	//   //   result, err := tx.Table("users").
	//   //     Update().
	//   //     BatchStruct(ctx, &user1, &user2, &user3)  // 可能导致嵌套事务
	//   //   return result, err
	//   // })
	//
	// 注意:
	//   - 所有 beans 必须属于同一张表
	//   - 每个 bean 必须有有效的主键值，用于定位要更新的记录
	//   - 方法内部使用事务，任何一个更新失败都会导致全部回滚
	//   - 不要在外部事务中调用此方法，避免嵌套事务问题
	//   - 如需在事务中批量更新，请手动循环调用 Struct 方法
	//   - 该方法会自动释放 Builder 资源（调用 Free()），不可重复使用
	//   - 大量数据批量更新时，建议分批处理，避免长事务
	//   - 返回的 sql.Result 是最后一次更新的结果，不是总和
	BatchStruct(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error)
}

type update struct {
	*orm
}

type expr struct {
	colName string
	arg     any
}

// Update 更新器
func (o *orm) Update(x ...*DB) Updater {
	if len(x) > 0 {
		o.db = x[0]
	}
	return &update{
		orm: o,
	}
}

// Exec 执行更新
func (u *update) Exec(ctx context.Context) (sql.Result, error) {
	defer u.Free()
	if u.err != nil {
		return nil, u.err
	}

	lens := len(u.cols) + len(u.exprCols)
	if lens == 0 {
		return nil, dialect.ErrCreateEmpty
	}

	d := u.db.Dialect()
	u.command.WriteString("UPDATE " + d.Quote(u.table) + " SET ")
	_cols := make([]string, 0, lens)
	for _, col := range u.cols {
		_cols = append(_cols, col.Quote(d)+" = "+d.Placeholder(&u.paramIndex))
	}
	for _, col := range u.exprCols {
		_cols = append(_cols, col.colName)
		if col.arg != nil {
			u.params = append(u.params, col.arg)
		}
	}
	u.command.WriteString(strings.Join(_cols, ","))

	where, params, e := u.parseCond(u.cond)
	if e != nil {
		u.err = e
		return nil, e
	}
	u.whereParams = params
	if where.Len() > 0 {
		u.command.WriteString(" WHERE " + where.String())
	}

	// 只返回SQL语句，不执行
	if u.debug || u.db.Debug() {
		log.Info(u.String())
		return &noRows{}, Err_ToSql
	}

	stmt, err := u.db.PrepareContext(ctx, u.command.String())
	if err != nil {
		return nil, err
	}
	if u.db.IsDB() {
		defer stmt.Close()
	}

	u.params = append(u.params, u.whereParams...)

	return stmt.ExecContext(ctx, u.params...)
}

// Update 更新单个结构体记录
// 该方法是对 Struct 方法的别名调用，用于根据 bean 的主键值更新对应的数据库记录
//
// 参数说明:
//   - ctx: 上下文对象，用于控制请求的生命周期和取消操作
//   - bean: 实现了 dialect.Modeler 接口的结构体指针，包含要更新的字段值和主键信息
//
// 返回值说明:
//   - sql.Result: SQL 执行结果，可通过 RowsAffected() 获取受影响的行数
//   - error: 错误信息，如果更新失败则返回具体错误，成功则为 nil
// func (u *update) Update(ctx context.Context, bean dialect.Modeler) (sql.Result, error) {
// 	return u.Struct(ctx, bean)
// }

// Struct 更新一个结构体
// 该方法是对 Struct 方法的别名调用，用于根据 bean 的主键值更新对应的数据库记录
//
// 参数说明:
//   - ctx: 上下文对象，用于控制请求的生命周期和取消操作
//   - bean: 实现了 dialect.Modeler 接口的结构体指针，包含要更新的字段值和主键信息
//
// 返回值说明:
//   - sql.Result: SQL 执行结果，可通过 RowsAffected() 获取受影响的行数
//   - error: 错误信息，如果更新失败则返回具体错误，成功则为 nil
func (u *update) Struct(ctx context.Context, bean dialect.Modeler) (sql.Result, error) {
	defer u.Free()
	if u.err != nil {
		return nil, u.err
	}

	d := u.db.Dialect()
	u.command.WriteString("UPDATE " + d.Quote(u.table) + " SET ")
	cols, vals := bean.AssignValues(d, u.cols...)
	for i, col := range cols {
		if i > 0 {
			u.command.WriteString(",")
		}
		u.command.WriteString(col + " = " + d.Placeholder(&u.paramIndex))
	}
	u.params = append(u.params, vals...)
	//
	keys, values := bean.AssignKeys()
	u.Where(keys.Eq(values))

	where, params, e := u.parseCond(u.cond)
	if e != nil {
		u.err = e
	}
	u.whereParams = params
	if where.Len() > 0 {
		u.command.WriteString(" WHERE " + where.String())
	}

	// 只返回SQL语句，不执行
	if u.debug || u.db.Debug() {
		log.Info(u.String())
		return &noRows{}, Err_ToSql
	}

	stmt, err := u.db.PrepareContext(ctx, u.command.String())
	if err != nil {
		return nil, err
	}
	if u.db.IsDB() {
		defer stmt.Close()
	}

	u.params = append(u.params, u.whereParams...)

	return stmt.ExecContext(ctx, u.params...)
}

// UpdateBatch 批量更新多个模型对象
// 该方法是 BatchStruct 的别名，提供语义化的命名
// 注意：该方法内部会使用事务处理，请不要在外部事务中调用
//
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期和取消操作
//   - beans: 可变参数，要更新的模型对象切片，每个对象必须实现 dialect.Modeler 接口
//
// 返回值:
//   - sql.Result: SQL执行结果，包含受影响的行数等信息
//   - error: 错误信息，如果更新失败则返回具体错误
// func (u *update) UpdateBatch(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error) {
// 	return u.BatchStruct(ctx, beans...)
// }

// BatchStruct 执行批量更新,请不要在事务中使用
// 该方法是 BatchStruct 的别名，提供语义化的命名
// 注意：该方法内部会使用事务处理，请不要在外部事务中调用
//
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期和取消操作
//   - beans: 可变参数，要更新的模型对象切片，每个对象必须实现 dialect.Modeler 接口
//
// 返回值:
//   - sql.Result: SQL执行结果，包含受影响的行数等信息
//   - error: 错误信息，如果更新失败则返回具体错误
func (u *update) BatchStruct(ctx context.Context, beans ...dialect.Modeler) (sql.Result, error) {
	defer u.Free()
	if u.err != nil {
		return nil, u.err
	}

	lens := len(beans)
	if lens == 0 {
		return nil, dialect.ErrCreateEmpty
	}

	d := u.db.Dialect()
	u.command.WriteString("UPDATE " + d.Quote(u.table) + " SET ")
	cols, vals := beans[0].RawAssignValues(d, u.cols...)
	for i, col := range cols {
		if i > 0 {
			u.command.WriteString(",")
		}
		u.command.WriteString(col + " = " + d.Placeholder(&u.paramIndex))
	}
	u.params = append(u.params, vals...)
	//
	keys, values := beans[0].AssignKeys()
	u.Where(keys.Eq(values))

	where, params, e := u.parseCond(u.cond)
	if e != nil {
		u.err = e
		return nil, e
	}
	u.whereParams = params
	if where.Len() > 0 {
		u.command.WriteString(" WHERE " + where.String())
	}

	// 只返回SQL语句，不执行
	if u.debug || u.db.Debug() {
		log.Info(u.String())
		return &noRows{}, Err_ToSql
	}

	u.params = append(u.params, u.whereParams...)
	// 启动事务批量执行更新
	ret, err := u.db.Transaction(ctx, func(tx *Tx) (any, error) {
		stmt, err := tx.PrepareContext(ctx, u.command.String())
		if err != nil {
			return nil, err
		}
		if u.db.IsDB() {
			defer stmt.Close()
		}

		result, err := stmt.ExecContext(ctx, u.params...)
		if err != nil {
			return nil, err
		}

		for i := 1; i < lens; i++ {
			bean := beans[i]
			if bean == nil {
				return nil, dialect.ErrBeanEmpty
			}
			_, vals = bean.RawAssignValues(d, u.cols...)
			u.params = u.params[:0]
			u.params = append(u.params, vals...)
			//
			_, values = bean.AssignKeys()
			u.params = append(u.params, values)

			result, err = stmt.ExecContext(ctx, u.params...)
			if err != nil {
				return nil, err
			}
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
