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

func TestParseAllInfo(t *testing.T) {
	sqls := []string{
		"SELECT u.id, o.name,count(o.id) as num FROM users as u left join orders as o on o.user = u.id WHERE u.age > :act.Age group by u.id having u.id > ? ORDER BY u.id DESC LIMIT ?, ?",
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
