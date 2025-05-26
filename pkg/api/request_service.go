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

package api

import (
	"context"
	"time"
)

// GetResult 调用service处理get请求，并返回结果数据
func GetResult[A, B any](ctx Context,
	fn func(ctx context.Context, req *A, resp *B) error) (*B, error) {
	var (
		req  A
		resp B
	)

	return serviceContext(ctx, &req, &resp, readGetRequest[A], fn)
}

// PostResult 调用service处理post请求，并返回结果数据
func PostResult[A, B any](ctx Context,
	fn func(ctx context.Context, req *A, resp *B) error) (*B, error) {
	var (
		req  A
		resp B
	)
	return serviceContext(ctx, &req, &resp, readPostRequest[A], fn)
}

// StreamResult 调用service处理post请求，并返回结果数据
func StreamResult[A, B any](ctx Context,
	fn func(ctx Context, req *A, resp *B) error) (*B, error) {
	var (
		req  A
		resp B
	)
	return service(ctx, &req, &resp, readPostRequest[A], fn)
}

// serviceContext 逻辑处理,serviceContext会超时
func serviceContext[A, B any](ctx Context, req *A, resp *B,
	read func(ctx Context, req *A) error,
	fn func(ctx context.Context, req *A, resp *B) error) (*B, error) {

	if e := read(ctx, req); e != nil {
		return resp, Param_Invalid.SetInfo(e)
	}
	if e := Validate(req); e != nil {
		return resp, e
	}

	// if e := Visit(ctx, req); e != nil {
	// 	return resp, e
	// }

	_ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	if e := fn(_ctx, req, resp); e != nil {
		return resp, e
	}

	return resp, nil
}

// service 逻辑处理,service不会超时
func service[A, B any](ctx Context, req *A, resp *B,
	read func(ctx Context, req *A) error,
	fn func(ctx Context, req *A, resp *B) error) (*B, error) {

	if e := read(ctx, req); e != nil {
		return resp, Param_Invalid.SetInfo(e)
	}
	if e := Validate(req); e != nil {
		return resp, e
	}

	// if e := Visit(ctx, req); e != nil {
	// 	return resp, e
	// }

	if e := fn(ctx, req, resp); e != nil {
		return resp, e
	}

	return resp, nil
}

// readGetRequest 读取get请求
func readGetRequest[A any](ctx Context, req *A) error {
	Initiate(ctx, req)
	if ctx.Request().URL.RawQuery == "" {
		return ReadForm(ctx, req)
	} else {
		return ReadQuery(ctx, req)
	}
}

// readPostRequest 读取post请求
func readPostRequest[A any](ctx Context, req *A) error {
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
