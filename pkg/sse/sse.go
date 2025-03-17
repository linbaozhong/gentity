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
	"github.com/r3labs/sse/v2"
	"net/http"
)

var (
	manager *ServerManager
)

func Start() error {
	manager = NewServerManager()
	return nil
}

func Publish(theme string) *sse.Stream {
	return manager.server.CreateStream(theme)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	manager.server.ServeHTTP(w, r)
}
