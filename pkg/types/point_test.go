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
	"math"
	"testing"
)

// TestPointValue 测试 Point 的 Value 方法（SRID 4326）
func TestPointValue(t *testing.T) {
	// 测试正常的 Point 值
	p := Point{
		X: 39.915,  // 北京天安门纬度
		Y: 116.404, // 北京天安门经度
	}

	value, err := p.Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}

	bytes, ok := value.([]byte)
	if !ok {
		t.Fatalf("Value() returned non-bytes type")
	}

	// 验证 EWKB 格式长度（25 字节）
	if len(bytes) != 25 {
		t.Errorf("Expected 25 bytes, got %d", len(bytes))
	}

	// 验证字节序
	if bytes[0] != 1 {
		t.Errorf("Expected little-endian byte order (1), got %d", bytes[0])
	}

	// 验证几何类型（POINT + SRID 标志位 = 0x20000001）
	geoType := binary.LittleEndian.Uint32(bytes[1:5])
	if geoType != 0x20000001 {
		t.Errorf("Expected geometry type 0x20000001, got 0x%08x", geoType)
	}

	// 验证 SRID 4326
	srid := binary.LittleEndian.Uint32(bytes[5:9])
	if srid != 4326 {
		t.Errorf("Expected SRID 4326, got %d", srid)
	}

	// 验证 X 坐标
	x := math.Float64frombits(binary.LittleEndian.Uint64(bytes[9:17]))
	if x != p.X {
		t.Errorf("Expected X = %f, got %f", p.X, x)
	}

	// 验证 Y 坐标
	y := math.Float64frombits(binary.LittleEndian.Uint64(bytes[17:25]))
	if y != p.Y {
		t.Errorf("Expected Y = %f, got %f", p.Y, y)
	}

	t.Logf("Point Value() test passed: X=%f, Y=%f, SRID=4326", p.X, p.Y)
}

// TestPointValueZero 测试零值 Point
func TestPointValueZero(t *testing.T) {
	p := Point{}

	value, err := p.Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}

	// 零值应该返回 nil
	if value != nil {
		t.Errorf("Expected nil for zero Point, got %v", value)
	}

	t.Log("Point zero value test passed")
}

// TestPointScan 测试 Point 的 Scan 方法（SRID 4326）
func TestPointScan(t *testing.T) {
	// 构造 EWKB 格式的测试数据（包含 SRID 4326）
	buf := make([]byte, 25)
	buf[0] = 1                                                           // 小端序
	binary.LittleEndian.PutUint32(buf[1:5], 0x20000001)                  // POINT + SRID 标志位
	binary.LittleEndian.PutUint32(buf[5:9], 4326)                        // SRID 4326
	binary.LittleEndian.PutUint64(buf[9:17], math.Float64bits(39.915))   // X = 纬度
	binary.LittleEndian.PutUint64(buf[17:25], math.Float64bits(116.404)) // Y = 经度

	var p Point
	err := p.Scan(buf)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	// 验证解析结果
	if p.X != 39.915 {
		t.Errorf("Expected X = 39.915, got %f", p.X)
	}

	if p.Y != 116.404 {
		t.Errorf("Expected Y = 116.404, got %f", p.Y)
	}

	t.Logf("Point Scan() test passed: X=%f (纬度), Y=%f (经度)", p.X, p.Y)
}

// TestPointScanNil 测试扫描 nil 值
func TestPointScanNil(t *testing.T) {
	var p Point
	err := p.Scan(nil)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	// nil 值不应该改变 Point
	if p.X != 0 || p.Y != 0 {
		t.Errorf("Point should remain unchanged for nil input")
	}

	t.Log("Point nil scan test passed")
}

// TestPointScanInvalidLength 测试无效长度
func TestPointScanInvalidLength(t *testing.T) {
	var p Point
	err := p.Scan([]byte{1, 2, 3})
	if err == nil {
		t.Error("Expected error for invalid length, got nil")
	}

	t.Logf("Point invalid length test passed: %v", err)
}

// TestPointRoundTrip 测试 Value 和 Scan 的往返转换
func TestPointRoundTrip(t *testing.T) {
	original := Point{
		X: 39.915,  // 纬度
		Y: 116.404, // 经度
	}

	// Value
	value, err := original.Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}

	// Scan
	var result Point
	err = result.Scan(value)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	// 验证往返转换是否正确
	if result.X != original.X {
		t.Errorf("X mismatch: original=%f, result=%f", original.X, result.X)
	}

	if result.Y != original.Y {
		t.Errorf("Y mismatch: original=%f, result=%f", original.Y, result.Y)
	}

	t.Logf("Point round-trip test passed: X=%f, Y=%f", result.X, result.Y)
}

// TestPointGeographicCoordinates 测试地理坐标（经纬度）
func TestPointGeographicCoordinates(t *testing.T) {
	// 测试一些著名的地理坐标
	// X = 纬度, Y = 经度
	testCases := []struct {
		name string
		x    float64 // 纬度
		y    float64 // 经度
	}{
		{"北京天安门", 39.915, 116.404},
		{"上海东方明珠", 31.239, 121.499},
		{"纽约时代广场", 40.758, -73.985},
		{"伦敦大本钟", 51.507, -0.127},
		{"东京塔", 35.658, 139.745},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := Point{X: tc.x, Y: tc.y}

			// Value
			value, err := p.Value()
			if err != nil {
				t.Fatalf("Value() error = %v", err)
			}

			// Scan
			var result Point
			err = result.Scan(value)
			if err != nil {
				t.Fatalf("Scan() error = %v", err)
			}

			// 验证（允许浮点数误差）
			if math.Abs(result.X-tc.x) > 1e-10 {
				t.Errorf("X mismatch: expected=%f, got=%f", tc.x, result.X)
			}

			if math.Abs(result.Y-tc.y) > 1e-10 {
				t.Errorf("Y mismatch: expected=%f, got=%f", tc.y, result.Y)
			}

			t.Logf("%s: X=%f, Y=%f (SRID 4326)", tc.name, result.X, result.Y)
		})
	}
}

// TestPointValueImplementsDriverValuer 验证 Point 实现了 driver.Valuer 接口
func TestPointValueImplementsDriverValuer(t *testing.T) {
	var _ driver.Valuer = Point{}
	t.Log("Point implements driver.Valuer interface")
}

// TestPointScanImplementsSqlScanner 验证 Point 实现了 sql.Scanner 接口
func TestPointScanImplementsSqlScanner(t *testing.T) {
	var p Point
	var _ interface{ Scan(interface{}) error } = &p
	t.Log("Point implements sql.Scanner interface")
}
