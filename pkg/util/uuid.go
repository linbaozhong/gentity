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
	"encoding/hex"

	"github.com/google/uuid"
)

// GetUUID 返回去除连接线(-)的32位字符的uuid字符串
func GetUUID() string {
	id := uuid.New()
	buf := make([]byte, 32)
	hex.Encode(buf[:], id[:])
	return string(buf[:])
}
