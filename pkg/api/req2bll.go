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
	"github.com/linbaozhong/gentity/pkg/api/iface"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

var (
	Param_Invalid = types.NewError(620, "参数无效")
	UnKnown       = types.NewError(610, "未知错误")
)

func Post[A, B any](
	ctx Context,
	fn func(ctx context.Context, req *A, resp *B) error,
) {
	var req A
	if e := ctx.ReadJSON(&req); e != nil {
		log.Error(e)
		Fail(ctx, Param_Invalid)
		return
	}
	if e := iface.Validate(&req); e != nil {
		Fail(ctx, e)
		return
	}

	var resp B
	if e := fn(ctx, &req, &resp); e != nil {
		log.Error(e)
		Fail(ctx, e)
		return
	}
	Ok(ctx, resp)
}

func Get[A, B any](
	ctx Context,
	fn func(ctx context.Context, req *A, resp *B) error,
) {
	var req A
	if e := ctx.ReadQuery(&req); e != nil {
		log.Error(e)
		Fail(ctx, Param_Invalid)
		return
	}
	if e := iface.Validate(&req); e != nil {
		Fail(ctx, e)
		return
	}

	var resp B
	if e := fn(ctx, &req, &resp); e != nil {
		log.Error(e)
		Fail(ctx, e)
		return
	}
	Ok(ctx, resp)
}
