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

package sqlserver

import (
	_ "github.com/denisenkom/go-mssqldb"
)

// const (
// 	Quote_Char_Left  = "["           // SQL Server 使用方括号
// 	Quote_Char_Right = "]"           // SQL Server 使用方括号
// 	PrimaryKey       = "PRIMARY KEY" //
// 	AutoInc          = "IDENTITY"    // SQL Server 自增关键字
// 	UniqueKey        = "UNIQUE"      //
// )
//
// func Limit(offset, limit uint) string {
// 	if offset > 0 {
// 		// SQL Server 2012+ 使用 OFFSET FETCH
// 		return fmt.Sprintf(" OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offset, limit)
// 	}
// 	// SQL Server 使用 TOP
// 	return fmt.Sprintf(" TOP %d", limit)
// }
//
// func Placeholder(index *uint8) string {
// 	*index = *index + 1
// 	return fmt.Sprintf("@p%d", *index) // SQL Server 参数化查询使用 @p1, @p2...
// }
//
// func Quote(name string) string {
// 	return Quote_Char_Left + name + Quote_Char_Right
// }
