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

package dto

import (
	"github.com/asaskevich/govalidator"
	"github.com/linbaozhong/gentity/pkg/types"
	"reader/internal/constant/err"
)

func (u *UserRegisterReq) Check() error {
	if u.UserName == "" {
		return err.ErrUserName
	}
	if u.Password == "" {
		return err.ErrPassword
	}
	if !govalidator.IsEmail(u.Email) {
		return types.NewError(610, "email格式错误")
	}
	if govalidator.ParamTagRegexMap["range"].MatchString(u.Password) {
	}
	return nil
}

func (u *GetUserReq) Check() error {
	if u.UserID == 0 {
		return err.ErrUserID
	}
	return nil
}
