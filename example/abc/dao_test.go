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
	"github.com/vetcher/go-astra"
	"reflect"
	"strings"
	"testing"
	"vitess.io/vitess/go/vt/sqlparser"
)

type Dao struct {
	InterfaceName string
	Description   string
	Namespace     string
	Methods       []*Method
}
type Method struct {
	Name        string
	Description string
	Statement   string
	Args        []Arg
	Results     []Result
}
type Arg struct {
	Name string
	Type string
}
type Result struct {
	Name string
	Type string
}

func TestParseDao(t *testing.T) {
	file, e := astra.ParseFile("./internal/model/dao/gentity_model_dao_interface.go")
	if e != nil {
		t.Error(e)
		return
	}
	daos := make([]*Dao, 0)

	for _, iface := range file.Interfaces {
		dao := &Dao{
			InterfaceName: iface.Name,
			Methods:       make([]*Method, 0),
		}
		//
		for _, doc := range iface.Docs {
			_doc := strings.TrimSpace(strings.TrimLeft(doc, "/"))
			if pos := strings.Index(_doc, "@Namespace"); pos > -1 {
				dao.Namespace = _doc[pos+len("@Namespace"):]
			} else if pos = strings.Index(_doc, iface.Name); pos > -1 {
				dao.Description = _doc[pos+len(iface.Name):]
			}
		}
		// 处理方法
		for _, method := range iface.Methods {
			_method := &Method{
				Name:    method.Name,
				Args:    make([]Arg, 0),
				Results: make([]Result, 0),
			}
			_method.Name = method.Name
			// 处理注释
			for _, doc := range method.Docs {
				_doc := strings.TrimSpace(strings.TrimLeft(doc, "/"))
				if pos := strings.Index(_doc, "@Statement"); pos > -1 {
					_method.Statement = _doc[pos+len("@Statement"):]
				} else if pos = strings.Index(_doc, method.Name); pos > -1 {
					_method.Description = _doc[pos+len(method.Name):]
				}
			}
			// 处理参数
			for _, param := range method.Args {
				_method.Args = append(_method.Args, Arg{
					Name: param.Name,
					Type: param.Type.String(),
				})
			}
			// 处理返回值
			for _, result := range method.Results {
				_method.Results = append(_method.Results, Result{
					Name: result.Name,
					Type: result.Type.String(),
				})
			}

			dao.Methods = append(dao.Methods, _method)
		}
		daos = append(daos, dao)
	}
	for _, dao := range daos {
		t.Log("接口：", dao.InterfaceName)
		fmt.Println("---- 接口描述：", dao.Description)
		fmt.Println("---- 接口命名空间：", dao.Namespace)
		for _, method := range dao.Methods {
			t.Log("方法：", method.Name)
			fmt.Println("---- 方法描述：", method.Description)
			fmt.Println("---- 查询语句：", method.Statement)
			fmt.Println("---- 方法参数：")
			for _, arg := range method.Args {
				fmt.Println("---- ---- ", arg.Name, arg.Type)
			}
			fmt.Println("---- 方法返回值：")
			for _, result := range method.Results {
				fmt.Println("---- ---- ", result.Name, result.Type)
			}
			fmt.Println("---- SQL语句解析开始：")
			parseSQLInfo(method.Statement)
			fmt.Println("---- SQL语句解析结束：")
			fmt.Println()
			fmt.Println()
		}
	}
}

// ///////////////////////
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
		parseSelect(s)

	case *sqlparser.Insert:
		parseInsert(s)
	case *sqlparser.Update:
		parseUpdate(s)

	case *sqlparser.Delete:
		parseDelete(s)
	default:
		fmt.Println("未知类型的 SQL 语句")
	}
}

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

func parseUpdate(s *sqlparser.Update) {
	fmt.Println("语句类型: UPDATE")
	fmt.Print("表名: ")
	getTables(s.TableExprs)

	fmt.Println()

	// 解析 SET 子句
	fmt.Print("更新内容: ")
	for _, updateExpr := range s.Exprs {
		fmt.Printf("%s = %v ", updateExpr.Name.Name, getExprInfo(updateExpr.Expr))
	}
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

func parseInsert(s *sqlparser.Insert) {
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
	// 解析插入的值
	fmt.Print("插入值: ")
	switch rows := s.Rows.(type) {
	case sqlparser.Values:
		for _, tuple := range rows {
			for _, val := range tuple {
				fmt.Print(getExprInfo(val), " ")
			}
		}
	case *sqlparser.Select:
		fmt.Print("子查询插入: ", sqlparser.String(rows))
	}
	fmt.Println()
}
func getTableName(tables []Table, alias string) string {
	if len(tables) == 0 {
		return ""
	}
	if len(tables) == 1 {
		return tables[0].Name
	}
	for _, table := range tables {
		if table.Name == alias || table.Alias == alias {
			return table.Name
		}
	}
	return tables[0].Name
}

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
			if cond, ok := node.(*sqlparser.ComparisonExpr); ok {
				val := getExprInfo(cond.Right)
				fmt.Printf("%v %v %v ", getExprInfo(cond.Left), cond.Operator.ToString(), val)
				cond.Right = getplaceholder(cond.Operator, val)
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

func getplaceholder(op sqlparser.ComparisonExprOperator, val string) *sqlparser.Literal {
	switch op {
	case sqlparser.InOp:
		return sqlparser.NewBitLiteral("(%s)")
	default:
		return sqlparser.NewBitLiteral("?")
	}
}

// getTables 从 TableExpr 中提取表名
func getTables(ts []sqlparser.TableExpr) []Table {
	tables := make([]Table, 0, len(ts))
	for _, tableExpr := range ts {
		table := Table{}
		switch expr := tableExpr.(type) {
		case *sqlparser.JoinTableExpr:
			sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
				if aliasTableName, ok := node.(*sqlparser.AliasedTableExpr); ok {
					if tableName, ok := aliasTableName.Expr.(sqlparser.TableName); ok {
						table.Name = tableName.Name.String()
					}
					table.Alias = aliasTableName.TableNameString()
				}
				return true, nil
			}, expr.LeftExpr)
			sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
				if aliasTableName, ok := node.(*sqlparser.AliasedTableExpr); ok {
					if tableName, ok := aliasTableName.Expr.(sqlparser.TableName); ok {
						table.Name = tableName.Name.String()
					}
					table.Alias = aliasTableName.TableNameString()
				}
				return true, nil
			}, expr.RightExpr)
			fmt.Println()
			fmt.Println(expr.Join.ToString(), getExprInfo(expr.Condition.On))
		default:
			sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
				if aliasTableName, ok := node.(*sqlparser.AliasedTableExpr); ok {
					if tableName, ok := aliasTableName.Expr.(sqlparser.TableName); ok {
						table.Name = tableName.Name.String()
					}
					table.Alias = aliasTableName.TableNameString()
				}
				return true, nil
			}, tableExpr)
		}
		tables = append(tables, table)
	}
	return tables
}

// getFieldName 从 sqlparser.Expr 中解析字段名
func getExprInfo(expr sqlparser.Expr) string {
	switch node := expr.(type) {
	case *sqlparser.ColName:
		// 如果是列名，直接返回列名
		tableName := node.Qualifier.Name.String()
		if tableName == "" {
			return node.Name.String()
		}
		return tableName + "." + node.Name.String()
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
	case *sqlparser.Argument:
		// 如果是参数，返回参数名称
		return node.Name
	case *sqlparser.FuncExpr:
		// 如果是函数表达式，返回函数名称和参数信息
		return fmt.Sprintf("函数: %s(%s)", node.Name, node.Exprs)
	case sqlparser.ValTuple:
		for _, val := range node {
			return getExprInfo(val)
		}
	}
	fmt.Println(reflect.TypeOf(expr))
	return "未知"
}
