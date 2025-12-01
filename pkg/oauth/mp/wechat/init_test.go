package wechat

import (
	"github.com/silenceper/wechat/v2/miniprogram/urllink"
	"github.com/silenceper/wechat/v2/miniprogram/urlscheme"
	"testing"
)

func TestMini(t *testing.T) {
	mini, e := Programe("wxecb6a1717661f6c1", "5a9251dacc72df464dae338943703706")
	if e != nil {
		t.Error(e)
	}

	ul := mini.GetURLLink()
	res, e := ul.Generate(&urllink.ULParams{
		Path:  "pages/login/login",
		Query: "state=1234567",
	})
	if e != nil {
		t.Error(e)
	}
	t.Log(res)

	scheme := mini.GetSURLScheme()
	res, e = scheme.Generate(&urlscheme.USParams{
		JumpWxa: &urlscheme.JumpWxa{
			Path:  "pages/login/login",
			Query: "state=1234567",
		},
		IsExpire:       true,
		ExpireType:     1,
		ExpireInterval: 30,
	})
	if e != nil {
		t.Error(e)
	}
	t.Log(res)
}
