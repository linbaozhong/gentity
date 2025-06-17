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
	"testing"
	"vitess.io/vitess/go/vt/sqlparser"
)

func TestParse(t *testing.T) {
	sql := "SELECT id, name FROM users WHERE age > 18"
	parser, e := sqlparser.New(sqlparser.Options{})
	if e != nil {
		fmt.Println("解析失败:", e)
		return
	}
	stmt, err := parser.Parse(sql)
	if err != nil {
		fmt.Println("解析失败:", err)
		return
	}
	tableNames := getTableNames(stmt)
	if len(tableNames) == 0 {
		t.Error("未解析出表名")
	} else {
		for _, tableName := range tableNames {
			fmt.Println("解析出的表名:", tableName)
		}
	}
}

// getTableNames 从 SQL 语句中提取表名
func getTableNames(stmt sqlparser.Statement) []string {
	var tableNames []string
	switch s := stmt.(type) {
	case *sqlparser.Select:
		// 处理 SELECT 语句
		for _, tableExpr := range s.From {
			getTableNamesFromTableExpr(tableExpr, &tableNames)
		}
	case *sqlparser.Insert:
		// 处理 INSERT 语句
		getTableNamesFromTableExpr(s.Table, &tableNames)
	case *sqlparser.Update:
		// 处理 UPDATE 语句
		for _, tableExpr := range s.TableExprs {
			getTableNamesFromTableExpr(tableExpr, &tableNames)
		}
	case *sqlparser.Delete:
		// 处理 DELETE 语句
		for _, tableExpr := range s.TableExprs {
			getTableNamesFromTableExpr(tableExpr, &tableNames)
		}
	}
	return tableNames
}

// getTableNamesFromTableExpr 从 TableExpr 中提取表名
func getTableNamesFromTableExpr(tableExpr sqlparser.TableExpr, tableNames *[]string) {
	sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		if tableName, ok := node.(*sqlparser.TableName); ok {
			*tableNames = append(*tableNames, tableName.Name.String())
		}
		if tableName, ok := node.(*sqlparser.AliasedTableExpr); ok {
			*tableNames = append(*tableNames, tableName.TableNameString())
		}
		return true, nil
	}, tableExpr)
}

// /////////////////

// parseSQLInfo 解析 SQL 语句的全部信息
func parseSQLInfo(sql string) {
	parser, e := sqlparser.New(sqlparser.Options{})
	if e != nil {
		fmt.Println("解析失败:", e)
		return
	}

	stmt, err := parser.Parse(sql)
	if err != nil {
		fmt.Println("解析失败:", err)
		return
	}

	switch s := stmt.(type) {
	case *sqlparser.Select:
		fmt.Println("语句类型: SELECT")
		// 解析 FROM 子句中的表名
		fmt.Print("表名: ")
		for _, tableExpr := range s.From {
			switch expr := tableExpr.(type) {
			case *sqlparser.JoinTableExpr:
				sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
					if tableName, ok := node.(*sqlparser.AliasedTableExpr); ok {
						fmt.Print(tableName.TableNameString(), " ")
					}
					return true, nil
				}, expr.LeftExpr)
				sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
					if tableName, ok := node.(*sqlparser.AliasedTableExpr); ok {
						fmt.Print(tableName.TableNameString(), " ")
					}
					return true, nil
				}, expr.RightExpr)
				fmt.Println()
				fmt.Println(expr.Join.ToString(), getExprInfo(expr.Condition.On))
			default:
				sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
					if tableName, ok := node.(*sqlparser.AliasedTableExpr); ok {
						fmt.Print(tableName.TableNameString(), " ")
					}
					return true, nil
				}, tableExpr)
			}
		}
		fmt.Println()

		// 解析 SELECT 子句中的列名
		fmt.Print("列名: ")
		for _, selectExpr := range s.SelectExprs.Exprs {
			if colName, ok := selectExpr.(*sqlparser.AliasedExpr); ok {
				if col, ok := colName.Expr.(*sqlparser.ColName); ok {
					fmt.Print(col.Name.String(), " ")
				}
			}
		}
		fmt.Println()

		// 解析 WHERE 子句中的条件
		if s.Where != nil {
			fmt.Print("条件: ")
			sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
				if cond, ok := node.(*sqlparser.ComparisonExpr); ok {
					fmt.Printf("%v %v %v ", getExprInfo(cond.Left), cond.Operator.ToString(), getExprInfo(cond.Right))
				}
				return true, nil
			}, s.Where.Expr)
			fmt.Println()
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

	case *sqlparser.Insert:
		fmt.Println("语句类型: INSERT")
		fmt.Print("表名: ")
		sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
			if tableName, ok := node.(*sqlparser.AliasedTableExpr); ok {
				fmt.Print(tableName.TableNameString())
			}
			return true, nil
		}, s.Table)
		fmt.Println()

		// 解析插入的列名
		fmt.Print("插入列名: ")
		for _, col := range s.Columns {
			fmt.Print(col.String(), " ")
		}
		fmt.Println()

	case *sqlparser.Update:
		fmt.Println("语句类型: UPDATE")
		fmt.Print("表名: ")
		for _, tableExpr := range s.TableExprs {
			sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
				if tableName, ok := node.(*sqlparser.AliasedTableExpr); ok {
					fmt.Print(tableName.TableNameString(), " ")
				}
				return true, nil
			}, tableExpr)
		}
		fmt.Println()

		// 解析 SET 子句
		fmt.Print("更新内容: ")
		for _, updateExpr := range s.Exprs {
			fmt.Printf("%s = %v ", updateExpr.Name.Name, getExprInfo(updateExpr.Expr))
		}
		fmt.Println()

	case *sqlparser.Delete:
		fmt.Println("语句类型: DELETE")
		fmt.Print("表名: ")
		for _, tableExpr := range s.TableExprs {
			sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
				if tableName, ok := node.(*sqlparser.AliasedTableExpr); ok {
					fmt.Print(tableName.TableNameString(), " ")
				}
				return true, nil
			}, tableExpr)
		}
		fmt.Println()

		// 解析 WHERE 子句中的条件
		if s.Where != nil {
			fmt.Print("条件: ")
			sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
				if cond, ok := node.(*sqlparser.ComparisonExpr); ok {
					fmt.Printf("%s %v %v ", getExprInfo(cond.Left), cond.Operator.ToString(), getExprInfo(cond.Right))
				}
				return true, nil
			}, s.Where.Expr)
			fmt.Println()
		}
	default:
		fmt.Println("未知类型的 SQL 语句")
	}
}

// getFieldName 从 sqlparser.Expr 中解析字段名
func getExprInfo(expr sqlparser.Expr) string {
	switch node := expr.(type) {
	case *sqlparser.ColName:
		// 如果是列名，直接返回列名
		return node.Name.String()
	case *sqlparser.Literal:
		// 如果是 SQL 固定值，返回值
		return node.Val
	case *sqlparser.Subquery:
		// 如果是子查询，返回子查询信息
		return fmt.Sprintf("子查询: %s", sqlparser.String(node.Select))
	case *sqlparser.BinaryExpr:
		// 如果是二元表达式，返回操作符和左右表达式信息
		return fmt.Sprintf("二元表达式: %s %s %s", getExprInfo(node.Left), node.Operator.ToString(), getExprInfo(node.Right))
	case *sqlparser.ComparisonExpr:
		// 如果是比较表达式，返回左右表达式和操作符信息
		return fmt.Sprintf("比较表达式: %s %s %s", getExprInfo(node.Left), node.Operator.ToString(), getExprInfo(node.Right))

	}

	return "未知"
}

func TestParseAllInfo(t *testing.T) {
	sqls := []string{
		"SELECT orders.id, orders.name,count(orders.id) as num FROM users left join orders on orders.user = users.id WHERE age > 18 ORDER BY id DESC",
		"INSERT INTO users (id, name) VALUES (1, 'Alice')",
		"UPDATE users SET name = 'Bob' WHERE id = 1",
		"DELETE FROM users WHERE age < 18",
	}

	for _, sql := range sqls {
		fmt.Printf("SQL 语句: %s\n", sql)
		parseSQLInfo(sql)
		fmt.Println()
	}
}
