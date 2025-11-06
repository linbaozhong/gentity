package wechat

import (
	"github.com/silenceper/wechat/v2/miniprogram/urllink"
	"testing"
)

func TestMini(t *testing.T) {
	mini, e := Programe("wxdd761eff21a0b34f", "c5b9574cf96157c3c3cd7a0490d41e11")
	if e != nil {
		t.Error(e)
	}

	ul := mini.GetURLLink()
	res, e := ul.Generate(&urllink.ULParams{
		Path:  "pages/index/index",
		Query: "state=1234567",
	})
	if e != nil {
		t.Error(e)
	}
	t.Log(res)
}
