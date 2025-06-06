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

package zip_aes

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	password = []byte("my-32-byte-password-123456789012") // AES-256 需要 32 字节

)

func TestCtr(t *testing.T) {
	root, e := os.Getwd()
	if e != nil {
		t.Error(e)
	}

	if err := CtrCompressAndEncrypt(filepath.Join(root, "readme.md"), filepath.Join(root, "readme.bin"), password); err != nil {
		t.Error("压缩加密失败:", err)
	}
	if err := CtrDecryptAndDecompress(filepath.Join(root, "readme.bin"), filepath.Join(root, "readme1.md"), password); err != nil {
		t.Error("解密解压失败:", err)
	}
}

func TestGcm(t *testing.T) {
	root, e := os.Getwd()
	if e != nil {
		t.Error(e)
	}

	if err := GcmCompressAndEncrypt(filepath.Join(root, "readme.md"), filepath.Join(root, "readme.bin"), password); err != nil {
		t.Error("压缩加密失败:", err)
	}
	if err := GcmDecryptAndDecompress(filepath.Join(root, "readme.bin"), filepath.Join(root, "readme1.md"), password); err != nil {
		t.Error("解密解压失败:", err)
	}

}

func BenchmarkCtrCompressAndEncrypt(b *testing.B) {
	root, e := os.Getwd()
	if e != nil {
		b.Error(e)
	}
	for i := 0; i < b.N; i++ {
		if err := CtrCompressAndEncrypt(filepath.Join(root, "readme.md"), filepath.Join(root, "readme.bin"), password); err != nil {
			b.Error("压缩加密失败:", err)
		}
		if err := CtrDecryptAndDecompress(filepath.Join(root, "readme.bin"), filepath.Join(root, "readme1.md"), password); err != nil {
			b.Error("解密解压失败:", err)
		}
	}
}

func BenchmarkGcmCompressAndEncrypt(b *testing.B) {
	root, e := os.Getwd()
	if e != nil {
		b.Error(e)
	}
	for i := 0; i < b.N; i++ {
		if err := GcmCompressAndEncrypt(filepath.Join(root, "readme.md"), filepath.Join(root, "readme.bin"), password); err != nil {
			b.Error("压缩加密失败:", err)
		}
		if err := GcmDecryptAndDecompress(filepath.Join(root, "readme.bin"), filepath.Join(root, "readme1.md"), password); err != nil {
			b.Error("解密解压失败:", err)
		}

	}
}
