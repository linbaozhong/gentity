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

package types

import (
	"sync/atomic"
)

var atomic_uint64 uint64

type (
	AceModel struct {
		ace_uuid uint64 `json:"-"` // 内部留用，禁止外部赋值
	}

	AceModeler interface {
		UUID() uint64
		Reset()
	}
)

func (a *AceModel) UUID() uint64 {
	if a.ace_uuid == 0 {
		a.ace_uuid = atomic.AddUint64(&atomic_uint64, 1)
	}
	return a.ace_uuid
}

func (p *AceModel) Lock()   {}
func (p *AceModel) Unlock() {}
