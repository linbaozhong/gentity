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

package util

import (
	"strings"
	"unicode"
)

func ParseField(fieldName string) string {
	var (
		isUpper bool
		upper   strings.Builder
	)
	for i, ch := range fieldName {
		if i == 0 {
			upper.WriteRune(unicode.ToUpper(ch))
		} else {
			if ch == 95 {
				isUpper = true
				continue
			}
			if isUpper {
				upper.WriteRune(unicode.ToUpper(ch))
				isUpper = false
			} else {
				upper.WriteRune(ch)
			}
		}
	}
	return upper.String()
}

func ParseFieldType(tp string, size int) string {
	switch strings.ToUpper(tp) {
	case "INT":
		return "int"
	case "VARCHAR":
		return "string"
	case "TINYINT":
		if size == 1 {
			return "bool"
		}
		return "int8"
	case "BIT":
		return "bool"
	case "TIMESTAMP", "DATETIME":
		return "time.Time"
	default:
		return "any" // 对于未明确映射的类型，使用接口类型作为占位符
	}
}
