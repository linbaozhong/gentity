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

package mysql

import "fmt"

type MySQL struct{}

func (m *MySQL) Name() string { return "mysql" }

func (m *MySQL) Quote(name string) string {
	return "`" + name + "`"
}

func (m *MySQL) Placeholder(index *uint16) string {
	*index = *index + 1
	return "?"
}

func (m *MySQL) Limit(offset, limit uint) string {
	if offset > 0 {
		return fmt.Sprintf(" LIMIT %d,%d", offset, limit)
	}
	return fmt.Sprintf(" LIMIT %d", limit)
}

func (m *MySQL) AutoIncrement() string { return "AUTO_INCREMENT" }
func (m *MySQL) PrimaryKey() string    { return "PRI" }
func (m *MySQL) UniqueKey() string     { return "UNI" }
func (m *MySQL) Null(expr string) string {
	return " IFNULL(" + expr + ")"
}
