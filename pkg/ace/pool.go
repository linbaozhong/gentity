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

package ace

import (
	"github.com/linbaozhong/gentity/pkg/ace/types"
	"sync"
	"time"
)

const timeout = time.Minute

// Pool 是一个通用的对象池
type Pool struct {
	mu sync.RWMutex
	// 存储对象的map，key是对象的UUID
	pool map[uint64]element
	// 存储对象的key的切片
	keys []uint64

	timeout time.Duration
	once    sync.Once

	New func() types.AceModeler

	cleanCh <-chan time.Time // 用于触发清理的通道
	doneCh  chan bool        // 用于停止清理循环的通道
}

type element struct {
	obj types.AceModeler
	t   time.Time
}

func (p *Pool) init() {
	p.once.Do(func() {
		p.pool = make(map[uint64]element)
		p.keys = make([]uint64, 0)
		p.timeout = timeout
		p.cleanCh = time.After(p.timeout)
		p.doneCh = make(chan bool)

		go p.startCleaner()
	})
}

func (p Pool) new() types.AceModeler {
	obj := p.New()
	obj.UUID() // 池元素唯一标识
	return obj
}

// Get 从池中获取一个对象
func (p *Pool) Get() types.AceModeler {
	p.mu.RLock()
	defer p.mu.RUnlock()

	p.init()

	l := len(p.keys)
	if l > 0 {
		// 取出最后一个key元素
		k := p.keys[l-1]
		// 从 pool 中取出对象并返回
		if obj, ok := p.pool[k]; ok {
			delete(p.pool, k)
			p.keys = p.keys[:l-1]
			return obj.obj
		}
	}
	// 如果池中没有对象，可以在这里创建并返回一个新的对象
	if p.New != nil {
		return p.new()
	}
	return nil
}

// Put 将一个对象放回池中
func (p *Pool) Put(obj types.AceModeler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.init()

	k := obj.UUID()
	// 如果对象已经存在，则直接返回
	if _, ok := p.pool[k]; ok {
		return
	}
	// 将对象放入pool中
	p.keys = append(p.keys, k)
	p.pool[k] = element{
		obj: obj,
		t:   time.Now().Add(p.timeout),
	}
	// 如果清理器已经停止，重新启动它
	if p.cleanCh == nil {
		p.cleanCh = time.After(p.timeout)
	}
}

// Len 返回栈中元素的数量
func (p *Pool) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.pool)
}

// Clean 清理长时间未使用的对象
func (p *Pool) cleanup() {
	if !p.mu.TryLock() {
		return
	}
	defer p.mu.Unlock()

	now := time.Now()
	i := 0

	for k, obj := range p.pool {
		if obj.t.Before(now) {
			delete(p.pool, k)
			continue
		}
		p.keys[i] = k
		i++
	}
	p.keys = p.keys[:i]
}

// startCleaner 启动协程定期清理超时元素
func (p *Pool) startCleaner() {
	for {
		select {
		case <-p.cleanCh:
			p.cleanup()
			if len(p.pool) > 0 {
				p.cleanCh = time.After(p.timeout)
			} else {
				p.cleanCh = nil
			}
		case <-p.doneCh:
			p.cleanCh = nil
			return
		}
	}
}

func (p *Pool) Stop() {
	p.doneCh <- true
}
