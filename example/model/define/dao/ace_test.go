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

package dao

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/linbaozhong/gentity/example/model/db"
	obj "github.com/linbaozhong/gentity/example/model/define/table/companytbl"
	"github.com/linbaozhong/gentity/internal/schema"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/conv"
	"github.com/linbaozhong/gentity/pkg/log"
	"testing"
	"time"
)

var (
	dbx *ace.DB
)

func init() {
	var err error
	dbx, err = ace.Connect("mysql",
		"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	dbx.SetMaxOpenConns(50)
	dbx.SetMaxIdleConns(25)
	dbx.SetDebug(true)
	dbx.SetCache(ace.CacheTypeSyncMap, 0)
	log.RegisterLogger(false)
}

func TestCreateSet(t *testing.T) {
	defer dbx.Close()
	x := Company(dbx).C()
	result, err := x.Set(
		obj.CorpName.Set("aaaaaa"),
		obj.IsAuthenticated.Set(true),
	).Exec(context.Background())

	t.Log(x)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestCreateCols(t *testing.T) {
	defer dbx.Close()
	x := Company(dbx).C()
	result, err := x.Cols(
		obj.CorpName,
		obj.IsAuthenticated,
		obj.State,
		obj.CreatedTime,
	).Struct(context.Background(), &db.Company{
		FullCorpName:     "m1",
		IsEcologicalCorp: false,
		State:            0,
		CreatedTime:      time.Now(),
	}, &db.Company{
		FullCorpName:     "m2",
		IsEcologicalCorp: true,
		State:            1,
		CreatedTime:      time.Now(),
	}, &db.Company{
		FullCorpName:     "m3",
		IsEcologicalCorp: true,
		State:            0,
		CreatedTime:      time.Now(),
	})

	t.Log(x)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestInsert(t *testing.T) {
	defer dbx.Close()
	//
	id, err := Company(dbx).InsertOne(context.Background(), &db.Company{
		FullCorpName:     "x1",
		IsEcologicalCorp: false,
		State:            0,
		CreatedTime:      time.Now(),
	})
	t.Log(id)

	if err != nil {
		t.Fatal(err)
	}

	n, err := Company(dbx).InsertBatch(context.Background(), []*db.Company{{
		FullCorpName:     "x1",
		IsEcologicalCorp: false,
		State:            0,
		CreatedTime:      time.Now(),
	}, {
		FullCorpName:     "x2",
		IsEcologicalCorp: true,
		State:            1,
		CreatedTime:      time.Now(),
	}, {
		FullCorpName:     "x3",
		IsEcologicalCorp: true,
		State:            0,
		CreatedTime:      time.Now(),
	}})

	t.Log(n)

	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateSet(t *testing.T) {
	defer dbx.Close()
	x := Company(dbx).U()
	result, err := x.
		Set(
			obj.CorpName.Set("ttt"),
			obj.IsAuthenticated.Set(true),
		).
		SetExpr(obj.State.Incr()).
		Where(
			obj.Id.Eq(1),
		).
		Exec(context.Background())

	t.Log(x)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestUpdateCols(t *testing.T) {
	defer dbx.Close()

	x := Company(dbx).U()
	result, err := x.Cols(
		obj.FullCorpName,
		obj.IsEcologicalCorp,
		obj.State,
		obj.CreatedTime,
	).Where(
		obj.Id.Eq(1),
	).Struct(context.Background(), &db.Company{
		Id:               2,
		FullCorpName:     "ccc",
		IsEcologicalCorp: false,
		State:            1,
		CreatedTime:      time.Now(),
	})

	t.Log(x)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestUpdate(t *testing.T) {
	defer dbx.Close()

	ok, err := Company(dbx).UpdateBatch(context.Background(), []*db.Company{&db.Company{
		Id:               2,
		FullCorpName:     "ccc",
		IsEcologicalCorp: false,
		State:            1,
		CreatedTime:      time.Now(),
	}})

	t.Log(ok)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	defer dbx.Close()

	x := Company(dbx).D()
	result, err := x.Where(
		obj.Id.Eq(1),
	).Exec(context.Background())
	t.Log(x)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestSelect(t *testing.T) {
	defer dbx.Close()
	obj, _, err := Company(dbx).Find4Cols(context.Background(), 0, 0, []dialect.Field{
		obj.Id,
		obj.FullCorpName,
		obj.State,
	},
		obj.Id.Eq(2),
	)

	if err != nil {
		t.Fatal(err)
	}
	// t.Log(obj)
	for _, o := range obj {
		t.Log(o)
	}
}

func TestGenStruct(t *testing.T) {
	buf, err := schema.SqlFile2Struct("../../../data.sql", "db")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf))
}

func TestConvert(t *testing.T) {
	type Abc struct {
		Name string
	}
	var num = make([]string, 2)
	num[0] = "1"
	num[1] = "kds界定"

	r, _ := conv.Interface2Bytes(num)
	t.Log(r)

	var n []any
	conv.Bytes2Interface(r, &n)
	t.Log(n)

}
