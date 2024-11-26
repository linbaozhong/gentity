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
gentity model路径 [ "数据库驱动" "SQL文件路径" |  "数据库驱动" "数据库连接字符串"] [flags]
```
#### 示例
```
gentity
gentity .\do
gentity .\do mysql "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
gentity .\do mysql .\database.sql
```
###### 1. 进入model目录
```
cd model
gentity
```
###### 2. 或者指定model目录
```
gentity ./model
```
###### 3. 指定数据库驱动、包含 Create Table 的 SQL文件，先在当前目录生成struct文件，再生成相应的ORM文件
```
gentity . mysql ./database.sql
```
###### 4. 指定数据库驱动、连接字符串，先在当前目录生成struct文件，再生成相应的ORM文件
```
gentity . mysql "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
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
	x := Company(dbx).C()
	result, err := x.Set(
		obj.CorpName.Set("aaaaaa"),
		obj.IsAuthenticated.Set(true),
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
	x := Company(dbx).C()
	result, err := x.Cols(
		obj.CorpName,
		obj.IsAuthenticated,
		obj.State,
		obj.CreatedTime,
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
	id, err := Company(dbx).InsertOne(context.Background(), &db.Company{
		FullCorpName:     "x1",
		IsEcologicalCorp: false,
		State:            0,
		CreatedTime:      time.Now(),
	})
	t.Log(id)

	if err != nil {
		t.Fatal(err)
	}

	n, err := Company(dbx).InsertBatch(context.Background(), []*db.Company{{
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
	x := Company(dbx).U()
	result, err := x.
		Set(
			obj.CorpName.Set("ttt"),
			obj.IsAuthenticated.Set(true),
		).
		SetExpr(obj.State.Incr()).
		Where(
			obj.Id.Eq(1),
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

	x := Company(dbx).U()
	result, err := x.Cols(
		obj.FullCorpName,
		obj.IsEcologicalCorp,
		obj.State,
		obj.CreatedTime,
	).Where(
		obj.Id.Eq(1),
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

	ok, err := Company(dbx).UpdateBatch(context.Background(), []*db.Company{&db.Company{
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

	x := Company(dbx).D()
	result, err := x.Where(
		obj.Id.Eq(1),
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
	obj, _, err := Company(dbx).Find4Cols(context.Background(), 0, 0, []dialect.Field{
		obj.Id,
		obj.FullCorpName,
		obj.State,
	},
		obj.Id.Eq(2),
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