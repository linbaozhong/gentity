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

import "strconv"

// Deprecated: 请使用 HashString 代替
// MemHashString 使用类 FNV-1a 算法计算哈希值
func MemHashString[T string | []byte](key T) string {
	return strconv.FormatUint(MemHash(key), 10)
}

// MemHash 使用类 FNV-1a 算法计算哈希值
//
// 注意：对 string 类型按 rune 遍历（非标准 FNV-1a 按字节），ASCII 字符串结果一致，
//
//	多字节 UTF-8 字符串与 []byte 版本结果不同，但不影响确定性
func MemHash[T string | []byte](key T) uint64 {
	const (
		offset64 = 14695981039346656037
		prime64  = 1099511628211
	)

	hash := uint64(offset64)
	for _, c := range []byte(key) {
		hash ^= uint64(c)
		hash *= prime64
	}
	return hash
}

// HashString Deprecated 实现 FNV-1a 哈希函数
func HashString[T string | []byte](key T) string {
	return strconv.FormatUint(MemHash(key), 10)
}
