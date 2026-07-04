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
	"github.com/gin-gonic/gin"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
	"net/http"
	"time"
)

type (
	Application  = *gin.Engine
	Context      = *gin.Context
	Party        = *gin.RouterGroup
	Handler      = gin.HandlerFunc
	ErrorHandler interface {
		HandleContextError(ctx Context, err error)
	}
	// ErrorHandlerFunc a function shortcut for ErrorHandler interface.
	ErrorHandlerFunc func(ctx Context, err error)
)

func NewApplication(name, version string) Application {
	app := gin.New()
	// .Configure(gin.WithRemoteAddrHeader(
	// 	"X-Forwarded-For",
	// ))
	// 中间件
	app.Use(Recovery())
	app.Use(Logger())

	// 调试服务
	app.GET("/", debug(name, version))
	app.HEAD("/", debug(name, version))
	// 错误处理
	app.NoMethod(NoMethod)
	app.NoRoute(NoRoute)

	return app
}

func NewParty(app Party, relativePath string) Party {
	return app.Group(relativePath)
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
