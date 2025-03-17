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
	"sync"
)

// ClientGroup 管理客户端组
type ClientGroup struct {
	clientIDs []string
	mu        sync.Mutex
}

// NewClientGroup 创建一个新的客户端组
func NewClientGroup() *ClientGroup {
	return &ClientGroup{
		clientIDs: make([]string, 0),
	}
}

// AddClient 向组中添加客户端 ID
func (g *ClientGroup) AddClient(clientID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, id := range g.clientIDs {
		if id == clientID {
			return
		}
	}
	g.clientIDs = append(g.clientIDs, clientID)
}

// RemoveClient 从组中移除客户端 ID
func (g *ClientGroup) RemoveClient(clientID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i, id := range g.clientIDs {
		if id == clientID {
			g.clientIDs = append(g.clientIDs[:i], g.clientIDs[i+1:]...)
			break
		}
	}
}
