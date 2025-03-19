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

package crypto

import (
	"encoding/hex"
	crypter "github.com/sekrat/aescrypter"
)

func AesEncrypt(src, key string) (string, error) {
	crypter := crypter.New()
	buf, err := crypter.Encrypt(key, []byte(src))
	return hex.EncodeToString(buf), err
}

func AesDecrypt(src, key string) (string, error) {
	buf, err := hex.DecodeString(src)
	if err != nil {
		return "", err
	}
	crypter := crypter.New()
	buf, err = crypter.Decrypt(key, buf)
	return string(buf), err
}

func AesEncryptBytes(src []byte, key string) (string, error) {
	crypter := crypter.New()
	buf, err := crypter.Encrypt(key, src)
	return hex.EncodeToString(buf), err
}
func AesDecryptBytes(src, key string) ([]byte, error) {
	buf, err := hex.DecodeString(src)
	if err != nil {
		return nil, err
	}
	crypter := crypter.New()
	buf, err = crypter.Decrypt(key, buf)
	return buf, err
}
