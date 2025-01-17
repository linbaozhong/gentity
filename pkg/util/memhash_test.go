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
	"github.com/dchest/siphash"
	"github.com/linbaozhong/gentity/pkg/conv"
	"testing"
)

func TestHashKey(t *testing.T) {
	prefix := "test"
	t.Log(MemHashString(prefix))
	t.Log(MemHashString(prefix + "abc"))
	t.Log(MemHashString(prefix + "123"))
	t.Log(MemHashString(prefix + "eter"))
	t.Log(MemHashString(prefix + "4579"))
}

func BenchmarkHash(b *testing.B) {
	buf := conv.String2Bytes("testtesttesttesttest")
	for i := 0; i < b.N; i++ {
		siphash.New(buf).Sum64()
	}
}

func BenchmarkMemHash(b *testing.B) {
	buf := conv.String2Bytes("testtesttesttesttest")
	for i := 0; i < b.N; i++ {
		MemHash(buf)
	}
}

func TestIIF(t *testing.T) {
	var b = false
	t.Log(IIF(b, 0.33, 4.0))
}
