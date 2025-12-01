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
	"sync"
)

type (
	// 定义一个私有类型，用于禁止拷贝
	noCopy struct{}

	Model struct {
		_       noCopy
		mu      sync.Mutex
		ace_put bool // 内部留用，禁止外部赋值
	}

	PoolModeler interface {
		Reset()
		// Put2Pool 方法，用于判断是否已经放入池中
		// 如果已经放入池中，返回 true，否则返回 false
		// 并将 ace_put 设置为 true，防止重复放入池中
		Put2Pool() bool
		// Get4Pool 方法，用于从池中取出对象时
		// 将 ace_put 设置为 false，防止重复放入池中
		Get4Pool()
	}
)

// Put2Pool 方法，用于判断是否已经放入池中
// 如果已经放入池中，返回 true，否则返回 false
// 并将 ace_put 设置为 true，防止重复放入池中
func (a *Model) Put2Pool() bool {
	if a.ace_put {
		return true
	}
	a.mu.Lock()
	a.ace_put = true
	a.mu.Unlock()
	return false
}

// Get4Pool 方法，用于从池中取出对象时
// 将 ace_put 设置为 false，防止重复放入池中
func (a *Model) Get4Pool() {
	if !a.ace_put {
		return
	}
	a.mu.Lock()
	a.ace_put = false
	a.mu.Unlock()
}

// 实现一个 Lock 方法，让 noCopy 实现 sync.Locker 接口
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
