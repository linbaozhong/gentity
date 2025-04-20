# gentity

### 1. 生成api框架
### 2. 生成 ORM 代码
* 根据结构体(struct)文件，生成相应的ORM代码
* 数据库表结构文件(.sql)，生成相应的ORM代码
* 数据库连接字符串(dns)，生成相应的ORM代码
### 3. 根据前端传入的request对象，生成相应的validator代码

# 使用方法
### 将项目git到本地
```
go get -u github.com/linbaozhong/gentity@latest
```
### 安装命令
```
go install github.com/linbaozhong/gentity@latest
```
### 生成 API 项目
```
gentity api [项目名称]
```
### 生成ORM代码
#### 定义一个struct，添加注释 tablename 数据库表名
```go
// tablename user
type User struct {
	ID         types.ID `json:"id" db:"'id' pk auto"`
	Avatar     string   `json:"avatar" db:"avatar"`
	Nickname   string   `json:"nickname" db:"nickname"`
	Status     int8     `json:"status" db:"status"`
	CreatedTime int64    `json:"created_time" db:"created_time <-"`
}
```
注释：
- `db:"'id' pk auto"` 表示数据库字段名，主键，自增
- `db:"created_time <-"` 表示数据库字段名，只读。 `<-`: 只读;`->`: 只写;`-`:不可读写
#### 用法：
```
gentity command [model路径] [ "数据库驱动" "SQL文件路径" |  "数据库驱动" "数据库连接字符串"] [flags]
```
#### 示例
```
// 生成api框架
gentity api project_name

// 根据当前目录的struct，生成相应的dao层
gentity dao

// 根据 .\do 路径的struct，生成相应的dao层
gentity dao .\do

// 根据指定的数据库引擎名称以及连接字符串在 .\do 路径，生成相应的struct和dao层
gentity db .\do mysql "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

// 根据指定的数据库引擎名称以及数据库建表sql文件在 .\do 路径，生成相应的struct和dao层
gentity sql .\do mysql .\database.sql

// 检查 .\dto 目录下所有request对象，生成相应的 validator 方法
gentity check .\dto
```
###### 1. 进入model目录
```
cd model
gentity dao
```
###### 2. 或者指定model目录
```
gentity dao ./model
```
###### 3. 指定数据库驱动、包含 Create Table 的 SQL文件，先在当前目录生成struct文件，再生成相应的ORM文件
```
gentity sql . mysql ./database.sql
```
###### 4. 指定数据库驱动、连接字符串，先在当前目录生成struct文件，再生成相应的ORM文件
```
gentity db . mysql "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
```
###### 5. 检查 .\dto 目录下所有request对象，生成相应的 validator 方法
```
gentity check .\dto
```
###### 6. 生成api框架
```
gentity api project_name
```
##### 注意：命令参数的顺序
## 示例
1. 导包
```go
import (
	"context"
	"github.com/linbaozhong/gentity/pkg/log"
	_ "github.com/go-sql-driver/mysql"
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
	"github.com/linbaozhong/gentity/pkg/ace"
)
```
2. 定义一个连接
```go
var (
	dbx *ace.DB
)
```
3. 初始化连接
```go
	var err error
	dbx, err = ace.Connect("mysql",
		"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	dbx.SetMaxOpenConns(50)
	dbx.SetMaxIdleConns(25)
	// 开启调试模式，打印sql语句
	// 注意：开启调试模式会影响性能，建议在开发环境开启，生产环境关闭
	// 如果关闭调试模式，也可以在dao方法中链式调用 ToSql()方法，打印单行sql语句
	dbx.SetDebug(true)
```

4. 使用
```go

// TestCreateSet 测试函数，用于测试不同方式插入数据的功能。
// 该函数会分别使用生成的 DAO 方法、DAO 构建器以及 ace 包提供的通用方法来插入数据，
// 并记录插入操作的结果和可能出现的错误。
func TestCreateSet(t *testing.T) {
	// 确保在函数执行结束后关闭数据库连接
	defer dbx.Close()
	// 使用生成的dao方法插入数据
	// 注意:
	//  1. 插入数据时，需要传入context.Context
	//  2. 插入数据时，需要传入dialect.Setter类型的参数
	// 调用 daocompany 包的 New 函数创建一个新的 DAO 实例，并调用其 Insert 方法插入数据
	// id 为插入数据的主键值，err 为可能出现的错误
	// 这里的 .ToSql()方法是为了不执行sql命令只打印sql语句，在需要检查sql语句是否正确时使用
	id, err := dao.Company(dbx).ToSql().Insert(context.Background(),
		// 设置公司的长名称为 "aaaaaa"
		tblcompany.LongName.Set("aaaaaa"),
		// 设置公司的状态为 1
		tblcompany.State.Set(1),
	)

	// 记录插入数据的主键值和可能出现的错误
	t.Log(id, err)

	// 若插入操作出现错误，则终止测试并输出错误信息
	if err != nil {
		t.Fatal(err)
	}
	
	// 使用 ace 包的通用方法插入数据
	// 调用 ace 包的 Table 函数指定要操作的表，设置插入的数据
	// 然后调用 Create 方法创建执行器并执行插入操作
	// result 为执行结果，包含插入的 ID 等信息，err 为可能出现的错误
	result, err = ace.
		Table(db.CompanyTableName).
		Set(
			// 设置公司的长名称为 "aaaaaa"
			tblcompany.LongName.Set("aaaaaa"),
			// 设置公司的状态为 1
			tblcompany.State.Set(1),
			// 设置公司的地址为 "beijing"
			tblcompany.Address.Set("beijing"),
		).ToSql().
		Create(dbx).
		Exec(context.Background())
	// 若插入操作出现错误，则终止测试并输出错误信息
	if err != nil {
		t.Fatal(err)
	}
	// 记录插入操作返回的最后插入的 ID
	t.Log(result.LastInsertId())
}

// TestCreateStruct 测试函数，用于测试不同方式通过结构体插入数据的功能。
// 该函数会分别使用生成的 DAO 批量插入方法、DAO 构建器以及 ace 包提供的通用方法，
// 通过结构体插入数据，并记录插入操作的结果和可能出现的错误。
func TestCreateStruct(t *testing.T) {
	// 确保在函数执行结束后关闭数据库连接
	defer dbx.Close()

	// 使用生成的 DAO 批量插入方法，通过结构体切片批量插入数据
	// n 为插入操作影响的行数，err 为可能出现的错误
	n, err := dao.Company(dbx).
		InsertBatch(context.Background(),
			[]*db.Company{
				{
					LongName: "m1",
					Address:  "北京",
					Email:    "1@2.com",
					State:    1,
				},
				{
					LongName: "m2",
					Address:  "上海",
					Email:    "",
					State:    1,
				},
			},
			tblcompany.LongName, tblcompany.Address, tblcompany.Email, tblcompany.State)

	t.Log(n, err)

	if err != nil {
		t.Fatal(err)
	}

	// 使用 ace 包的通用方法，通过单个结构体插入数据
	// result 为执行结果，包含插入的 ID 等信息，err 为可能出现的错误
	result, err = ace.
		Table(db.CompanyTableName).
		Cols(
			tblcompany.LongName,
			tblcompany.Address,
			tblcompany.Email,
		).
		Create(dbx).
		Struct(context.Background(),
			&db.Company{
				LongName: "m1",
				Address:  "北京",
				Email:    "1@2.com",
			})
	// 若插入操作出现错误，则终止测试并输出错误信息
	if err != nil {
		t.Fatal(err)
	}
	// 记录插入操作返回的最后插入的 ID
	t.Log(result.LastInsertId())
}

// TestUpdateSet 测试函数，用于测试不同方式更新数据的功能。
// 该函数会分别使用生成的 DAO 方法、DAO 构建器以及 ace 包提供的通用方法来更新数据，
// 并记录更新操作的结果和可能出现的错误。
func TestUpdateSet(t *testing.T) {
	// 确保在函数执行结束后关闭数据库连接
	defer dbx.Close()
	// 使用生成的 DAO 方法更新数据
	// 调用 daocompany 包的 New 函数创建一个新的 DAO 实例，并调用其 Update 方法更新数据
	// n 为更新操作影响的行数，err 为可能出现的错误
	// 第一个参数为 context 上下文，第二个参数为要更新的字段设置，第三个参数为更新条件
	n, err := dao.Company(dbx).
		Update(context.Background(),
			ace.Sets(
				tblcompany.LongName.Set("aaaaaa"),
				tblcompany.State.Set(1),
			).ToSlice(),
			tblcompany.Id.Eq(1),
		)
	t.Log(n, err)
	if err != nil {
		t.Fatal(err)
	}

	// 使用 ace 包的通用方法更新数据
	// 调用 ace 包的 Table 函数指定要操作的表，设置要更新的数据
	// 并通过 Where 方法指定更新条件，然后调用 Update 方法创建执行器并执行更新操作
	// result 为执行结果，err 为可能出现的错误
	result, err = ace.
		Table(db.CompanyTableName).
		Set(
			tblcompany.LongName.Set("aaaaaa"),
			tblcompany.State.Set(1),
		).
		Where(tblcompany.Id.Eq(1)).
		Update(dbx).
		Exec(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
}

// TestSelect 测试函数，用于测试不同方式从数据库中查询单条公司数据的功能。
// 该函数会分别使用生成的 DAO 方法、DAO 构建器以及 ace 包提供的通用方法来查询数据，
// 并根据查询结果记录信息或终止测试。
func TestSelect(t *testing.T) {
	// 确保在函数执行结束后关闭数据库连接
	defer dbx.Close()

	// 使用生成的 DAO 方法查询数据
	// 调用 daocompany 包的 New 函数创建一个新的 DAO 实例，并调用其 Get 方法进行查询
	// obj 为查询结果对象，has 表示是否查询到数据，err 为可能出现的错误
	// 这里指定查询 tblcompany 表的 Id 和 LongName 列，查询条件为 Id 等于 1
	obj, has, err := dao.Company(dbx).
		Get(context.Background(),
			ace.Cols(tblcompany.Id, tblcompany.LongName).
				Where(tblcompany.Id.Eq(1)),
		)
	// 若查询操作出现错误，则终止测试并输出错误信息
	if err != nil {
		t.Fatal(err)
	}
	// 若未查询到数据，则终止测试并输出提示信息
	if !has {
		t.Fatal("not found")
	}
	// 确保在函数结束时释放查询结果对象的资源
	defer obj.Free()

	// 记录查询结果对象
	t.Log(obj)

	// 创建一个新的 Company 实例
	obj = db.NewCompany()

	// 使用 ace 包的通用方法查询数据
	// 调用 ace 包的 Table 函数指定要操作的表，设置查询条件为 Id 等于 1
	// 然后调用 Select 方法指定数据库连接，再调用 Get 方法将查询结果填充到 obj 中
	// err 为可能出现的错误
	err = ace.Table(db.CompanyTableName).
		Where(tblcompany.Id.Eq(1)).
		Select(dbx).
		Get(context.Background(), &obj)
	// 根据查询结果处理不同情况
	switch err {
	case nil:
		// 若查询成功，记录查询结果对象
		t.Log(obj)
	case sql.ErrNoRows:
		// 若未查询到数据，终止测试并输出提示信息
		t.Fatal("没找到")
	default:
		// 若出现其他错误，终止测试并输出错误信息
		t.Fatal(err)
	}
}

```
提示：查询条件ace.Or、ace.And 和 dao.Or、dao.And等同
5. 实验性查询

* 结构体嵌套在生成dao代码时，会自动生成关联查询。
* 关联查询的标签为：db:"ref:关联字段|关联表对应字段"
* dao相应查询方法为：FindX、GetX ......

例如：
```go
// tablename user_log
type UserLog struct {
	ace.Model
	Id         types.BigInt `json:"id,omitempty" db:"'id' pk auto size(20)"`           //
	UserId     types.BigInt `json:"user_id,omitempty" db:"'user_id' size(20)"`         //
	LoginTime  types.Time   `json:"login_time,omitempty" db:"'login_time'"`            // 登录时间
	Device     types.String `json:"device,omitempty" db:"'device' size(255)"`          // 登录终端参数
	Os         types.String `json:"os,omitempty" db:"'os' size(10)"`                   //
	OsVersion  types.String `json:"os_version,omitempty" db:"'os_version' size(10)"`   //
	AppName    types.String `json:"app_name,omitempty" db:"'app_name' size(10)"`       //
	AppVersion types.String `json:"app_version,omitempty" db:"'app_version' size(10)"` //
	Ip         types.String `json:"ip,omitempty" db:"'ip' size(50)"`                   // ip地址

	User User `json:"user,omitempty" db:"ref:UserId|Id"`
}

// tablename user
type User struct {
	ace.Model
	Id    types.BigInt `json:"id,omitempty" db:"'id' pk auto size(20)"` //
	Uuid  types.String `json:"uuid,omitempty" db:"'uuid' size(45)"`     // 用户识别码
	Ctime types.Time   `json:"ctime,omitempty" db:"'ctime'"`            //

	UserLogs []UserLog `json:"user_logs,omitempty" db:"ref:Id|UserId"` //
}

```
## 生成Checker接口
检查指定目录下所有request对象，生成相应的 validator 方法
```
gentity check [dir]
```

### 示例
注意：
* 如果是请求结构体，一定要要在struct上添加注释：// request
* 如果是响应结构体，一定要要在struct上添加注释：// response

#### Post请求
* 如果 Content-Type：application/json，req结构体的字段tag为json
* 如果 Content-Type: application/x-www-form-urlencoded，req结构体的字段tag为form
* 如果 Content-Type: multipart/form-data，req结构体的字段tag为form
```
// request
type ArticleReq struct {
	Title    string `json:"title" valid:"alphanum,required"`
	Emain    string `json:"email" valid:"email~email格式错误,required"`
	Content  string `json:"content" valid:"maxstringlength(50)"`
	AuthorIP string `json:"author_ip" valid:"ipv4"`
	Date     string `json:"date" valid:"-"`
}

// response
type ArticleReq struct {
	Title    string `json:"title" valid:"alphanum,required"`
	Emain    string `json:"email" valid:"email~email格式错误,required"`
	Content  string `json:"content" valid:"maxstringlength(50)"`
	AuthorIP string `json:"author_ip" valid:"ipv4"`
	Date     string `json:"date" valid:"-"`
}

// request
type ChannelReq struct {
	ID int `form:"id" valid:"int,required"`
}
```
#### Get请求
* 首先尝试读取query，req结构体的字段 tag 为 url 或者 param。
* 如果query为空，则尝试读取form，req结构体的字段tag为form。
```
// request
type ArticleReq struct {
	ID       int    `param:"id" valid:"int,required"`
	Title    string `url:"title" valid:"alphanum,required"`
	Emain    string `url:"email" valid:"email~email格式错误,required"`
}
// request
type ChannelReq struct {
	ID int `form:"id" valid:"int,required"`
}
```
### valid 支持的 struct filed 校验器
```
"email":              // email
"url":                // url
"dialstring":         // dial string
"requrl":             // request url
"requri":             // request uri
"alpha":              // alpha
"utfletter":          // utf letter
"alphanum":           // alpha numeric
"utfletternum":       // utf letter numeric
"numeric":            // numeric
"utfnumeric":         // utf numeric
"utfdigit":           // utf digit
"hexadecimal":        // hexadecimal
"hexcolor":           // hex color
"rgbcolor":           // rgb color
"lowercase":          // lowercase
"uppercase":          // uppercase
"haslowercase":       // has lowercase
"hasuppercase":       // has uppercase
"whitespace":         // whitespace
"haswhitespace":      // has whitespace
"int":                // int
"float":              // float
"null":               // null
"uuid":               // uuid
"uuidv3":             // uuidv3
"uuidv4":             // uuidv4
"uuidv5":             // uuidv5
"creditcard":         // credit card
"isbn10":             // isbn10
"isbn13":             // isbn13
"json":               // json
"multibyte":          // multi byte
"ascii":              // ascii
"printableascii":     // printable ascii
"fullwidth":          // full width
"halfwidth":          // half width
"variablewidth":      // variable width
"base64":             // base64
"datauri":            // data uri
"ip":                 // ip
"port":               // port
"ipv4":               // ipv4
"ipv6":               // ipv6
"dns":                // dns
"host":               // host
"mac":                // mac
"latitude":           // latitude
"longitude":          // longitude
"ssn":                // ssn
"semver":             // semver
"rfc3339":            // rfc3339
"rfc3339WithoutZone": // rfc3339 without zone
"ISO3166Alpha2":      // ISO3166 Alpha2
"ISO3166Alpha3":      // ISO3166 Alpha3
"ulid":               // ulid
"mobile":             // 手机号码
```
### 校验器可用参数
```
"range(min|max)"                    // 数值范围
"length(min|max)"                   // 字节长度
"runelength(min|max)"               // Rune串长度
"stringlength(min|max)"             // 字符串长度
"matches(pattern)"                  // 正则表达式
"in(string1|string2|...|stringN)"   // 枚举
"rsapub(keylength)"                 // RSA公钥
"minstringlength(int)               // 字符串最小长度
"maxstringlength(int)               // 字符串最大长度
```
### 自定义错误信息
在 valid 注解的参数后面添加：~错误信息

参照下面示例代码 struct 中的 email 字段

### 示例
注意：一定要要在struct上添加注释：// checker
```
// request
type PostRequest struct {
	Title    string `json:"title" valid:"alphanum,required"`
	Emain    string `json:"email" valid:"email~email格式错误,required"`
	Content  string `json:"content" valid:"maxstringlength(50)"`
	AuthorIP string `json:"author_ip" valid:"ipv4"`
	Date     string `json:"date" valid:"-"`
}
// request
type GetRequest struct {
	ID int `url:"id" valid:"int,required"`
}

```