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
	"fmt"
	"testing"
)

func TestGetInviteCode(t *testing.T) {
	// 1. 创建生成器，传入你的密钥盐值 (上线后不可更改)
	coder := NewHashID("my_secret_salt_2023")

	fmt.Println("--- 连续 ID 混淆效果测试 ---")
	// 测试连续的 UID，看是否还有规律
	testUIDs := []uint64{1, 2, 3, 4, 5, 100, 1000, 123456789}

	for _, uid := range testUIDs {
		code, _ := coder.Encode(uid)
		decodeUID, _ := coder.Decode(code)
		fmt.Printf("UID: %-10d -> 邀请码: %-10s -> 解码: %-10d (一致: %v)\n", uid, code, decodeUID, uid == decodeUID)
	}

	fmt.Println("\n--- 盐值不同，结果完全不同 ---")
	coder2 := NewHashID("another_salt")
	code1, _ := coder.Encode(1)
	code2, _ := coder2.Encode(1)
	fmt.Printf("UID=1 在盐值1下的邀请码: %s\n", code1)
	fmt.Printf("UID=1 在盐值2下的邀请码: %s\n", code2)

	fmt.Println("\n--- 非法字符防御测试 ---")
	_, err := coder.Decode("!!!@@@") // 输入乱码
	fmt.Printf("输入乱码解码错误: %v\n", err)
}
