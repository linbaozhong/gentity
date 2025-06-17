package define

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/log"
	"model/define/dao"
	"model/define/table/tblcompany"
	"model/do"
	"testing"
)

var (
	dbx *ace.DB
)

func init() {
	var err error
	dbx, err = ace.Connect("mysql",
		"ssld_dev:Cu83&sr66@tcp(123.56.5.53:13306)/dispatch?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	dbx.SetMaxOpenConns(50)
	dbx.SetMaxIdleConns(25)
	// dbx.SetDebug(true)
}

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

	// 使用 daocompany 构建器插入数据
	// 调用 daocompany 包的 Builder 函数创建一个构建器实例，设置插入的数据
	// 然后调用 Create 方法创建执行器并执行插入操作
	// result 为执行结果，包含插入的 ID 等信息，err 为可能出现的错误
	result, err := ace.Table(do.CompanyTableName).Set(
		// 设置公司的长名称为 "aaaaaa"
		tblcompany.LongName.Set("aaaaaa"),
		// 设置公司的状态为 1
		tblcompany.State.Set(1),
	).
		Create(dbx).
		Exec(context.Background())
	// 若插入操作出现错误，则终止测试并输出错误信息
	if err != nil {
		t.Fatal(err)
	}
	// 记录插入操作返回的最后插入的 ID
	t.Log(result.LastInsertId())

	// 使用 ace 包的通用方法插入数据
	// 调用 ace 包的 Table 函数指定要操作的表，设置插入的数据
	// 然后调用 Create 方法创建执行器并执行插入操作
	// result 为执行结果，包含插入的 ID 等信息，err 为可能出现的错误
	result, err = ace.
		Table(do.CompanyTableName).
		Set(
			// 设置公司的长名称为 "aaaaaa"
			tblcompany.LongName.Set("aaaaaa"),
			// 设置公司的状态为 1
			tblcompany.State.Set(1),
			// 设置公司的地址为 "beijing"
			tblcompany.Address.Set("beijing"),
		).
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
			[]*do.Company{
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

	// 使用 DAO 构建器，通过单个结构体插入数据
	// result 为执行结果，包含插入的 ID 等信息，err 为可能出现的错误
	result, err := ace.Table(do.CompanyTableName).
		Cols(
			tblcompany.LongName,
			tblcompany.Address,
			tblcompany.Email,
			tblcompany.State,
		).
		Create(dbx).
		Struct(context.Background(),
			&do.Company{
				LongName: "m1",
				Address:  "北京",
				Email:    "1@2.com",
				State:    1,
			})
	// 若插入操作出现错误，则终止测试并输出错误信息
	if err != nil {
		t.Fatal(err)
	}
	// 记录插入操作返回的最后插入的 ID
	t.Log(result.LastInsertId())

	// 使用 ace 包的通用方法，通过单个结构体插入数据
	// result 为执行结果，包含插入的 ID 等信息，err 为可能出现的错误
	result, err = ace.
		Table(do.CompanyTableName).
		Cols(
			tblcompany.LongName,
			tblcompany.Address,
			tblcompany.Email,
		).
		Create(dbx).
		Struct(context.Background(),
			&do.Company{
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

	// 使用 daocompany 构建器更新数据
	// 调用 daocompany 包的 Builder 函数创建一个构建器实例，设置要更新的数据
	// 并通过 Where 方法指定更新条件，然后调用 Update 方法创建执行器并执行更新操作
	// result 为执行结果，err 为可能出现的错误
	sss := ace.Table(do.CompanyTableName).
		Set(
			tblcompany.LongName.Set("aaaaaa"),
			// 设置公司的状态为 1
			tblcompany.State.Set(1),
		).
		Where(tblcompany.Id.Eq(1)).Clone()

	result, err := sss.ToSql().Update(dbx).
		Exec(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(result.LastInsertId())

	// 使用生成的 DAO 方法更新数据
	// 调用 daocompany 包的 New 函数创建一个新的 DAO 实例，并调用其 Update 方法更新数据
	// n 为更新操作影响的行数，err 为可能出现的错误
	// 第一个参数为 context 上下文，第二个参数为要更新的字段设置，第三个参数为更新条件
	n, err := dao.Company(dbx).ToSql().
		Update(context.Background(),
			ace.Sets(
				tblcompany.LongName.Set("bbbbbbb"),
				// 设置公司的状态为 1
				tblcompany.State.Set(2),
			).ToSlice(),
			tblcompany.Id.Eq(2),
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
		Table(do.CompanyTableName).
		Set(
			tblcompany.LongName.Set("aaaaaa"),
			// 设置公司的状态为 1
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
	obj, has, err := dao.Company(dbx).ToSql().
		Get(context.Background(),
			ace.Cols(tblcompany.Id, tblcompany.LongName).
				Where(tblcompany.Id.In(nil)),
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
	return

	// 创建一个新的 Company 实例
	obj = do.NewCompany()

	// 使用 DAO 构建器查询数据
	// 调用 daocompany 包的 Builder 函数创建一个构建器实例，设置查询条件为 Id 等于 1
	// 然后调用 Select 方法指定数据库连接，再调用 Get 方法将查询结果填充到 obj 中
	// err 为可能出现的错误
	err = ace.Table(do.CompanyTableName).ToSql().
		Cols(tblcompany.Id, tblcompany.State, tblcompany.Address).
		Where(tblcompany.Id.Eq(2)).
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

	// 使用 ace 包的通用方法查询数据
	// 调用 ace 包的 Table 函数指定要操作的表，设置查询条件为 Id 等于 1
	// 然后调用 Select 方法指定数据库连接，再调用 Get 方法将查询结果填充到 obj 中
	// err 为可能出现的错误
	err = ace.Table(do.CompanyTableName).ToSql().
		Where(tblcompany.Id.Eq(3)).
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

func TestPool(t *testing.T) {
	a1 := do.NewAccount()
	a1.Id = 1
	a1.LoginName = "a1"
	a1.Free()
	a1.Free()

	a2 := do.NewAccount()
	a3 := do.NewAccount()
	a2.Id = 2
	a2.LoginName = "a2"
	t.Logf("a3: %+v", a3)
	t.Logf("a2: %+v", a2)
}

func TestCopy(t *testing.T) {
	var ii []int
	copy(ii, []int{1, 2, 3})
	t.Log(len(ii))
}
