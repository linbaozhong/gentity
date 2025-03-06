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
	"github.com/linbaozhong/alipay/v3"
	"github.com/linbaozhong/gentity/pkg/conv"
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

func NewAlipay(appid, privateKey, publidKey, redirectUrl, notifyUrl string) Loginer {
	return &ali{
		appid:      appid,
		privateKey: privateKey,
		publicKey:  publidKey,
		notifyUrl:  notifyUrl,
		returnUrl:  redirectUrl,
	}
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

// GetLoginURL
func (a *ali) GetLoginURL(ctx context.Context, state string) (string, error) {
	stateCache.Save(ctx, state, struct{}{}, 10)
	u, e := a.client().PublicAppAuthorize([]string{"auth_user"}, a.returnUrl, state)
	if e != nil {
		return "", e
	}
	return u.String(), nil
}

// Login 登录验证
func (a *ali) Login(ctx context.Context, code, state string) (*SystemOauthTokenRsp, error) {
	// 检查state是否存在,并在使用完成后移除
	if _, e := stateCache.Fetch(ctx, state); e != nil {
		return nil, types.NewError(conv.String2Int(string(alipay.CodeMissingParam)),
			"state is not exist")
	}
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
		return &SystemOauthTokenRsp{
			AccessToken:  _res.AccessToken,
			ExpiresIn:    _res.ExpiresIn,
			RefreshToken: _res.RefreshToken,
			ReExpiresIn:  _res.ReExpiresIn,
			UserId:       _res.UserId,
			AuthStart:    _res.AuthStart,
			OpenId:       _res.OpenId,
			UnionId:      _res.UnionId,
		}, nil
	}
	return nil, types.NewError(conv.String2Int(string(_res.Code)),
		_res.Error.Error())
}

// GetUserInfo 获取用户信息
func (a *ali) GetUserInfo(ctx context.Context, token string) (*UserInfoShareRsp, error) {
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
		return &UserInfoShareRsp{
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

func (a *ali) Notify(ctx context.Context, req *alipay.Notification) (bool, error) {
	return false, nil
}
