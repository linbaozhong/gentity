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
	"github.com/linbaozhong/gentity/example/model/db"
	"github.com/linbaozhong/gentity/example/model/define/table/tblaccount"
	"testing"
)

func TestUpdater(t *testing.T) {
	Set(tblaccount.State.Set(1)).Where(tblaccount.Id.Eq(1)).Update().Exec(context.Background())

}

func TestSelecter(t *testing.T) {
	d := db.NewAccount()
	e := ormSelectbuilder(Where(tblaccount.Id.Eq(1))).Get(context.Background(), d)
	if e != nil {
		t.Fatal(e)
	}
	t.Log(d)
}

func ormSelectbuilder(s SelectBuilder) Selecter {
	return s.Select()
}

func TestNew(t *testing.T) {
	New(
		WithWhere(
			tblaccount.Id.Eq(1),
			Or(tblaccount.State.Eq(1),
				tblaccount.State.Eq(2))),
		WithOrderBy(Asc(tblaccount.Id)),
		WithSet(tblaccount.Id.Set(1)))
}
