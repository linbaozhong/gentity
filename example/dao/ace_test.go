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
	"github.com/linbaozhong/gentity/example/model"
	"github.com/linbaozhong/gentity/pkg/ace"
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/schema"
	"testing"
	"time"
)

var (
	dbx *ace.DB
	obj = model.UserTbl
)

func init() {
	var err error
	dbx, err = ace.Connect("mysql",
		"snow_dev:harrylin_68@tcp(39.107.252.66:13306)/snow_ganji?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	dbx.SetMaxOpenConns(50)
	dbx.SetMaxIdleConns(25)
	dbx.SetDebug(true)
	log.RegisterLogger(false)
}

func TestCreateSet(t *testing.T) {
	defer dbx.Close()
	x := User(dbx).C()

	result, err := x.Set(
		obj.Name.Set("aaaaaa"),
		obj.IsAllow.Set(true),
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
	x := User(dbx).C()
	result, err := x.Cols(
		obj.Name,
		obj.IsAllow,
		obj.Status,
		obj.CreatedTime,
	).Struct(context.Background(), &model.User{
		Name:        "m1",
		IsAllow:     false,
		Status:      0,
		CreatedTime: time.Now(),
	}, &model.User{
		Name:        "m2",
		IsAllow:     true,
		Status:      1,
		CreatedTime: time.Now(),
	}, &model.User{
		Name:        "m3",
		IsAllow:     true,
		Status:      0,
		CreatedTime: time.Now(),
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
	id, err := User(dbx).InsertOne(context.Background(), &model.User{
		Name:        "x1",
		IsAllow:     false,
		Status:      0,
		CreatedTime: time.Now(),
	})
	t.Log(id)

	if err != nil {
		t.Fatal(err)
	}

	n, err := User(dbx).InsertMulti(context.Background(), []*model.User{&model.User{
		Name:        "x1",
		IsAllow:     false,
		Status:      0,
		CreatedTime: time.Now(),
	}, &model.User{
		Name:        "x2",
		IsAllow:     true,
		Status:      1,
		CreatedTime: time.Now(),
	}, &model.User{
		Name:        "x3",
		IsAllow:     true,
		Status:      0,
		CreatedTime: time.Now(),
	}})

	t.Log(n)

	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateSet(t *testing.T) {
	defer dbx.Close()
	x := User(dbx).U()
	result, err := x.
		Set(
			obj.Name.Set("ttt"),
			obj.IsAllow.Set(true),
		).
		SetExpr(obj.Status.Incr()).
		Where(
			obj.ID.Eq(1),
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

	x := User(dbx).U()
	result, err := x.Cols(
		obj.Name,
		obj.IsAllow,
		obj.Status,
		obj.CreatedTime,
	).Where(
		obj.ID.Eq(1),
	).Struct(context.Background(), &model.User{
		ID:          2,
		Name:        "ccc",
		IsAllow:     false,
		Status:      1,
		CreatedTime: time.Now(),
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

	ok, err := User(dbx).UpdateMulti(context.Background(), []*model.User{&model.User{
		ID:          2,
		Name:        "ccc",
		IsAllow:     false,
		Status:      1,
		CreatedTime: time.Now(),
	}})

	t.Log(ok)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	defer dbx.Close()

	x := User(dbx).D()
	result, err := x.Where(
		obj.ID.Eq(1),
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
	obj, err := User(dbx).Find4Cols(context.Background(), 0, 0, []atype.Field{
		obj.ID,
		obj.Name,
		obj.Status,
	},
		obj.ID.Eq(2),
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
	buf, err := schema.SqlFile2Struct("./sql.sql")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf))
}
