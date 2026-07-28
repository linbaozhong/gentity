// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package ack

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/linbaozhong/gentity/pkg/ack/internal/core"

	"github.com/kataras/iris/v12"
)

// ctxAdapter adapts iris.Context to core.Context interface.
type ctxAdapter struct {
	c iris.Context
}

func adapt(c iris.Context) core.Context {
	return &ctxAdapter{c: c}
}

// === 请求信息 ===

func (a *ctxAdapter) Path() string                { return a.c.Path() }
func (a *ctxAdapter) Method() string              { return a.c.Method() }
func (a *ctxAdapter) RemoteAddr() string          { return a.c.RemoteAddr() }
func (a *ctxAdapter) GetHeader(key string) string { return a.c.GetHeader(key) }
func (a *ctxAdapter) Param(key string) string     { return a.c.Params().Get(key) }

// === 请求对象 ===

func (a *ctxAdapter) Request() *core.HttpRequest {
	r := a.c.Request()
	return &core.HttpRequest{
		Body:   r.Body,
		Method: r.Method,
		URL: &core.HttpURL{
			Path:     r.URL.Path,
			RawQuery: r.URL.RawQuery,
			Values:   r.URL.Query(),
		},
		UserAgent: r.UserAgent,
	}
}

// === 请求数据读取 ===

// readPathParams 把路由路径参数补充绑定到 ptr。
// 实现了 UnmarshalValueser 时按 json 字段名匹配（与生成的 UnmarshalValues 约定一致），
// 否则回退到 iris 自身的路径参数绑定。
func (a *ctxAdapter) readPathParams(ptr any) error {
	if uv, ok := ptr.(core.UnmarshalValueser); ok {
		if params := a.c.Params(); params.Len() > 0 {
			values := make(map[string][]string, params.Len())
			params.Visit(func(key string, value string) {
				values[key] = []string{value}
			})
			return uv.UnmarshalValues(values)
		}
		return nil
	}

	if e := a.c.ReadParams(ptr); e != nil && !iris.IsErrPath(e) {
		return e
	}
	return nil
}

func (a *ctxAdapter) ReadJSON(ptr any) error {
	body, e := io.ReadAll(a.c.Request().Body)
	if e != nil {
		return e
	}
	defer a.c.Request().Body.Close()

	if len(body) == 0 {
		return errors.New("请求体为空")
	}

	if x, ok := ptr.(json.Unmarshaler); ok {
		if e := x.UnmarshalJSON(body); e != nil {
			return e
		}
	} else if e := json.Unmarshal(body, ptr); e != nil {
		return e
	}
	// body 已解析，再合并路径参数
	return a.readPathParams(ptr)
}

func (a *ctxAdapter) ReadForm(ptr any) error {
	if uv, ok := ptr.(core.UnmarshalValueser); ok {
		if fv := a.c.FormValues(); len(fv) > 0 {
			if e := uv.UnmarshalValues(fv); e != nil {
				return e
			}
		}
		// form 已解析，再合并路径参数
		return a.readPathParams(ptr)
	}
	if e := a.c.ReadForm(ptr); e != nil && !iris.IsErrPath(e) {
		return e
	}
	return a.readPathParams(ptr)
}

func (a *ctxAdapter) ReadQuery(ptr any) error {
	if uv, ok := ptr.(core.UnmarshalValueser); ok {
		if qv := a.c.Request().URL.Query(); len(qv) > 0 {
			if e := uv.UnmarshalValues(qv); e != nil {
				return e
			}
		}
		// query 已解析，再合并路径参数
		return a.readPathParams(ptr)
	}
	if e := a.c.ReadQuery(ptr); e != nil && !iris.IsErrPath(e) {
		return e
	}
	return a.readPathParams(ptr)
}

// iris: 获取 Content-Type 用 GetContentTypeRequested
func (a *ctxAdapter) ContentType() string             { return a.c.GetContentTypeRequested() }
func (a *ctxAdapter) FormValues() map[string][]string { return a.c.FormValues() }

// === 响应写入 ===

func (a *ctxAdapter) JSON(v any) error                 { return a.c.JSON(v) }
func (a *ctxAdapter) StatusCode(status int)            { a.c.StatusCode(status) }
func (a *ctxAdapter) Header(key, value string)         { a.c.Header(key, value) }
func (a *ctxAdapter) SetContentType(ct string)         { a.c.ContentType(ct) }
func (a *ctxAdapter) Redirect(url string)              { a.c.Redirect(url) }
func (a *ctxAdapter) SendFile(path, name string) error { return a.c.SendFile(path, name) }

// === 流式响应 ===

func (a *ctxAdapter) ResponseWriter() io.Writer { return a.c.ResponseWriter() }

// === 中断控制 ===

func (a *ctxAdapter) StopWithStatus(code int) { a.c.StopWithStatus(code) }

// === 上下文值存储 ===

func (a *ctxAdapter) Values() core.Map { return &irisMapAdapter{c: a.c} }
func (a *ctxAdapter) Set(key string, val any) {
	a.c.Values().Set(key, val)
}
func (a *ctxAdapter) Get(key string) (any, bool) {
	v := a.c.Values().Get(key)
	return v, v != nil
}

// === context.Context ===

func (a *ctxAdapter) Deadline() (time.Time, bool) { return a.c.Deadline() }
func (a *ctxAdapter) Done() <-chan struct{}       { return a.c.Done() }
func (a *ctxAdapter) Err() error                  { return a.c.Err() }
func (a *ctxAdapter) Value(key any) any           { return a.c.Value(key) }

// === Middleware ===

func (a *ctxAdapter) Next() { a.c.Next() }

// irisMapAdapter implements core.Map for iris.Values.
type irisMapAdapter struct {
	c iris.Context
}

func (ma *irisMapAdapter) Set(key string, val any) { ma.c.Values().Set(key, val) }
func (ma *irisMapAdapter) Get(key string) any      { return ma.c.Values().Get(key) }
