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
	"fmt"
)

type SQLite struct{}

func (s *SQLite) Name() string { return "sqlite" }

// Quote 引用标识符：SQLite 使用双引号（与 PostgreSQL 相同）
func (s *SQLite) Quote(name string) string {
	return `"` + name + `"`
}

// Placeholder 参数占位符：SQLite 使用 ?（与 MySQL 相同）
func (s *SQLite) Placeholder(index *uint8) string {
	*index = *index + 1
	return "?"
}

// Limit 分页语句：SQLite 使用 LIMIT ... OFFSET ...（与 PostgreSQL 格式相同）
func (s *SQLite) Limit(offset, limit uint) string {
	if offset > 0 {
		return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	return fmt.Sprintf(" LIMIT %d", limit)
}

// AutoIncrement 自增关键字：SQLite 使用 AUTOINCREMENT
func (s *SQLite) AutoIncrement() string { return "AUTOINCREMENT" }

// PrimaryKey 主键标识：从 PRAGMA table_info 的 pk 列获取（值为 0 或 1）
func (s *SQLite) PrimaryKey() string { return "pk" }

// UniqueKey 唯一键标识：通过 PRAGMA index_list 获取 unique=1 的索引
func (s *SQLite) UniqueKey() string { return "unique" }
