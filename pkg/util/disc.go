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
	"regexp"
	"strings"
)

type iif interface {
	~bool | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string | any
}

// IIF returns trueValue if condition is true, otherwise returns falseValue.
func IIF[T iif](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// IsUrl 检查url是否合法
// 如果url没有前缀，会自动添加https前缀或者scheme指定的前缀
// scheme: http, https, ftp
func IsUrl(url string, scheme ...string) (string, bool) {
	regex := regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
	if regex.MatchString(url) {
		return url, true
	}
	regex = regexp.MustCompile(`^(?:\/\/)?[^\s/$.?#].[^\s]*$`)
	if regex.MatchString(url) {
		prefix := "https:"
		if len(scheme) > 0 {
			prefix = scheme[0] + ":"
		}
		if strings.HasPrefix(url, "//") {
			return prefix + url, true
		}
		return prefix + "//" + url, true
	}
	return url, false
}
