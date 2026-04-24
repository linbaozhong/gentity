package model

import (
	"context"
	"github.com/linbaozhong/gentity/example/model/define/table/tblaccount"
	"github.com/linbaozhong/gentity/example/model/do"
	"github.com/linbaozhong/gentity/pkg/ace"
	"testing"
)

func TestName(t *testing.T) {
	_, e := ace.Table(do.AccountTableName).
		Where(tblaccount.Id.MBRContains(116.461, 39.908, 5000)).
		ToSql(true).Select().QueryRow(context.Background())
	t.Log(e)
}
