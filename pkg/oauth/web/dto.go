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

import "github.com/linbaozhong/gentity/pkg/types"

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

// PagePayReq 网页支付请求
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
}

// NotifyResp 回调通知数据
type NotifyResp struct {
	Code        int         `json:"code"`
	Id          string      `json:"id"`                    // 回调通知的唯一编号
	CreateTime  string      `json:"create_time"`           // 本次回调通知的创建时间
	EventType   string      `json:"event_type"`            // 支付成功通知的类型为TRANSACTION.SUCCESS
	Summary     string      `json:"summary,omitempty"`     // 内容摘要备注
	Transaction Transaction `json:"transaction,omitempty"` // 业务数据
}

// Transaction 业务数据
type Transaction struct {
	Amount         Amount `json:"amount,omitempty"`
	Appid          string `json:"appid,omitempty"` // App ID
	Attach         string `json:"attach,omitempty"`
	BankType       string `json:"bank_type,omitempty"`
	Mchid          string `json:"mchid,omitempty"`
	OutTradeNo     string `json:"out_trade_no,omitempty"`
	Payer          Payer  `json:"payer,omitempty"`
	SuccessTime    string `json:"success_time,omitempty"`
	TradeState     string `json:"trade_state,omitempty"`
	TradeStateDesc string `json:"trade_state_desc,omitempty"`
	TradeType      string `json:"trade_type,omitempty"`
	TransactionId  string `json:"transaction_id,omitempty"`
}

// Amount 金额
type Amount struct {
	Currency      string `json:"currency,omitempty"`       // 货币
	PayerCurrency string `json:"payer_currency,omitempty"` // 支付货币
	PayerTotal    int64  `json:"payer_total,omitempty"`    // 支付金额
	Total         int64  `json:"total,omitempty"`          // 订单金额
}

// Payer 支付方
type Payer struct {
	Openid string `json:"openid,omitempty"`
}
