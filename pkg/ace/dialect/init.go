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
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect/mysql"
	"github.com/linbaozhong/gentity/pkg/ace/dialect/postgres"
	"github.com/linbaozhong/gentity/pkg/ace/dialect/sqlite"
	"github.com/linbaozhong/gentity/pkg/sqlparser"
)

type (
	limit       func(offset, limit uint) string
	placeholder func(index int) string
	getTables   func(db *ace.DB, dbName string) ([]*sqlparser.Table, error)
)

var (
	Placeholder placeholder
	Quote_Char  = "`"
	PrimaryKey  = ""
	AutoInc     = ""
	UniqueKey   = ""
	Limit       limit
	GetTables   getTables
)

func Register(driverName string) {
	switch driverName {
	case "mysql":
		Placeholder = mysql.Placeholder
		Quote_Char = mysql.Quote_Char
		PrimaryKey = mysql.PrimaryKey
		AutoInc = mysql.AutoInc
		UniqueKey = mysql.UniqueKey
		Limit = mysql.Limit
		GetTables = mysql.GetTables
	case "sqlite":
		Placeholder = sqlite.Placeholder
		Quote_Char = sqlite.Quote_Char
		PrimaryKey = sqlite.PrimaryKey
		AutoInc = sqlite.AutoInc
		UniqueKey = sqlite.UniqueKey
		Limit = sqlite.Limit
		GetTables = sqlite.GetTables
	case "postgres":
		Placeholder = postgres.Placeholder
		Quote_Char = postgres.Quote_Char
		PrimaryKey = postgres.PrimaryKey
		AutoInc = postgres.AutoInc
		UniqueKey = postgres.UniqueKey
		Limit = postgres.Limit
		GetTables = postgres.GetTables
	}
}
