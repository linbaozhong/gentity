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

package ack

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/linbaozhong/gentity/pkg/ack/core"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
	"net/http"
	"time"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/view"
)

type (
	Application = *iris.Application
	Context     = iris.Context
	Party       = iris.Party
	Handler     = iris.Handler
	// ErrorHandler interface {
	// 	HandleContextError(ctx *Context, err error)
	// }
	// // ErrorHandlerFunc a function shortcut for ErrorHandler interface.
	// ErrorHandlerFunc func(ctx *Context, err error)
)

func NewApplication(name, version string) Application {
	app := iris.New().Configure(iris.WithRemoteAddrHeader(
		"X-Forwarded-For",
	))
	// 中间件
	app.Use(Logger(), Recovery())

	// 调试服务
	app.Get("/", debug(name, version))
	app.Head("/", debug(name, version))
	// 错误处理
	app.OnErrorCode(iris.StatusNotFound, NoMethod)
	app.OnErrorCode(iris.StatusMethodNotAllowed, NoRoute)

	return app
}

func NewParty(app Party, relativePath string) Party {
	return app.Party(relativePath)
}

// Server 封装 iris 的 HTTP 服务器，提供统一的 Run/Shutdown 生命周期接口。
// 与 gin 包的 Server 接口一致，切换框架只需改动 import 路径。
type Server struct {
	app  *iris.Application
	addr string
}

// NewServer 创建一个 Server 实例。
func NewServer(app Application, addr string) *Server {
	return &Server{app: app, addr: addr}
}

// Run 启动服务（阻塞直到 Shutdown 被调用）。
func (s *Server) Run() error {
	return s.app.Listen(s.addr)
}

// Shutdown 优雅关闭服务，等待正在处理的请求完成或超时。
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.Shutdown(ctx)
}

func Logger() Handler {
	return logger.New(logger.Config{
		Status: true,
		IP:     true,
		Method: true,
		Path:   true,
		Query:  true,
	})
}

// HtmlView 设置 HTML 模板引擎
// dir: 模板文件目录，如 "./views"
// extension: 模板文件扩展名，如 ".html"
// reload: 是否实时加载模板文件
func HtmlView(dir, extension string, reload bool) *view.HTMLEngine {
	return iris.HTML(dir, extension).Reload(reload)
}

// StaticWeb 设置静态文件服务
// urlPath: URL 访问路径，如 "/static"
// dir: 静态文件目录，如 "./public"
func StaticWeb(party Party, urlPath, dir string) {
	party.HandleDir(urlPath, iris.Dir(dir))
}

func debug(name, version string) Handler {
	return func(c Context) {
		c.JSON(iris.Map{
			"name":    name,
			"version": version,
			"time":    time.Now().Format(time.DateTime),
		})
		return
	}
}

func NoRoute(c Context) {
	core.Fail(adapt(c), types.NewError(http.StatusMethodNotAllowed, "方法不允许"))
}

func NoMethod(c Context) {
	core.Fail(adapt(c), types.NewError(http.StatusNotFound, "方法未找到"))
}

func Recovery() Handler {
	return func(c Context) {
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(error); ok {
					log.Panic("Recovery", err)
					return
				}
				log.Panic("Recovery", fmt.Errorf("%v", e))
			}
		}()
		c.Next()
	}
}
