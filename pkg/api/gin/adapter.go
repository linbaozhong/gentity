// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package api

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linbaozhong/gentity/pkg/api/core"
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

func (a *ctxAdapter) ReadJSON(ptr any) error {
	defer a.c.Request.Body.Close()
	body, err := io.ReadAll(a.c.Request.Body)
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return errors.New("请求体为空")
	}
	if x, ok := ptr.(json.Unmarshaler); ok {
		return x.UnmarshalJSON(body)
	}
	return json.Unmarshal(body, ptr)
}

func (a *ctxAdapter) ReadForm(ptr any) error {
	a.c.Request.ParseForm()
	if uv, ok := ptr.(core.UnmarshalValueser); ok {
		values := a.c.Request.Form
		if len(values) == 0 {
			return nil
		}
		return uv.UnmarshalValues(values)
	}
	e := a.c.ShouldBind(ptr)
	if e != nil {
		return e
	}
	return nil
}

func (a *ctxAdapter) ReadQuery(ptr any) error {
	if uv, ok := ptr.(core.UnmarshalValueser); ok {
		values := a.c.Request.URL.Query()
		if len(values) == 0 {
			return nil
		}
		return uv.UnmarshalValues(values)
	}
	return a.c.ShouldBindQuery(ptr)
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
