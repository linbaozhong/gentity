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

func SliceContains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// 计算两个slice的差集
func SliceDiff[T comparable](slice1, slice2 []T) []T {
	var diff []T
	var set = make(map[T]struct{})
	for _, s := range slice2 {
		set[s] = struct{}{}
	}

	for _, s := range slice1 {
		if _, ok := set[s]; !ok {
			diff = append(diff, s)
		}
	}

	return diff
}
