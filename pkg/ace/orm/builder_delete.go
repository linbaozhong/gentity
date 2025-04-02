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

package orm

import (
	"context"
	"database/sql"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

// Deleter 删除器
type Deleter interface {
	Exec(ctx context.Context) (sql.Result, error)
}

// delete 删除器
type delete struct {
	*orm
}

// Delete 删除器
func (o *orm) Delete(x ...ace.Executer) Deleter {
	o.connect(x...)
	return &delete{
		orm: o,
	}
}

// Exec 执行删除
func (d *delete) Exec(ctx context.Context) (sql.Result, error) {
	defer d.Free()

	// if d.err != nil {
	// 	return nil, d.err
	// }

	d.command.WriteString("DELETE FROM " + dialect.Quote_Char + d.table + dialect.Quote_Char)
	// WHERE
	if d.where.Len() > 0 {
		d.command.WriteString(" WHERE " + d.where.String())
	}

	stmt, err := d.db.PrepareContext(ctx, d.command.String())
	if err != nil {
		return nil, err
	}
	if d.db.IsDB() {
		defer stmt.Close()
	}

	return stmt.ExecContext(ctx, d.whereParams...)
}
