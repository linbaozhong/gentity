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
	"github.com/linbaozhong/gentity/pkg/conv"
	"runtime"
	"strconv"
)

type (
	Int64 int64
	Int32 int32
	Int16 int16
	Int8  int8
	Int   int
)

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

func (i8 Int8) Bytes() []byte {
	return []byte{byte(i8)}
}

func (i8 Int8) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i8), 10)), nil
}

func (i8 *Int8) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

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

func (i16 Int16) Bytes() []byte {
	return conv.Base2Bytes(i16)
}

func (i16 Int16) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i16), 10)), nil
}

func (i16 *Int16) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

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

func (i32 Int32) Bytes() []byte {
	return conv.Base2Bytes(i32)
}

func (i32 Int32) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i32), 10)), nil
}

func (i32 *Int32) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

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

func (i64 Int64) Bytes() []byte {
	return conv.Base2Bytes(i64)
}

func (i64 Int64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i64), 10)), nil
}

func (i64 *Int64) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

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

func (i Int) Bytes() []byte {
	if runtime.GOARCH == "arm64" || runtime.GOARCH == "amd64" {
		return conv.Base2Bytes(int64(i))
	} else {
		return conv.Base2Bytes(int32(i))
	}
}

func (i Int) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i), 10)), nil
}

func (i *Int) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

	if c == "" {
		*i = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i = Int(tem)
	return e
}
