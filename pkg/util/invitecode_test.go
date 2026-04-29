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
	// 测试连续的 UID，看是否还有规律
	testUIDs := []uint64{1, 2, 3, 4, 5, 100, 1000, 123456789}

	for _, uid := range testUIDs {
		code, _ := GetInviteCode(uid)
		decodeUID, _ := GetIDFromInviteCode(code)
		fmt.Printf("UID: %-10d -> 邀请码: %-10s -> 解码: %-10d (一致: %v)\n", uid, code, decodeUID, uid == decodeUID)
	}

}
