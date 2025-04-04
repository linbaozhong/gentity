package dao

import (
	"fmt"
	"github.com/linbaozhong/gentity/example/model/define/table/tblaccount"
	"testing"
)

func TestAce(t *testing.T) {
	cond := Where(tblaccount.Id.Eq(1),
		tblaccount.LoginName.Eq("linbaozhong"),
		Or(tblaccount.Id.Eq(2),
			tblaccount.LoginName.Eq("linbaozhong2")),
	)
	for _, c := range cond {
		s, v := c()
		t.Log(fmt.Sprintf("%s %+v", s, v))
	}

	ords := Order(Asc(tblaccount.Id, tblaccount.LoginName), Desc(tblaccount.State, tblaccount.Ctime))
	for _, o := range ords {
		s, v := o()
		t.Log(fmt.Sprintf("%s %+v", s, v))
	}
}
