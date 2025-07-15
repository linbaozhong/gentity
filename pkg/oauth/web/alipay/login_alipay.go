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

package alipay

import (
	"context"
	"fmt"
	"github.com/linbaozhong/alipay/v3"
	"github.com/linbaozhong/gentity/pkg/conv"
	"github.com/linbaozhong/gentity/pkg/oauth/web"
	"github.com/linbaozhong/gentity/pkg/types"
)

type ali struct {
	appid      string
	privateKey string
	publicKey  string
	returnUrl  string
	notifyUrl  string
}

var (
	aliClient *alipay.Client
)

type option func(a *ali)

func WithAppId(appid string) option {
	return func(a *ali) {
		a.appid = appid
	}
}
func WithPrivateKey(privateKey string) option {
	return func(a *ali) {
		a.privateKey = privateKey
	}
}
func WithPublicKey(publicKey string) option {
	return func(a *ali) {
		a.publicKey = publicKey
	}
}
func WithReturnUrl(returnUrl string) option {
	return func(a *ali) {
		a.returnUrl = returnUrl
	}
}
func WithNotifyUrl(notifyUrl string) option {
	return func(a *ali) {
		a.notifyUrl = notifyUrl
	}
}

func New(opts ...option) web.Platformer {
	a := &ali{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *ali) client() *alipay.Client {
	if aliClient == nil {
		var e error
		aliClient, e = alipay.New(a.appid, a.privateKey, true)
		if e != nil {
			panic(e)
		}
		aliClient.LoadAliPayPublicKey(a.publicKey)
	}

	return aliClient
}

// Authorize 生成支付宝授权链接
func (a *ali) Authorize(ctx context.Context, state string) (string, error) {
	u, e := a.client().PublicAppAuthorize([]string{"auth_user"}, a.returnUrl, web.Alipay.String()+":"+state)
	if e != nil {
		return "", e
	}
	return u.String(), nil
}

// Callback alipay 用户授权回调，使用授权码换取 access_token
func (a *ali) Callback(ctx context.Context, code, state string) (*web.OauthTokenRsp, error) {
	fmt.Println(code, state)

	_res, e := a.client().SystemOauthToken(ctx,
		alipay.SystemOauthToken{
			Code:      code,
			GrantType: "authorization_code",
		})
	if e != nil {
		return nil, types.NewError(conv.String2Int(string(alipay.CodeUnknowError)),
			e.Error())
	}
	if _res.IsSuccess() {
		_token := &web.OauthTokenRsp{
			AccessToken:  _res.AccessToken,
			ExpiresIn:    _res.ExpiresIn,
			RefreshToken: _res.RefreshToken,
			ReExpiresIn:  _res.ReExpiresIn,
			UserId:       _res.UserId,
			AuthStart:    _res.AuthStart,
			OpenId:       _res.OpenId,
			UnionId:      _res.UnionId,
		}

		return _token, nil
	}
	return nil, types.NewError(conv.String2Int(string(_res.Code)),
		_res.Error.Error())
}

// GetUserInfo 获取用户信息
func (a *ali) GetUserInfo(ctx context.Context, token, openid string) (*web.UserInfoRsp, error) {
	_res, e := a.client().UserInfoShare(ctx,
		alipay.UserInfoShare{
			AuthToken: token,
		})
	if e != nil {
		return nil, &alipay.Error{
			Code: alipay.CodeUnknowError,
			Msg:  e.Error(),
		}
	}
	if _res.IsSuccess() {
		return &web.UserInfoRsp{
			AuthNo:             _res.AuthNo,
			UserId:             _res.UserId,
			OpenId:             _res.OpenId,
			UnionId:            _res.UnionId,
			Avatar:             _res.Avatar,
			Province:           _res.Province,
			City:               _res.City,
			NickName:           _res.NickName,
			IsStudentCertified: _res.IsStudentCertified,
			UserType:           _res.UserType,
			UserStatus:         _res.UserStatus,
			IsCertified:        _res.IsCertified,
			Gender:             _res.Gender,
			Username:           _res.Username,
			CertNo:             _res.CertNo,
			CertType:           _res.CertType,
			Mobile:             _res.Mobile,
		}, nil
	}
	return nil, types.NewError(conv.String2Int(string(_res.Code)),
		_res.Error.Error())
}

func (a *ali) GetPlatform() string {
	return "alipay"
}
