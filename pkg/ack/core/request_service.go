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

package core

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/linbaozhong/gentity/pkg/types"
)

// Initiate 初始化请求上下文，设置 IP、UserAgent、Authorization 和 OperationID。
func Initiate(ctx Context, arg any) {
	ctx.Values().Set(IpKey, ctx.RemoteAddr())
	ctx.Values().Set(UserAgent, ctx.Request().UserAgent())
	ctx.Values().Set(Authorization, ctx.GetHeader(Authorization))

	id := ctx.GetHeader(OperationID)
	if len(id) == 0 {
		ctx.Values().Set(OperationID, fmtUnixMilli())
	} else {
		ctx.Values().Set(OperationID, id)
	}

	if ier, ok := arg.(Initializer); ok {
		ier.Init()
	}
}

// Validate 校验参数，如果 arg 实现了 Checker 接口则调用 Check 方法。
func Validate(arg any) error {
	if checker, ok := arg.(Checker); ok {
		return checker.Check()
	}
	return nil
}

func fmtUnixMilli() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}

// GetResult handles a GET request and returns result data.
func GetResult[A, B any](ctx Context,
	callService func(context.Context, *A, *B) error) (*B, error) {
	var req A
	var resp B
	return serviceContext(ctx, &req, &resp, readGetRequest[A], callService)
}

// PostResult handles a POST request and returns result data.
func PostResult[A, B any](ctx Context,
	callService func(context.Context, *A, *B) error) (*B, error) {
	var req A
	var resp B
	return serviceContext(ctx, &req, &resp, readPostRequest[A], callService)
}

// StreamResult handles a streaming POST request and returns result data.
func StreamResult[A, B any](ctx Context,
	callService func(context.Context, *A, *B) error) (*B, error) {
	var req A
	var resp B
	return service(ctx, &req, &resp, readPostRequest[A], callService)
}

// serviceContext executes business logic with timeout support.
func serviceContext[A, B any](ctx Context, req *A, resp *B,
	read func(Context, *A) error,
	callService func(context.Context, *A, *B) error) (*B, error) {

	defer func() {
		if e := recover(); e != nil {
			var err error
			if eerr, ok := e.(error); ok {
				err = eerr
			} else {
				err = fmt.Errorf("%v", e)
			}
			Fail(ctx, types.NewError(http.StatusInternalServerError, "内部服务器错误").Join(err))
		}
	}()

	if e := read(ctx, req); e != nil {
		return resp, types.NewError(http.StatusBadRequest, "反序列化参数错误", "serviceContext.read").Join(e)
	}
	if e := Validate(req); e != nil {
		return resp, e
	}

	_ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	if e := callService(_ctx, req, resp); e != nil {
		return resp, e
	}

	return resp, nil
}

// service executes business logic without timeout.
func service[A, B any](ctx Context, req *A, resp *B,
	read func(Context, *A) error,
	callService func(context.Context, *A, *B) error) (*B, error) {

	defer func() {
		if e := recover(); e != nil {
			var err error
			if eerr, ok := e.(error); ok {
				err = eerr
			} else {
				err = fmt.Errorf("%v", e)
			}
			Fail(ctx, types.NewError(http.StatusInternalServerError, "内部服务器错误").Join(err))
		}
	}()

	if e := read(ctx, req); e != nil {
		return resp, types.NewError(http.StatusBadRequest, "反序列化参数错误", "service.read").Join(e)
	}
	if e := Validate(req); e != nil {
		return resp, e
	}

	if e := callService(ctx, req, resp); e != nil {
		return resp, e
	}

	return resp, nil
}

// readGetRequest 读取 GET 请求参数
func readGetRequest[A any](ctx Context, req *A) error {
	Initiate(ctx, req)
	if ctx.Request().URL.RawQuery == "" {
		return ctx.ReadForm(req)
	}
	return ctx.ReadQuery(req)
}

// readPostRequest 读取 POST 请求参数
func readPostRequest[A any](ctx Context, req *A) error {
	Initiate(ctx, req)
	switch ctx.ContentType() {
	case "application/json":
		return ctx.ReadJSON(req)
	case "application/x-www-form-urlencoded", "multipart/form-data":
		return ctx.ReadForm(req)
	default:
		if ctx.Request().URL.RawQuery == "" {
			return ctx.ReadForm(req)
		}
		return ctx.ReadQuery(req)
	}
}
