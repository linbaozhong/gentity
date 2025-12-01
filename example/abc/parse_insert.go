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

func parseInsert(s *sqlparser.Insert) {
	stmt := &Statement{
		Type: InsertType,
	}
	fmt.Println("语句类型: INSERT")
	fmt.Print("表名: ")
	sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		if tableName, ok := node.(*sqlparser.AliasedTableExpr); ok {
			stmt.Table = append(stmt.Table, Table{
				Name: tableName.TableNameString(),
			})
		}
		return true, nil
	}, s.Table)
	fmt.Println("表名: ", stmt.Table)

	// 解析插入的列名
	fmt.Print("插入列名: ")
	for _, col := range s.Columns {
		stmt.Columns = append(stmt.Columns, Column{
			Name: col.String(),
		})
	}
	fmt.Println(stmt.Columns)
	// 解析插入的值
	fmt.Print("插入值: ")
	switch rows := s.Rows.(type) {
	case sqlparser.Values:
		for _, tuple := range rows {
			for _, val := range tuple {
				// fmt.Print(getExprInfo(val), " ")
				stmt.PlaceHolders = append(stmt.PlaceHolders, PlaceHolder{
					Name: getExprInfo(val),
				})
			}
		}
		// case *sqlparser.Select:
		// 	fmt.Print("子查询插入: ", sqlparser.String(rows))
	}
	fmt.Println(stmt.PlaceHolders)
}
