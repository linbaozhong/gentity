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
	"os"
)

func CtrCompressAndEncrypt(inputPath, outputPath string, password []byte) error {
	// 打开原始文件
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// 创建输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 创建 AES 块
	block, err := aes.NewCipher(password)
	if err != nil {
		return err
	}

	// 生成 nonce
	nonce := make([]byte, aes.BlockSize)
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	// 将 nonce 写入输出文件
	if _, err := outputFile.Write(nonce); err != nil {
		return err
	}

	// 使用 Pipe 连接压缩和加密
	pr, pw := io.Pipe()

	// 压缩协程
	go func() {
		defer pw.Close()
		xzWriter, err := xz.NewWriter(pw)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		defer xzWriter.Close()
		if _, err := io.Copy(xzWriter, inputFile); err != nil {
			pw.CloseWithError(err)
		}
	}()

	// 创建 CTR 加密器
	stream := cipher.NewCTR(block, nonce)
	// 创建加密流
	encryptedStream := &cipher.StreamWriter{S: stream, W: outputFile}

	// 将压缩数据加密并写入输出文件
	_, err = io.Copy(encryptedStream, pr)
	if err != nil {
		return err
	}

	// log.Printf("CTR 模式压缩加密完成: %s → %s\n", inputPath, outputPath)
	return nil
}

func CtrDecryptAndDecompress(inputPath, outputPath string, password []byte) error {
	// 打开加密文件
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// 创建输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 创建 AES 块
	block, err := aes.NewCipher(password)
	if err != nil {
		return err
	}

	// 读取 nonce
	nonce := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inputFile, nonce); err != nil {
		return err
	}

	// 使用 Pipe 连接解密和解压
	pr, pw := io.Pipe()

	// 解密协程
	go func() {
		defer pw.Close()
		// 创建 CTR 解密器
		stream := cipher.NewCTR(block, nonce)
		// 创建解密流
		decryptedStream := &cipher.StreamReader{S: stream, R: inputFile}
		_, err := io.Copy(pw, decryptedStream)
		if err != nil {
			pw.CloseWithError(err)
		}
	}()

	// 解压数据
	xzReader, err := xz.NewReader(pr)
	if err != nil {
		return err
	}
	defer xzReader.Close()

	// 将解压后的数据写入输出文件
	if _, err := io.Copy(outputFile, xzReader); err != nil {
		return err
	}

	// log.Printf("CTR 模式解密解压完成: %s → %s\n", inputPath, outputPath)
	return nil
}

// ////////////////////

//export CtrCompressAndEncryptWasm
func CtrCompressAndEncryptWasm(input []byte, password []byte) ([]byte, error) {
	// 创建 AES 块
	block, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}

	// 生成 nonce
	nonce := make([]byte, aes.BlockSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// 使用 Pipe 连接压缩和加密
	pr, pw := xzPipe()

	// 压缩协程
	go func() {
		defer pw.Close()
		xzWriter, err := xz.NewWriter(pw)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		defer xzWriter.Close()
		if _, err := xzWriter.Write(input); err != nil {
			pw.CloseWithError(err)
		}
	}()

	// 创建 CTR 加密器
	stream := cipher.NewCTR(block, nonce)
	// 创建加密流
	var encryptedData []byte
	encryptedStream := &cipher.StreamWriter{S: stream, W: &bytesBuffer{&encryptedData}}

	// 将 nonce 写入输出
	encryptedData = append(encryptedData, nonce...)

	// 将压缩数据加密并写入输出
	if _, err := io.Copy(encryptedStream, pr); err != nil {
		return nil, err
	}

	return encryptedData, nil
}

//export CtrDecryptAndDecompressWasm
func CtrDecryptAndDecompressWasm(input []byte, password []byte) ([]byte, error) {
	// 创建 AES 块
	block, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}

	// 读取 nonce
	nonce := input[:aes.BlockSize]
	input = input[aes.BlockSize:]

	// 创建 CTR 解密器
	stream := cipher.NewCTR(block, nonce)
	// 创建解密流
	decryptedStream := &cipher.StreamReader{S: stream, R: bytes.NewReader(input)}

	// 解压数据
	xzReader, err := xz.NewReader(decryptedStream)
	if err != nil {
		return nil, err
	}

	var decryptedData []byte
	decryptedBuffer := &bytesBuffer{&decryptedData}
	if _, err := io.Copy(decryptedBuffer, xzReader); err != nil {
		return nil, err
	}

	return decryptedData, nil
}

// xzPipe 创建一个用于 XZ 压缩的管道
func xzPipe() (*io.PipeReader, *io.PipeWriter) {
	return io.Pipe()
}

// bytesBuffer 实现 io.Writer 接口，用于将数据写入字节切片
type bytesBuffer struct {
	buf *[]byte
}

func (b *bytesBuffer) Write(p []byte) (n int, err error) {
	*b.buf = append(*b.buf, p...)
	return len(p), nil
}
