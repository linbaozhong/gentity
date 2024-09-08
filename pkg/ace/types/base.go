// Copyright © 2023 Linbaozhong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"context"
	"fmt"
)

const (
	Inner_Join JoinType = " INNER"
	Left_Join  JoinType = " LEFT"
	Right_Join JoinType = " RIGHT"

	Operator_and = " AND "
	Operator_or  = " OR "

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
	//// NullBool is an alias to sql.NullBool.
	//NullBool = sql.NullBool
	//// NullInt64 is an alias to sql.NullInt64.
	//NullInt64 = sql.NullInt64
	//// NullInt32 is an alias to sql.NullInt32.
	//NullInt32 = sql.NullInt32
	//// NullInt16 is an alias to sql.NullInt16.
	//NullInt16 = sql.NullInt16
	//// NullByte is an alias to sql.NullByte.
	//NullByte = sql.NullByte
	//// NullString is an alias to sql.NullString.
	//NullString = sql.NullString
	//// NullFloat64 is an alias to sql.NullFloat64.
	//NullFloat64 = sql.NullFloat64
	//// NullTime represents a time.Time that may be null.
	//NullTime = sql.NullTime
	//// UnknownType is a named type to any indicates the info
	//// needs to be extracted from the underlying rows.
	//UnknownType any

	Modeler interface {
		TableName() string
		AssignValues(args ...Field) ([]string, []any)
		AssignKeys() ([]Field, []any)
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
