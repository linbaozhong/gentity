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

// Point 表示 MySQL POINT 类型（SRID 4326 - WGS84）
type Point struct {
	X float64 // 纬度 (Latitude)
	Y float64 // 经度 (Longitude)
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

	// MySQL EWKB 格式（包含 SRID 4326）
	// 总长度 25 字节：
	// 字节 0: 字节序标记
	// 字节 1-4: 几何类型 (POINT + SRID 标志位)
	// 字节 5-8: SRID
	// 字节 9-17: X 坐标
	// 字节 17-25: Y 坐标
	if len(bytes) < 25 {
		return errors.New("invalid point data: expected 25 bytes, got " + string(rune(len(bytes))))
	}

	// 检查字节序
	if bytes[0] != 1 {
		return errors.New("unsupported byte order, only little-endian is supported")
	}

	// 解析几何类型（包含 SRID 标志位 0x20000000）
	// POINT 类型 + SRID 标志位 = 0x20000001
	geoType := binary.LittleEndian.Uint32(bytes[1:5])
	if geoType != 0x20000001 {
		return errors.New("unsupported geometry type or missing SRID")
	}

	// 解析 SRID
	srid := binary.LittleEndian.Uint32(bytes[5:9])
	if srid != 4326 {
		return errors.New("unsupported SRID, only SRID 4326 is supported")
	}

	// 解析 X 坐标（纬度）
	p.X = math.Float64frombits(binary.LittleEndian.Uint64(bytes[9:17]))
	// 解析 Y 坐标（经度）
	p.Y = math.Float64frombits(binary.LittleEndian.Uint64(bytes[17:25]))

	return nil
}

// Value 实现 driver.Valuer 接口
func (p Point) Value() (driver.Value, error) {
	if p.X == 0 && p.Y == 0 {
		return nil, nil
	}

	// 构建 EWKB 格式（包含 SRID 4326）
	// 总长度 25 字节
	buf := make([]byte, 25)

	buf[0] = 1 // 小端序
	// 几何类型：POINT (1) + SRID 标志位 (0x20000000)
	binary.LittleEndian.PutUint32(buf[1:5], 0x20000001)
	// SRID 4326 (WGS84)
	binary.LittleEndian.PutUint32(buf[5:9], 4326)
	// X 坐标（纬度）
	binary.LittleEndian.PutUint64(buf[9:17], math.Float64bits(p.X))
	// Y 坐标（经度）
	binary.LittleEndian.PutUint64(buf[17:25], math.Float64bits(p.Y))

	return buf, nil
}
