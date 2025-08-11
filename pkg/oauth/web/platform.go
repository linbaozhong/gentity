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

package web

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/types"
	"net/http"
	"strings"
)

type Platformer interface {
	Authorize(ctx context.Context, state string, isMobile bool) (string, error)
	// Callback code: 授权码，注意：微信返回的是code，支付宝返回的是auth_code
	Callback(ctx context.Context, code, state string) (*OauthTokenRsp, error)
	GetUserInfo(ctx context.Context, token, openid string) (*UserInfoRsp, error)
	GetPlatform() string
	PagePay(ctx context.Context, req *PagePayReq) (types.Smap, error)
	Notify(ctx context.Context, req *http.Request, resp *NotifyResp) error
}

// 第三方平台
type Platform string

const (
	// 微信
	Wechat Platform = "wechat"
	// 支付宝
	Alipay Platform = "alipay"
)

func (p Platform) String() string {
	return string(p)
}

func (p Platform) Int() uint8 {
	switch p {
	case Wechat:
		return 1
	case Alipay:
		return 2
	default:
		return 0
	}
}

func (p Platform) CombinState(state string) string {
	return p.String() + ":" + state
}

func ToPlatform(i int) Platform {
	switch i {
	case 1:
		return Wechat
	case 2:
		return Alipay
	default:
		return Wechat
	}
}

func SplitState(state string) (Platform, string) {
	pos := strings.Index(state, ":")
	if pos == -1 {
		return "", state
	}
	return Platform(state[:pos]), state[pos+1:]
}

const Passbackchar = "__"
