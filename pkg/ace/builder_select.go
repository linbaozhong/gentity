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
	// SetDB 设置数据库连接对象
	//
	// 该方法用于为 Builder 实例设置或更换数据库连接。通常在以下场景使用：
	// 1. Builder 实例化时未传入 DB，后续需要设置
	// 2. 需要在不同的数据库连接之间切换
	//
	// 参数说明:
	//   - d: 数据库连接对象（*DB），必须是有效的已连接实例
	//
	// 返回值说明:
	//   - Builder: 返回构建器实例本身，支持链式调用
	//
	// 使用示例:
	//   // 示例1: 实例化时未传 DB，后续设置
	//   builder := ace.New(nil)
	//   builder.SetDB(db).Table("users").Select().Gets(ctx, &users)
	//
	//   // 示例2: 切换数据库连接
	//   builder.Table("users").SetDB(anotherDB).Select().Gets(ctx, &users)
	//
	// 注意:
	//   - 如果 Builder 已有 DB 连接，SetDB 会覆盖原有连接
	//   - 确保传入的 DB 对象是有效的、已连接的实例
	//   - 该方法主要用于灵活的场景，通常建议在 ace.New() 时直接传入 DB
	SetDB(d *DB) Builder

	// Table 设置查询的表名或数据源
	//
	// 该方法是构建 SQL 查询的起点，用于指定要查询的数据来源。
	// 支持多种类型的参数，可以灵活地指定表名、结构体或子查询。
	//
	// 参数说明:
	//   - a: 数据源参数，支持以下类型:
	//          * string: 直接的数据库表名，如 "users"、"orders"
	//          * dialect.TableNamer: 实现了 TableName() 方法的对象
	//            （通常是自动生成的 DO 对象），会自动调用 TableName() 获取表名
	//          * Builder: 另一个 Builder 实例，表示子查询，会被包装为 (...) 形式
	//          * 其他类型: 通过反射获取类型名称作为表名（较少使用）
	//   - as: 可选的别名参数，最多一个字符串
	//         * 对于普通表：设置表别名，如 Table("users", "u") → "users AS u"
	//         * 对于子查询：必须提供别名，如 Table(subQuery, "sub") → "(SELECT ...) AS sub"
	//
	// 返回值说明:
	//   - Builder: 返回构建器实例，支持链式调用
	//
	// 使用示例:
	//   // 示例1: 使用字符串表名
	//   db.Table("users").Select().Gets(ctx, &users)
	//
	//   // 示例2: 使用带别名的表名
	//   db.Table("users", "u").
	//     Where(tblUsers.ID.Eq(1)).
	//     Select().
	//     Get(ctx, &user)
	//
	// 注意:
	//   - 每个 Builder 只能设置一个主表（最后一次调用会覆盖之前的设置）
	//   - 多表关联应该使用 Join/LeftJoin/RightJoin 方法，而不是多次调用 Table
	//   - 使用子查询时，必须提供别名（as 参数）
	//   - 推荐使用自动生成的 DO 对象（如 tblUsers），具有类型安全性
	//   - 表名和别名会被自动添加引号（根据数据库方言）
	Table(a any, as ...string) Builder

	Columner
	Wherer
	Orderer
	Grouper

	// Join 添加 JOIN 连接查询
	//
	// 该方法用于在查询中添加表连接操作，支持各种 JOIN 类型（INNER JOIN、LEFT JOIN、RIGHT JOIN 等）。
	// 通过指定连接条件和额外的过滤条件，可以灵活地组合多个表的数据。
	//
	// 参数说明:
	//   - joinType: 连接类型，可选值包括:
	//                 * dialect.Inner_Join: 内连接（默认，只返回匹配的行）
	//                 * dialect.Left_Join: 左连接（返回左表所有行，右表无匹配则为 NULL）
	//                 * dialect.Right_Join: 右连接（返回右表所有行，左表无匹配则为 NULL）
	//                 * dialect.Full_Join: 全外连接（返回两表所有行）
	//                 * dialect.Cross_Join: 交叉连接（笛卡尔积）
	//   - left: 连接条件的左侧字段（通常是当前表的字段）
	//           例如: tblUsers.Id、F("u.id")
	//   - right: 连接条件的右侧字段（通常是目标表的字段）
	//            该字段同时决定了 JOIN 的目标表名
	//            例如: tblOrders.UserId、F("o.user_id")
	//   - fns: 可选的额外连接条件，用于添加更复杂的 ON 条件
	//          例如: tblOrders.Status.Eq(1)、tblOrders.CreateTime.Gt(startTime)
	//
	// 返回值说明:
	//   - Builder: 返回构建器实例，支持链式调用
	//
	// 使用示例:
	//   // 示例1: 基本的 INNER JOIN
	//   db.Table("users", "u").
	//     Join(dialect.Inner_Join, tblUsers.Id, tblOrders.UserId).
	//     Select().
	//     Gets(ctx, &results)
	//
	//   // 示例2: LEFT JOIN（推荐写法）
	//   db.Table(tblUsers).
	//     LeftJoin(tblUsers.Id, tblOrders.UserId).
	//     Select().
	//     Gets(ctx, &users)
	//
	//   // 示例3: 带额外条件的 JOIN
	//   db.Table(tblUsers).
	//     Join(dialect.Inner_Join, tblUsers.Id, tblOrders.UserId,
	//          tblOrders.Status.Eq(1),
	//          tblOrders.CreateTime.Gt(startTime)).
	//     Select().
	//     Gets(ctx, &results)
	//
	//   // 示例4: 多表 JOIN
	//   db.Table(tblUsers).
	//     LeftJoin(tblUsers.Id, tblOrders.UserId).
	//     LeftJoin(tblOrders.Id, tblOrderItems.OrderId).
	//     Select().
	//     Gets(ctx, &results)
	//
	// 注意:
	//   - right 字段决定了 JOIN 的目标表，确保该字段包含正确的表信息
	//   - 可以多次调用 Join 来连接多个表
	//   - 额外的条件（fns）会通过 AND 连接到 ON 子句中
	//   - 对于简单的左连接，建议使用 LeftJoin() 快捷方法
	//   - 对于简单的右连接，建议使用 RightJoin() 快捷方法
	//   - JOIN 条件中的字段应该使用正确的表前缀或别名
	Join(joinType dialect.JoinType, left, right dialect.Field, fns ...dialect.Condition) Builder

	// LeftJoin 添加 LEFT JOIN 左连接查询
	//
	// 该方法是 Join 的快捷版本，专门用于添加 LEFT JOIN 操作。
	// LEFT JOIN 会返回左表（主表）的所有记录，即使右表中没有匹配的记录。
	// 对于右表中无匹配的记录，其字段值为 NULL。
	//
	// 参数说明:
	//   - left: 连接条件的左侧字段（通常是左表/主表的字段）
	//           例如: tblUsers.Id、F("u.id")
	//   - right: 连接条件的右侧字段（通常是右表/从表的字段）
	//            该字段同时决定了 JOIN 的目标表名
	//            例如: tblOrders.UserId、F("o.user_id")
	//   - fns: 可选的额外连接条件，用于添加更复杂的 ON 条件
	//          例如: tblOrders.Status.Eq(1)、tblOrders.Amount.Gt(100)
	//
	// 返回值说明:
	//   - Builder: 返回构建器实例，支持链式调用
	//
	// 使用示例:
	//   // 示例1: 基本左连接
	//   db.Table(tblUsers).
	//     LeftJoin(tblUsers.Id, tblOrders.UserId).
	//     Select().
	//     Gets(ctx, &users)
	//   // 生成 SQL: SELECT * FROM users LEFT JOIN orders ON users.id = orders.user_id
	//
	//   // 示例2: 带额外条件的左连接
	//   db.Table(tblUsers).
	//     LeftJoin(tblUsers.Id, tblOrders.UserId,
	//              tblOrders.Status.Eq(1),
	//              tblOrders.CreateTime.Gt(startTime)).
	//     Select().
	//     Gets(ctx, &users)
	//   // 生成 SQL: SELECT * FROM users
	//   //          LEFT JOIN orders ON users.id = orders.user_id
	//   //          AND orders.status = 1 AND orders.create_time > ?
	//
	//   // 示例3: 多个左连接
	//   db.Table(tblUsers).
	//     LeftJoin(tblUsers.Id, tblOrders.UserId).
	//     LeftJoin(tblUsers.DepartmentId, tblDepartments.Id).
	//     LeftJoin(tblUsers.RoleId, tblRoles.Id).
	//     Select().
	//     Gets(ctx, &users)
	//
	//   // 示例4: 左连接 + WHERE 条件
	//   db.Table(tblUsers).
	//     LeftJoin(tblUsers.Id, tblOrders.UserId).
	//     Where(tblUsers.Status.Eq(1)).
	//     Order(tblUsers.CreateTime.Desc()).
	//     Select().
	//     Gets(ctx, &activeUsers)
	//
	//   // 示例5: 左连接 + 聚合函数
	//   db.Table(tblUsers).
	//     LeftJoin(tblUsers.Id, tblOrders.UserId).
	//     Func(dialect.Count(tblOrders.Id).As("order_count")).
	//     Group(tblUsers.Id).
	//     Select().
	//     Gets(ctx, &userStats)
	//
	// 注意:
	//   - LEFT JOIN 是最常用的连接类型之一，特别适合查询"主表 + 可选的关联数据"
	//   - 如果只需要匹配的记录，使用 Join (INNER JOIN) 更高效
	//   - 可以在 WHERE 子句中使用右表字段进行过滤，但要注意 NULL 值的处理
	//   - 多个 LeftJoin 会按调用顺序依次添加到 SQL 中
	//   - 该方法等价于: Join(dialect.Left_Join, left, right, fns...)
	LeftJoin(left, right dialect.Field, fns ...dialect.Condition) Builder

	// RightJoin 添加 RIGHT JOIN 右连接查询
	//
	// 该方法是 Join 的快捷版本，专门用于添加 RIGHT JOIN 操作。
	// RIGHT JOIN 会返回右表的所有记录，即使左表中没有匹配的记录。
	// 对于左表中无匹配的记录，其字段值为 NULL。
	//
	// 在实际应用中，RIGHT JOIN 使用频率低于 LEFT JOIN。
	// 大多数情况下，可以通过交换表的位置使用 LEFT JOIN 替代。
	//
	// 参数说明:
	//   - left: 连接条件的左侧字段（通常是左表的字段）
	//           例如: tblUsers.Id、F("u.id")
	//   - right: 连接条件的右侧字段（通常是右表/主表的字段）
	//            该字段同时决定了 JOIN 的目标表名
	//            例如: tblOrders.UserId、F("o.user_id")
	//   - fns: 可选的额外连接条件，用于添加更复杂的 ON 条件
	//          例如: tblUsers.Status.Eq(1)
	//
	// 返回值说明:
	//   - Builder: 返回构建器实例，支持链式调用
	//
	// 使用示例:
	//   // 示例1: 基本右连接
	//   db.Table(tblUsers).
	//     RightJoin(tblUsers.Id, tblOrders.UserId).
	//     Select().
	//     Gets(ctx, &results)
	//   // 生成 SQL: SELECT * FROM users RIGHT JOIN orders ON users.id = orders.user_id
	//   // 等价于: SELECT * FROM orders LEFT JOIN users ON orders.user_id = users.id
	//
	//   // 示例2: 带额外条件的右连接
	//   db.Table(tblProducts).
	//     RightJoin(tblProducts.Id, tblOrderItems.ProductId,
	//               tblProducts.Status.Eq(1)).
	//     Select().
	//     Gets(ctx, &results)
	//
	//   // 示例3: 使用 LEFT JOIN 替代 RIGHT JOIN（推荐）
	//   // 不推荐：
	//   db.Table(tblUsers).
	//     RightJoin(tblUsers.Id, tblOrders.UserId).
	//     Select()
	//
	//   // 推荐（交换表位置，使用 LEFT JOIN）：
	//   db.Table(tblOrders).
	//     LeftJoin(tblOrders.UserId, tblUsers.Id).
	//     Select()
	//
	// 注意:
	//   - RIGHT JOIN 在实际开发中使用较少，大多数场景可以用 LEFT JOIN 替代
	//   - 如果需要频繁使用 RIGHT JOIN，考虑重新设计查询逻辑，改用 LEFT JOIN
	//   - RIGHT JOIN 的性能与 LEFT JOIN 相当，取决于数据库优化器
	//   - 该方法等价于: Join(dialect.Right_Join, left, right, fns...)
	RightJoin(left, right dialect.Field, fns ...dialect.Condition) Builder

	// Page 设置分页查询参数
	//
	// 该方法用于对查询结果进行分页，基于页码和每页大小计算偏移量。
	// 页码从 1 开始（不是从 0 开始），系统会自动转换为 SQL 的 LIMIT/OFFSET 语法。
	//
	// 参数说明:
	//   - pageIndex: 页码，从 1 开始
	//                  * 如果传入 0 或负数，会自动修正为 1（第一页）
	//                  * 例如: 1 表示第一页，2 表示第二页
	//   - pageSize: 每页记录数
	//                 * 如果小于 1，会返回空结果（LIMIT 0）
	//                 * 建议设置在 10-100 之间，避免单次查询过多数据
	//
	// 返回值说明:
	//   - Builder: 返回构建器实例，支持链式调用
	//
	// 计算逻辑:
	//   OFFSET = (pageIndex - 1) * pageSize
	//   LIMIT = pageSize
	//
	// 使用示例:
	//   // 示例1: 查询第 1 页，每页 10 条
	//   db.Table(tblUsers).
	//     Page(1, 10).
	//     Select().
	//     Gets(ctx, &users)
	//   // 生成 SQL: SELECT * FROM users LIMIT 10 OFFSET 0
	//
	//   // 示例2: 查询第 3 页，每页 20 条
	//   db.Table(tblUsers).
	//     Page(3, 20).
	//     Order(tblUsers.CreateTime.Desc()).
	//     Select().
	//     Gets(ctx, &users)
	//   // 生成 SQL: SELECT * FROM users ORDER BY create_time DESC LIMIT 20 OFFSET 40
	//
	//   // 示例3: 带条件的分页查询
	//   db.Table(tblOrders).
	//     Where(tblOrders.Status.Eq(1)).
	//     Page(pageIndex, pageSize).
	//     Order(tblOrders.CreateTime.Desc()).
	//     Select().
	//     Gets(ctx, &orders)
	//
	//   // 示例4: 分页 + 统计总数
	//   var users []User
	//   db.Table(tblUsers).
	//     Page(1, 10).
	//     Select().
	//     Gets(ctx, &users)
	//
	//   total, _ := db.Table(tblUsers).Count(ctx)
	//   totalPages := (total + 9) / 10  // 计算总页数
	//
	//   // 示例5: 边界情况处理
	//   db.Table(tblUsers).Page(0, 10).Select()   // 自动修正为 Page(1, 10)
	//   db.Table(tblUsers).Page(1, 0).Select()    // 返回空结果（LIMIT 0）
	//   db.Table(tblUsers).Page(-1, 10).Select()  // 自动修正为 Page(1, 10)
	//
	// 注意:
	//   - 页码从 1 开始，不是从 0 开始
	//   - 深分页（大偏移量）性能较差，建议限制最大页码或使用游标分页
	//   - 对于大数据集，建议使用 PageByBookmark 进行游标分页
	//   - Page 方法内部调用 Limit，不要同时使用 Page 和 Limit
	//   - 不同数据库的分页语法可能不同，系统会根据方言自动生成
	Page(pageIndex, pageSize uint) Builder

	// PageByBookmark 基于书签的游标分页查询
	//
	// 该方法实现了高效的游标分页（也称为键集分页或书签分页），
	// 特别适用于大数据集和无限滚动场景。相比传统的 OFFSET 分页，
	// 游标分页在深分页时性能更好，且不会出现数据重复或遗漏的问题。
	//
	// 工作原理:
	//   使用上一页最后一条记录的主键值（或其他唯一字段）作为"书签"，
	//   查询下一页时，只获取书签之后的记录。
	//
	// 参数说明:
	//   - size: 每页记录数（页大小）
	//             * 如果小于 1，会返回空结果
	//             * 建议设置在 10-100 之间
	//   - bm: 书签条件，定义如何筛选下一页数据
	//         * 正序查询（ASC）：使用大于条件（Gt），如 Id.Gt(lastId)
	//         * 倒序查询（DESC）：使用小于条件（Lt），如 Id.Lt(lastId)
	//         * 必须是有效的 Condition 对象
	//
	// 返回值说明:
	//   - Builder: 返回构建器实例，支持链式调用
	//
	// 使用示例:
	//   // 示例1: 正序游标分页（ID 递增）
	//   // 第一页
	//   var users []User
	//   db.Table(tblUsers).
	//     Order(tblUsers.Id.Asc()).
	//     PageByBookmark(10, nil).  // 第一页不需要书签
	//     Select().
	//     Gets(ctx, &users)
	//
	//   // 获取最后一条记录的 ID
	//   lastId := users[len(users)-1].Id
	//
	//   // 第二页及后续页
	//   db.Table(tblUsers).
	//     Order(tblUsers.Id.Asc()).
	//     PageByBookmark(10, tblUsers.Id.Gt(lastId)).
	//     Select().
	//     Gets(ctx, &users)
	//
	//   // 示例2: 倒序游标分页（ID 递减，最新在前）
	//   // 第一页
	//   db.Table(tblArticles).
	//     Order(tblArticles.Id.Desc()).
	//     PageByBookmark(20, nil).
	//     Select().
	//     Gets(ctx, &articles)
	//
	//   // 后续页
	//   lastId := articles[len(articles)-1].Id
	//   db.Table(tblArticles).
	//     Order(tblArticles.Id.Desc()).
	//     PageByBookmark(20, tblArticles.Id.Lt(lastId)).
	//     Select().
	//     Gets(ctx, &articles)
	//
	//   // 示例3: 基于时间字段的游标分页
	//   // 按创建时间倒序，获取最新动态
	//   db.Table(tblPosts).
	//     Where(tblPosts.Status.Eq(1)).
	//     Order(tblPosts.CreateTime.Desc()).
	//     PageByBookmark(15, tblPosts.CreateTime.Lt(lastCreateTime)).
	//     Select().
	//     Gets(ctx, &posts)
	//
	//   // 示例4: 复合条件游标（时间 + ID）
	//   // 当时间相同时，用 ID 保证顺序一致性
	//   db.Table(tblLogs).
	//     Order(tblLogs.CreateTime.Desc(), tblLogs.Id.Desc()).
	//     PageByBookmark(50,
	//       tblLogs.CreateTime.Lt(lastTime).
	//         Or(tblLogs.CreateTime.Eq(lastTime).And(tblLogs.Id.Lt(lastId)))).
	//     Select().
	//     Gets(ctx, &logs)
	//
	//   // 示例5: 前端传递书签的实现
	//   type PageRequest struct {
	//     Size   uint   `json:"size"`
	//     Cursor string `json:"cursor"`  // 游标（base64 编码的最后一条记录 ID）
	//   }
	//
	//   func getList(req PageRequest) ([]User, error) {
	//     builder := db.Table(tblUsers).
	//       Order(tblUsers.Id.Asc()).
	//       PageByBookmark(req.Size, nil)
	//
	//     if req.Cursor != "" {
	//       lastId, _ := decodeCursor(req.Cursor)
	//       builder = db.Table(tblUsers).
	//         Order(tblUsers.Id.Asc()).
	//         PageByBookmark(req.Size, tblUsers.Id.Gt(lastId))
	//     }
	//
	//     var users []User
	//     err := builder.Select().Gets(ctx, &users)
	//     return users, err
	//   }
	//
	// 优势对比（vs 传统 OFFSET 分页）:
	//   1. 性能更好：不使用 OFFSET，避免了扫描和跳过大量记录
	//   2. 数据一致：插入或删除记录不会导致数据重复或遗漏
	//   3. 适合深分页：无论翻到第几页，性能都稳定
	//   4. 适合无限滚动：完美支持下拉加载更多的场景
	//
	// 劣势:
	//   1. 不能跳转到指定页码（只能一页一页翻）
	//   2. 需要客户端保存游标状态
	//   3. 实现稍复杂
	//
	// 注意:
	//   - 第一页查询时，bm 参数可以传 nil 或不设置条件
	//   - 后续页必须提供正确的书签条件
	//   - 排序方向必须与书签条件匹配：
	//     * 正序（ASC）用 Gt（大于）
	//     * 倒序（DESC）用 Lt（小于）
	//   - 用于游标的字段应该有索引，以保证查询性能
	//   - 建议使用主键或唯一字段作为游标
	//   - 如果排序字段可能重复，需要添加次要排序字段（如 ID）保证唯一性
	//   - 该方法内部调用 Where 和 Limit，不要重复设置
	PageByBookmark(size uint, bm dialect.Condition) Builder

	// Limit 设置查询结果的数量限制
	//
	// 该方法用于限制查询返回的记录数量，可以直接指定限制条数和起始位置。
	// 不同数据库的 LIMIT 语法可能不同，系统会根据数据库方言自动生成正确的 SQL。
	//
	// 参数说明:
	//   - size: 限制返回的记录数量
	//             * 如果为 0，会清除 LIMIT 子句（返回所有记录）
	//             * 建议设置合理的上限，避免一次性加载过多数据
	//   - start: 可选的起始位置（偏移量），从 0 开始
	//              * 如果不传，默认为 0（从第一条记录开始）
	//              * 例如: Limit(10, 20) 表示从第 21 条开始，取 10 条
	//              * 相当于 SQL: LIMIT 10 OFFSET 20
	//
	// 返回值说明:
	//   - Builder: 返回构建器实例，支持链式调用
	//
	// 使用示例:
	//   // 示例1: 限制返回 10 条记录
	//   db.Table(tblUsers).
	//     Limit(10).
	//     Select().
	//     Gets(ctx, &users)
	//   // 生成 SQL: SELECT * FROM users LIMIT 10
	//
	//   // 示例2: 从第 21 条开始，取 10 条
	//   db.Table(tblUsers).
	//     Limit(10, 20).
	//     Select().
	//     Gets(ctx, &users)
	//   // 生成 SQL: SELECT * FROM users LIMIT 10 OFFSET 20
	//
	//   // 示例3: 结合排序使用
	//   db.Table(tblArticles).
	//     Order(tblArticles.ViewCount.Desc()).
	//     Limit(5).
	//     Select().
	//     Gets(ctx, &topArticles)
	//   // 获取阅读量最高的 5 篇文章
	//
	//   // 示例4: 获取最新的一条记录
	//   db.Table(tblOrders).
	//     Where(tblOrders.UserId.Eq(userId)).
	//     Order(tblOrders.CreateTime.Desc()).
	//     Limit(1).
	//     Select().
	//     Get(ctx, &latestOrder)
	//
	//   // 示例5: 清除 LIMIT 限制
	//   db.Table(tblUsers).
	//     Limit(10).      // 先设置限制
	//     Limit(0).       // 再清除限制
	//     Select().
	//     Gets(ctx, &allUsers)
	//   // 最终不会添加 LIMIT 子句
	//
	//   // 示例6: 不同数据库的语法差异
	//   // MySQL:    LIMIT 10 OFFSET 20
	//   // SQLite:   LIMIT 10 OFFSET 20
	//   // PostgreSQL: LIMIT 10 OFFSET 20
	//   // SQL Server: OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY
	//
	// 注意:
	//   - start 参数从 0 开始，不是从 1 开始
	//   - Limit(0) 会清除之前设置的 LIMIT 子句
	//   - 不要同时使用 Page 和 Limit，它们会相互覆盖
	//   - 对于大数据集，建议始终设置合理的 LIMIT，避免内存溢出
	//   - 如果需要分页功能，建议使用 Page 方法（更直观）
	//   - Limit 是底层方法，Page 内部调用了 Limit
	Limit(size uint, start ...uint) Builder

	// Select 创建查询执行器
	//
	// 该方法是查询构建的最后一步，将当前的 Builder 转换为 Selecter 执行器，
	// 用于执行 SELECT 查询并获取结果。Selecter 提供了多种结果获取方法，
	// 可以灵活地将查询结果映射到不同的数据结构中。
	//
	// 参数说明:
	//   - db: 可选的数据库连接对象
	//         * 如果 Builder 实例化时已经传入了 DB，则此处可以省略
	//         * 如果实例化时未传入 DB（如 ace.New(nil)），则此处必须传入
	//         * 如果传入新的 DB，会覆盖原有的数据库连接
	//
	// 返回值说明:
	//   - Selecter: 返回查询执行器接口，提供以下方法:
	//                 * Query(ctx): 返回 *sql.Rows，手动处理结果集
	//                 * QueryRow(ctx): 返回 *sql.Row，单行查询
	//                 * Get(ctx, dest): 获取单条记录到结构体
	//                 * Gets(ctx, dest): 获取多条记录到切片
	//                 * Map(ctx): 获取单条记录为 map[string]any
	//                 * Maps(ctx): 获取多条记录为 []map[string]any
	//                 * Slice(ctx): 获取单条记录为 []any
	//                 * Slices(ctx): 获取多条记录为 [][]any
	//                 * Count(ctx): 获取记录总数
	//                 * Sum/Avg/Max/Min: 聚合函数
	//                 * RawQuery*: 原生 SQL 查询方法
	//
	// 使用示例:
	//   // 示例1: 基本查询（DB 已在实例化时传入）
	//   var users []User
	//   db.Table(tblUsers).
	//     Where(tblUsers.Status.Eq(1)).
	//     Select().  // 创建查询执行器
	//     Gets(ctx, &users)
	//
	//   // 示例2: 查询单条记录
	//   var user User
	//   db.Table(tblUsers).
	//     Where(tblUsers.Id.Eq(1)).
	//     Select().
	//     Get(ctx, &user)
	//
	//   // 示例3: 如果实例化时未传入 DB，则此处必须传入
	//   builder := ace.New(nil)
	//   var users []User
	//   builder.Table(tblUsers).
	//     Select(db).  // 这里必须传入 db
	//     Gets(ctx, &users)
	//
	//   // 示例4: 查询结果为 map
	//   result, err := db.Table(tblUsers).
	//     Where(tblUsers.Id.Eq(1)).
	//     Select().
	//     Map(ctx)
	//   // result 是 map[string]any 类型
	//
	//   // 示例5: 统计数量
	//   count, err := db.Table(tblUsers).
	//     Where(tblUsers.Status.Eq(1)).
	//     Select().
	//     Count(ctx)
	//
	//   // 示例6: 聚合查询
	//   stats, err := db.Table(tblOrders).
	//     Where(tblOrders.UserId.Eq(userId)).
	//     Select().
	//     Sum(ctx, []dialect.Field{tblOrders.Amount})
	//   // stats["amount"] 包含总金额
	//
	//   // 示例7: 自定义列查询
	//   db.Table(tblUsers).
	//     Cols(tblUsers.Id, tblUsers.Name, tblUsers.Email).
	//     Where(tblUsers.Status.Eq(1)).
	//     Select().
	//     Gets(ctx, &users)
	//
	//   // 示例8: 复杂查询
	//   db.Table(tblUsers).
	//     Cols(tblUsers.Id, tblUsers.Name).
	//     Func(dialect.Count(tblOrders.Id).As("order_count")).
	//     LeftJoin(tblUsers.Id, tblOrders.UserId).
	//     Where(tblUsers.Status.Eq(1)).
	//     Group(tblUsers.Id, tblUsers.Name).
	//     Having(dialect.Count(tblOrders.Id).Gt(5)).
	//     Order(dialect.F("order_count").Desc()).
	//     Page(1, 20).
	//     Select().
	//     Gets(ctx, &userStats)
	//
	// 注意:
	//   - Select() 是查询构建的最后一步，之后不能再调用 Table、Where 等方法
	//   - 调用 Select() 后，必须调用 Selecter 的执行方法（Get、Gets 等）才能真正执行查询
	//   - 如果 Builder 没有 DB 连接，必须在 Select() 时传入，否则执行时会出错
	//   - Select() 会锁定当前的查询配置，后续的修改不会影响本次查询
	//   - 每次调用 Select() 都会创建一个新的 Selecter 实例
	//   - Selecter 执行完成后会自动释放 Builder 资源
	Select(...*DB) Selecter
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
	// Avg 返回平均值
	Avg(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error)
	// Max 返回最大值
	Max(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error)
	// Min 返回最小值
	Min(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error)
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

// Select 创建查询器
func (o *orm) Select(x ...*DB) Selecter {
	if len(x) > 0 {
		o.db = x[0]
	}
	return &read{
		orm: o,
	}
}

// // Sub 子查询
// func (o *orm) Sub(b Builder, as ...string) Builder {
// 	cmd, params := b.parse()
// 	o.table = "(" + cmd.String() + ")"
// 	o.whereParams = append(o.whereParams, params...)
// 	if len(as) > 0 {
// 		o.table = fmt.Sprintf("%s AS %s", o.table, as[0])
// 	}
// 	return o
// }

// Query
func (s *read) Query(ctx context.Context) (*sql.Rows, error) {
	defer s.Free()
	if s.err != nil {
		return nil, s.err
	}

	return s.query(ctx)
}

// QueryRow
func (s *read) QueryRow(ctx context.Context) (*sql.Row, error) {
	defer s.Free()
	if s.err != nil {
		return nil, s.err
	}

	cmd, params, e := s.parse()
	if e != nil {
		return nil, e
	}

	return s.row(ctx, cmd.String(), params...)
}

// Get 返回单个数据，dest 必须是指针
func (s *read) Get(ctx context.Context, dest any) error {
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
func (s *read) Map(ctx context.Context) (map[string]any, error) {
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
	err = r.MapScan(dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

// Maps 返回 map[string]any 的切片 []map[string]any，用于列数未知的情况
func (s *read) Maps(ctx context.Context) ([]map[string]any, error) {
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
func (s *read) Slice(ctx context.Context) ([]any, error) {
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
func (s *read) Slices(ctx context.Context) ([][]any, error) {
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
func (s *read) Count(ctx context.Context, cond ...dialect.Condition) (int64, error) {
	defer s.Free()
	if s.err != nil {
		return 0, s.err
	}
	//
	// 保存原始 limit 和 command，Count 不应受分页限制
	savedLimit := s.limit
	s.limit = "" // Count 忽略 LIMIT/PAGE

	s.command.Reset()

	s.Where(cond...)

	s.command.WriteString("SELECT COUNT(*)")

	if e := s.buildFromClause(cond); e != nil {
		return 0, e
	}

	row, err := s.row(ctx, s.command.String(), s.mergeParams()...)
	if err != nil {
		return 0, err
	}

	// 恢复原始 limit（虽然 Free 后无意义，保持语义完整）
	s.limit = savedLimit

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
	if s.err != nil {
		return nil, s.err
	}

	for _, col := range cols {
		s.Func(col.Sum())
	}
	return s.aggregateQuery(ctx, cols, cond...)
}

// Avg 返回平均值
func (s *read) Avg(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error) {
	defer s.Free()
	if s.err != nil {
		return nil, s.err
	}

	for _, col := range cols {
		s.Func(col.Avg())
	}
	return s.aggregateQuery(ctx, cols, cond...)
}

// Max 返回最大值
func (s *read) Max(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error) {
	defer s.Free()
	if s.err != nil {
		return nil, s.err
	}

	for _, col := range cols {
		s.Func(col.Max())
	}
	return s.aggregateQuery(ctx, cols, cond...)
}

// Min 返回最小值
func (s *read) Min(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error) {
	defer s.Free()
	if s.err != nil {
		return nil, s.err
	}

	for _, col := range cols {
		s.Func(col.Min())
	}
	return s.aggregateQuery(ctx, cols, cond...)
}

// aggregateQuery 聚合查询的公共部分：构建SQL、执行查询、扫描结果
// 调用前需要先设置 s.funcs 和 s.where
func (s *read) aggregateQuery(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error) {
	//
	s.command.Reset()

	s.Where(cond...)
	s.command.WriteString("SELECT ")
	s.command.WriteString(strings.Join(s.parseFunc(s.funcs), ","))

	if e := s.buildFromClause(cond); e != nil {
		return nil, e
	}

	// LIMIT
	if s.limit != "" {
		s.command.WriteString(s.limit)
	}

	row, err := s.row(ctx, s.command.String(), s.mergeParams()...)
	if err != nil {
		return nil, err
	}

	var results = make([]any, len(cols))
	err = row.Scan(results...)
	if err != nil {
		return nil, err
	}

	resultMap := make(map[string]any, len(cols))
	for i := range results {
		resultMap[cols[i].Name] = results[i]
	}
	return resultMap, nil
}

// buildFromClause 构建公共的 FROM + JOIN + WHERE 子句
// 用于 Count / Sum / Avg / Max / Min 等聚合查询
func (s *read) buildFromClause(cond []dialect.Condition) error {
	d := s.db.Dialect()
	// FROM TABLE
	if s.table == "" {
		return Err_TableName
	}
	s.command.WriteString(" FROM " + d.Quote(s.table))

	if len(s.join) > 0 {
		joinStr, params, e := s.parseJoin(s.join)
		if e != nil {
			return e
		}
		s.joinParams = params
		if joinStr.Len() > 0 {
			s.command.WriteString(joinStr.String())
		}
	}

	if len(s.cond) > 0 {
		where, params, e := s.parseCond(s.cond)
		if e != nil {
			return e
		}
		s.whereParams = params
		if where.Len() > 0 {
			s.command.WriteString(" WHERE " + where.String())
		}
	}
	return nil
}

// Select 执行原生的 SQL 查询
// 此方法接受一个上下文、原生 SQL 语句和对应的参数，返回查询结果和可能的错误
func (s *read) RawQuery(ctx context.Context, sqlStr string, args ...any) (*sql.Rows, error) {
	defer s.Free()
	if s.err != nil {
		return nil, s.err
	}

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
	err = r.MapScan(dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
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
