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
	"github.com/linbaozhong/gentity/pkg/ace/data"
	"testing"
)

var (
	dns = ""
)

func TestOrmCond(t *testing.T) {
	db, _ := Connect("mysql", dns)

	o := newOrm()
	o.Table("company").Cols(data.Id, data.LongName) // .
	// 	Join(dialect.Left_Join, data.Id, data.Id, data.Status.Eq(9), data.State.Gte(99))
	// o.Join(dialect.Right_Join, data.Id, data.Id, data.Status.Eq(8), data.State.Gte(88))
	o.Where(data.LongName.Eq("test")).And(data.Id.Eq(1))
	o.AndOr(data.Id.Eq(99), data.Id.Eq(98))
	o.OrAnd(data.Id.Eq(3), data.Id.Eq(4))
	o.RawWhere("id=?", 5)
	o.Debug(true).Select(db).QueryRow(context.Background())

	// where,params,e:=o.parseCond(o.cond)
	// t.Log(o.String())
}
