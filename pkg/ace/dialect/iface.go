package dialect

import (
	"context"
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

	Modeler interface {
		TableName() string
		AssignPtr(args ...Field) []any
		AssignValues(args ...Field) ([]string, []any)
		AssignKeys() (Field, any)
		AssignPrimaryKeyValues(result sql.Result) error
	}

	Daoer interface {
		// Exists 是否存在符合条件的数据
		Exists(ctx context.Context, cond ...Condition) (bool, error)
		// Sum 获取指定列的总和
		Sum(ctx context.Context, col []Field, cond ...Condition) (map[string]any, error)
		// Count 获取符合条件的数据总数
		Count(ctx context.Context, cond ...Condition) (int64, error)
		// Delete 删除符合条件的数据
		Delete(ctx context.Context, cond ...Condition) (bool, error)
		// Update 更新符合条件的数据
		Update(ctx context.Context, sets []Setter, cond ...Condition) (bool, error)
		// Insert 插入数据
		Insert(ctx context.Context, sets ...Setter) (int64, error)
	}
)
