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

package validator

import (
	"github.com/linbaozhong/gentity/pkg/conv"
)

// InRangeInt returns true if value lies between left and right border
func InRangeInt(value, left, right interface{}) bool {
	value64 := conv.Any2Int64(value)
	left64 := conv.Any2Int64(left)
	right64 := conv.Any2Int64(right)
	if left64 > right64 {
		left64, right64 = right64, left64
	}
	return value64 >= left64 && value64 <= right64
}

// InRangeFloat32 returns true if value lies between left and right border
func InRangeFloat32(value, left, right float32) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

// InRangeFloat64 returns true if value lies between left and right border
func InRangeFloat64(value, left, right float64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

// InRange returns true if value lies between left and right border, generic type to handle int, float32, float64 and string.
// All types must the same type.
// False if value doesn't lie in range or if it incompatible or not comparable
func InRange(value interface{}, left interface{}, right interface{}) bool {
	switch value.(type) {
	case int:
		intValue := conv.Any2Int64(value)
		intLeft := conv.Any2Int64(left)
		intRight := conv.Any2Int64(right)
		return InRangeInt(intValue, intLeft, intRight)
	case float32, float64:
		intValue := conv.Any2Float64(value)
		intLeft := conv.Any2Float64(left)
		intRight := conv.Any2Float64(right)
		return InRangeFloat64(intValue, intLeft, intRight)
	case string:
		return value.(string) >= left.(string) && value.(string) <= right.(string)
	default:
		return false
	}
}
