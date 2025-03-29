package model

import (
	"context"
	"github.com/linbaozhong/gentity/example/model/define/table/tblusers"
	"github.com/linbaozhong/gentity/example/model/do"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/orm"
	"github.com/linbaozhong/gentity/pkg/log"
	"testing"
)

func TestBuild(t *testing.T) {
	db, e := ace.Connect("mysql", "snow_dev:harrylin_68@tcp(123.56.5.53:13306)/snow_cipher?charset=utf8mb4&parseTime=True&loc=Local")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	//db.SetDebug(true)

	bld := orm.New()

	_, e = bld.Set(tblusers.State.Set(0)).
		Where(tblusers.Id.Eq(1), tblusers.State.Eq(4)).
		And(tblusers.State.Eq(2), tblusers.Nick.Eq("test")).
		Or(tblusers.Id.Eq(3)).And(tblusers.State.Eq(4)).
		Desc(tblusers.Id).
		Update(do.UsersTableName).
		Exec(context.Background())
	if e != nil {
		t.Fatal(e)
	}
	log.Info(bld.String())
}
