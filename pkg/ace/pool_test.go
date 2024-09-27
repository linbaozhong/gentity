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

package ace

import (
	"fmt"
	"github.com/linbaozhong/gentity/pkg/ace/types"
	"runtime"
	"sync"
	"testing"
	"time"
)

var (
	appplMap sync.Map
	apool    = &sync.Pool{New: func() any {
		return &Object{}
	}}
)

type Object struct {
	types.AceModel
	Name string
}

func (o *Object) Free() {
	uuid := o.UUID()
	if _, ok := appplMap.Load(uuid); ok {
		return
	}

	appplMap.Store(uuid, struct{}{})
	apool.Put(o)
}

func NewA() *Object {
	o := apool.Get().(*Object)
	uuid := o.UUID()
	appplMap.Delete(uuid)
	return o
}

func TestPut(t *testing.T) {

	obj := NewA()
	runtime.SetFinalizer(obj, func(obj *Object) {
		fmt.Println("finalizer")
	})

	obj = nil

	runtime.GC()

	time.Sleep(time.Second * 2)
	fmt.Println("Main function has finished.")

	// obj.Free()
	// obj.Free()
	//
	// obj1 := NewA()
	// obj2 := NewA()
	// obj1.Free()
	// obj2.Free()
	//
	// obj3 := NewA()
	// fmt.Println(obj1, obj2, obj3)
}
