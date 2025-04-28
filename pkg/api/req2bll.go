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
	"time"
)

var (
	Param_Invalid = types.NewError(620, "参数无效")
	UnKnown       = types.NewError(610, "未知错误")
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
		read = func(ctx Context, req *A) error {
			Initiate(ctx, req)

			if ctx.Request().URL.RawQuery == "" {
				return ReadForm(ctx, req)
			} else {
				return ReadQuery(ctx, req)
			}
		}
	)

	return logicProcessing(ctx, &req, &resp, read, fn)
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
		read = func(ctx Context, req *A) error {
			Initiate(ctx, req)

			switch ctx.GetContentTypeRequested() {
			case "application/json":
				return ReadJSON(ctx, req)
			case "application/x-www-form-urlencoded", "multipart/form-data":
				return ReadForm(ctx, req)
			default:
				if ctx.Request().URL.RawQuery == "" {
					return ReadForm(ctx, req)
				} else {
					return ReadQuery(ctx, req)
				}
			}
		}
	)

	return logicProcessing(ctx, &req, &resp, read, fn)
}

// logicProcessing 逻辑处理
func logicProcessing[A, B any](ctx Context, req *A, resp *B,
	read func(ctx Context, req *A) error,
	fn func(ctx context.Context, req *A, resp *B) error) error {

	if e := read(ctx, req); e != nil {
		Fail(ctx, Param_Invalid)
		log.Error(e)
		return e
	}
	if e := Validate(req); e != nil {
		Fail(ctx, e)
		log.Error(e)
		return e
	}

	if e := Visit(ctx, req); e != nil {
		Fail(ctx, e)
		log.Error(e)
		return e
	}

	_ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	if e := fn(_ctx, req, resp); e != nil {
		Fail(ctx, e)
		log.Error(e)
		return e
	}

	return Ok(ctx, resp)
}

func GetWithCache[A, B any](
	ctx Context,
	fn func(ctx context.Context, req *A, resp *B) error,
) error {
	var (
		req  A
		resp B
		read = func(ctx Context, req *A) error {
			Initiate(ctx, req)

			if ctx.Request().URL.RawQuery == "" {
				return ReadForm(ctx, req)
			} else {
				return ReadQuery(ctx, req)
			}
		}
	)

	return logicProcessing(ctx, &req, &resp, read, fn)
}

func Redirect[A any](ctx Context, fn func(ctx context.Context, req *A, resp *string) error) error {
	var (
		req  A
		resp string
		read = func(ctx Context, req *A) error {
			Initiate(ctx, req)
			if ctx.Request().URL.RawQuery == "" {
				return ReadForm(ctx, req)
			} else {
				return ReadQuery(ctx, req)
			}
		}
	)
	if e := read(ctx, &req); e != nil {
		log.Error(e)
		return e
	}
	if e := Validate(&req); e != nil {
		log.Error(e)
		return e
	}

	_ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	if e := fn(_ctx, &req, &resp); e != nil {
		log.Error(e)
		return e
	}
	ctx.Redirect(resp)
	return nil
}

func Stream[A, B any](
	ctx Context,
	fn func(ctx Context, req *A, resp *B) error,
) error {
	var (
		req  A
		resp B
		read = func(ctx Context, req *A) error {
			Initiate(ctx, req)

			switch ctx.GetContentTypeRequested() {
			case "application/json":
				return ReadJSON(ctx, req)
			case "application/x-www-form-urlencoded", "multipart/form-data":
				return ReadForm(ctx, req)
			default:
				if ctx.Request().URL.RawQuery == "" {
					return ReadForm(ctx, req)
				} else {
					return ReadQuery(ctx, req)
				}
			}
		}
	)
	if e := read(ctx, &req); e != nil {
		Fail(ctx, Param_Invalid)
		log.Error(e)
		return e
	}
	if e := Validate(&req); e != nil {
		Fail(ctx, e)
		log.Error(e)
		return e
	}

	if e := Visit(ctx, &req); e != nil {
		Fail(ctx, e)
		log.Error(e)
		return e
	}

	if e := fn(ctx, &req, &resp); e != nil {
		Fail(ctx, e)
		log.Error(e)
		return e
	}

	return Ok(ctx, &resp)
}
