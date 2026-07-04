# 错误处理最佳实践

## 概述

为了避免在 DAO 层记录日志（这是自动生成的代码），我们采用以下分层错误处理策略：

1. **DAO 层**：直接返回原始错误，不记录日志
2. **Service 层**：包装错误，添加操作上下文信息
3. **Handler 层**：记录错误日志，返回格式化的错误响应

## 使用方法

### 1. Service 层错误包装

在 Service 层使用 `apperr.WrapDB()` 或 `apperr.New()` 包装数据库错误：

```go
import apperr "Mingtianjian/internal/pkg/errors"

// 插入数据库
ok, e := dao.Merchants().InsertOne(c, _merchant)
if e != nil {
    // 包装错误，添加操作信息
    return apperr.WrapDB(e, "CreateMerchant.InsertOne")
}
if !ok {
    // 创建自定义错误
    return apperr.New(apperr.ErrorTypeUnknown, "CreateMerchant.InsertOne", nil, "插入商家失败")
}
```

### 2. 错误类型

```go
const (
    ErrorTypeUnknown      // 未知错误
    ErrorTypeDB          // 数据库错误
    ErrorTypeValidation  // 验证错误
    ErrorTypePermission  // 权限错误
    ErrorTypeNotFound    // 未找到错误
)
```

### 3. 可用的包装函数

```go
// 包装数据库错误
apperr.WrapDB(err error, op string) *AppError

// 创建验证错误
apperr.NewValidation(op string, err error, message string) *AppError

// 创建权限错误
apperr.NewPermission(op string, err error, message string) *AppError

// 创建未找到错误
apperr.NewNotFound(op string, err error, message string) *AppError

// 通用包装
apperr.Wrap(err error, op string) *AppError
```

### 4. Handler 层错误处理

在 Handler 中使用中间件统一处理错误：

```go
import "Mingtianjian/internal/handler"

func Init() api.Application {
    _app := api.NewApplication("Mingtianjian", "0.1")

    // 添加错误处理中间件
    _app.Use(handler.ErrorMiddleware())

    // 注册自定义错误处理器
    _app.SetErrorHandler(handler.ErrorHandler())

    return _app
}
```

## 日志输出示例

当发生数据库错误时，Handler 层会输出以下格式日志：

```
【数据库错误】POST /v1/merchant/create - IP: 127.0.0.1 - 操作: CreateMerchant.InsertOne - 错误: duplicate key value violates unique constraint "merchants_pkey"
```

## 错误响应格式

### 数据库错误
```json
{
  "code": 500,
  "message": "数据库操作失败",
  "op": "CreateMerchant.InsertOne"
}
```

### 验证错误
```json
{
  "code": 400,
  "message": "用户名不能为空",
  "op": "CreateUser.Validate"
}
```

### 权限错误
```json
{
  "code": 403,
  "message": "权限不足",
  "op": "UpdateMerchant.CheckPermission"
}
```

### 未找到错误
```json
{
  "code": 404,
  "message": "商家不存在",
  "op": "GetMerchant.Query"
}
```

## 注意事项

1. **不要在自动生成的 DAO 文件中修改代码**，因为重新生成时会丢失修改
2. **始终在 Service 层包装错误**，添加有意义的操作名称
3. **操作名称格式**：`函数名.操作描述`，如 `CreateMerchant.InsertOne`
4. **只在 Handler 层记录日志**，确保错误能够被追踪和调试
5. **对于已知的业务错误**（如用户不存在），使用 `apperr.NewNotFound()` 等函数创建
6. **对于未知的错误**，使用 `apperr.WrapDB()` 包装，保留原始错误信息

## 完整示例

### Service 层

```go
func CreateMerchant(c context.Context, in *dto.CreateMerchantReq, out *dto.CreateMerchantResp) error {
    vis := getVisitor(c)

    // 查询商家是否存在
    _existing, has, e := dao.Merchants().GetByID(c, vis.Id)
    if e != nil {
        return apperr.WrapDB(e, "CreateMerchant.CheckExist")
    }
    if has {
        return apperr.NewValidation("CreateMerchant.CheckExist", nil, "商家已存在")
    }
    _existing.Free()

    // 创建商家
    _merchant := do.NewMerchants()
    // ... 设置字段 ...

    ok, e := dao.Merchants().InsertOne(c, _merchant)
    if e != nil {
        return apperr.WrapDB(e, "CreateMerchant.InsertOne")
    }
    if !ok {
        return apperr.New(apperr.ErrorTypeUnknown, "CreateMerchant.InsertOne", nil, "插入失败")
    }

    out.MerchantID = _merchant.Id
    return nil
}
```

### Handler 层

```go
func Init() api.Application {
    _app := api.NewApplication("Mingtianjian", "0.1")

    // 添加错误处理中间件
    _app.Use(handler.ErrorMiddleware())

    // 注册自定义错误处理器
    _app.SetErrorHandler(handler.ErrorHandler())

    return _app
}
```

通过这种方式，我们实现了清晰的错误分层处理：
- DAO 层：只返回错误
- Service 层：包装错误，添加上下文
- Handler 层：记录日志，格式化响应
