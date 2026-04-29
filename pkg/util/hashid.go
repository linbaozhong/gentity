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
	"errors"
	"fmt"
)

const (
	Base                 uint64 = 30
	baseNumber           uint64 = 100000000
	ERR_UID_INVALID             = "uid 必须大于 0"
	ERR_CODE_INVALID            = "非法的邀请码"
	ERR_CODE_PARSE_ERROR        = "邀请码解析异常"
)

// 基础字母表 (去除了 0OIL 等易混淆字符，共 30 个)
var baseAlphabet = []rune{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V',
	'W', 'X', 'Y', '3', '4', '5', '6', '7', '8', '9',
}

// hashid 混淆码生成器
type hashid struct {
	salt             string
	shuffledAlphabet []rune
}

// NewHashID 创建实例
func NewHashID(salt string) *hashid {
	coder := &hashid{
		salt: salt,
	}
	coder.shuffledAlphabet = shuffle(baseAlphabet, salt)
	return coder
}

func magicID(uid uint64) uint64 {
	return (uid % 10) * baseNumber
}

// Encode 根据用户id获取混淆码
func (c *hashid) Encode(uid uint64) (string, error) {
	if uid == 0 {
		return "", errors.New(ERR_UID_INVALID)
	}

	targetUID := uid + magicID(uid)

	currentAlphabet := make([]rune, len(c.shuffledAlphabet))
	copy(currentAlphabet, c.shuffledAlphabet)

	result := make([]rune, 0, 8)
	for targetUID > 0 {
		mod := targetUID % Base
		targetUID = targetUID / Base

		result = append(result, currentAlphabet[mod])

		// 动态移位
		currentAlphabet = shiftAlphabet(currentAlphabet, int(mod))
	}

	reverseRunes(result)

	return string(result), nil
}

// Decode 根据混淆码获取id
func (c *hashid) Decode(hid string) (uint64, error) {
	if hid == "" {
		return 0, errors.New(ERR_CODE_INVALID)
	}

	runes := []rune(hid)

	currentAlphabet := make([]rune, len(c.shuffledAlphabet))
	copy(currentAlphabet, c.shuffledAlphabet)

	var targetUID uint64 = 0
	var power uint64 = 1 // 权重从 Base^0 (即1) 开始

	// 从后往前遍历，先还原低位，再还原高位
	for i := len(runes) - 1; i >= 0; i-- {
		char := runes[i]

		idx := indexOf(currentAlphabet, char)
		if idx == -1 {
			return 0, fmt.Errorf("%s: 包含非法字符 %c", ERR_CODE_INVALID, char)
		}

		// 修复点：使用正确的权重累加
		targetUID += uint64(idx) * power
		power *= Base

		// 动态移位顺序必须与编码一致
		currentAlphabet = shiftAlphabet(currentAlphabet, idx)
	}

	magic := magicID(targetUID)
	if targetUID < magic {
		return 0, errors.New(ERR_CODE_PARSE_ERROR)
	}

	return targetUID - magic, nil
}

// --- 辅助函数 ---

func shuffle(alphabet []rune, salt string) []rune {
	result := make([]rune, len(alphabet))
	copy(result, alphabet)

	if salt == "" {
		return result
	}

	v := 0
	for i := len(result) - 1; i > 0; i-- {
		v += int(salt[i%len(salt)])
		j := v % (i + 1)
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func shiftAlphabet(alphabet []rune, n int) []rune {
	if n <= 0 {
		return alphabet
	}
	length := len(alphabet)
	pos := n % length
	if pos == 0 {
		return alphabet
	}
	return append(alphabet[pos:], alphabet[:pos]...)
}

func indexOf(alphabet []rune, char rune) int {
	for i, c := range alphabet {
		if c == char {
			return i
		}
	}
	return -1
}

func reverseRunes(r []rune) {
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
}
