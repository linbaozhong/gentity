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

package sse

import (
	"github.com/linbaozhong/gentity/pkg/api"
	"github.com/linbaozhong/sse/v2"
)

var (
	_sseServer *sse.Server
)

func Start() error {
	_sseServer = sse.New()
	return nil
}

func Close() {
	if _sseServer == nil {
		return
	}
	_sseServer.Close()
}

// ServeHTTP 服务端推送
// streamID: 流ID
// lastEventId: 上次的event id
func ServeHTTP(ctx api.Context, streamID, lastEventId string) {
	if _sseServer == nil {
		Start()
	}

	ctx.Header("Access-Control-Allow-Origin", "*")

	r := ctx.Request()
	query := r.URL.Query()
	query.Set(sse.StreamKey, streamID)
	query.Set(sse.LastEventIdKey, lastEventId)
	r.URL.RawQuery = query.Encode()

	_sseServer.ServeHTTP(ctx.ResponseWriter(), r)
}
