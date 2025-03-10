// Copyright © 2023 Linbaozhong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/oauth/web"
	"github.com/linbaozhong/gentity/pkg/types"
	"io"
	"net/http"
)

type wx struct {
	appid     string
	appSecret string
	lang      string
}
type option func(w *wx)

// 微信 API 地址
const (
	wechatAuthURL     = "https://open.weixin.qq.com/connect/qrconnect"
	wechatTokenURL    = "https://api.weixin.qq.com/sns/oauth2/access_token"
	wechatUserInfoURL = "https://api.weixin.qq.com/sns/userinfo"
	redirectURI       = "http://your-domain.com/wx/callback" // 替换为你的回调地址
)

func WithAppId(appid string) option {
	return func(w *wx) {
		w.appid = appid
	}
}
func WithLang(lang string) option {
	return func(w *wx) {
		w.lang = lang
	}
}
func WithAppSecret(appSecret string) option {
	return func(w *wx) {
		w.appSecret = appSecret
	}
}

func New(opts ...option) web.Loginer {
	w := &wx{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *wx) GetLoginURL(ctx context.Context, state string) (string, error) {
	web.StateCache.Save(ctx, state, struct{}{}, 10)
	return fmt.Sprintf("%s?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect",
		wechatAuthURL, w.appid, redirectURI, state), nil
}
func (w *wx) Login(ctx context.Context, code, state string) (*web.OauthTokenRsp, error) {
	// 检查state是否存在,并在使用完成后移除
	if _, e := web.StateCache.Fetch(ctx, state); e != nil {
		return nil, types.NewError(400001,
			"state is not exist")
	}
	url := fmt.Sprintf("%s?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		wechatTokenURL, w.appid, w.appSecret, code)
	resp, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()

	body, e := io.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}

	var result web.OauthTokenRsp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &web.OauthTokenRsp{}, nil
}

func (w *wx) GetUserInfo(ctx context.Context, token string) (*web.UserInfoRsp, error) {
	return &web.UserInfoRsp{}, nil
}
