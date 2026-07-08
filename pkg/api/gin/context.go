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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

type (
	Application = *gin.Engine
	Context     = *gin.Context
	Party       = *gin.RouterGroup
	Handler     = gin.HandlerFunc
	// ErrorHandler interface {
	// 	HandleContextError(ctx Context, err error)
	// }
	// // ErrorHandlerFunc a function shortcut for ErrorHandler interface.
	// ErrorHandlerFunc func(ctx Context, err error)
)

func NewApplication(name, version string) Application {
	app := gin.New()
	// 中间件
	app.Use(Logger(), Recovery())

	// 调试服务
	app.GET("/", debug(name, version))
	app.HEAD("/", debug(name, version))
	// 错误处理
	app.NoMethod(NoMethod)
	app.NoRoute(NoRoute)

	return app
}

func NewParty(app gin.IRouter, relativePath string) Party {
	return app.Group(relativePath)
}

// Server 封装 gin 的 HTTP 服务器，提供统一的 Run/Shutdown 生命周期接口。
// 与 iris 包的 Server 接口一致，切换框架只需改动 import 路径。
type Server struct {
	srv *http.Server
}

// NewServer 创建一个 Server 实例。
func NewServer(app Application, addr string) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: app,
		},
	}
}

// Run 启动服务（阻塞直到 Shutdown 被调用）。
func (s *Server) Run() error {
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown 优雅关闭服务，等待正在处理的请求完成或超时。
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func Logger() Handler {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// HtmlView 加载 HTML 模板用于网页渲染。
// 与 iris 包同名函数签名一致，返回类型适配 gin（*template.Template）。
// 注意: reload 参数保留用于 API 统一，gin 不内置热加载。
//
// 使用方式: app.SetHTMLTemplate(HtmlView("views", ".html", false))
func HtmlView(dir, extension string, reload bool) *template.Template {
	_ = reload
	return template.Must(template.ParseGlob(filepath.Join(dir, "*"+extension)))
}

// StaticWeb 设置静态文件服务
// urlPath: URL 访问路径，如 "/static"
// dir: 静态文件目录，如 "./public"
func StaticWeb(party Party, urlPath, dir string) {
	party.Static(urlPath, dir)
}

func debug(name, version string) Handler {
	return func(c Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    name,
			"version": version,
			"time":    time.Now().Format(time.DateTime),
		})
		return
	}
}

func NoRoute(c Context) {
	Fail(c, types.NewError(http.StatusMethodNotAllowed, "方法不允许"))
}

func NoMethod(c Context) {
	Fail(c, types.NewError(http.StatusNotFound, "方法未找到"))
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
