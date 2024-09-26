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
	"testing"
	"time"
)

var (
	apool = &Pool{New: func() any {
		return &Object{}
	}}
)

type Object struct {
	Name string
}

func NewA() *Object {
	obj := apool.Get().(*Object)
	fmt.Print("\t\t")
	return obj
}

func TestPut(t *testing.T) {

	obj := NewA()
	apool.Put(obj)
	obj1 := NewA()
	obj1.Name = "hello"
	// fmt.Println(apool.Len(), "长度")
	obj2 := NewA()
	obj2.Name = "linbaozhong"
	obj3 := NewA()
	// apool.Put(obj1)
	// apool.Put(obj2)
	// fmt.Println(apool.Len(), "长度")
	obj1 = nil
	obj2 = nil
	obj3 = nil
	fmt.Println(obj1, obj2, obj3)
	// fmt.Println(reflect.ValueOf(obj).Pointer(), reflect.ValueOf(obj2).Pointer())
	// apool.Put(obj2)
	// apool.Put(obj3)
	// fmt.Println(apool.Len(), "长度")

	time.Sleep(time.Minute * 5)
	fmt.Println(apool.Len(), "长度")
}
