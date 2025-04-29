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
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

var (
	Param_Invalid = types.NewError(-620, "参数无效")
	UnKnown       = types.NewError(-610, "未知错误")
)

// Get get请求：
// 读取query，req结构体的字段 tag 为 url。
func Get[A, B any](
	ctx Context,
	fn func(ctx context.Context, req *A, resp *B) error,
) error {
	var (
		req  A
		resp B
	)

	_, e := serviceContext(ctx, &req, &resp, readGetRequest[A], fn)
	if e != nil {
		Fail(ctx, e)
		log.Error(e)
		return e
	}
	return Ok(ctx, resp)
}

// Post post请求
// Content-Type：application/json，req结构体的字段tag为json
// Content-Type: application/x-www-form-urlencoded，req结构体的字段tag为form
// Content-Type: multipart/form-data，req结构体的字段tag为form
func Post[A, B any](
	ctx Context,
	fn func(ctx context.Context, req *A, resp *B) error,
) error {
	var (
		req  A
		resp B
	)

	_, e := serviceContext(ctx, &req, &resp, readPostRequest[A], fn)
	if e != nil {
		Fail(ctx, e)
		log.Error(e)
		return e
	}
	return Ok(ctx, resp)
}

// Redirect 重定向
func Redirect[A any](ctx Context,
	fn func(ctx context.Context, req *A, resp *string) error) error {
	var (
		req  A
		resp string
	)
	_, e := serviceContext(ctx, &req, &resp, readPostRequest[A], fn)
	if e != nil {
		Fail(ctx, e)
		log.Error(e)
		return e
	}
	ctx.Redirect(resp)
	return nil
}

func Stream[A, B any](
	ctx Context,
	fn func(ctx Context, req *A, resp *B) error) error {
	var (
		req  A
		resp B
	)
	_, e := service(ctx, &req, &resp, readPostRequest[A], fn)
	if e != nil {
		Fail(ctx, e)
		log.Error(e)
		return e
	}
	return Ok(ctx, resp)
}
