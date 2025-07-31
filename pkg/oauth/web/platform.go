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

type OauthTokenRsp struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`  // 微信
	ExpiresIn    int64  `json:"expires_in"`    // 微信
	RefreshToken string `json:"refresh_token"` // 微信
	ReExpiresIn  int64  `json:"re_expires_in"`
	AuthStart    string `json:"auth_start"`
	OpenId       string `json:"openid"`  // 微信openid
	UnionId      string `json:"unionid"` // 微信unionid
	Scope        string `json:"scope"`   // 微信
}

type UserInfoRsp struct {
	AuthNo             string `json:"auth_no"`
	UserId             string `json:"user_id"`
	OpenId             string `json:"openid"`     // 微信
	UnionId            string `json:"unionid"`    // 微信
	Avatar             string `json:"headimgurl"` // 微信
	Province           string `json:"province"`   // 微信
	City               string `json:"city"`       // 微信
	Country            string `json:"country"`    // 微信
	NickName           string `json:"nickname"`   // 微信
	IsStudentCertified string `json:"is_student_certified"`
	UserType           string `json:"user_type"`
	UserStatus         string `json:"user_status"`
	IsCertified        string `json:"is_certified"`
	Sex                int8   `json:"sex"`    // 微信
	Gender             string `json:"gender"` // 支付宝
	Username           string `json:"user_name"`
	CertNo             string `json:"cert_no"`
	CertType           string `json:"cert_type"`
	Mobile             string `json:"mobile"`
}

type Platformer interface {
	Authorize(ctx context.Context, state string, isMobile bool) (string, error)
	// Callback code: 授权码，注意：微信返回的是code，支付宝返回的是auth_code
	Callback(ctx context.Context, code, state string) (*OauthTokenRsp, error)
	GetUserInfo(ctx context.Context, token, openid string) (*UserInfoRsp, error)
	GetPlatform() string
	PagePay(ctx context.Context, req *PagePayReq) (types.Smap, error)
	Notify(ctx context.Context, req *http.Request) (int, []byte)
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

func SplitState(state string) (Platform, string) {
	pos := strings.Index(state, ":")
	if pos == -1 {
		return "", state
	}
	return Platform(state[:pos]), state[pos+1:]
}

const Passbackchar = "__"

type PagePayReq struct {
	Bill     types.BigInt // 交易ID
	Sku      types.BigInt // 商品ID
	Sharer   types.BigInt // 分享人ID
	Buyer    types.BigInt // 买家ID
	Seller   types.BigInt // 卖家ID
	OpenID   types.String // 支付宝或微信用户openid
	Name     types.String // 商品名称
	Desc     types.String // 商品描述
	Amount   types.Money  // 交易金额
	Currency types.String // 货币名称
	// NotifyUrl types.String // 通知地址
	// ReturnUrl types.String // 返回地址
}
