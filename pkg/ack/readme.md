# ack — gentity 的 Web / API 应用层

`ack` 是 gentity 框架的 HTTP 应用层，封装了**请求读取、参数校验、业务调用、响应写入、路由注册、缓存、幂等**等通用能力。
它采用「泛型 API + 框架适配器」的设计：

- 提供 `iris` 与 `gin` 两套适配器，二者**导出的包名都是 `ack`，API 完全一致**；
- 切换 Web 框架时**只需改动 import 路径，业务代码一行不动**；
- 框架无关的核心逻辑集中在 `internal/core`，通过 `core.Context` 接口抽象 `iris.Context` / `gin.Context`。

> **命名来源：**
> 
> Ack = API Construction Kit 
> * API：专注于构建 HTTP API 服务
> * Construction：像搭积木一样构建路由、中间件、处理器
> * Kit：工具箱，轻量、实用、即拿即用
> 
> 同时致敬 TCP 协议中的 ACK（确认包） —— 每一个请求，Ack 都给你一个确定的回应。

---

## 1. 包结构

| 路径 | 说明 |
|------|------|
| `pkg/ack/iris` | 基于 iris 的适配器，包名 `ack` |
| `pkg/ack/gin` | 基于 gin 的适配器，包名 `ack` |
| `pkg/ack/internal/core` | 框架无关核心：Context 接口、请求处理、响应、缓存、幂等 |

> 注意：`pkg/ack` 下若出现了 `iris - 副本`、`gin - 副本` 目录，通常为误提交的重复副本，可清理。

---

## 2. 快速开始

以 iris 为例（gin 仅需替换 import 路径）：

    import (
        ack "github.com/linbaozhong/gentity/pkg/ack/iris"
        _ "your_project/internal/handler" // 触发路由注册
    )

    func main() {
        app := ack.NewApplication("myapp", "0.1")
        v1 := ack.NewParty(app, "/v1")
        ack.RegisterRouter(v1) // 注册所有 handler 中声明的路由

        server := ack.NewServer(app, ":8080")
        if err := server.Run(); err != nil {
            log.Fatal(err)
        }
    }

`Server` 提供 `Run()`（阻塞启动）与 `Shutdown(ctx)`（优雅关闭），iris / gin 接口一致。

---

## 3. 核心概念

| 类型 / 符号 | 说明 |
|------|------|
| `ack.Context` | 类型别名，等于底层框架的 Context（`iris.Context` / `gin.Context`） |
| `ack.Application` / `ack.Party` / `ack.Handler` | 分别对应框架的 app / router group / handler 类型 |
| `core.Context` | `internal/core` 中定义的抽象接口，屏蔽框架差异 |
| `ack.Get` / `ack.Post` | 泛型请求处理入口，自动完成「读参 → 校验 → 调服务 → 写响应」 |

请求处理链路（以 `ack.Post` 为例）：

1. `adapt(ctx)` 把框架 Context 适配为 `core.Context`；
2. `readPostRequest` 按 `Content-Type` 读取参数（json / form / query），并调用 `Initiate` 设置 IP、UserAgent、RequestId；
3. 若请求结构体实现了 `Checker` 接口，自动调用 `Check()` 校验；
4. 调用 `callService` 执行业务（带 3 秒超时）；
5. 成功则 `Ok` 写 JSON，失败则 `Fail` 写错误 JSON。

---

## 4. 定义路由

实现 `IRegisterRoute` 接口，在 `init()` 中注册即可被 `RegisterRouter` 统一挂载：

    type user struct{}

    func init() {
        ack.RegisterRoute(&user{})
    }

    func (u *user) RegisterRoute(group ack.Party) {
        g := ack.NewParty(group, "/user")
        g.Post("/register", u.register) // 无需登录
        g.Use(lib.AuthMiddleware())
        g.Get("/get", u.get)            // 需要登录
    }

| 符号 | 说明 |
|------|------|
| `IRegisterRoute` | 路由注册接口，含 `RegisterRoute(Party)` |
| `RegisterRoute(r)` | 把一个路由处理器加入全局注册表 |
| `RegisterRouter(group)` | 遍历注册表，把全部路由挂到指定 Party 上 |
| `NewParty(app, path)` | 创建路由分组（等价于 `app.Party(path)`） |

---

## 5. 处理请求（泛型 API）

`ack.Get` / `ack.Post` 是最常用的两个入口，业务只需实现一个 `callService` 函数：

    // 业务函数：A 为请求结构体，B 为响应结构体
    func GetUser(ctx context.Context, req *GetUserReq, resp *GetUserResp) error {
        // ... 调用 ace 的 DAO 等
        return nil
    }

    func (u *user) get(c ack.Context) {
        ack.Get(c, GetUser)
    }

    func (u *user) register(c ack.Context) {
        ack.Post(c, UserRegister)
    }

`callService` 签名统一为 `func(ctx context.Context, req *A, resp *B) error`。
`Get` 读取 URL query，`Post` 按 `Content-Type` 读取 JSON / 表单。

可选 `after` 回调，用于在写响应前对 `resp` 做二次加工：

    ack.Get(c, GetUser, func(c ack.Context, resp *GetUserResp) error {
        resp.Secret = "" // 脱敏
        return nil
    })

其他入口：

| 函数 | 说明 |
|------|------|
| `ack.Get[A,B]` / `ack.Post[A,B]` | 标准 GET / POST，自动写 JSON 响应 |
| `ack.Redirect[A]` | POST 后做 HTTP 重定向（resp 为 `*string` 目标地址） |
| `ack.Stream[A,B]` | 流式响应，无超时限制 |
| `ack.GetResult[A,B]` / `ack.PostResult[A,B]` | 同上，但返回 `(*B, error)` 由调用方自行处理响应 |

---

## 6. 响应与错误

    // 业务出错时，在 handler / middleware 中手动写错误响应
    func AuthMiddleware() ack.Handler {
        return func(c ack.Context) {
            if c.GetHeader("Authorization") == "" {
                ack.Fail(c, constant.ErrAuthorizationNotFound)
                return
            }
            c.Next()
        }
    }

| 函数 | 说明 |
|------|------|
| `ack.Fail(c, err, args...)` | 写错误 JSON（自动识别 `types.Error` 的 code/message） |
| `ack.Ok(c, data...)` | 写成功 JSON；GET 请求且启用缓存时自动写入响应缓存 |
| `ack.SendLocalFile(c, path, name)` | 以附件形式发送本地文件 |
| `ack.SendUrlFile(c, url, name)` | 下载远端 URL 并以附件形式发送 |

---

## 7. GET 响应缓存

在 GET handler 开头调用 `ack.ReadCache`，命中则直接返回缓存并中止后续处理：

    func (u *user) get(c ack.Context) {
        if ack.ReadCache(c, time.Second*30) {
            return // 命中缓存，响应已写出
        }
        ack.Get(c, GetUser)
    }

缓存 key 为 `路径 + 查询参数`，默认有效期 30 秒，可传参覆盖。

---

## 8. 幂等请求

对写接口做幂等保护（相同请求体在有效期内只执行一次）：

    var idempotencyConfig = ack.DefaultIdempotencyConfig(mmap.New())

    func (u *user) register(c ack.Context) {
        ack.PostIdempotent(c, idempotencyConfig, UserRegister)
    }

| 符号 | 说明 |
|------|------|
| `ack.DefaultIdempotencyConfig(cache)` | 返回默认配置（默认 24h，键 = Hash(请求体)） |
| `ack.PostIdempotent[A,B](c, config, svc, after...)` | 带幂等保护的 POST |

---

## 9. 中间件与静态资源

    app.Use(ack.Logger(), ack.Recovery()) // 日志 + 异常恢复（NewApplication 已默认挂载）

    ack.StaticWeb(party, "/static", "./public")           // 静态文件目录
    engine := ack.HtmlView("./views", ".html", true)      // HTML 模板引擎

| 符号 | 说明 |
|------|------|
| `ack.Logger()` | 请求日志中间件 |
| `ack.Recovery()` | panic 恢复中间件 |
| `ack.StaticWeb(party, url, dir)` | 注册静态文件服务 |
| `ack.HtmlView(dir, ext, reload)` | 创建 HTML 模板引擎 |

---

## 10. 切换 Web 框架

因为 `iris` 与 `gin` 两个适配器导出**完全相同的 `ack` 包 API，业务代码无需改动，只改 import：

    // 用 iris
    import ack "github.com/linbaozhong/gentity/pkg/ack/iris"

    // 改用 gin（其余代码一字不改）
    import ack "github.com/linbaozhong/gentity/pkg/ack/gin"

---

## 11. 设计要点

- **泛型 + 适配器**：请求/响应结构体由泛型参数 `A`/`B` 约束，业务无需手写 JSON 解析；框架差异被 `core.Context` 接口隔离。
- **自动 Initiate / Validate**：`Initiate` 注入 IP、UserAgent、RequestId；若 `req` 实现 `Checker` 接口则自动 `Check()`，保证参数校验不遗漏。
- **统一错误模型**：`Fail` 识别 `types.Error`（含业务 code），其余错误归为未知错误，响应格式一致。
- **超时保护**：`Get` / `Post` 默认带 3 秒上下文超时；`Stream` 不限时，适合长连接/流式场景。
- **GET 缓存**：`ReadCache` + `Ok` 配合，对只读接口做响应级缓存，降低后端压力。
