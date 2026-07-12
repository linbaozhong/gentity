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
)

// Get handles a GET request: reads query parameters, calls service, writes response.
func Get[A, B any](
	ctx Context,
	callService func(context.Context, *A, *B) error,
	after ...func(Context, *B) error,
) error {
	var req A
	var resp B
	_, e := serviceContext(ctx, &req, &resp, readGetRequest[A], callService)
	if e != nil {
		return Fail(ctx, e)
	}
	for _, fn := range after {
		if e := fn(ctx, &resp); e != nil {
			return e
		}
	}
	return Ok(ctx, &resp)
}

// Post handles a POST request: reads body (JSON/form), calls service, writes response.
//
// Content-Type: application/json → json tag
// Content-Type: application/x-www-form-urlencoded → form tag
// Content-Type: multipart/form-data → form tag
func Post[A, B any](
	ctx Context,
	callService func(context.Context, *A, *B) error,
	after ...func(Context, *B) error,
) error {
	var req A
	var resp B
	_, e := serviceContext(ctx, &req, &resp, readPostRequest[A], callService)
	if e != nil {
		return Fail(ctx, e)
	}
	for _, fn := range after {
		if e := fn(ctx, &resp); e != nil {
			return e
		}
	}
	return Ok(ctx, &resp)
}

// Redirect performs an HTTP redirect after calling service.
func Redirect[A any](ctx Context,
	callService func(context.Context, *A, *string) error,
) error {
	var req A
	var resp string
	_, e := serviceContext(ctx, &req, &resp, readPostRequest[A], callService)
	if e != nil {
		return Fail(ctx, e)
	}
	ctx.Redirect(resp)
	return nil
}

// Stream handles a streaming response (no timeout).
func Stream[A, B any](
	ctx Context,
	callService func(context.Context, *A, *B) error,
	after ...func(Context, *B) error,
) error {
	var req A
	var resp B
	_, e := service(ctx, &req, &resp, readPostRequest[A], callService)
	if e != nil {
		return Fail(ctx, e)
	}
	for _, fn := range after {
		if e := fn(ctx, &resp); e != nil {
			return e
		}
	}
	return Ok(ctx, &resp)
}
