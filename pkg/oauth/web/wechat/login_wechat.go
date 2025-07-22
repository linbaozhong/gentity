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
	"io"
	"net/http"
	"net/url"
)

type wx struct {
	appid       string
	appSecret   string
	lang        string
	redirectURI string // 回调地址
}
type option func(w *wx)

// 微信 API 地址
const (
	wechatAuthURL     = "https://open.weixin.qq.com/connect/qrconnect"
	wechatTokenURL    = "https://api.weixin.qq.com/sns/oauth2/access_token"
	wechatUserInfoURL = "https://api.weixin.qq.com/sns/userinfo"
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

func WithRedirectURI(redirectURI string) option {
	return func(w *wx) {
		w.redirectURI = redirectURI
	}
}

func New(opts ...option) web.Platformer {
	w := &wx{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *wx) Authorize(ctx context.Context, state string, isMobile bool) (string, error) {
	params := url.Values{}
	params.Set("appid", w.appid)
	params.Set("redirect_uri", w.redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "snsapi_login")                // 授权范围
	params.Set("state", web.Wechat.String()+":"+state) // 防CSRF令牌

	if isMobile {
		return fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?%s#wechat_redirect", params.Encode()), nil
	}
	return fmt.Sprintf("%s?%s#wechat_redirect", wechatAuthURL, params.Encode()), nil
}

func (w *wx) Callback(ctx context.Context, code, state string) (*web.OauthTokenRsp, error) {
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
	return &result, nil
}

func (w *wx) GetUserInfo(ctx context.Context, token, openid string) (*web.UserInfoRsp, error) {
	// 构造正确的 URL
	url := fmt.Sprintf("%s?access_token=%s&openid=%s&lang=%s",
		wechatUserInfoURL, token, openid, w.lang)

	// 发起 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析 JSON 响应
	var userInfo web.UserInfoRsp
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (w *wx) GetPlatform() string {
	return "wechat"
}
