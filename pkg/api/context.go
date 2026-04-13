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
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/linbaozhong/gentity/pkg/types"
	"net/http"
	"time"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/view"
)

type (
	Application  = *iris.Application
	Context      = iris.Context
	Party        = iris.Party
	Handler      = iris.Handler
	ErrorHandler interface {
		HandleContextError(ctx *Context, err error)
	}
	// ErrorHandlerFunc a function shortcut for ErrorHandler interface.
	ErrorHandlerFunc func(ctx *Context, err error)
)

func NewApplication(name, version string) Application {
	app := iris.New().Configure(iris.WithRemoteAddrHeader(
		"X-Forwarded-For",
	))
	// 中间件
	app.Use(Recovery())
	app.Use(Logger())

	// 调试服务
	app.Get("/", debug(name, version))
	app.Head("/", debug(name, version))
	// 错误处理
	app.OnErrorCode(iris.StatusNotFound, NoMethod)
	app.OnErrorCode(iris.StatusMethodNotAllowed, NoRoute)

	return app
}

func OnInterrupt(fn func()) {
	iris.RegisterOnInterrupt(fn)
}

func NewParty(app Party, relativePath string) Party {
	return app.Party(relativePath)
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
func GetHtmlView(dir, extension string, reload bool) *view.HTMLEngine {
	return iris.HTML(dir, extension).Reload(reload)
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
					Fail(c, types.NewError(http.StatusInternalServerError, "内部服务器错误").
						SetOp("Recovery").Join(err))
					return
				}
				Fail(c, types.NewError(http.StatusInternalServerError, "内部服务器错误").
					Join(fmt.Errorf("%v", e)))
			}
		}()
		c.Next()
	}
}
