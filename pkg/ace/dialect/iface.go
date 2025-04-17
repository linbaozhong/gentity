package dialect

import (
	"database/sql"
	"fmt"
)

const (
	Inner_Join JoinType = " INNER"
	Left_Join  JoinType = " LEFT"
	Right_Join JoinType = " RIGHT"

	Operator_and  = " AND "
	Operator_or   = " OR "
	Operator_Asc  = " ASC"
	Operator_Desc = " DESC"

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
	JoinType string

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
