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

package wechat

import (
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	"sync"
)

var (
	programs = sync.Map{}
)

func Programe(id, secret string) (*miniprogram.MiniProgram, error) {
	mi, ok := programs.Load(id)
	if ok {
		return mi.(*miniprogram.MiniProgram), nil
	}

	p := miniprogram.NewMiniProgram(&config.Config{
		AppID:     id,
		AppSecret: secret,
		Cache:     cache.NewMemory(),
	})
	programs.Store(id, p)
	return p, nil
}
