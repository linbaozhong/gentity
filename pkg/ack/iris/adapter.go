// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package ack

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/linbaozhong/gentity/pkg/ack/core"
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

func (a *ctxAdapter) ReadJSON(ptr any) error {
	body, err := io.ReadAll(a.c.Request().Body)
	if err != nil {
		return err
	}
	defer a.c.Request().Body.Close()
	if len(body) == 0 {
		return errors.New("请求体为空")
	}
	if x, ok := ptr.(json.Unmarshaler); ok {
		return x.UnmarshalJSON(body)
	}
	return json.Unmarshal(body, ptr)
}

func (a *ctxAdapter) ReadForm(ptr any) error {
	if uv, ok := ptr.(core.UnmarshalValueser); ok {
		values := a.c.FormValues()
		if len(values) == 0 {
			return nil
		}
		return uv.UnmarshalValues(values)
	}
	e := a.c.ReadForm(ptr)
	if e != nil && !iris.IsErrPath(e) {
		return e
	}
	return nil
}

func (a *ctxAdapter) ReadQuery(ptr any) error {
	if uv, ok := ptr.(core.UnmarshalValueser); ok {
		values := a.c.Request().URL.Query()
		if len(values) == 0 {
			return nil
		}
		return uv.UnmarshalValues(values)
	}
	return a.c.ReadQuery(ptr)
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
