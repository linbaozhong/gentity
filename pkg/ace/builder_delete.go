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

package ace

import (
	"context"
	"database/sql"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/log"
)

type DeleteBuilder interface {
	Table(a any) Builder
	GetTableName() string
	Wherer
	Delete(x ...Executer) Deleter
	// ToSql 不传参数或者参数为 true 时，仅打印SQL语句，不执行。
	ToSql(...bool) Builder
}

// Deleter 删除器
type Deleter interface {
	Exec(ctx context.Context) (sql.Result, error)
}

// delete 删除器
type delete struct {
	*orm
}

// Delete 删除器
func (o *orm) Delete(x ...Executer) Deleter {
	o.connect(x...)
	return &delete{
		orm: o,
	}
}

// Exec 执行删除
func (d *delete) Exec(ctx context.Context) (sql.Result, error) {
	defer d.Free()

	d.command.WriteString("DELETE FROM " + dialect.Quote_Char + d.table + dialect.Quote_Char)
	// WHERE
	if d.where.Len() > 0 {
		d.command.WriteString(" WHERE " + d.where.String())
	}

	// 只返回SQL语句，不执行
	if d.toSql {
		log.Info(d.String())
		return &noRows{}, Err_ToSql
	}

	stmt, err := d.db.PrepareContext(ctx, d.command.String())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if d.db.IsDB() {
		defer stmt.Close()
	}

	r, err := stmt.ExecContext(ctx, d.whereParams...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return r, nil
}
