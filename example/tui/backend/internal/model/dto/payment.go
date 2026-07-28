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

import "github.com/linbaozhong/gentity/pkg/types"

// PagePayReq 支付宝页面支付请求
// @request
type PagePayReq struct {
	Mobile *bool         `json:"mobile"`
	From   *string       `json:"from"`
	ID     *types.BigInt `json:"id" valid:"required,min(1)~无效的id值"` // 文档Sku id
	Sharer *types.BigInt `json:"sharer" valid:"int"`                // 分享者id
	Vis    *types.String `json:"vis" valid:"required"`              // 访问者id
}

// PagePayResp 支付宝页面支付响应
// @response
type PagePayResp struct {
	Link   string `json:"link,omitempty"` // 统一下单url
	Scheme string `json:"scheme,omitempty"`
	Query  string `json:"query,omitempty"`
}

// CheckPaidReq 检查订单是否支付请求
// @request
type CheckPaidReq struct {
	ID *types.BigInt `json:"id" valid:"required,min(1)"` // 文档Sku id
}

// CheckPaidResp 检查订单是否支付响应
// @response
type CheckPaidResp struct {
	Key *types.String `json:"key"` // 文档安全码
	// Status *types.String `json:"status"` // 订单状态
}
