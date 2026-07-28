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
	"github.com/linbaozhong/gentity/pkg/types"
)

// OrderLog 订单日志
// @response
// @request
type OrderLog struct {
	Status      *types.Uint16 `json:"status"`      // 状态码
	StatusText  string        `json:"status_text"` // 状态
	StatusTime  *types.Time   `json:"status_time"` // 状态时间
	Description string        `json:"description"` // 描述
	Event       string        `json:"event"`       // 触发状态的事件
	Trigger     string        `json:"trigger"`     // 触发状态的人
}

// ListOrdersReq 获取订单列表请求
// @request
type ListOrdersReq struct {
	Page *int `json:"page" valid:"numeric,min(1)"`
	Size *int `json:"size" valid:"numeric,min(1)"`
}

// ListOrdersResp 获取订单列表响应
// @response
type ListOrdersResp struct {
	Id           *types.BigInt `json:"id"`            // 订单id
	Sku          *types.BigInt `json:"sku"`           // 商品id
	Name         *types.String `json:"name"`          // 商品名称
	Price        *types.Money  `json:"price"`         // 商品单价
	Amount       *types.Money  `json:"amount"`        // 订单金额
	Currency     *types.String `json:"currency"`      // 订单币种
	AuthorAmount *types.Money  `json:"author_amount"` // 作者金额
	SharerAmount *types.Money  `json:"sharer_amount"` // 分享者金额
	Status       *types.Uint16 `json:"status"`        // 订单状态
	StatusText   string        `json:"status_text"`   // 订单状态文本
	Ctime        *types.Time   `json:"ctime"`         // 订单时间
	RefundBill   *types.BigInt `json:"refund_bill"`   // 退款账单id
	Bill         *types.BigInt `json:"bill"`          // 支付账单id
	AllowRefund  bool          `json:"allow_refund"`  // 是否运行退款
}

// RefundReq 退款请求
// @request
type RefundReq struct {
	Id         *types.BigInt `json:"id" valid:"required"` // 订单id
	Bill       *types.BigInt `json:"bill"`                // 支付账单id
	TransId    *types.String `json:"trans_id"`            // 支付交易id
	RefundBill *types.BigInt `json:"refund_bill"`         // 退款账单id
	Amount     *types.Money  `json:"amount"`              // 退款金额
	Currency   *types.String `json:"currency"`            // 退款币种
	Reason     *types.String `json:"reason"`              // 退款原因
	Platform   string        `json:"platform"`            // 退款平台
	NotifyUrl  *types.String `json:"notify_url"`          // 退款通知地址
}

// RefundResp 退款响应
// @response
type RefundResp struct {
	Id         *types.BigInt `json:"id"`          // 订单id
	RefundBill *types.BigInt `json:"refund_bill"` // 退款账单id
	Status     *types.Uint16 `json:"status"`      // 订单状态
	StatusText string        `json:"status_text"` // 订单状态文本
	StatusTime *types.Time   `json:"status_time"` // 订单状态时间
}
