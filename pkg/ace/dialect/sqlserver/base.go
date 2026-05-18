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
	"fmt"
)

type SQLServer struct{}

func (s *SQLServer) Name() string { return "sqlserver" }

// Quote 引用标识符：SQL Server 使用方括号 []
func (s *SQLServer) Quote(name string) string {
	return "[" + name + "]"
}

// Placeholder 参数占位符：
// SQL Server go-mssqldb 驱动支持 ?（自动转换）
// 原生格式为 @p1, @p2...
func (s *SQLServer) Placeholder(index *uint8) string {
	*index = *index + 1
	// 使用 ? 兼容 go-mssqldb 驱动
	// 如需原生 @p1 格式可改为: fmt.Sprintf("@p%d", *index)
	return "?"
}

// Limit 分页语句：SQL Server 使用 OFFSET ... FETCH NEXT ... ROWS ONLY
// 注意：OFFSET FETCH 必须配合 ORDER BY 使用（SQL Server 2012+）
func (s *SQLServer) Limit(offset, limit uint) string {
	return fmt.Sprintf(" OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offset, limit)
}

// AutoIncrement 自增关键字：SQL Server 使用 IDENTITY
func (s *SQLServer) AutoIncrement() string { return "IDENTITY" }

// PrimaryKey 主键标识：从 information_schema.table_constraints 获取
func (s *SQLServer) PrimaryKey() string { return "PRIMARY KEY" }

// UniqueKey 唯一键标识
func (s *SQLServer) UniqueKey() string { return "UNIQUE" }
