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

package pool

import (
	"reflect"
	"sync"
	"time"
)

const timeout = time.Minute

// Pool 是一个通用的对象池
type Pool struct {
	mu sync.Mutex
	// 存储对象的map，key是对象的类型，value是对象的切片
	pool    []*element
	timeout time.Duration

	once    sync.Once
	chclean chan bool
	ticker  *time.Ticker

	New func() any
}

type element struct {
	obj any
	t   time.Time
}

func (p *Pool) init() {
	p.once.Do(func() {
		p.timeout = timeout
		go p.startCleaner()
	})
}

// Get 从池中获取一个对象
func (p *Pool) Get() any {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 如果池中没有对象，可以在这里创建一个新的对象
	if len(p.pool) == 0 {
		return p.New()
	}

	el := p.pool[len(p.pool)-1]
	p.pool = p.pool[:len(p.pool)-1]

	p.chclean <- true

	return el.obj
}

// Put 将一个对象放回池中
func (p *Pool) Put(obj any) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.init()

	for _, el := range p.pool {
		if reflect.ValueOf(obj).Pointer() == reflect.ValueOf(el.obj).Pointer() {
			return
		}
	}

	el := &element{
		obj: obj,
		t:   time.Now().Add(p.timeout * 2),
	}
	p.pool = append(p.pool, el)

	p.chclean <- true
}

// Len 返回栈中元素的数量
func (p *Pool) Len() int {
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

	for _, obj := range p.pool {
		if obj.t.After(now) {
			p.pool[i] = obj
			i++
		}
	}
	p.pool = p.pool[:i]

	p.ticker.Reset(p.timeout)
}

// startCleaner 启动协程定期清理超时元素
func (p *Pool) startCleaner() {
	p.ticker = time.NewTicker(p.timeout)
	for {
		select {
		case <-p.ticker.C:
			p.cleanup()
		default:

		}
	}
}
