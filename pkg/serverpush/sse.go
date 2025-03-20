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

package serverpush

import (
	"github.com/linbaozhong/gentity/pkg/api"
	"github.com/linbaozhong/sse/v2"
)

var (
	_sseServer *sse.Server
)

// option 服务端推送配置
type option func(*sse.Server)

// WithAutoStream 是否自动创建流, 即当客户端连接后是否自动创建新的流
func WithAutoStream(autoStream bool) option {
	return func(s *sse.Server) {
		s.AutoStream = autoStream
	}
}

// WithAutoReplay 是否自动重放, 即当客户端断开并重新连接后是否自动重放事件
func WithAutoReplay(autoReplay bool) option {
	return func(s *sse.Server) {
		s.AutoReplay = autoReplay
	}
}

// Start 启动服务
// 初始化sse服务
func Start(opts ...option) error {
	_sseServer = sse.New()
	for _, opt := range opts {
		opt(_sseServer)
	}
	return nil
}

func Close() {
	if _sseServer == nil {
		return
	}
	_sseServer.Close()
}

// CreateStream 创建流
// streamID: 流ID
func CreateStream(streamID string) {
	if _sseServer == nil {
		Start(WithAutoStream(true), WithAutoReplay(true))
	}
	_sseServer.CreateStream(streamID)
}

// Push 推送事件
// streamID: 流ID
// event: 事件
func Push(streamID string, event *sse.Event) {
	if _sseServer == nil {
		Start(WithAutoStream(true), WithAutoReplay(true))
	}
	_sseServer.Publish(streamID, event)
}

// Boardcast 广播事件
// event: 事件
func Boardcast(event *sse.Event) {
	if _sseServer == nil {
		Start(WithAutoStream(true), WithAutoReplay(true))
	}
	_sseServer.Boardcast(event)
}

// ServeHTTP 服务端推送
// streamID: 流ID
// lastEventId: 上次的event id
func ServeHTTP(ctx api.Context, streamID, lastEventId string) {
	if _sseServer == nil {
		Start(WithAutoStream(true), WithAutoReplay(true))
	}

	ctx.Header("Access-Control-Allow-Origin", "*")

	r := ctx.Request()
	query := r.URL.Query()
	query.Set(sse.StreamKey, streamID)
	query.Set(sse.LastEventIdKey, lastEventId)
	r.URL.RawQuery = query.Encode()

	_sseServer.ServeHTTP(ctx.ResponseWriter(), r)
}
