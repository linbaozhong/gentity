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

package conv

import (
	"encoding/binary"
	"testing"
	"time"
)

func TestAny2Bytes(t *testing.T) {
	var s time.Time = time.Now()
	t.Log(s)
	b, _ := Any2Bytes(s)
	t.Log(b)
	_ = Bytes2Any(b, &s)
	t.Log(s)
}

func TestByte(t *testing.T) {
	b := Uint64ToBytes(123456)
	t.Log(b)
	t.Log(binary.BigEndian.Uint64(b))
}

// Uint64ToBytes 将 uint64 类型转换为字节切片
func Uint64ToBytes(num uint64) []byte {
	// 创建一个长度为 8 的字节切片，因为 uint64 占 8 个字节
	buf := make([]byte, 8)
	// 使用 BigEndian 字节序将 uint64 写入字节切片
	binary.BigEndian.PutUint64(buf, num)
	return buf
}
