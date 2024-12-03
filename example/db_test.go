package example

import (
	"context"
	"fmt"
	"github.com/linbaozhong/gentity/example/model/define/dao"
	"github.com/linbaozhong/gentity/example/model/define/table/tblapp"
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
		"")
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
