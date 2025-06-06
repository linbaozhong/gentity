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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/ulikunitz/xz"
	"io"
	"log"
	"os"
)

func GcmCompressAndEncrypt(inputPath, outputPath string, password []byte) error {
	// 1. 打开原始文件
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// 2. 创建输出文件（加密后）
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 3. 初始化 AES-256 加密器
	block, err := aes.NewCipher(password)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return err
	}
	// 将 nonce 写入输出文件
	if _, err = outputFile.Write(nonce); err != nil {
		return err
	}

	// 4. 压缩数据
	var compressedData bytes.Buffer
	xzWriter, err := xz.NewWriter(&compressedData)
	if err != nil {
		return err
	}
	if _, err = io.Copy(xzWriter, inputFile); err != nil {
		xzWriter.Close()
		return err
	}
	if err = xzWriter.Close(); err != nil {
		return err
	}

	// 5. 加密压缩后的数据
	ciphertext := gcm.Seal(nil, nonce, compressedData.Bytes(), nil)
	if _, err := outputFile.Write(ciphertext); err != nil {
		return err
	}

	// log.Printf("压缩加密完成: %s → %s\n", inputPath, outputPath)
	return nil
}

func GcmDecryptAndDecompress(inputPath, outputPath string, password []byte) error {
	// 1. 打开加密文件
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Printf("打开加密文件失败: %v", err)
		return err
	}
	defer inputFile.Close()

	// 2. 创建输出文件（解密后）
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Printf("创建输出文件失败: %v", err)
		return err
	}
	defer outputFile.Close()

	// 3. 初始化 AES-256 解密器
	block, err := aes.NewCipher(password)
	if err != nil {
		log.Printf("初始化 AES 块失败: %v", err)
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("初始化 GCM 失败: %v", err)
		return err
	}

	// 4. 读取 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(inputFile, nonce); err != nil {
		log.Printf("读取 nonce 失败: %v", err)
		return err
	}
	// log.Printf("成功读取 nonce，长度: %d，内容: %x", len(nonce), nonce)

	// 5. 读取加密数据
	ciphertext, err := io.ReadAll(inputFile)
	if err != nil {
		log.Printf("读取加密数据失败: %v", err)
		return err
	}

	// 6. 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Printf("解密数据失败，认证错误: %v", err)
		return err
	}

	// 7. 解压数据
	xzReader, err := xz.NewReader(bytes.NewReader(plaintext))
	if err != nil {
		log.Printf("创建 xz 读取器失败: %v", err)
		return err
	}
	defer xzReader.Close()

	if _, err := io.Copy(outputFile, xzReader); err != nil {
		log.Printf("解压数据失败: %v", err)
		return err
	}

	// log.Printf("解密解压完成: %s → %s\n", inputPath, outputPath)
	return nil
}
