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
	"bytes"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/conv"
	"strconv"
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
	AceByte    []byte
	AceFloat64 float64
	AceFloat32 float32
	AceBool    bool
	AceTime    struct {
		time.Time
	}
)

// ////////////////////////////
// AceByte
func (b *AceByte) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*b = []byte{}
	case []byte:
		*b = v
	default:
		return fmt.Errorf("unsupported scan type for AceByte: %T", src)
	}
	return nil
}

//func (b AceByte) Value() (driver.Value, error) {
//	return b, nil
//}

// ////////////////////////////
// AceString
func (s *AceString) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*s = ""
	case []byte:
		*s = AceString(conv.Bytes2String(v))
	case string:
		*s = AceString(v)
	default:
		return fmt.Errorf("unsupported scan type for AceString: %T", src)
	}
	return nil
}

//	func (s AceString) Value() (driver.Value, error) {
//		return s, nil
//	}
func (s *AceString) String() string {
	return string(*s)
}

func (s AceString) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(`"` + string(s) + `"`), nil
}

func (s *AceString) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))
	*s = AceString(c)
	return nil
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

//	func (i8 AceInt8) Value() (driver.Value, error) {
//		return i8, nil
//	}
func (i8 *AceInt8) Int8() int8 {
	return int8(*i8)
}

func (i8 AceInt8) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i8), 10)), nil
}

func (i8 *AceInt8) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i8 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i8 = AceInt8(tem)
	return e
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

func (i16 AceInt16) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i16), 10)), nil
}

func (i16 *AceInt16) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i16 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i16 = AceInt16(tem)
	return e
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

func (i32 AceInt32) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i32), 10)), nil
}

func (i32 *AceInt32) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i32 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i32 = AceInt32(tem)
	return e
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

func (i64 AceInt64) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i64), 10)), nil
}

func (i64 *AceInt64) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i64 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i64 = AceInt64(tem)
	return e
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

func (i AceInt) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i), 10)), nil
}

func (i *AceInt) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i = AceInt(tem)
	return e
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

func (i8 AceUint8) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i8), 10)), nil
}

func (i8 *AceUint8) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i8 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i8 = AceUint8(tem)
	return e
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

func (i16 AceUint16) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i16), 10)), nil
}

func (i16 *AceUint16) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i16 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i16 = AceUint16(tem)
	return e
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

func (i32 AceUint32) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i32), 10)), nil
}

func (i32 *AceUint32) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i32 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i32 = AceUint32(tem)
	return e
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

func (i64 AceUint64) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i64), 10)), nil
}

func (i64 *AceUint64) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i64 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i64 = AceUint64(tem)
	return e
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

func (i AceUint) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i), 10)), nil
}

func (i *AceUint) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i = AceUint(tem)
	return e
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

func (f32 AceFloat32) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatFloat(float64(f32), 'f', -1, 32)), nil
}

func (f32 *AceFloat32) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*f32 = 0
		return nil
	}
	f, e := strconv.ParseFloat(c, 32)
	*f32 = AceFloat32(f)
	return e
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

func (f64 AceFloat64) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatFloat(float64(f64), 'f', -1, 64)), nil
}

func (f64 *AceFloat64) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))
	if c == "" {
		*f64 = 0
		return nil
	}
	f, e := strconv.ParseFloat(c, 64)
	*f64 = AceFloat64(f)
	return e
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

func (b *AceBool) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatBool(bool(*b))), nil
}

func (b *AceBool) UnmarshalJSON(bs []byte) error {
	c := conv.Bytes2String(bytes.Trim(bs, "\""))
	if c == "" {
		*b = false
		return nil
	}
	tem, e := strconv.ParseBool(c)
	*b = AceBool(tem)
	return e
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

func (t *AceTime) String() string {
	return t.Format(time.DateTime)
}

func (t *AceTime) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(`"` + t.String() + `"`), nil
}

func (t *AceTime) UnmarshalJSON(b []byte) error {
	c := string(bytes.Trim(b, "\""))

	if c == "" {
		*t = AceTime{}
		return nil
	}
	tem, e := time.Parse(time.DateTime, c)
	*t = AceTime{tem}
	return e
}
