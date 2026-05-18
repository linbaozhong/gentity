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

package sqlite

import (
	_ "github.com/mattn/go-sqlite3"
)

// const (
// 	Quote_Char_Left  = "\""
// 	Quote_Char_Right = "\""
// 	PrimaryKey       = "PK"
// 	AutoInc          = "AUTOINCREMENT"
// 	UniqueKey        = "UNIQUE"
// )
//
// func Limit(offset, limit uint) string {
// 	if offset > 0 {
// 		return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
// 	}
// 	return fmt.Sprintf(" LIMIT %d", limit)
// }
//
// func Placeholder(index *uint8) string {
// 	return "?"
// }
//
// func Quote(name string) string {
// 	return Quote_Char_Left + name + Quote_Char_Right
// }
