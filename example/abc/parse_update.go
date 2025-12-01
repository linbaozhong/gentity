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

package abc

import (
	"fmt"
	"vitess.io/vitess/go/vt/sqlparser"
)

func parseUpdate(s *sqlparser.Update) {
	stmt := &Statement{
		Type: UpdateType,
	}
	fmt.Println("语句类型: UPDATE")
	fmt.Print("表名: ")
	stmt.Table = getTables(s.TableExprs)
	fmt.Println(stmt.Table)

	// 解析 SET 子句
	fmt.Print("更新内容: ")
	for _, updateExpr := range s.Exprs {
		fmt.Printf("%s = %v ", updateExpr.Name.Name, getExprInfo(updateExpr.Expr))
		stmt.PlaceHolders = append(stmt.PlaceHolders, PlaceHolder{
			Name: getExprInfo(updateExpr.Expr),
		})
	}
	fmt.Println()
	// 解析 WHERE 子句中的条件
	if s.Where != nil {
		fmt.Println("条件: ")
		sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
			placeholder := PlaceHolder{}
			if cond, ok := node.(*sqlparser.ComparisonExpr); ok {
				val := getExprInfo(cond.Right)
				// fmt.Printf("%s %v %v ", getExprInfo(cond.Left), cond.Operator.ToString(), val)
				// cond.Right = sqlparser.NewIntLiteral("?")
				placeholder.Name = val
				placeholder.Operator = OperatorType(cond.Operator)
				stmt.PlaceHolders = append(stmt.PlaceHolders, placeholder)
			}
			return true, nil
		}, s.Where.Expr)
	}
	fmt.Println(stmt.PlaceHolders)
}
