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

package dao

import (
	"fmt"
)

type DataAccessInterface struct {
	Name        string // 方法名
	Method      string // DAL方法名，如Get,List,Save,Update,Delete等，默认Get
	Title       string // 方法标题
	Description string // 方法描述
	Table       string // 表名
	Input       []any  // 方法输入
	Output      []any  // 方法输出
}

var dais = make(map[string]DataAccessInterface)

func RegisterDpi(dai DataAccessInterface) {
	dais[dai.Name] = dai
}

func Run() {
	for _, dai := range dais {
		if dai.Method == "" {
			dai.Method = "Get"
		}
		fmt.Println(`//`, dai.Title)
		fmt.Println(`//`, dai.Description)
		fmt.Println(`func `, dai.Name, `(ctx context.Context,`, `)`)
	}
}
