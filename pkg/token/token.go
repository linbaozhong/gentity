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

package token

import (
	"fmt"
	"github.com/linbaozhong/gentity/pkg/crypto"
	"strings"
)

const Token_Secret_Key = "snow19921115love"

// EncryptToken 加密token
func EncryptToken(tk string) (string, error) {
	return crypto.AesEncrypt(tk, Token_Secret_Key)
}

// DecryptToken 解密token
func DecryptToken(tk string) (string, error) {
	return crypto.AesDecrypt(tk, Token_Secret_Key)
}

// GenToken 生成Token
func GenToken(id string, bean string) (string, error) {
	return EncryptToken(fmt.Sprintf("%s_%s", bean, id))
}

// ParseToken 从token中解析出id和token
func ParseToken(tk string) (id string, token string) {
	pos := strings.Index(tk, "_")
	if pos < 1 {
		return
	}
	token = tk[:pos]
	id = tk[pos+1:]
	return
}

// GetIDAndTokenFromCipher 从密文中获取id和token
func GetIDAndTokenFromCipher(tk string) (id string, token string, err error) {
	// 解密tk
	token, err = DecryptToken(tk)
	if err != nil {
		return
	}
	// 解析tk
	id, token = ParseToken(token)
	return
}
