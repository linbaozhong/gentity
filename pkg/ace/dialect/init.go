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

import (
	"database/sql"
	"github.com/linbaozhong/gentity/pkg/ace/dialect/mysql"
	"github.com/linbaozhong/gentity/pkg/ace/dialect/postgres"
	"github.com/linbaozhong/gentity/pkg/ace/dialect/sqlite"
	"github.com/linbaozhong/gentity/pkg/ace/dialect/sqlserver"
	"github.com/linbaozhong/gentity/pkg/sqlparser"
)

type (
	limit       func(offset, limit uint) string
	placeholder func(index *uint8) string
	getTables   func(db *sql.DB, dbName string) ([]*sqlparser.Table, error)
	getColumns  func(db *sql.DB, dbName string) (map[string][]*sqlparser.Column, error)
	quote       func(name string) string
)

var (
	Placeholder      placeholder
	Quote_Char_Left  = "`"
	Quote_Char_Right = "`"
	PrimaryKey       = ""
	// AutoInc     = ""
	// UniqueKey   = ""
	Limit   limit
	Tables  getTables
	Columns getColumns
	Quote   quote
)

// Register 注册数据库方言
func Register(driverName string) {
	switch driverName {
	case "mysql":
		Placeholder = mysql.Placeholder
		Quote_Char_Left = mysql.Quote_Char_Left
		Quote_Char_Right = mysql.Quote_Char_Right
		PrimaryKey = mysql.PrimaryKey
		// AutoInc = mysql.AutoInc
		// UniqueKey = mysql.UniqueKey
		Limit = mysql.Limit
		Tables = mysql.GetTables
		Columns = mysql.GetColumns
		Quote = mysql.Quote
	case "sqlite":
		Placeholder = sqlite.Placeholder
		Quote_Char_Left = sqlite.Quote_Char_Left
		Quote_Char_Right = sqlite.Quote_Char_Right
		PrimaryKey = sqlite.PrimaryKey
		// AutoInc = sqlite.AutoInc
		// UniqueKey = sqlite.UniqueKey
		Limit = sqlite.Limit
		Tables = sqlite.GetTables
		Columns = sqlite.GetColumns
		Quote = sqlite.Quote
	case "postgres":
		Placeholder = postgres.Placeholder
		Quote_Char_Left = postgres.Quote_Char_Left
		Quote_Char_Right = postgres.Quote_Char_Right
		PrimaryKey = postgres.PrimaryKey
		// AutoInc = postgres.AutoInc
		// UniqueKey = postgres.UniqueKey
		Limit = postgres.Limit
		Tables = postgres.GetTables
		Columns = postgres.GetColumns
		Quote = postgres.Quote
	case "sqlserver":
		Placeholder = sqlserver.Placeholder
		Quote_Char_Left = sqlserver.Quote_Char_Left
		Quote_Char_Right = sqlserver.Quote_Char_Right
		PrimaryKey = sqlserver.PrimaryKey
		// AutoInc = sqlserver.AutoInc
		// UniqueKey = sqlserver.UniqueKey
		Limit = sqlserver.Limit
		Tables = sqlserver.GetTables
		Columns = sqlserver.GetColumns
		Quote = sqlserver.Quote
	}
}

type Driverer interface {
	GetTables(db *sql.DB, dbName string) ([]*sqlparser.Table, error)
}

func GetTables(db *sql.DB, dbName string) ([]*sqlparser.Table, error) {
	// 表名,表注释
	ts, err := Tables(db, dbName)
	if err != nil {
		return nil, err
	}

	// 表字段信息
	ms, err := Columns(db, dbName)
	if err != nil {
		return nil, err
	}

	for _, t := range ts {
		t.ColumnsX = ms[t.Name]
	}
	return ts, nil
}
