package example

import (
	"context"
	"fmt"
	"github.com/linbaozhong/gentity/example/model/define/dao"
	"github.com/linbaozhong/gentity/example/model/define/table/tblapp"
	"github.com/linbaozhong/gentity/example/model/do"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"log"
	"testing"
)

var (
	dbx *ace.DB
)

func init() {
	var err error
	dbx, err = ace.Connect(ace.Context, "mysql",
		"snow_dev:@tcp(123.56.5.53:3306)/snow_reminder?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	dbx.SetMaxOpenConns(50)
	dbx.SetMaxIdleConns(25)
	dbx.SetDebug(true)
	dbx.SetCache(ace.CacheTypeSyncMap, 0)
}

func TestName(t *testing.T) {
	defer dbx.Close()
	obj, _, err := dao.App(dbx).Find4Cols(context.Background(), 0, 0, []dialect.Field{
		tblapp.Id,
		tblapp.Arch,
		tblapp.Version,
	},
		tblapp.Id.Eq(2),
	)

	if err != nil {
		t.Fatal(err)
	}
	// t.Log(obj)
	for _, o := range obj {
		t.Log(fmt.Sprintf("%+v", o))
	}

}

func TestInsert(t *testing.T) {
	defer dbx.Close()

	r, e := dbx.C(do.AppTableName).Set(
		tblapp.Arch.Set("test"),
		tblapp.Id.Set(4),
		tblapp.Version.Set("v1.0.0"),
		tblapp.Url.Set("https://www.baidu.com"),
	).Exec(context.Background())
	if e != nil {
		t.Fatal(e)
	}
	t.Log(r.LastInsertId())
	t.Log(r.RowsAffected())
}
