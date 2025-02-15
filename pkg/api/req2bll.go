// Copyright © 2023 SnowIM. All rights reserved.
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

package api

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/api/broker"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

var (
	Param_Invalid = types.NewError(620, "参数无效")
	UnKnown       = types.NewError(610, "未知错误")
)

// Post post请求
// Content-Type：application/json，req结构体的字段tag为json
// Content-Type: application/x-www-form-urlencoded，req结构体的字段tag为form
// Content-Type: multipart/form-data，req结构体的字段tag为form
func Post[A, B any](
	ctx Context,
	fn func(ctx context.Context, req *A, resp *B) error,
) {
	var (
		req A
		e   error
	)
	Initiate(ctx, &req)

	switch ctx.GetContentTypeRequested() {
	case "application/json":
		e = ReadJSON(ctx, &req)
	case "application/x-www-form-urlencoded", "multipart/form-data":
		e = ctx.ReadForm(&req)
	default:
		if ctx.Request().URL.RawQuery == "" {
			e = ctx.ReadForm(&req)
		} else {
			e = ReadQuery(ctx, &req)
		}
	}

	if e != nil {
		Fail(ctx, Param_Invalid)
		log.Error(e)
		return
	}
	if e := broker.Validate(&req); e != nil {
		Fail(ctx, e)
		log.Error(e)
		return
	}

	var resp B
	if e := fn(ctx, &req, &resp); e != nil {
		Fail(ctx, e)
		log.Error(e)
		return
	}
	Ok(ctx, resp)
}

// Get get请求：
// 首先尝试读取query，req结构体的字段 tag 为 url 或者 param。
// 如果query为空，则尝试读取form，req结构体的字段tag为form。
func Get[A, B any](
	ctx Context,
	fn func(ctx context.Context, req *A, resp *B) error,
) {
	var (
		req A
		e   error
	)
	Initiate(ctx, &req)

	if ctx.Request().URL.RawQuery == "" {
		e = ctx.ReadForm(&req)
	} else {
		e = ReadQuery(ctx, &req)
	}
	if e != nil {
		Fail(ctx, Param_Invalid)
		log.Error(e)
		return
	}
	//
	if e = broker.Validate(&req); e != nil {
		Fail(ctx, e)
		log.Error(e)
		return
	}

	var resp B
	if e = fn(ctx, &req, &resp); e != nil {
		Fail(ctx, e)
		log.Error(e)
		return
	}
	Ok(ctx, resp)
}
