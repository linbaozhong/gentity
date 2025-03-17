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
	"encoding/json"
	"github.com/r3labs/sse/v2"
	"log"
	"sync"
)

// ServerManager 管理 SSE 服务器和客户端组
type ServerManager struct {
	server *sse.Server
	groups map[string]*ClientGroup
	mu     sync.Mutex
}

// NewServerManager 创建一个新的服务器管理器
func NewServerManager() *ServerManager {
	srv := sse.New()
	srv.CreateStream("broadcast")
	return &ServerManager{
		server: srv,
		groups: make(map[string]*ClientGroup),
	}
}

// GetGroup 获取指定名称的客户端组，如果不存在则创建
func (m *ServerManager) GetGroup(groupName string) *ClientGroup {
	m.mu.Lock()
	defer m.mu.Unlock()
	if group, ok := m.groups[groupName]; ok {
		return group
	}
	group := NewClientGroup()
	m.groups[groupName] = group
	return group
}

func (m *ServerManager) Subscribe(theme string) *sse.Stream {
	return m.server.CreateStream(theme)
}

// Broadcast 向所有客户端广播消息
func (m *ServerManager) Broadcast(eventType []byte, data any) {
	message, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	m.server.Publish("broadcast", &sse.Event{
		Event: eventType,
		Data:  message,
	})
}

// SendToClient 向单个客户端发送消息
func (m *ServerManager) SendToClient(clientID string, eventType []byte, data any) {
	message, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	m.server.Publish(clientID, &sse.Event{
		Event: eventType,
		Data:  message,
	})
}

// SendToGroup 向指定组的所有客户端发送消息
func (m *ServerManager) SendToGroup(groupName string, eventType []byte, data interface{}) {
	group := m.GetGroup(groupName)
	if group == nil {
		return
	}
	message, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}
	group.mu.Lock()
	defer group.mu.Unlock()
	for _, clientID := range group.clientIDs {
		m.server.Publish(clientID, &sse.Event{
			Event: eventType,
			Data:  message,
		})
	}
}
