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

package api

import (
	"sync"
	"testing"
)

type User struct {
	Id   int
	Name string
}

func TestPost(t *testing.T) {
	var u User
	Init(&u)
}

func Init[A any](a *A) {
	InitiateX(nil, a)
}

func TestPool(t *testing.T) {
	type User struct {
		Id   int
		Name string
	}
	// 初始化一个对象池
	var pool = sync.Pool{
		New: func() any {
			return &User{}
		},
	}

	// 从对象池中获取一个对象
	user := pool.Get().(*User)
	user.Id = 1
	user.Name = "Hello World"

	// 将对象放回对象池中
	pool.Put(user)
	// 将同一对象重复放回对象池
	pool.Put(user)

	// 从对象池中获取一个对象
	userOne := pool.Get().(*User)
	userOne.Id = 2               // 修改Id的值
	t.Log("userOne = ", userOne) // 输出：userOne = &{2 Hello World}
	// 从对象池中再获取一个对象
	userTwo := pool.Get().(*User)
	t.Log("userTwo = ", userTwo) // 输出：userTwo = &{2 Hello World}
}
