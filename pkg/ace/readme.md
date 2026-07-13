# ace — gentity 的 ORM 内核

`ace` 是 gentity 框架的数据库访问内核，提供一套**链式、类型安全**的 SQL 构建器（Builder）与执行器（Executer）。
它不依赖反射做字段映射，生成的 DAO 代码与结构体均已**池化**，对 GC 友好。

- 支持 MySQL / PostgreSQL / SQLite / SQL Server（按驱动名自动选择方言）
- `*DB` 与 `*Tx` 都实现 `Executer` 接口，因此**同一个 Builder 既能在普通连接上执行，也能在事务内执行**
- 条件、赋值、聚合函数均以 `dialect.Field` / `dialect.Function` 的形式由代码生成器产出，类型安全

> **命名来源：**
> 
> ace 取自扑克牌中的 "A"（王牌），寓意它是 gentity 框架数据访问层的"王牌引擎"；其缩写可解读为 Access Connection for Entity —— 实体访问连接层。在框架分层中，ace 负责底层 SQL 的构建与执行，与上层应用包 ack（"收到请求、给出响应"）形成 ace → ack 的数据流呼应。
---

## 1. 连接数据库

    import "github.com/linbaozhong/gentity/pkg/ace"

    db, err := ace.Connect("mysql", "user:pass@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

`Connect` 内部会 `sql.Open` + `Ping`，并按驱动名自动装配方言（`mysql`/`postgres`/`sqlite`/`sqlserver`）。
连接对象通过 `app.RegisterServiceCloser` 注册到应用生命周期，进程退出时自动关闭。

---

## 2. 核心类型

| 类型 | 说明 |
|------|------|
| `DB` | 普通数据库连接，内嵌 `*sql.DB`，实现 `Executer` |
| `Tx` | 事务连接，内嵌 `*sql.Tx`，同样实现 `Executer` |
| `Executer` | 执行器接口（`PrepareContext`/`QueryContext`/`ExecContext`/`Table`/`Where`/`Dialect`/`IsDB` 等），`*DB` 与 `*Tx` 都满足 |
| `Builder` | 链式 SQL 构建器，聚合了 `SelectBuilder`/`CreateBuilder`/`UpdateBuilder`/`DeleteBuilder` 四个子接口 |

> 因为 `*Tx` 也实现了 `Executer`，所以只要把 `*Tx` 当作 `Executer` 传入 DAO / Builder，
> 后续所有 `SELECT/INSERT/UPDATE/DELETE` 都会**在事务内**执行（详见第 6 节）。

---

## 3. 查询（SELECT）

生成的列字段（如 `data.Id`、`data.LongName`）自带 `Eq/In/Like/...` 等方法，返回 `dialect.Condition`：

    var companies []Company
    has, err := db.Table("company").
        Where(data.Status.Eq(1), data.LongName.Like("%科技%")).
        Order(data.Ctime.Desc()).
        Page(1, 20).
        Select().
        Gets(ctx, &companies)

常用查询入口（均在 `Select()` 之后调用）：

| 方法 | 说明 |
|------|------|
| `Get(ctx, &dest)` | 取首行到结构体 |
| `Gets(ctx, &slice)` | 取多行到切片 |
| `Map(ctx)` / `Maps(ctx)` | 取为 `map[string]any` / `[]map[string]any` |
| `Count(ctx)` | 总数 |
| `Sum/Avg/Max/Min(ctx, fields)` | 聚合函数 |
| `Exist(ctx)` | 是否存在 |

支持 `Join`/`LeftJoin`/`RightJoin`、`Group` + `Having`、`Func`（聚合函数）、`Distinct`、`Cols`/`Omit`（指定/排除列）、`Page` / `PageByBookmark`（游标分页）。

---

## 4. 插入（INSERT）

    // 方式一：显式指定列与值
    res, err := db.Table("company").
        Set(
            data.LongName.Set("示例科技"),
            data.Status.Set(1),
            data.Ctime.Set(time.Now()),
        ).
        Create().
        Exec(ctx)

    // 方式二：通过结构体（自动按字段赋值）
    err := db.Table("company").Create().Struct(ctx, &company)

    // 方式三：批量插入（注意：批量插入不建议在事务中使用）
    err := db.Table("company").Create().BatchStruct(ctx, &c1, &c2)

字段赋值除了 `Set(v)` 普通赋值，还支持表达式：

| 方法 | 生成的 SQL 片段 |
|------|----------------|
| `f.Incr(n)` | `col = col + n`（默认 1） |
| `f.Decr(n)` | `col = col - n`（默认 1） |
| `f.Replace(old, new)` | `col = REPLACE(col, old, new)` |
| `f.Expr("col + 1")` | `col = col + 1` |

---

## 5. 更新与删除（UPDATE / DELETE）

    // 更新：视图计数 +1，并修改状态
    _, err := db.Table("company").
        Where(data.Id.Eq(1001)).
        Set(data.Status.Set(2), data.ViewCount.Incr(1)).
        Update().
        Exec(ctx)

    // 删除
    res, err := db.Table("company").
        Where(data.Id.Eq(1001)).
        Delete().
        Exec(ctx)

---

## 6. 事务（Transaction）

`*DB.Transaction` 自动 `Begin` → 执行业务 → `Commit` / `Rollback`：

    result, err := db.Transaction(ctx, func(tx *ace.Tx) (any, error) {
        if _, e := tx.Table("company").
            Set(data.Status.Set(2)).
            Where(data.Id.Eq(1001)).
            Update().Exec(ctx); e != nil {
            return nil, e
        }
        if _, e := tx.Table("company").
            Set(data.LongName.Set("new")).
            Where(data.Id.Eq(1002)).
            Update().Exec(ctx); e != nil {
            return nil, e // 返回错误会自动 Rollback
        }
        return nil, nil
    })

事务语义保证：`*Tx` 创建出的 Builder 内部执行器就是 `*Tx` 自身，
所有 SQL 经由 `*sql.Tx` 执行，不会绕回连接池。

---

## 7. 调试模式

    // 全局：仅打印生成的 SQL，不执行
    db.Debug(true)

    // 单次 Builder：仅打印本条 SQL，不执行
    db.Table("company").Where(data.Id.Eq(1)).Debug().Select().Gets(ctx, &list)

调试模式下 `Exec`/`Struct`/`BatchStruct` 等方法会返回 `ace.Err_ToSql` 错误，并打印最终 SQL 与参数。

---

## 8. 子包（subpackages）

| 子包 | 职责 |
|------|------|
| `dialect` | 核心：`Field`/`Condition`/`Function`/`Setter` 构建器 + `Dialect` 方言接口；含 `mysql`、`postgres`、`sqlite`、`sqlserver` 四套实现 |
| `data` | 由代码生成器产出的**列字段常量**（如 `data.Id`、`data.WritableFields`），提供给 Builder 使用 |
| `dao` | `DataAccessInterface` 元数据结构体，描述一个 DAO 方法的命名空间、表名、输入输出等，供代码生成与文档使用 |
| `pool` | 基于 `sync.Pool` 的对象池（`pool.New[T]`）与 `Model` 基类，提供 `Reset/Put/Get` 的防重复入池机制 |
| `reflectx` | 结构体字段映射（`reflectx.Mapper`），用于读取 struct tag 与字段值 |

---

## 9. 设计要点

- **对象池化**：`orm`（构建器内部状态）与生成的 DAO 对象均从 `sync.Pool` 复用，降低 GC 压力。
- **无反射热路径**：SQL 拼接与参数绑定走代码生成的 `Field`/`Setter`，执行期不依赖反射。
- **统一执行器**：`Executer` 接口让 `*DB` 与 `*Tx` 无缝互换，事务接入对业务代码零侵入。
- **方言隔离**：不同数据库的引号、占位符（`?` / `$1` / `@p1`）、分页语法由 `Dialect` 实现各自处理。

---

## 10. 常见错误

| 错误 | 含义 |
|------|------|
| `ace.Err_TableName` | 未指定表名（`Table` 为空） |
| `ace.Err_ToSql` | 处于 Debug 模式，SQL 仅打印未执行 |
| `dialect.ErrCreateEmpty` | `Create` 未提供任何列 |
| `dialect.ErrBeanEmpty` | `Struct`/`BatchStruct` 传入空对象 |
