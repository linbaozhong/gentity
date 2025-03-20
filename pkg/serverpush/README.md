# SSE 使用说明

## 1. 启动服务
在系统启动时，需要启动SSE推送服务：
```go
// 启动推送服务
serverpush.Start(
    serverpush.WithAutoStream(true),
    serverpush.WithAutoReplay(true),
)
```

## 2. 关闭推送服务
在系统关闭前，请手动关闭推送服务：
```go
// 关闭推送服务
serverpush.Close()
```
## 3. 推送消息
```go
// 推送消息
serverpush.Push("test", "hello world")
```
## 4. 广播消息
```go
// 广播消息
serverpush.Broadcast("hello world")
```
## 提供前端连接的api接口
```go
package handler

import (
	"github.com/linbaozhong/gentity/pkg/api"
	"github.com/linbaozhong/gentity/pkg/serverpush"
	"github.com/linbaozhong/sse/v2"
	"time"
)

type sevent struct{}

func init() {
	api.RegisterRoute(&sevent{})
}

func (s *sevent) RegisterRoute(group api.Party) {
	g := api.NewParty(group, "/sse")
	g.Get("/connect", s.connect)
}

func (s *sevent) connect(c api.Context) {
	var _clientId string
	
	values := c.Request().URL.Query()
	// 从查询参数中获取名为 "tk" 的值，并从加密的令牌中解析出 ID 和令牌
	_tk := values.Get("tk")
	if _tk != "" {
		_clientId, _, _ = token.GetIDAndTokenFromCipher(_tk)
	}
	// 从查询参数中获取名为 "event_id" 的值，
	_lastEventId := values.Get("event_id")
	// 启动一个 goroutine，在 1 秒后向客户端模拟推送事件
	go func() {
		time.Sleep(time.Second)
		serverpush.Push(_clientId, &sse.Event{
			Data: []byte("hello world"),
		})
		serverpush.Push(_clientId, &sse.Event{
			Event: []byte("login"),
			Data:  []byte("welcome"),
		})
	}()
	// 调用 serverpush 包的 ServeHTTP 函数，处理 SSE 请求
	serverpush.ServeHTTP(c, _clientId, _lastEventId)
}

```