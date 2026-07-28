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

// Package core provides framework-agnostic HTTP handler utilities.
// It defines the Context interface that abstracts over gin and iris,
// and contains all shared business logic for request handling.
package core

import (
	"context"
	"io"
	"net/url"
)

// Context defines the HTTP context interface that abstracts framework-specific
// implementations (gin.Context, iris.Context).
// It embeds context.Context so it can be used wherever a context.Context is expected.
type Context interface {
	context.Context

	// === 请求信息 ===
	Path() string
	Method() string
	RemoteAddr() string
	GetHeader(key string) string
	// Param 读取路由参数（如 /users/:id 中的 id）
	Param(key string) string

	// === 请求对象 ===
	Request() *HttpRequest

	// === 请求数据读取（由框架 adapter 实现）===
	ReadJSON(ptr any) error
	ReadForm(ptr any) error
	ReadQuery(ptr any) error
	// ReadParams(ptr any) error
	ContentType() string
	FormValues() map[string][]string

	// === 响应写入 ===
	JSON(v any) error
	StatusCode(status int)
	Header(key, value string)
	SetContentType(contentType string)
	Redirect(url string)
	SendFile(path, name string) error

	// === 流式响应 ===
	ResponseWriter() io.Writer

	// === 中断控制 ===
	StopWithStatus(code int)

	// === 上下文值存储 ===
	Values() Map
	Set(key string, val any)
	Get(key string) (any, bool)

	// === Middleware ===
	Next()
}

// HttpRequest wraps http.Request with unified access patterns.
type HttpRequest struct {
	Body      io.ReadCloser
	URL       *HttpURL
	Method    string
	UserAgent func() string
}

// HttpURL wraps URL access for both frameworks.
type HttpURL struct {
	Path     string
	RawQuery string
	Values   url.Values
}

// Map is a key-value store interface (abstracts iris.Values / gin.Context.Get/Set)
type Map interface {
	Set(key string, val any)
	Get(key string) any
}
