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
	"fmt"
	"sync"
	"testing"
)

type A struct {
	Name string
}

var pool = sync.Pool{
	New: func() interface{} {
		fmt.Println("New")
		return A{}
	},
}

func NewA() *A {
	obj := pool.Get().(A)
	return &obj
}
func Dispose(x *A) {
	x = nil
	//fmt.Println(x == nil)
}
func (a *A) Free() {
	if a == nil {
		return
	}
	a.Name = ""
	Dispose(a)
	pool.Put(*a)
}

func TestType(t *testing.T) {
	a := NewA()
	fmt.Println(a)
	a.Free()
	a.Free()
	a.Free()
	a.Free()
	a.Name = "world"
	a = NewA()
	fmt.Println(1, a)
	a = NewA()
	fmt.Println(2, a)
	a = NewA()
	fmt.Println(3, a)
	a = NewA()
	fmt.Println(a)
	a = NewA()
	fmt.Println(a)
	a = NewA()
	fmt.Println(a)
	a = NewA()
	fmt.Println(a)
	a = NewA()
	fmt.Println(a)
	a = NewA()
	fmt.Println(a)
}
