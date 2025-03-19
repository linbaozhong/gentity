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

type Payload struct {
	Token       string `json:"token"`
	LastEventId string `json:"lastEventId"`
}

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

func ServeHTTP(ctx api.Context, clientID, lastEventId string) {
	if _sseServer == nil {
		Start()
	}
	w := ctx.ResponseWriter()
	r := ctx.Request()

	query := r.URL.Query()
	query.Set("stream", clientID)
	r.URL.RawQuery = query.Encode()
	r.Header.Set("Last-Event-ID", lastEventId)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	_sseServer.ServeHTTP(w, r)
}
