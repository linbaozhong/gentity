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
	New: func() any {
		fmt.Print("New \t")
		return &A{}
	},
}

func NewA() *A {
	obj := pool.Get().(*A)
	// time.Sleep(time.Millisecond)
	// obj.Name = time.Now().String()
	fmt.Print("\t\t")

	return obj
}
func Dispose(x *A) {
	pool.Put(x)
	x = nil
}
func (a *A) Free() {
	Dispose(a)
}

func TestType(t *testing.T) {
	a := NewA()
	a.Name = "hello"
	fmt.Println(a)
	a.Free()
	fmt.Println(1)
	a.Free()
	fmt.Println(2)
	a.Free()
	fmt.Println(3)
	a.Name = "world"
	a.Free()
	b := NewA()
	fmt.Println(1, b, &b)
	c := NewA()
	fmt.Println(2, c, &c)
	d := NewA()
	fmt.Println(3, d, &d)
	a = NewA()
	fmt.Println(a, &a)
	a = NewA()
	fmt.Println(a, &a)
	a = NewA()
	fmt.Println(a, &a)
	a = NewA()
	fmt.Println(a, &a)
	a = NewA()
	fmt.Println(a, &a)
	a = NewA()
	fmt.Println(a, &a)
}

func TestNil(t *testing.T) {
	var a = new(A)
	fmt.Println(a == nil)
	a.Name = "hello"
	a = nil
	fmt.Println(a == nil)
}
