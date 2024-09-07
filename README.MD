# 使用方法
## 将项目git到本地
```
go get -u github.com/linbaozhong/gentity@latest
```
## 安装命令
```
go install github.com/linbaozhong/gentity/cmd/gentity@latest
```
## 定义一个struct，添加注释 tablename 数据库表名
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
## 生成代码
### 用法：
```
gentity model路径 [ "数据库驱动" "SQL文件路径" |  "数据库驱动" "数据库连接字符串"] [flags]
```
### 示例
```
gentity

gentity .\db

gentity .\db mysql "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

gentity . mysql .\database.sql
```
##### 1. 进入model目录
```
cd model
gentity
```
##### 2. 或者指定model目录
```
gentity ./model
```
##### 3. 指定数据库驱动、包含 Create Table 的 SQL文件，先在当前目录生成struct文件，再生成相应的ORM文件
```
gentity . mysql ./database.sql
```
##### 4. 指定数据库驱动、连接字符串，先在当前目录生成struct文件，再生成相应的ORM文件
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
    tt  = model.UserTbl
)
```
3. 初始化连接
```go
var err error
db, err = ace.Connect("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
if err != nil {
	panic(err)
}
defer db.Close()

db.SetMaxOpenConns(50)
db.SetMaxIdleConns(25)
db.SetDebug(true)
log.RegisterLogger(false)
```
4. 使用
```go
func TestCreateSet(t *testing.T) {
	x := Test(dbx).C()  // Create

	result, err := x.Set(
		tt.Name.Set("aaaaaa"),
		tt.IsAllow.Set(true),
	).Do(context.Background())

	t.Log(x)  // 打印SQL语句

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestCreateCols(t *testing.T) {
	x := Test(dbx).C()  // Create
	result, err := x.Cols(
		tt.Name,
		tt.IsAllow,
		tt.Type,
		tt.CreatedTime,
	).Struct(context.Background(), &db.Test{
		Name:        "m1",
		IsAllow:     false,
		Type:        0,
		CreatedTime: time.Now(),
	}, &db.Test{
		Name:        "m2",
		IsAllow:     true,
		Type:        1,
		CreatedTime: time.Now(),
	}, &db.Test{
		Name:        "m3",
		IsAllow:     true,
		Type:        0,
		CreatedTime: time.Now(),
	})

	t.Log(x)  // 打印SQL语句

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestInsert(t *testing.T) {
	// 批量创建
	n, err := Test(dbx).InsertMulti(context.Background(), []*db.Test{&db.Test{
		Name:        "x1",
		IsAllow:     false,
		Type:        0,
		CreatedTime: time.Now(),
	}, &db.Test{
		Name:        "x2",
		IsAllow:     true,
		Type:        1,
		CreatedTime: time.Now(),
	}, &db.Test{
		Name:        "x3",
		IsAllow:     true,
		Type:        0,
		CreatedTime: time.Now(),
	}})

	t.Log(n)

	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateSet(t *testing.T) {
	x := Test(dbx).U()  // Update
	result, err := x.
		Set(
			tt.Name.Set("ttt"),
			tt.IsAllow.Set(true),
		).
		SetExpr(tt.Type.Incr()).
		Where(
			tt.ID.Eq(1),
		).
		Do(context.Background())

	t.Log(x)   // 打印SQL语句
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestUpdateCols(t *testing.T) {
	x := Test(dbx).U()  // Update
	result, err := x.Cols(
		tt.Name,
		tt.IsAllow,
		tt.Type,
		tt.CreatedTime,
	).Where(
		tt.ID.Eq(1),
	).Struct(context.Background(), &db.Test{
		ID:          2,
		Name:        "ccc",
		IsAllow:     false,
		Type:        1,
		CreatedTime: time.Now(),
	})

	t.Log(x)   // 打印SQL语句

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())

}

func TestUpdate(t *testing.T) {
    // 批量更新
	ok, err := Test(dbx).UpdateMulti(context.Background(), []*db.Test{&db.Test{
		ID:          2,
		Name:        "ccc",
		IsAllow:     false,
		Type:        1,
		CreatedTime: time.Now(),
	}})

	t.Log(ok)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	x := Test(dbx).D()  // Delete
	result, err := x.Where(
		tt.ID.Eq(1),
	).Do(context.Background())
	
	t.Log(x)    // 打印SQL语句

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestSelect(t *testing.T) {
	// 查询指定列
	obj, err := Test(dbx).Find4Cols(context.Background(), 0, 0, []atype.Field{
		tt.ID,
		tt.Name,
		tt.Type,
	},
		tt.ID.Eq(2),
	)

	if err != nil {
		t.Fatal(err)
	}
	
	for _, o := range obj {
		t.Log(o)
	}
}
```