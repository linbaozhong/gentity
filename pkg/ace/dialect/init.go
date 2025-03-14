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

import "github.com/linbaozhong/gentity/pkg/ace/dialect/mysql"

var (
	Placeholder = "?"
	Quote_Char  = "`"
	PrimaryKey  = ""
	AutoInc     = ""
	UniqueKey   = ""
)

func Register(driverName string) {
	switch driverName {
	case "mysql":
		Placeholder = mysql.Mysql_Placeholder
		Quote_Char = mysql.Mysql_Quote_Char
		PrimaryKey = mysql.Mysql_PrimaryKey
		AutoInc = mysql.Mysql_AutoInc
		UniqueKey = mysql.Mysql_UniqueKey
	}
}
