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
	obj "github.com/linbaozhong/gentity/example/define/table/companytbl"
	"github.com/linbaozhong/gentity/example/model"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/schema"
	"testing"
	"time"
)

var (
	dbx *ace.DB
)

func init() {
	var err error
	dbx, err = ace.Connect("mysql",
		"ssld_dev:Cu83&sr66@tcp(39.107.252.66:13306)/assessment?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	dbx.SetMaxOpenConns(50)
	dbx.SetMaxIdleConns(25)
	dbx.SetDebug(true)
	log.RegisterLogger(false)
}

func TestScheme(t *testing.T) {
	rows, err := dbx.Query(`SELECT table_name,column_name,column_default,data_type,ifnull(character_maximum_length,0),column_key,extra,column_comment FROM information_schema.COLUMNS WHERE table_schema = ?`, "assessment")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	ms := make(map[string][]dialect.Column)
	if rows.Next() {
		var tableName string
		col := dialect.Column{}
		err = rows.Scan(&tableName, &col.Name, &col.Default, &col.Type, &col.Size, &col.Key, &col.Extra, &col.Comment)
		if err != nil {
			t.Fatal(err)
		}

		if cols, ok := ms[tableName]; ok {
			ms[tableName] = append(cols, col)
		} else {
			ms[tableName] = []dialect.Column{col}
		}
	}
	if err = rows.Err(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateSet(t *testing.T) {
	defer dbx.Close()
	x := Company(dbx).C()
	result, err := x.Set(
		obj.CorpName.Set("aaaaaa"),
		obj.IsAuthenticated.Set(true),
	).Do(context.Background())

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
	).Struct(context.Background(), &model.Company{
		FullCorpName:     "m1",
		IsEcologicalCorp: false,
		State:            0,
		CreatedTime:      time.Now(),
	}, &model.Company{
		FullCorpName:     "m2",
		IsEcologicalCorp: true,
		State:            1,
		CreatedTime:      time.Now(),
	}, &model.Company{
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
	id, err := Company(dbx).InsertOne(context.Background(), &model.Company{
		FullCorpName:     "x1",
		IsEcologicalCorp: false,
		State:            0,
		CreatedTime:      time.Now(),
	})
	t.Log(id)

	if err != nil {
		t.Fatal(err)
	}

	n, err := Company(dbx).InsertMulti(context.Background(), []*model.Company{&model.Company{
		FullCorpName:     "x1",
		IsEcologicalCorp: false,
		State:            0,
		CreatedTime:      time.Now(),
	}, &model.Company{
		FullCorpName:     "x2",
		IsEcologicalCorp: true,
		State:            1,
		CreatedTime:      time.Now(),
	}, &model.Company{
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
		Do(context.Background())

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
	).Struct(context.Background(), &model.Company{
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

	ok, err := Company(dbx).UpdateMulti(context.Background(), []*model.Company{&model.Company{
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
	).Do(context.Background())
	t.Log(x)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.LastInsertId())
	t.Log(result.RowsAffected())
}

func TestSelect(t *testing.T) {
	defer dbx.Close()
	obj, _, err := Company(dbx).Find4Cols(context.Background(), 0, 0, []atype.Field{
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
	buf, err := schema.SqlFile2Struct("./data.sql", "db")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf))
}
