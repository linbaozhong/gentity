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

package postgres

import "fmt"

type PostgreSQL struct{}

func (m *PostgreSQL) Name() string { return "postgres" }

// Quote 引用标识符：PostgreSQL 使用双引号
func (m *PostgreSQL) Quote(name string) string {
	return `"` + name + `"`
}

// Placeholder 参数占位符：PostgreSQL 使用 $1, $2, $3...
func (m *PostgreSQL) Placeholder(index *uint8) string {
	*index = *index + 1
	return fmt.Sprintf("$%d", *index)
}

// Limit 分页语句：PostgreSQL 使用 LIMIT ... OFFSET ...
func (m *PostgreSQL) Limit(offset, limit uint) string {
	if offset > 0 {
		return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	return fmt.Sprintf(" LIMIT %d", limit)
}

// AutoIncrement 自增关键字：PostgreSQL 使用 SERIAL 或 IDENTITY
func (m *PostgreSQL) AutoIncrement() string { return "SERIAL" }

// PrimaryKey 主键标识：从 information_schema.table_constraints 获取
func (m *PostgreSQL) PrimaryKey() string { return "p" }

// UniqueKey 唯一键标识
func (m *PostgreSQL) UniqueKey() string { return "u" }
