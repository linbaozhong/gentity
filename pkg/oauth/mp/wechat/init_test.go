package wechat

import (
	"github.com/silenceper/wechat/v2/miniprogram/urllink"
	"testing"
)

func TestMini(t *testing.T) {
	mini, e := Programe("wxecb6a1717661f6c1", "5a9251dacc72df464dae338943703706")
	if e != nil {
		t.Error(e)
	}

	ul := mini.GetURLLink()
	res, e := ul.Generate(&urllink.ULParams{
		Path:       "pages/index/index",
		Query:      "state=1234567",
		EnvVersion: "trial",
	})
	if e != nil {
		t.Error(e)
	}
	t.Log(res)
}
