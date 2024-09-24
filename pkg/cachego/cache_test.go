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

package cachego

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/dchest/siphash"
	"github.com/linbaozhong/gentity/pkg/conv"
	"testing"
)

func TestConvert(t *testing.T) {
	key := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	hash := siphash.New(conv.String2Bytes(key)).Sum64()
	// 将哈希值转换为字节切片
	hashBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(hashBytes, hash)

	buf := hex.EncodeToString(hashBytes)
	t.Log(buf)
	buf = "company:" + buf
	t.Log(buf)
}
