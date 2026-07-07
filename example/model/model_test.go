package model

import (
	"context"
	"github.com/linbaozhong/gentity/example/model/define/dao"
	"github.com/linbaozhong/gentity/example/model/define/table/tblaccount"
	"github.com/linbaozhong/gentity/example/model/do"
	"github.com/linbaozhong/gentity/pkg/ace"
	"os"
	"path/filepath"
	"testing"
)

func TestName(t *testing.T) {
	ctx := context.Background()
	_, e := ace.Table(do.AccountTableName).
		Where(tblaccount.Id.MBRContains(116.461, 39.908, 5000)).
		Debug(true).Select().QueryRow(ctx)
	t.Log(e)

	db := new(ace.DB)
	db.Transaction(ctx, func(tx *ace.Tx) (any, error) {
		_, e := dao.CompanyMan(tx).Exists(ctx)
		if e != nil {
			return nil, e
		}
		a, _, e := dao.Account(tx).GetByID(ctx, 1)
		if e != nil {
			return nil, e
		}
		return a, nil
	})
}

func TestArgs(t *testing.T) {
	t.Log(filepath.Base(os.Args[1]))
}
