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
	"github.com/linbaozhong/gentity/pkg/conv"
	"unsafe"
)

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func MemHashString(s string) uintptr {
	return MemHash(conv.String2Bytes(s))
}

// MemHash
func MemHash(b []byte) uintptr {
	s := *(*stringStruct)(unsafe.Pointer(&b))
	return memhash(s.str, 0, uintptr(s.len))
}

// Hashfnv32 实现 FNV-1a 32 位哈希函数
func Hashfnv32(key string) uint32 {
	const (
		offset32 = 2166136261
		prime32  = 16777619
	)
	hash := uint32(offset32)
	for i := 0; i < len(key); i++ {
		hash ^= uint32(key[i])
		hash *= prime32
	}
	return hash
}
