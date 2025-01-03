# gentity
### 根据结构体(struct)文件/数据库表结构文件(.sql)/数据库连接字符串(dns)，生成对应的ORM代码

# 使用方法
## 将项目git到本地
```
go get -u github.com/linbaozhong/gentity@latest
```
## 安装命令
```
go install github.com/linbaozhong/gentity@latest
```
## 生成 API 项目
```
gentity init [项目名称]
```
## 生成ORM代码
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
gentity command model路径 [ "数据库驱动" "SQL文件路径" |  "数据库驱动" "数据库连接字符串"] [flags]
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
## 创建连接
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
		tblcompany.Id.Eq(2),
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

## validator
### 支持的 struct filed 校验器
```
"email":              IsEmail,
"url":                IsURL,
"dialstring":         IsDialString,
"requrl":             IsRequestURL,
"requri":             IsRequestURI,
"alpha":              IsAlpha,
"utfletter":          IsUTFLetter,
"alphanum":           IsAlphanumeric,
"utfletternum":       IsUTFLetterNumeric,
"numeric":            IsNumeric,
"utfnumeric":         IsUTFNumeric,
"utfdigit":           IsUTFDigit,
"hexadecimal":        IsHexadecimal,
"hexcolor":           IsHexcolor,
"rgbcolor":           IsRGBcolor,
"lowercase":          IsLowerCase,
"uppercase":          IsUpperCase,
"int":                IsInt,
"float":              IsFloat,
"null":               IsNull,
"uuid":               IsUUID,
"uuidv3":             IsUUIDv3,
"uuidv4":             IsUUIDv4,
"uuidv5":             IsUUIDv5,
"creditcard":         IsCreditCard,
"isbn10":             IsISBN10,
"isbn13":             IsISBN13,
"json":               IsJSON,
"multibyte":          IsMultibyte,
"ascii":              IsASCII,
"printableascii":     IsPrintableASCII,
"fullwidth":          IsFullWidth,
"halfwidth":          IsHalfWidth,
"variablewidth":      IsVariableWidth,
"base64":             IsBase64,
"datauri":            IsDataURI,
"ip":                 IsIP,
"port":               IsPort,
"ipv4":               IsIPv4,
"ipv6":               IsIPv6,
"dns":                IsDNSName,
"host":               IsHost,
"mac":                IsMAC,
"latitude":           IsLatitude,
"longitude":          IsLongitude,
"ssn":                IsSSN,
"semver":             IsSemver,
"rfc3339":            IsRFC3339,
"rfc3339WithoutZone": IsRFC3339WithoutZone,
"ISO3166Alpha2":      IsISO3166Alpha2,
"ISO3166Alpha3":      IsISO3166Alpha3,
"ulid":               IsULID,
```
### 校验器可用参数
```
"range(min|max)"                    Range,
"length(min|max)"                   ByteLength,
"runelength(min|max)"               RuneLength,
"stringlength(min|max)"             StringLength,
"matches(pattern)"                  StringMatches,
"in(string1|string2|...|stringN)"   IsIn,
"rsapub(keylength)"                 IsRsaPub,
"minstringlength(int)               MinStringLength,
"maxstringlength(int)               MaxStringLength,
```
### 自定义错误信息
在 valid 注解的参数后面添加：~错误信息
### 示例
```go
type Post struct {
	Title    string `valid:"alphanum,required"`
	Emain    string `valid:"email~email格式错误,required"`
	Content  string `valid:"maxstringlength(50)"`
	AuthorIP string `valid:"ipv4"`
	Date     string `valid:"-"`
}

```