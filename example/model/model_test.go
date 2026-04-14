package model

import (
	"github.com/linbaozhong/gentity/example/model/define/table/tblaccount"
	"testing"
)

func TestName(t *testing.T) {
	s := tblaccount.Id.MBRContains(116.461, 39.908, 5000)()
	t.Log(s)
}
