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
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"testing"
)

type App struct {
	Id      BigInt  `json:"id,omitempty" db:"'id' pk auto"`   //
	Arch    float64 `json:"arch,omitempty" db:"'arch'"`       // 操作系统架构
	Version Bool    `json:"version,omitempty" db:"'version'"` // 版本号
	Url     String  `json:"url,omitempty" db:"'url'"`         // 应用下载地址
	State   Int     `json:"state,omitempty" db:"'state'"`     //
	Force   Money   `json:"force,omitempty" db:"'force'"`     //
	Ctime   Time    `json:"ctime,omitempty" db:"'ctime'"`     //
	Data    []Smap  `json:"data"`
}

func TestBase(t *testing.T) {
	// a := new(App)
	// a.Id = 1234567
	// a.Arch = 3.14159265358979323846
	// a.Version = Bool(-1)
	// a.Url = "https://www.baidu.com"
	// a.State = 1
	// a.Force = 11256
	// a.Ctime = Now()
	//
	// b, e := json.Marshal(a)
	// if e != nil {
	// 	t.Error(e)
	// }
	// s := string(b)
	// t.Log(s)
	//
	// n := NewSmap(3).
	// 	Set("id", a.Id).
	// 	Set("arch", a.Arch).
	// 	Set("version", a.Version).
	// 	Set("force", a.Force).
	// 	Set("url", a.Url).
	// 	Set("ctime", a.Ctime)
	// m := NewSmap(3).
	// 	Set("data", []Smap{n, n})
	//
	// r := NewResult()
	// r.Data = m
	// b, e = json.Marshal(r)
	// if e != nil {
	// 	t.Error(e)
	// }
	// s = string(b)
	// t.Log(s)
	//
	var a2 App
	a2.Id = math.MaxUint64
	e := json.Unmarshal([]byte(`{"id":"undefined",  "arch":3.14159265358979323846,"version":"true","url":"undefined","state":"undefined","force":112.56,"ctime":"12:03:35"}`), &a2)
	if e != nil {
		t.Error(e)
	}

	t.Log(a2)
}

func TestError(t *testing.T) {
	e1 := NewError(1, "error1")
	e3 := NewError(1, "error1")
	e2 := fmt.Errorf("error2:,%w", e1)
	t.Log(e1)
	t.Log(e2)
	t.Log(errors.Is(e1, e3))
}

func TestConv(t *testing.T) {
	s := "as12"
	i, e := strconv.ParseInt(s, 10, 64)
	t.Log(i, e)
}

func TestJsonResult(t *testing.T) {
	var r = NewResult()
	b, e := json.Marshal(r)
	if e != nil {
		t.Error(e)
	}
	t.Log(string(b))
}

func TestTime(t *testing.T) {

}
