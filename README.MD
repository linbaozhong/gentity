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
	dbx, err = ace.Connect(ace.Context, "mysql",
		"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	dbx.SetMaxOpenConns(50)
	dbx.SetMaxIdleConns(25)
	dbx.SetDebug(true)
	dbx.SetCache(ace.CacheTypeSyncMap, 0)
```
4. 使用
```go
func TestCreateSet(t *testing.T) {
	defer dbx.Close()
	x := dao.Company(dbx).C()
	result, err := x.Set(
		tblcompany.CorpName.Set("aaaaaa"),
		tblcompany.IsAuthenticated.Set(true),
	).Exec(context.Background())

	t.Log(x)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestCreateCols(t *testing.T) {
	defer dbx.Close()
	x := dao.Company(dbx).C()
	result, err := x.Cols(
		tblcompany.CorpName,
		tblcompany.IsAuthenticated,
		tblcompany.State,
		tblcompany.CreatedTime,
	).StructBatch(context.Background(), &db.Company{
		FullCorpName:     "m1",
		IsEcologicalCorp: false,
		State:            0,
		CreatedTime:      time.Now(),
	}, &db.Company{
		FullCorpName:     "m2",
		IsEcologicalCorp: true,
		State:            1,
		CreatedTime:      time.Now(),
	}, &db.Company{
		FullCorpName:     "m3",
		IsEcologicalCorp: true,
		State:            0,
		CreatedTime:      time.Now(),
	})

	t.Log(x)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestInsert(t *testing.T) {
	defer dbx.Close()
	//
	id, err := dao.Company(dbx).InsertOne(context.Background(), &db.Company{
		FullCorpName:     "x1",
		IsEcologicalCorp: false,
		State:            0,
		CreatedTime:      time.Now(),
	})
	t.Log(id)

	if err != nil {
		t.Fatal(err)
	}

	n, err := dao.Company(dbx).InsertBatch(context.Background(), []*db.Company{{
		FullCorpName:     "x1",
		IsEcologicalCorp: false,
		State:            0,
		CreatedTime:      time.Now(),
	}, {
		FullCorpName:     "x2",
		IsEcologicalCorp: true,
		State:            1,
		CreatedTime:      time.Now(),
	}, {
		FullCorpName:     "x3",
		IsEcologicalCorp: true,
		State:            0,
		CreatedTime:      time.Now(),
	}})

	t.Log(n)

	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateSet(t *testing.T) {
	defer dbx.Close()
	x := dao.Company(dbx).U()
	result, err := x.
		Set(
			tblcompany.CorpName.Set("ttt"),
			tblcompany.IsAuthenticated.Set(true),
		).
		SetExpr(tblcompany.State.Incr()).
		Where(
			tblcompany.Id.Eq(1),
		).
		Exec(context.Background())

	t.Log(x)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestUpdateCols(t *testing.T) {
	defer dbx.Close()

	x := dao.Company(dbx).U()
	result, err := x.Cols(
		tblcompany.FullCorpName,
		tblcompany.IsEcologicalCorp,
		tblcompany.State,
		tblcompany.CreatedTime,
	).Where(
		tblcompany.Id.Eq(1),
	).Struct(context.Background(), &db.Company{
		Id:               2,
		FullCorpName:     "ccc",
		IsEcologicalCorp: false,
		State:            1,
		CreatedTime:      time.Now(),
	})

	t.Log(x)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestUpdate(t *testing.T) {
	defer dbx.Close()

	ok, err := dao.Company(dbx).UpdateBatch(context.Background(), []*db.Company{&db.Company{
		Id:               2,
		FullCorpName:     "ccc",
		IsEcologicalCorp: false,
		State:            1,
		CreatedTime:      time.Now(),
	}})

	t.Log(ok)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	defer dbx.Close()

	x := dao.Company(dbx).D()
	result, err := x.Where(
		tblcompany.Id.Eq(1),
	).Exec(context.Background())
	t.Log(x)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestSelect(t *testing.T) {
	defer dbx.Close()
	obj, _, err := dao.Company(dbx).Find4Cols(context.Background(), 0, 0, []dialect.Field{
		tblcompany.Id,
		tblcompany.FullCorpName,
		tblcompany.State,
	},
		tblcompany.Id.Eq(2),dao.Or(tblcompany.Id.Eq(3),tblcompany.State.Eq(1))
	)

	if err != nil {
		t.Fatal(err)
	}
	// t.Log(obj)
	for _, o := range obj {
		t.Log(o)
	}
}
```
提示：查询条件ace.Or、ace.And 和 dao.Or、dao.And等同
5. 实验性查询
结构体嵌套在生成dao代码时，会自动生成关联查询。
关联查询的标签为：db:"ref:关联字段|关联表对应字段"
dao相应查询方法为：FindX、GetX ......
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
注意：一定要要在struct上添加注释：// checker

#### Post请求
* 如果 Content-Type：application/json，req结构体的字段tag为json
* 如果 Content-Type: application/x-www-form-urlencoded，req结构体的字段tag为form
* 如果 Content-Type: multipart/form-data，req结构体的字段tag为form
```
// checker
type ArticleReq struct {
	Title    string `json:"title" valid:"alphanum,required"`
	Emain    string `json:"email" valid:"email~email格式错误,required"`
	Content  string `json:"content" valid:"maxstringlength(50)"`
	AuthorIP string `json:"author_ip" valid:"ipv4"`
	Date     string `json:"date" valid:"-"`
}
// checker
type ChannelReq struct {
	ID int `form:"id" valid:"int,required"`
}
```
#### Get请求
* 首先尝试读取query，req结构体的字段 tag 为 url 或者 param。
* 如果query为空，则尝试读取form，req结构体的字段tag为form。
```
// checker
type ArticleReq struct {
	ID       int    `param:"id" valid:"int,required"`
	Title    string `url:"title" valid:"alphanum,required"`
	Emain    string `url:"email" valid:"email~email格式错误,required"`
}
// checker
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
"alphanum":           
"utfletternum":       
"numeric":            
"utfnumeric":         
"utfdigit":           
"hexadecimal":        
"hexcolor":           
"rgbcolor":           
"lowercase":          
"uppercase":          
"int":                
"float":              
"null":               
"uuid":               
"uuidv3":             
"uuidv4":             
"uuidv5":             
"creditcard":         
"isbn10":             
"isbn13":             
"json":               
"multibyte":          
"ascii":              
"printableascii":     
"fullwidth":          
"halfwidth":          
"variablewidth":      
"base64":             
"datauri":            
"ip":                 
"port":               
"ipv4":               
"ipv6":               
"dns":                
"host":               
"mac":                
"latitude":           
"longitude":          
"ssn":                
"semver":             
"rfc3339":            
"rfc3339WithoutZone": 
"ISO3166Alpha2":      
"ISO3166Alpha3":      
"ulid":               
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
// checker
type PostRequest struct {
	Title    string `json:"title" valid:"alphanum,required"`
	Emain    string `json:"email" valid:"email~email格式错误,required"`
	Content  string `json:"content" valid:"maxstringlength(50)"`
	AuthorIP string `json:"author_ip" valid:"ipv4"`
	Date     string `json:"date" valid:"-"`
}
// checker
type GetRequest struct {
	ID int `url:"id" valid:"int,required"`
}

```