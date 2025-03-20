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

package handler

import (
	"github.com/linbaozhong/gentity/pkg/api"
	"github.com/linbaozhong/gentity/pkg/serverpush"
	"github.com/linbaozhong/gentity/pkg/token"
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

	_tk := values.Get("tk")
	if _tk != "" {
		_clientId, _, _ = token.GetIDAndTokenFromCipher(_tk)
	}

	_lastEventId := values.Get("event_id")
	go func() {
		time.Sleep(time.Second * 5)
		serverpush.Push(_clientId, &sse.Event{
			Data: []byte("hello world"),
		})
		serverpush.Push(_clientId, &sse.Event{
			Event: []byte("login"),
			Data:  []byte("welcome"),
		})
	}()
	serverpush.ServeHTTP(c, _clientId, _lastEventId)
}
