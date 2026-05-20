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

package dialect

import "errors"

var (
	Err_Expression_Empty_Param = errors.New("Expression parameter must have one value")
	Err_Condition_Empty_Param  = errors.New("Condition parameter must have one value")
)

type (
	SetOp int8 // 赋值运算符
	// Field 字段
	Field struct {
		Name      string
		Json      string
		OmitEmpty bool
		Table     string
		Type      string
	}
	// Function 聚合函数
	Function func(Dialect) string
	// Condition 条件
	Condition struct {
		Op        LogicalOperator
		Condition CondFunc
		Children  []Condition
	}
	CondFunc func(*uint8, Dialect) (string, any)
	// Order 排序
	Order func() (OrderType, []Field)
	// Setter 赋值
	Setter func() (Field, any, SetOp)
)

const (
	Op_Normal    SetOp = iota // insert 赋值
	Op_Increment              // update 自增
	Op_Decrement              // update 自减
	Op_Replace                // update 替换
	Op_Expr                   // update 其它表达式
)
