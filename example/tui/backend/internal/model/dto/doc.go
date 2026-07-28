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

// NewDocReq 新建文档请求
// @request
type NewDocReq struct {
	Filename *types.String `json:"filename" valid:"required"`
	Size     *types.Uint64 `json:"size,omitempty" valid:"min(1)"`
	Title    *types.String `json:"title,omitempty"`
	Descr    *types.String `json:"descr,omitempty"`
	Currency *types.Uint8  `json:"currency,omitempty"`
	Price    *types.Money  `json:"price" valid:"required,min(0.01)"`
	Bonus    *types.Money  `json:"bonus,omitempty"`
	Type     *types.String `json:"type,omitempty"`
	Ver      *types.Uint8  `json:"ver,omitempty"`
	Images   *types.String `json:"images,omitempty"`
}

// NewDocResp 新建文档响应
// @response
type NewDocResp struct {
	ID  *types.BigInt `json:"id"` // SKU ID
	Key *types.String `json:"key"`
}

// GetDocReq 获取文档请求
// @request
type GetDocReq struct {
	ID *types.BigInt `json:"id" valid:"required"`
}

// GetDocSkuReq 获取文档SKU请求
// @request
type GetDocSkuReq struct {
	ID *types.BigInt `json:"id" valid:"required"` // SKU ID
}

// GetDocSkuResp 获取文档SKU响应
// @response
type GetDocSkuResp struct {
	ID       *types.BigInt `json:"id"`
	Sharer   *types.BigInt `json:"sharer"`
	Title    *types.String `json:"title"`
	Descr    *types.String `json:"descr"`
	Type     *types.String `json:"type"`
	Ver      *types.Uint8  `json:"ver"`
	Images   *types.String `json:"images"`
	Currency *types.String `json:"currency"`
	Filename *types.String `json:"filename"`
	Ext      *types.String `json:"ext"`
	Length   *types.BigInt `json:"length"`
	Price    *types.Money  `json:"price"`
	Bonus    *types.Money  `json:"bonus"`
	State    *types.Int8   `json:"state"`
	Status   *types.Int8   `json:"status"`
	Key      *types.String `json:"key"` // 安全码
	Vis      *types.String `json:"vis"` // 访问者ID
}

// ListDocsReq 获取文档列表请求
// @request
type ListDocsReq struct {
	Page *int `json:"page" valid:"min(1)"`
	Size *int `json:"size" valid:"min(1)"`
}

// ListDocsResp 获取文档列表响应
// @response
type ListDocsResp struct {
	ID         *types.BigInt `json:"id"`
	Title      *types.String `json:"title,omitempty"`
	Filename   *types.String `json:"filename,omitempty"`
	Length     *types.BigInt `json:"length,omitempty"`
	Price      *types.Money  `json:"price,omitempty"`
	Bonus      *types.Money  `json:"bonus,omitempty"`
	State      *types.Int8   `json:"state"`
	StateText  string        `json:"state_text"`
	Status     *types.Int8   `json:"status"`
	StatusText string        `json:"status_text"`
	Ctime      *types.Time   `json:"ctime"`
	Utime      *types.Time   `json:"utime"`
}

// UpDownDocReq 上下架文档请求
// @request
type UpDownDocReq struct {
	ID     *types.BigInt `json:"id" valid:"required"` // doc id
	Status *types.Int8   `json:"status"`
}

// UpDownDocResp 上下架文档响应
// @response
type UpDownDocResp struct {
	ID     *types.BigInt `json:"id"`
	Status *types.Int8   `json:"status"`
}
