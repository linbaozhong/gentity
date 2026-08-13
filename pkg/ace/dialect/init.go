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
		Name       string // 字段名
		Json       string // json字段名
		OmitEmpty  bool   // 是否忽略空值
		Table      string // 表名
		Type       string // 字段类型
		IsRelation bool   // 是否关联字段
		as         string // 别名
	}
	// Function 聚合函数
	Function func(Dialect) string
	// Condition 条件
	Condition struct {
		Op        LogicalOperator
		Condition CondFunc
		Children  []Condition
	}
	CondFunc func(*uint16, Dialect) (string, any)
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

// Or 将当前条件与后续条件以 OR 组合，返回一个新的组合条件。
// 可作为单个 Condition 传入 Exists/Count/Where 等接收 ...Condition 的方法，
// 实现「外层 AND + 内层 OR」并存。
//
// 例:
//
//	tblstaff.Phone.Eq(p).Or(tblstaff.Email.Eq(p))
//	=> (Phone = ? OR Email = ?)
func (c Condition) Or(conds ...Condition) Condition {
	all := append([]Condition{c}, conds...)
	for i := range all {
		all[i].Op = Operator_or
	}
	return Condition{
		Op:       Operator_or,
		Children: all,
	}
}

// And 将当前条件与后续条件以 AND 组合，返回一个新的组合条件。
// 常用于在 OR 子组中再嵌套 AND。
//
// 例:
//
//	tblstaff.A.Eq(1).Or(tblstaff.B.Eq(2).And(tblstaff.C.Eq(3)))
//	=> (A = ? OR (B = ? AND C = ?))
func (c Condition) And(conds ...Condition) Condition {
	all := append([]Condition{c}, conds...)
	for i := range all {
		all[i].Op = Operator_and
	}
	return Condition{
		Op:       Operator_and,
		Children: all,
	}
}

// Or 将多个条件组合为 OR 关系的一组条件。
// 生成形如 (a = ? OR b = ?) 的 SQL，可作为单个 Condition 传入
// Exists / Count / Where 等接收 ...Condition 的方法，实现「外层 AND + 内层 OR」并存。
//
// 例:
//
//	dialect.Or(tbl.X.Eq(1), tbl.Y.Eq(2))
//	dao.Staff(db).Exists(ctx,
//	    tblstaff.Status.Eq(1),
//	    dialect.Or(tblstaff.Phone.Eq(p), tblstaff.Email.Eq(p)),
//	)
//	// => Status = ? AND (Phone = ? OR Email = ?)
func Or(conds ...Condition) Condition {
	for i := range conds {
		conds[i].Op = Operator_or
	}
	return Condition{
		Op:       Operator_or,
		Children: conds,
	}
}

// And 将多个条件组合为 AND 关系的一组条件（显式分组），
// 常用于在 OR 子组中再嵌套 AND。
// 例: dialect.Or(tbl.A.Eq(1), dialect.And(tbl.B.Eq(2), tbl.C.Eq(3)))
func And(conds ...Condition) Condition {
	for i := range conds {
		conds[i].Op = Operator_and
	}
	return Condition{
		Op:       Operator_and,
		Children: conds,
	}
}
