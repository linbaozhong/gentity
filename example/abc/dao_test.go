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

var daos = make([]*Dao, 0)

func TestParseDao(t *testing.T) {
	file, e := astra.ParseFile("./internal/model/do/gentity_model_dao_interface.go")
	if e != nil {
		t.Error(e)
		return
	}

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
		fmt.Println("// ", dao.InterfaceName, dao.Description)
		for _, method := range dao.Methods {
			fmt.Println("// ", method.Name, method.Description)
			fmt.Print("func (*", dao.Namespace, ")")
			fmt.Print(method.Name, "(")
			// fmt.Println("---- 查询语句：", method.Statement)

			for i, arg := range method.Args {
				if i > 0 {
					fmt.Print(",")
				}
				fmt.Printf("%s %s", arg.Name, arg.Type)
			}
			fmt.Print(") (")
			for i, result := range method.Results {
				if i > 0 {
					fmt.Print(",")
				}
				fmt.Printf("%s %s", result.Name, result.Type)
			}
			fmt.Print(") {\n")
			fmt.Println("---- SQL语句解析开始：")
			parseSQLInfo(method.Statement)
			fmt.Println("---- SQL语句解析结束：")
			fmt.Println("}")
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
