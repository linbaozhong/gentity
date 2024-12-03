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
	"fmt"
	"time"
)

type (
	AceString  string
	AceUint64  uint64
	AceUint32  uint32
	AceUint16  uint16
	AceUint8   uint8
	AceUint    uint
	AceInt64   int64
	AceInt32   int32
	AceInt16   int16
	AceInt8    int8
	AceInt     int
	AceByte    byte
	AceFloat64 float64
	AceFloat32 float32
	AceBool    bool
	AceTime    struct {
		time.Time
	}
)

// ////////////////////////////
// AceString
func (s *AceString) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*s = ""
		return nil
	case []byte:
		*s = AceString(v)
		return nil
	case string:
		*s = AceString(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceString: %T", src)
	}
}

func (s *AceString) String() string {
	return string(*s)
}

// ////////////////////////////
// AceByte
func (b *AceByte) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*b = 0
		return nil
	case int64:
		*b = AceByte(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceByte: %T", src)
	}
}

func (b *AceByte) Byte() byte {
	return byte(*b)
}

// /////////////////////////////
// AceInt8
func (i8 *AceInt8) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i8 = 0
		return nil
	case int64:
		*i8 = AceInt8(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceInt8: %T", src)
	}
}

func (i8 *AceInt8) Int8() int8 {
	return int8(*i8)
}

// /////////////////////////////
// AceInt16
func (i16 *AceInt16) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i16 = 0
		return nil
	case int64:
		*i16 = AceInt16(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceInt16: %T", src)
	}
}

func (i16 *AceInt16) Int16() int16 {
	return int16(*i16)
}

// ///////////////////////////////
// AceInt32
func (i32 *AceInt32) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i32 = 0
		return nil
	case int64:
		*i32 = AceInt32(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceInt32: %T", src)
	}
}

func (i32 *AceInt32) Int32() int32 {
	return int32(*i32)
}

// /////////////////////////////////////
// AceInt64
func (i64 *AceInt64) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i64 = 0
		return nil
	case int64:
		*i64 = AceInt64(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceInt64: %T", src)
	}
}

func (i64 *AceInt64) Int64() int64 {
	return int64(*i64)
}

// ///////////////////////////////////////
// AceInt
func (i *AceInt) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i = 0
		return nil
	case int64:
		*i = AceInt(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceInt: %T", src)
	}
}

func (i *AceInt) Int() int {
	return int(*i)
}

// /////////////////////////////
// AceUint8
func (i8 *AceUint8) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i8 = 0
		return nil
	case int64:
		*i8 = AceUint8(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceUint8: %T", src)
	}
}

func (i8 *AceUint8) Uint8() uint8 {
	return uint8(*i8)
}

// /////////////////////////////
// AceUint16
func (i16 *AceUint16) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i16 = 0
		return nil
	case int64:
		*i16 = AceUint16(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceInt16: %T", src)
	}
}

func (i16 *AceUint16) Uint16() uint16 {
	return uint16(*i16)
}

// ///////////////////////////////
// AceUint32
func (i32 *AceUint32) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i32 = 0
		return nil
	case int64:
		*i32 = AceUint32(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceUint32: %T", src)
	}
}

func (i32 *AceUint32) Uint32() uint32 {
	return uint32(*i32)
}

// /////////////////////////////////////
// AceUint64
func (i64 *AceUint64) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i64 = 0
		return nil
	case int64:
		*i64 = AceUint64(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceUint64: %T", src)
	}
}

func (i64 *AceUint64) Int64() uint64 {
	return uint64(*i64)
}

// ///////////////////////////////////////
// AceUint
func (i *AceUint) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i = 0
		return nil
	case int64:
		*i = AceUint(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceUint: %T", src)
	}
}

func (i *AceUint) Uint() uint {
	return uint(*i)
}

// //////////////////////////////////////
// AceFloat32
func (f32 *AceFloat32) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*f32 = 0
		return nil
	case float64:
		*f32 = AceFloat32(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceFloat32: %T", src)
	}
}

func (f32 *AceFloat32) Float32() float32 {
	return float32(*f32)
}

// ///////////////////////////
// AceFloat64
func (f64 *AceFloat64) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*f64 = 0
		return nil
	case float64:
		*f64 = AceFloat64(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceFloat64: %T", src)
	}
}

func (f64 *AceFloat64) Float64() float64 {
	return float64(*f64)
}

// //////////////////////////////////
// AceBool
func (b *AceBool) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*b = false
		return nil
	case bool:
		*b = AceBool(v)
		return nil
	case int64:
		*b = v != 0
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceBool: %T", src)
	}
}

func (b *AceBool) Bool() bool {
	return bool(*b)
}

// //////////////////////////////////
// AceTime
func (t *AceTime) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*t = AceTime{}
		return nil
	case time.Time:
		*t = AceTime{v}
		return nil
	case []byte:
		t2, err := time.Parse(time.DateTime, string(v))
		if err != nil {
			return err
		}
		*t = AceTime{t2}
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceTime: %T", src)
	}
}

//
// func (t *AceTime) Time() time.Time {
// 	return time.Time(*t)
// }

func (t AceTime) String() string {
	return t.Format(time.DateTime)
}
