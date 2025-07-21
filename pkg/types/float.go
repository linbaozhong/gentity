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
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

type (
	Float64 float64
	Float32 float32
)

// //////////////////////////////////////
// Float32
func (f32 *Float32) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*f32 = 0
		return nil
	case float64:
		*f32 = Float32(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Float32: %T", src)
	}
}

// IsNil 是否空值，注意空值!=零值
func (f32 Float32) IsNil() bool {
	return f32 == NilFloat32
}

func (f32 Float32) Float32() float32 {
	return float32(f32)
}

func (f32 Float32) String() string {
	return strconv.FormatFloat(float64(f32), 'f', -1, 32)
}

func (f32 Float32) Bytes() []byte {
	return binary.BigEndian.AppendUint32(nil, math.Float32bits(f32.Float32()))
}

func (f32 *Float32) FromBytes(b []byte) {
	*f32 = Float32(math.Float32frombits(binary.BigEndian.Uint32(b)))
}

func (f32 Float32) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(f32), 'f', -1, 32)), nil
}

func (f32 *Float32) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

	if c == "" {
		*f32 = 0
		return nil
	}
	f, e := strconv.ParseFloat(c, 32)
	*f32 = Float32(f)
	return e
}

// ///////////////////////////
// Float64
func (f64 *Float64) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*f64 = 0
		return nil
	case float64:
		*f64 = Float64(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Float64: %T", src)
	}
}

// IsNil 是否空值，注意空值!=零值
func (f64 Float64) IsNil() bool {
	return f64 == NilFloat64
}

func (f64 Float64) Float64() float64 {
	return float64(f64)
}

func (f64 Float64) String() string {
	return strconv.FormatFloat(float64(f64), 'f', -1, 64)
}

func (f64 Float64) Bytes() []byte {
	return binary.BigEndian.AppendUint64(nil, math.Float64bits(f64.Float64()))
}

func (f64 *Float64) FromBytes(b []byte) {
	*f64 = Float64(math.Float64frombits(binary.BigEndian.Uint64(b)))
}

func (f64 Float64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(f64), 'f', -1, 64)), nil
}

func (f64 *Float64) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)
	if c == "" {
		*f64 = 0
		return nil
	}
	f, e := strconv.ParseFloat(c, 64)
	*f64 = Float64(f)
	return e
}
