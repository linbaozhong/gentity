// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package ack

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/linbaozhong/gentity/pkg/ack/internal/core"

	"github.com/gin-gonic/gin"
)

// ctxAdapter adapts gin.Context to core.Context interface.
type ctxAdapter struct {
	c *gin.Context
}

func adapt(c *gin.Context) core.Context {
	return &ctxAdapter{c: c}
}

// === 请求信息 ===

func (a *ctxAdapter) Path() string                { return a.c.Request.URL.Path }
func (a *ctxAdapter) Method() string              { return a.c.Request.Method }
func (a *ctxAdapter) RemoteAddr() string          { return a.c.RemoteIP() }
func (a *ctxAdapter) GetHeader(key string) string { return a.c.GetHeader(key) }
func (a *ctxAdapter) Param(key string) string     { return a.c.Param(key) }

// === 请求对象 ===

func (a *ctxAdapter) Request() *core.HttpRequest {
	r := a.c.Request
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
// 否则回退到 gin 自身的路径参数绑定。
func (a *ctxAdapter) readPathParams(ptr any) error {
	if uv, ok := ptr.(core.UnmarshalValueser); ok {
		if params := a.c.Params; len(params) > 0 {
			values := make(map[string][]string, len(params))
			for _, p := range params {
				values[p.Key] = []string{p.Value}
			}
			return uv.UnmarshalValues(values)
		}
		return nil
	}

	if e := a.c.BindUri(ptr); e != nil {
		return e
	}
	return nil
}

func (a *ctxAdapter) ReadJSON(ptr any) error {
	defer a.c.Request.Body.Close()
	body, e := io.ReadAll(a.c.Request.Body)
	if e != nil {
		return e
	}
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
	a.c.Request.ParseForm()
	if uv, ok := ptr.(core.UnmarshalValueser); ok {
		if values := a.c.Request.Form; len(values) > 0 {
			if e := uv.UnmarshalValues(values); e != nil {
				return e
			}
		}

		// form 已解析，再合并路径参数
		return a.readPathParams(ptr)
	}

	if e := a.c.ShouldBind(ptr); e != nil {
		return e
	}
	return a.readPathParams(ptr)
}

func (a *ctxAdapter) ReadQuery(ptr any) error {
	if uv, ok := ptr.(core.UnmarshalValueser); ok {
		if values := a.c.Request.URL.Query(); len(values) > 0 {
			if e := uv.UnmarshalValues(values); e != nil {
				return e
			}
		}
		// query 已解析，再合并路径参数
		return a.readPathParams(ptr)
	}

	if e := a.c.ShouldBindQuery(ptr); e != nil {
		return e
	}
	return a.readPathParams(ptr)
}

func (a *ctxAdapter) ContentType() string             { return a.c.ContentType() }
func (a *ctxAdapter) FormValues() map[string][]string { return a.c.Request.PostForm }

// === 响应写入 ===

func (a *ctxAdapter) JSON(v any) error                 { a.c.JSON(200, v); return nil }
func (a *ctxAdapter) StatusCode(status int)            { a.c.Status(status) }
func (a *ctxAdapter) Header(key, value string)         { a.c.Header(key, value) }
func (a *ctxAdapter) SetContentType(ct string)         { a.c.Header("Content-Type", ct) }
func (a *ctxAdapter) Redirect(url string)              { a.c.Redirect(302, url) }
func (a *ctxAdapter) SendFile(path, name string) error { a.c.FileAttachment(path, name); return nil }

// === 流式响应 ===

func (a *ctxAdapter) ResponseWriter() io.Writer { return a.c.Writer }

// === 中断控制 ===

func (a *ctxAdapter) StopWithStatus(code int) { a.c.AbortWithStatus(code) }

// === 上下文值存储 ===

func (a *ctxAdapter) Values() core.Map           { return &ginMapAdapter{c: a.c} }
func (a *ctxAdapter) Set(key string, val any)    { a.c.Set(key, val) }
func (a *ctxAdapter) Get(key string) (any, bool) { return a.c.Get(key) }

// === context.Context ===

func (a *ctxAdapter) Deadline() (time.Time, bool) { return a.c.Deadline() }
func (a *ctxAdapter) Done() <-chan struct{}       { return a.c.Done() }
func (a *ctxAdapter) Err() error                  { return a.c.Err() }
func (a *ctxAdapter) Value(key any) any           { return a.c.Value(key) }

// === Middleware ===

func (a *ctxAdapter) Next() { a.c.Next() }

// ginMapAdapter adapts gin.Context to core.Map interface.
type ginMapAdapter struct {
	c *gin.Context
}

func (m *ginMapAdapter) Set(key string, val any) { m.c.Set(key, val) }
func (m *ginMapAdapter) Get(key string) any {
	v, _ := m.c.Get(key)
	return v
}
