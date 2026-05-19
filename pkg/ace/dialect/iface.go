package dialect

import (
	"database/sql"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/sqlparser"
)

const (
	Inner_Join JoinType = " INNER"
	Left_Join  JoinType = " LEFT"
	Right_Join JoinType = " RIGHT"

	Operator_and LogicalOperator = " AND "
	Operator_or  LogicalOperator = " OR "
	Operator_not LogicalOperator = " NOT "

	Operator_Asc  OrderType = " ASC"
	Operator_Desc OrderType = " DESC"

	MaxLimit uint = 1000
	PageSize uint = 20
)

var (
	ErrCreateEmpty        = fmt.Errorf("No data is created")
	ErrBeanEmpty          = fmt.Errorf("bean=nil 或者 len(beans)=0 或者 len(beans)>100")
	ErrNotFound           = fmt.Errorf("not found")
	ErrSetterEmpty        = fmt.Errorf("setter=nil 或者 len(setter)=0")
	ErrBeansEmpty         = fmt.Errorf("beans=nil 或者 len(beans)=0")
	ErrArgsNotMatch       = fmt.Errorf("args not match")
	ErrPrimaryKeyNotMatch = fmt.Errorf("primary key not match")
)

type (
	JoinType        string
	OrderType       string
	LogicalOperator string

	TableNamer interface {
		TableName() string
	}

	Modeler interface {
		TableNamer
		AssignPtr(args ...Field) []any
		// AssignValues 向数据库写入数据前，为表列赋值。
		// 如果 args 为空，则将非零值赋与可写字段
		// 如果 args 不为空，则只赋值 args 中的字段
		AssignValues(args ...Field) ([]string, []any)
		// RawAssignValues 向数据库写入数据前，为表列赋值。多用于批量插入和更新
		// 如果 args 为空，则赋值所有可写字段
		// 如果 args 不为空，则只赋值 args 中的字段
		RawAssignValues(args ...Field) ([]string, []any)
		AssignKeys() (Field, any)
		AssignPrimaryKeyValues(result sql.Result) error
	}

	// Daoer interface {
	// 	// Exists 是否存在符合条件的数据
	// 	Exists(ctx context.Context, cond ...Condition) (bool, error)
	// 	// Sum 获取指定列的总和
	// 	Sum(ctx context.Context, col []Field, cond ...Condition) (map[string]any, error)
	// 	// Count 获取符合条件的数据总数
	// 	Count(ctx context.Context, cond ...Condition) (int64, error)
	// 	// Delete 删除符合条件的数据
	// 	Delete(ctx context.Context, cond ...Condition) (bool, error)
	// 	// Update 更新符合条件的数据
	// 	Update(ctx context.Context, sets []Setter, cond ...Condition) (bool, error)
	// 	// Insert 插入数据
	// 	Insert(ctx context.Context, sets ...Setter) (int64, error)
	// }
)

type Dialect interface {
	// Name 返回数据库名称
	Name() string

	// Quote 引用标识符（表名、列名）
	// MySQL: `name`  PostgreSQL: "name"  SQL Server: [name]
	Quote(name string) string

	// Placeholder 返回参数占位符
	// MySQL: ?           PostgreSQL: $1, $2...
	// SQL Server: @p1    Oracle: :1
	Placeholder(index *uint8) string

	// Limit 生成分页语句
	// MySQL: LIMIT offset,limit      PostgreSQL: LIMIT limit OFFSET offset
	// SQL Server: OFFSET offset ROWS FETCH NEXT limit ROWS ONLY
	Limit(offset, limit uint) string

	// AutoIncrement 返回自增关键字
	AutoIncrement() string

	// PrimaryKey 返回主键标识
	PrimaryKey() string

	// UniqueKey 返回唯一键标识
	UniqueKey() string
	GetTables(db *sql.DB, dbName string) ([]*sqlparser.Table, error)
}

func (l LogicalOperator) String() string {
	return string(l)
}

func (o OrderType) String() string {
	return string(o)
}
