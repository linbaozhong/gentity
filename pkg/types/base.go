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
	"database/sql/driver"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/conv"
	"strconv"
	"time"
)

type (
	String  string
	Uint64  uint64
	Uint32  uint32
	Uint16  uint16
	Uint8   uint8
	Uint    uint
	Int64   int64
	Int32   int32
	Int16   int16
	Int8    int8
	Int     int
	Bytes   []byte
	Float64 float64
	Float32 float32
	Bool    bool
	Time    struct {
		time.Time
	}
)

// ////////////////////////////
// Byte
func (b *Bytes) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*b = []byte{}
	case []byte:
		*b = v
	default:
		return fmt.Errorf("unsupported scan type for Byte: %T", src)
	}
	return nil
}

func (b Bytes) String() string {
	return conv.Bytes2String(b)
}

// ////////////////////////////
// String
func (s *String) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*s = ""
	case []byte:
		*s = String(conv.Bytes2String(v))
	case string:
		*s = String(v)
	default:
		return fmt.Errorf("unsupported scan type for String: %T", src)
	}
	return nil
}

func (s String) String() string {
	return string(s)
}

func (s String) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(`"` + string(s) + `"`), nil
}

func (s *String) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))
	*s = String(c)
	return nil
}

// /////////////////////////////
// Int8
func (i8 *Int8) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i8 = 0
		return nil
	case int64:
		*i8 = Int8(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Int8: %T", src)
	}
}

func (i8 Int8) Int8() int8 {
	return int8(i8)
}

func (i8 Int8) String() string {
	return strconv.FormatInt(int64(i8), 10)
}
func (i8 Int8) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i8), 10)), nil
}

func (i8 *Int8) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i8 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i8 = Int8(tem)
	return e
}

// /////////////////////////////
// Int16
func (i16 *Int16) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i16 = 0
		return nil
	case int64:
		*i16 = Int16(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Int16: %T", src)
	}
}

func (i16 Int16) Int16() int16 {
	return int16(i16)
}

func (i16 Int16) String() string {
	return strconv.FormatInt(int64(i16), 10)
}
func (i16 Int16) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i16), 10)), nil
}

func (i16 *Int16) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i16 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i16 = Int16(tem)
	return e
}

// ///////////////////////////////
// Int32
func (i32 *Int32) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i32 = 0
		return nil
	case int64:
		*i32 = Int32(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Int32: %T", src)
	}
}

func (i32 Int32) Int32() int32 {
	return int32(i32)
}

func (i32 Int32) String() string {
	return strconv.FormatInt(int64(i32), 10)
}
func (i32 Int32) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i32), 10)), nil
}

func (i32 *Int32) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i32 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i32 = Int32(tem)
	return e
}

// /////////////////////////////////////
// Int64
func (i64 *Int64) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i64 = 0
		return nil
	case int64:
		*i64 = Int64(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Int64: %T", src)
	}
}

func (i64 Int64) Int64() int64 {
	return int64(i64)
}

func (i64 Int64) String() string {
	return strconv.FormatInt(int64(i64), 10)
}
func (i64 Int64) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i64), 10)), nil
}

func (i64 *Int64) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i64 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i64 = Int64(tem)
	return e
}

// ///////////////////////////////////////
// Int
func (i *Int) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i = 0
		return nil
	case int64:
		*i = Int(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Int: %T", src)
	}
}

func (i Int) Int() int {
	return int(i)
}

func (i Int) String() string {
	return strconv.FormatInt(int64(i), 10)
}
func (i Int) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i), 10)), nil
}

func (i *Int) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i = Int(tem)
	return e
}

// /////////////////////////////
// Uint8
func (i8 *Uint8) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i8 = 0
		return nil
	case int64:
		*i8 = Uint8(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Uint8: %T", src)
	}
}

func (i8 Uint8) Uint8() uint8 {
	return uint8(i8)
}

func (i8 Uint8) String() string {
	return strconv.FormatUint(uint64(i8), 10)
}

func (i8 Uint8) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i8), 10)), nil
}

func (i8 *Uint8) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i8 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i8 = Uint8(tem)
	return e
}

// /////////////////////////////
// Uint16
func (i16 *Uint16) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i16 = 0
		return nil
	case int64:
		*i16 = Uint16(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Int16: %T", src)
	}
}

func (i16 Uint16) Uint16() uint16 {
	return uint16(i16)
}

func (i16 Uint16) String() string {
	return strconv.FormatUint(uint64(i16), 10)
}

func (i16 Uint16) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i16), 10)), nil
}

func (i16 *Uint16) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i16 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i16 = Uint16(tem)
	return e
}

// ///////////////////////////////
// Uint32
func (i32 *Uint32) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i32 = 0
		return nil
	case int64:
		*i32 = Uint32(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Uint32: %T", src)
	}
}

func (i32 Uint32) Uint32() uint32 {
	return uint32(i32)
}

func (i32 Uint32) String() string {
	return strconv.FormatUint(uint64(i32), 10)
}

func (i32 Uint32) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i32), 10)), nil
}

func (i32 *Uint32) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i32 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i32 = Uint32(tem)
	return e
}

// /////////////////////////////////////
// Uint64
func (i64 *Uint64) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i64 = 0
		return nil
	case int64:
		*i64 = Uint64(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Uint64: %T", src)
	}
}

func (i64 Uint64) Uint64() uint64 {
	return uint64(i64)
}

func (i64 Uint64) String() string {
	return strconv.FormatUint(uint64(i64), 10)
}

func (i64 Uint64) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i64), 10)), nil
}

func (i64 *Uint64) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i64 = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i64 = Uint64(tem)
	return e
}

// ///////////////////////////////////////
// Uint
func (i *Uint) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i = 0
		return nil
	case int64:
		*i = Uint(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Uint: %T", src)
	}
}

func (i Uint) Uint() uint {
	return uint(i)
}

func (i Uint) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

func (i Uint) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatInt(int64(i), 10)), nil
}

func (i *Uint) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*i = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i = Uint(tem)
	return e
}

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

func (f32 Float32) Float32() float32 {
	return float32(f32)
}

func (f32 Float32) String() string {
	return strconv.FormatFloat(float64(f32), 'f', -1, 32)
}

func (f32 Float32) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatFloat(float64(f32), 'f', -1, 32)), nil
}

func (f32 *Float32) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

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

func (f64 Float64) Float64() float64 {
	return float64(f64)
}

func (f64 Float64) String() string {
	return strconv.FormatFloat(float64(f64), 'f', -1, 64)
}

func (f64 Float64) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatFloat(float64(f64), 'f', -1, 64)), nil
}

func (f64 *Float64) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))
	if c == "" {
		*f64 = 0
		return nil
	}
	f, e := strconv.ParseFloat(c, 64)
	*f64 = Float64(f)
	return e
}

// //////////////////////////////////
// Bool
func (b *Bool) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*b = false
		return nil
	case bool:
		*b = Bool(v)
		return nil
	case int64:
		*b = v != 0
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Bool: %T", src)
	}
}

func (b Bool) Bool() bool {
	return bool(b)
}

func (b Bool) String() string {
	return strconv.FormatBool(b.Bool())
}
func (b Bool) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(strconv.FormatBool(bool(b))), nil
}

func (b *Bool) UnmarshalJSON(bs []byte) error {
	c := conv.Bytes2String(bytes.Trim(bs, "\""))
	if c == "" {
		*b = false
		return nil
	}
	tem, e := strconv.ParseBool(c)
	*b = Bool(tem)
	return e
}

// //////////////////////////////////
// Time
func (t *Time) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*t = Time{}
		return nil
	case time.Time:
		*t = Time{v}
		// *t = Time(v)
		return nil
	case []byte:
		t2, err := time.Parse(time.DateTime, conv.Bytes2String(v))
		if err != nil {
			return err
		}
		*t = Time{t2}
		// *t = Time(t2)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Time: %T", src)
	}
}
func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

func (t Time) Now() Time {
	return Time{time.Now()}
}

func Now() Time {
	return Time{time.Now()}
}
func (t Time) String() string {
	return t.Format(time.DateTime)
}

func (t Time) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(`"` + t.String() + `"`), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	c := conv.Bytes2String(bytes.Trim(b, "\""))

	if c == "" {
		*t = Time{}
		return nil
	}
	tem, e := time.Parse(time.DateTime, c)
	*t = Time{tem}
	// *t = Time(tem)
	return e
}
