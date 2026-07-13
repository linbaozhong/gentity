# validator — gentity 的通用校验工具包

`validator` 是 gentity 框架的**纯函数式校验层**，提供一组无状态、并发安全的字符串 / 数值校验函数。
它不依赖任何 Web 框架，被 `ack` 的请求 `Checker`、代码生成器的结构体 tag 校验等场景复用。

- 全部校验函数签名统一为 `func(str string) bool`（或带参数的泛型函数），返回 `true/false`
- 正则全部在 `package` 级预编译（`rxXxx = regexp.MustCompile(...)`），零运行时编译开销
- 数值类校验（`InRange`/`Range`/`Min`/`Max`/`IsIn`）使用 Go 泛型，支持任意 `constraints.Ordered` 类型
- 提供 `TagMap` / `ParamTagMap` / `ParamTagRegexMap` 三张映射表，可把结构体 tag（如 `valid:"email"`）关联到具体校验函数
- 采用 Apache-2.0 许可，源自社区 validator 实现风格的裁剪与扩展

---

## 1. 文件结构

| 文件 | 职责 |
|------|------|
| `validator.go` | 核心：各类 `IsXXX` 格式校验、长度/范围/集合校验、RSA 公钥校验 |
| `numerics.go` | 数值范围与比较：`InRange` / `Range` / `Min` / `Max` |
| `patterns.go` | 正则常量与预编译正则（`Email`/`UUID`/`URL`/`CnMobile` ...） |
| `types.go` | `TagMap` / `ParamTagMap` / `ParamTagRegexMap` 映射表 + `ISO3166List` / `ISO4217List` 数据表 |
| `utils.go` | `Matches` 正则匹配辅助 |

---

## 2. 字符串格式校验（IsXXX）

最常用的格式校验函数（返回 `bool`）：

| 函数 | 校验内容 |
|------|----------|
| `IsEmail(str)` | 邮箱 |
| `IsMobile(str)` | 中国大陆手机号（`^1[3456789]\d{9}$`） |
| `IsURL(str)` / `IsRequestURL` / `IsRequestURI` | URL / 请求 URL / 请求 URI |
| `IsIP` / `IsIPv4` / `IsIPv6` / `IsPort` / `IsDNSName` / `IsHost` / `IsMAC` | 网络地址 |
| `IsUUID` / `IsUUIDv3/v4/v5` | UUID |
| `IsCreditCard(str)` | 信用卡号（Luhn 算法） |
| `IsISBN10` / `IsISBN13` / `IsISBN(str, ver)` | 书号 |
| `IsJSON(str)` | JSON 合法性 |
| `IsBase64(str)` / `IsDataURI(str)` | Base64 / Data URI |
| `IsLatitude` / `IsLongitude` | 经纬度 |
| `IsSemver(str)` | 语义化版本号 |
| `IsRFC3339` / `IsRFC3339WithoutZone` / `IsTime(str, format)` | 时间格式 |
| `IsISO3166Alpha2` / `IsISO3166Alpha3` | 国家代码 |
| `IsISO4217(str)` | 货币代码 |
| `IsIMEI(str)` / `IsSSN(str)` / `IsULID(str)` | IMEI / 美国社保号 / ULID |

示例：

    if !validator.IsEmail(req.Email) {
        return errors.New("邮箱格式不正确")
    }
    if !validator.IsMobile(req.Phone) {
        return errors.New("手机号格式不正确")
    }

---

## 3. 字符集与大小写校验

| 函数 | 说明 |
|------|------|
| `IsAlpha` / `IsUTFLetter` | 仅字母（ASCII / 全语言） |
| `IsAlphanumeric` / `IsUTFLetterNumeric` | 仅字母+数字 |
| `IsNumeric` / `IsUTFNumeric` / `IsUTFDigit` | 仅数字 |
| `IsHexadecimal` / `IsHexcolor` / `IsRGBcolor` | 十六进制 / 颜色 |
| `IsASCII` / `IsPrintableASCII` / `IsMultibyte` | 字符编码范围 |
| `IsLowerCase` / `IsUpperCase` / `HasLowerCase` / `HasUpperCase` | 大小写 |
| `IsFullWidth` / `IsHalfWidth` / `IsVariableWidth` | 全角/半角 |
| `HasWhitespace` / `IsWhitespaceOnly` | 空白字符 |
| `IsInt` / `IsFloat` | 整数/浮点字符串 |

---

## 4. 长度、范围与集合校验

| 函数 | 签名 | 说明 |
|------|------|------|
| `ByteLength(str, min, max)` | `func(string, ...int) bool` | 字节长度区间 |
| `StringLength(str, min, max)` | 泛型 | 字符（rune）长度区间 |
| `RuneLength(str, min, max)` | 泛型 | `StringLength` 别名 |
| `MinStringLength(str, min)` | `func(string, int) bool` | 最小字符长度 |
| `MaxStringLength(str, max)` | `func(string, int) bool` | 最大字符长度 |
| `Range(value, min, max)` | 泛型 | 数值落在 `[min, max]` |
| `Min(value, min)` / `Max(value, max)` | 泛型 | 数值下/上界 |
| `InRange(value, min, max)` | 泛型 | `Range` 底层实现 |
| `IsIn[T](str, params...)` | 泛型 | 属于给定集合 |
| `IsInRaw[T](str, params...)` | 泛型 | `IsIn` 别名 |

示例：

    // 用户名 3~20 个字符
    if !validator.StringLength(name, 3, 20) {
        return errors.New("用户名长度需为 3-20")
    }
    // 年龄在 1~120
    if !validator.Range(age, 1, 120) {
        return errors.New("年龄不合法")
    }
    // 状态必须是枚举值之一
    if !validator.IsIn(status, "active", "inactive", "banned") {
        return errors.New("状态非法")
    }

---

## 5. 正则与自定义匹配

    // 直接传入正则
    if !validator.Matches(code, `^\d{6}$`) {
        return errors.New("验证码必须为 6 位数字")
    }
    // 配合 tag 引擎使用（见第 7 节）
    validator.StringMatches(code, `^\d{6}$`)

`Matches` 内部用 `regexp.MatchString`，匹配失败或正则非法时返回 `false`。

---

## 6. 密钥与编码校验

| 函数 | 说明 |
|------|------|
| `IsRsaPub(str, keylen)` | 校验 RSA 公钥 PEM / Base64 且指定位长 |
| `IsRsaPublicKey(str, keylen)` | `IsRsaPub` 底层实现 |
| `IsBase64(str)` / `IsDataURI(str)` | 见第 2 节 |

示例：

    if !validator.IsRsaPub(pubKeyPEM, 2048) {
        return errors.New("RSA 公钥无效或位长不为 2048")
    }

---

## 7. 通过 Tag 进行结构体校验

`types.go` 中提供了三张映射表，把字符串 tag 关联到具体校验函数，供校验引擎（如 gentity 生成的 `Checker`）按结构体字段 tag 自动执行：

- `tagName = "valid"`：默认 tag 名
- `TagMap`：无参 tag → 函数名（如 `email → IsEmail`）
- `ParamTagMap`：带参 tag → 函数名（如 `length → ByteLength`）
- `ParamTagRegexMap`：带参 tag → 解析正则（如 `^length$(\d+)\|(\d+)$$`）

支持的内置 tag（节选）：

| tag | 对应函数 |
|-----|----------|
| `email` / `url` / `mobile` / `ip` / `uuid` ... | 见 `TagMap` |
| `length(min|max)` | `ByteLength` |
| `runelength(min|max)` / `stringlength(min|max)` | `StringLength` |
| `minstringlength(n)` / `maxstringlength(n)` | 字符长度上下界 |
| `range(min|max)` | `Range` |
| `min(n)` / `max(n)` | `Min` / `Max` |
| `in(a|b|c)` | `IsIn` |
| `matches(re)` | `StringMatches` |
| `rsapub(bits)` | `IsRsaPub` |

字段 tag 写法示例：

    type RegisterReq struct {
        Email string `valid:"email"`
        Name  string `valid:"stringlength(1|20),alphanum"`
        Age   int    `valid:"range(1|120)"`
        Role  string `valid:"in(admin|user|guest)"`
        Code  string `valid:"matches(^\\d{6}$)"`
    }

> 说明：`validator` 包本身只提供 tag → 函数的映射与底层校验函数；真正的"遍历结构体字段、解析 tag、收集错误"的引擎在 gentity 的校验框架中（如 `ack` 的 `Checker` 接口），本包是被其调用的能力底座。

---

## 8. 内置数据表

| 变量 | 说明 |
|------|------|
| `ISO3166List` | `[]ISO3166Entry`，含全球国家/地区的 Alpha2 / Alpha3 / Numeric 代码，供 `IsISO3166Alpha2/Alpha3` 使用 |
| `ISO4217List` | `[]string`，ISO 4217 货币代码，供 `IsISO4217` 使用 |

---

## 9. 设计要点

- **无状态纯函数**：所有校验函数不保存状态，可安全并发调用。
- **零运行时正则编译**：正则集中在 `patterns.go` 以 `package` 变量预编译，避免重复 `Compile` 开销。
- **泛型数值校验**：`InRange` / `Range` / `Min` / `Max` / `IsIn` 基于 `constraints.Ordered`，一套逻辑覆盖 int / float / string 等可比较类型。
- **tag 驱动**：`TagMap` 系列把"声明式 tag"映射到"命令式函数"，让结构体字段校验可配置化，配合代码生成器自动产出 `Check()` 逻辑。
- **可组合**：各 `IsXXX` 返回 `bool`，业务层可自由 `&&` / `||` 组合出复杂规则。

---

## 10. 命名来源

`validator` 字面即"校验器 / 验证器"，是 gentity 框架的**校验基础设施层**——它不关心 HTTP、不关心数据库，只负责"这个字符串/数值合不合法"。上层 `ack` 的请求 `Checker`、生成代码的字段约束都建立在这套原子校验能力之上，形成 `validator`（校验）→ `ack`（请求）→ `ace`（数据）的能力分层。
