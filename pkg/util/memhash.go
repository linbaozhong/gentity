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

// //go:noescape
// //go:linkname memhash runtime.memhash
// func memhash(p unsafe.Pointer, h, s uintptr) uintptr
//
// type stringStruct struct {
// 	str unsafe.Pointer
// 	len int
// }

func MemHashString(s string) uint {
	return HashString(s)
}

// MemHash
func MemHash(b []byte) uint {
	// s := *(*stringStruct)(unsafe.Pointer(&b))
	// return uint(memhash(s.str, 0, uintptr(s.len)))
	return HashByte(b)
}

// Hashfnv32 实现 FNV-1a 哈希函数
func HashString(key string) uint {
	const (
		offset64 = 14695981039346656037
		prime64  = 1099511628211
	)
	hash := uint64(offset64)
	for _, c := range key {
		hash ^= uint64(c)
		hash *= prime64
	}
	return uint(hash)
}

// Hashfnv32 实现 FNV-1a 哈希函数
func HashByte(key []byte) uint {
	const (
		offset64 = 14695981039346656037
		prime64  = 1099511628211
	)
	hash := uint64(offset64)
	for _, c := range key {
		hash ^= uint64(c)
		hash *= prime64
	}
	return uint(hash)
}
