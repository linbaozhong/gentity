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
)

type OauthTokenRsp struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ReExpiresIn  int64  `json:"re_expires_in"`
	AuthStart    string `json:"auth_start"`
	OpenId       string `json:"open_id"`
	UnionId      string `json:"union_id"`
}

type UserInfoRsp struct {
	AuthNo             string `json:"auth_no"`
	UserId             string `json:"user_id"`
	OpenId             string `json:"open_id"`
	UnionId            string `json:"union_id"`
	Avatar             string `json:"avatar"`
	Province           string `json:"province"`
	City               string `json:"city"`
	NickName           string `json:"nick_name"`
	IsStudentCertified string `json:"is_student_certified"`
	UserType           string `json:"user_type"`
	UserStatus         string `json:"user_status"`
	IsCertified        string `json:"is_certified"`
	Gender             string `json:"gender"`
	Username           string `json:"user_name"`
	CertNo             string `json:"cert_no"`
	CertType           string `json:"cert_type"`
	Mobile             string `json:"mobile"`
}

type Platformer interface {
	Authorize(ctx context.Context, state string) (string, error)
	Callback(ctx context.Context, code, state string) (*OauthTokenRsp, error)
	GetUserInfo(ctx context.Context, token string) (*UserInfoRsp, error)
	GetPlatform() string
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
