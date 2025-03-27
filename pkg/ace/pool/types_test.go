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
	"fmt"
	"github.com/linbaozhong/gentity/pkg/app"
	"testing"
)

type A struct {
	Model
	Name string
	Age  int
}

var poolA = New(app.Context, func() any {
	_obj := &A{}
	_obj.UUID()
	return _obj
})

func NewA() *A {
	obj := poolA.Get().(*A)
	// time.Sleep(time.Millisecond)
	// _obj.Name = time.Now().String()
	return obj
}
func (a *A) Reset() {
	a = nil
}
func (a *A) Free() {
	a.Reset()
	poolA.Put(a)
}

func (a *A) Clone() *A {
	_a := *a
	_a.UUID()
	return &_a
	// return &(*a)
}

func TestClone(t *testing.T) {
	_a := NewA()
	_a.Name = "hello"
	_b := _a.Clone()
	_b.Age = 123
	fmt.Println(&_a, _a, &_b, _b)
	_c := *_a
	_a.Free()
	_b.Free()
	_c.Free()
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

func TestSlice(t *testing.T) {
	a := []int{1, 2, 3}
	b := a[1:]
	fmt.Println(a, b)
	b[0] = 100
	fmt.Println(a, b)
	b = append([]int(nil), a...)
	fmt.Println(a, b)
	b[0] = 200
	fmt.Println(a, b)
}
