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

package handler

import (
	"github.com/linbaozhong/gentity/pkg/api"
	"reader/internal/constant/cst"
	"reader/internal/constant/err"
)

// AuthRequired 身份认证中间件
func AuthRequired(rule ...bool) api.Handler {
	return func(c api.Context) {
		// 从header读取access_token
		access := c.GetHeader(cst.Authorization_Key)

		if access == "" {
			api.Fail(c, err.Err_AuthToken_NotFound)
			return
		}

		// 路径鉴权
		if len(rule) == 0 || rule[0] {
			goto stop
		}

		//// 将用户简介植入上下文
		//c.Values().Set(cst.VisitorKey, types.Visitor{
		//	ID: 12345,
		//	IP: c.RemoteAddr(),
		//})

		c.Next()
		return
	stop:
		api.Fail(c, err.Err_Authorization_Limited)
		c.StopExecution()
	}
}
