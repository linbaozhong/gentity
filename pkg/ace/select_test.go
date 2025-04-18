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
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/linbaozhong/gentity/example/model/db"
	"github.com/linbaozhong/gentity/example/model/define/dao/daoaccount"
	"github.com/linbaozhong/gentity/example/model/define/table/tblaccount"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"testing"
)

func init() {
	// 注册数据库驱动
	dialect.Register("mysql")
	_, e := Connect("mysql", "ssld_dev:Cu83&sr66@tcp(123.56.5.53:13306)/dispatch?charset=utf8mb4&parseTime=True&loc=Local")
	if e != nil {
		fmt.Println(e)
	}
}

func TestFunc(t *testing.T) {
	conds := Conds(tblaccount.Id.Eq(1)).
		Or(tblaccount.State.Eq(2),
			tblaccount.State.Eq(3))
	conds.Conds(tblaccount.Id.Eq(4)).
		And(tblaccount.State.Eq(5),
			tblaccount.State.Eq(6))
	bld := Where(conds.ToSlice()...).(*orm)
	t.Log(bld.where.String(), bld.whereParams)
}

func TestBuilder(t *testing.T) {
	db, e := Connect("mysql", "ssld_dev:Cu83&sr66@tcp(123.56.5.53:13306)/dispatch?charset=utf8mb4&parseTime=True&loc=Local")
	if e != nil {
		t.Error(e)
	}
	bld := newOrm().connect(db).Table("company").(*orm)
	t.Log(bld.db != nil)

	bld.Free()
	t.Log(bld.db != nil)

}

func TestUpdater(t *testing.T) {
	// Table("account").Set(tblaccount.State.Set(1), tblaccount.LoginName.Set("aaaaaaaaaaa")).Where(tblaccount.Id.Eq(1)).Update().Exec(context.Background())
	daoaccount.New().Update(context.Background(), Sets(tblaccount.State.Set(1), tblaccount.LoginName.Set("aaaaaaaaaaa")).ToSlice(), Conds(tblaccount.Id.Eq(1)).ToSlice()...)
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
	// New(
	// 	WithWhere(
	// 		Conds(tblaccount.Id.Eq(1))...),
	// 	WithOrderBy(Asc(tblaccount.Id)),
	// 	WithSet(tblaccount.Id.Set(1)))
}
