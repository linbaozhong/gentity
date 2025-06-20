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

type StatementType string

const (
	InsertType StatementType = "INSERT"
	UpdateType               = "UPDATE"
	DeleteType               = "DELETE"
	SelectType               = "SELECT"
)

// Table 表信息
type Table struct {
	Name  string // 表名
	Alias string // 表别名
}

// Column 列信息
type Column struct {
	Table string // 表名
	Name  string // 列名
	Type  string // 列类型
}

type Statement struct {
	Type    StatementType // 语句类型
	Table   []Table       // 表名
	Columns []Column      // 列信息
}
