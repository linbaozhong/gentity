package dialect

import (
	"context"
)

type (
	Modeler[T BaseType] interface {
		TableName() string
		AssignValues(args ...Field[T]) ([]string, []any)
		AssignKeys() ([]Field[T], []any)
	}

	Daoer[T BaseType] interface {
		// Exists 是否存在符合条件的数据
		Exists(ctx context.Context, cond ...Condition[T]) (bool, error)
		// Sum 获取指定列的总和
		Sum(ctx context.Context, col Field[T], cond ...Condition[T]) (int64, error)
		// Count 获取符合条件的数据总数
		Count(ctx context.Context, cond ...Condition[T]) (int64, error)
		// Delete 删除符合条件的数据
		Delete(ctx context.Context, cond ...Condition[T]) (bool, error)
		// Update 更新符合条件的数据
		Update(ctx context.Context, sets []Setter[T], cond ...Condition[T]) (bool, error)
		// Insert 插入数据
		Insert(ctx context.Context, sets ...Setter[T]) (int64, error)
	}
)
