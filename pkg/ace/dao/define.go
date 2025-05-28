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
	"github.com/linbaozhong/gentity/pkg/ace"
)

type DataAccessInterface struct {
	NameSpace   string                // 命名空间
	Name        string                // 方法名
	Title       string                // 方法标题
	Description string                // 方法描述
	Table       string                // 表名
	Method      ace.MethodType        // 方法类型
	Input       ace.Builder           // 方法输入
	Output      any                   // 方法输出
	Children    []DataAccessInterface // 子方法
}
