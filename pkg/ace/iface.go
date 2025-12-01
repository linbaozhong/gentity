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
	"database/sql/driver"
	"github.com/linbaozhong/gentity/pkg/ace/reflectx"
)

type (
	Executer interface {
		Mapper() *reflectx.Mapper
		Transaction(ctx context.Context, f func(tx *Tx) (any, error)) (any, error)
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		Debug() bool
		//Cache(string) cachego.Cache
		IsDB() bool
	}
)

type noRows struct{}

var _ driver.Result = noRows{}

func (noRows) LastInsertId() (int64, error) {
	return 0, Err_ToSql
}

func (noRows) RowsAffected() (int64, error) {
	return 0, Err_ToSql
}
