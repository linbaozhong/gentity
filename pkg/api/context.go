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
	"github.com/kataras/iris/v12"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/view"
)

type (
	Application = *iris.Application
	Context     = iris.Context
	Party       = iris.Party
	Handler     = iris.Handler
)

var (
	// 全局context，用于支持外部调用。最好在程序启动时引用，否则可能造成panic。
	// 在程序退出时，需要调用Cancel()。
	AppContext, AppCancel = context.WithCancel(context.Background())
)

func NewApplication(name, version string) Application {
	app := iris.New().Configure(iris.WithRemoteAddrHeader(
		"X-Forwarded-For",
	))

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
			"app_name":    name,
			"app_version": version,
		})
		return
	}
}

func NoRoute(c Context) {
	Fail(c, types.NewError(iris.StatusMethodNotAllowed, "方法不允许"))
}

func NoMethod(c Context) {
	Fail(c, types.NewError(iris.StatusNotFound, "方法未找到"))
}

func Recovery() Handler {
	return func(c Context) {
		defer func() {
			if e := recover(); e != nil {
				log.Panic(e)
				Fail(c, types.NewError(iris.StatusInternalServerError, "内部服务器错误"))
			}
		}()
		c.Next()
	}
}
