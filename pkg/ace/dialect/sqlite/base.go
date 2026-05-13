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

// pkg/ace/dialect/sqlite/sqlite.go
package sqlite

import "fmt"

type SQLite struct{}

func (s *SQLite) Name() string { return "sqlite" }

func (s *SQLite) Quote(name string) string {
	return `"` + name + `"`
}

func (s *SQLite) Placeholder(index int) string {
	return "?"
}

func (s *SQLite) Limit(offset, limit uint) string {
	if offset > 0 {
		return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	return fmt.Sprintf(" LIMIT %d", limit)
}

func (s *SQLite) AutoIncrement() string { return "AUTOINCREMENT" }
func (s *SQLite) PrimaryKey() string    { return "pk" }
func (s *SQLite) UniqueKey() string     { return "unique" }
