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
	"sync/atomic"
)

var atomic_uint64 uint64

type (
	// 定义一个私有类型，用于禁止拷贝
	noCopy struct{}

	Model struct {
		_        noCopy
		ace_uuid uint64 `json:"-"` // 内部留用，禁止外部赋值
	}

	PoolModeler interface {
		UUID() uint64
		Reset()
	}
)

func (a *Model) UUID() uint64 {
	if a.ace_uuid == 0 {
		a.ace_uuid = atomic.AddUint64(&atomic_uint64, 1)
	}
	return a.ace_uuid
}

// 实现一个 Lock 方法，让 noCopy 实现 sync.Locker 接口
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
