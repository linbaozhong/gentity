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

func parseDelete(s *sqlparser.Delete) {
	fmt.Println("语句类型: DELETE")
	fmt.Print("表名: ")
	getTables(s.TableExprs)
	fmt.Println()

	// 解析 WHERE 子句中的条件
	if s.Where != nil {
		fmt.Print("条件: ")
		sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
			if cond, ok := node.(*sqlparser.ComparisonExpr); ok {
				fmt.Printf("%s %v %v ", getExprInfo(cond.Left), cond.Operator.ToString(), getExprInfo(cond.Right))
				cond.Right = sqlparser.NewIntLiteral("?")
			}
			return true, nil
		}, s.Where.Expr)
	}
	fmt.Println()
}
