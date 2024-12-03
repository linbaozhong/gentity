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
	"testing"
	"time"
)

type App struct {
	Id      BigInt     `json:"id,omitempty" db:"'id' pk auto"`   //
	Arch    AceFloat64 `json:"arch,omitempty" db:"'arch'"`       // 操作系统架构
	Version AceBool    `json:"version,omitempty" db:"'version'"` // 版本号
	Url     AceString  `json:"url,omitempty" db:"'url'"`         // 应用下载地址
	State   AceInt8    `json:"state,omitempty" db:"'state'"`     //
	Force   Money      `json:"force,omitempty" db:"'force'"`     //
	Ctime   AceTime    `json:"ctime,omitempty" db:"'ctime'"`     //
}

func TestBase(t *testing.T) {
	a := new(App)
	a.Id = 1234567
	a.Arch = 3.14159265358979323846
	a.Version = true
	a.Url = "https://www.baidu.com"
	a.State = 1
	a.Force = 11256
	a.Ctime = AceTime{
		time.Now(),
	}

	b, e := json.Marshal(a)
	if e != nil {
		t.Error(e)
	}
	s := string(b)

	t.Log(s)
	//
	var a2 App
	e = json.Unmarshal([]byte(`{"id":"1234567","arch":3.14159265358979323846,"version":"true","url":"https://www.baidu.com","state":1,"force":112.56,"ctime":"2024-12-03 15:59:30"}`), &a2)
	if e != nil {
		t.Error(e)
	}
	t.Log(fmt.Sprintf("%+v", a2))
}

func TestError(t *testing.T) {
	e1 := NewError(1, "error")

	t.Log(e1.SetInfo("haha"))
	t.Log(e1)
}
