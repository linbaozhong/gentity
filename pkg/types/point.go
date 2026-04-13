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

package types

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"math"
)

// Point 表示 MySQL POINT 类型
type Point struct {
	X float64 // 经度
	Y float64 // 纬度
}

// Scan 实现 sql.Scanner 接口
func (p *Point) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-bytes into Point")
	}

	// MySQL WKB 格式解析
	if len(bytes) < 25 {
		return errors.New("invalid point data")
	}

	// 跳过 SRID（4字节）和字节序标记
	// X = 经度，Y = 纬度
	p.X = math.Float64frombits(binary.LittleEndian.Uint64(bytes[5:13]))
	p.Y = math.Float64frombits(binary.LittleEndian.Uint64(bytes[13:21]))

	return nil
}

// Value 实现 driver.Valuer 接口
func (p Point) Value() (driver.Value, error) {
	if p.X == 0 && p.Y == 0 {
		return nil, nil
	}

	// 构建 WKB 格式
	buf := make([]byte, 21)
	buf[0] = 1 // 小端序
	buf[1] = 1 // POINT 类型

	binary.LittleEndian.PutUint64(buf[5:13], math.Float64bits(p.X))
	binary.LittleEndian.PutUint64(buf[13:21], math.Float64bits(p.Y))

	return buf, nil
}
