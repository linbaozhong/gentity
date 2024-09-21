package dialect

import (
	"context"
	"database/sql"
)

type (
	Modeler interface {
		TableName() string
		AssignPtr(args ...Field) []any
		AssignValues(args ...Field) ([]string, []any)
		AssignKeys() ([]Field, []any)
		AssignPrimaryKeyValues(result sql.Result) error
	}

	Daoer interface {
		// Exists 是否存在符合条件的数据
		Exists(ctx context.Context, cond ...Condition) (bool, error)
		// Sum 获取指定列的总和
		Sum(ctx context.Context, col Field, cond ...Condition) (int64, error)
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
