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

func parseSelect(s *sqlparser.Select) *Statement {
	stmt := &Statement{
		Type: SelectType,
	}
	fmt.Println("语句类型: SELECT")
	// 解析 FROM 子句中的表名
	stmt.Table = getTables(s.From)
	fmt.Println("表名: ", stmt.Table)

	// 解析 SELECT 子句中的列名
	fmt.Print("列名: ")
	for _, selectExpr := range s.SelectExprs.Exprs {
		column := Column{}
		switch colName := selectExpr.(type) {
		case *sqlparser.StarExpr:
			fmt.Print("*")
		case *sqlparser.AliasedExpr:
			switch col := colName.Expr.(type) {
			case *sqlparser.ColName:
				if col.Qualifier.Name.String() == "" {
					column.Table = stmt.Table[0].Name
				} else {
					column.Table = getTableName(stmt.Table, col.Qualifier.Name.String())
				}
				column.Name = col.Name.String()
			case *sqlparser.Count:
				fmt.Print(getExprInfo(col.Args[0]), " ")
			}
		}
		stmt.Columns = append(stmt.Columns, column)
	}
	fmt.Println(stmt.Columns)

	// 解析 WHERE 子句中的条件
	if s.Where != nil {
		fmt.Print("条件: ")
		sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
			placeholder := PlaceHolder{}
			if cond, ok := node.(*sqlparser.ComparisonExpr); ok {
				val := getExprInfo(cond.Right)
				fmt.Printf("%v %v %v ", getExprInfo(cond.Left), cond.Operator.ToString(), val)
				// cond.Right = getplaceholder(cond.Operator, val)
				placeholder.Name = val
				placeholder.Operator = OperatorType(cond.Operator)
				stmt.PlaceHolders = append(stmt.PlaceHolders, placeholder)
			}
			return true, nil
		}, s.Where.Expr)
		fmt.Println(stmt.PlaceHolders)
	}

	// 解析 ORDER BY 子句
	if s.OrderBy != nil {
		fmt.Print("排序规则: ")
		for _, order := range s.OrderBy {
			if col, ok := order.Expr.(*sqlparser.ColName); ok {
				fmt.Printf("%s %v ", col.Name.String(), order.Direction.ToString())
			}
		}
		fmt.Println()
	}
	// 解析 GROUP BY 子句
	if s.GroupBy != nil {
		fmt.Print("分组规则: ")
		for _, group := range s.GroupBy.Exprs {
			fmt.Printf("%v ", getExprInfo(group))
		}
		fmt.Println()
		if s.Having != nil {
			fmt.Print("分组条件: ")
			sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
				if cond, ok := node.(*sqlparser.ComparisonExpr); ok {
					fmt.Printf("%v %v %v ", getExprInfo(cond.Left), cond.Operator.ToString(), getExprInfo(cond.Right))
					cond.Right = getplaceholder(cond.Operator, getExprInfo(cond.Right))
				}
				return true, nil
			}, s.Having.Expr)
		}
		fmt.Println()
	}
	// 解析 LIMIT 子句
	if s.Limit != nil {
		fmt.Printf("分页信息: Offset %v , Size %v \n", getExprInfo(s.Limit.Offset), getExprInfo(s.Limit.Rowcount))
		s.Limit.Offset = sqlparser.NewIntLiteral("?")
		s.Limit.Rowcount = sqlparser.NewIntLiteral("?")
		fmt.Println()
	}
	fmt.Println(sqlparser.String(s))
	return stmt
}
