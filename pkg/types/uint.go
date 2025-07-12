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
	"runtime"
	"strconv"
)

type (
	Uint64 uint64
	Uint32 uint32
	Uint16 uint16
	Uint8  uint8
	Uint   uint
)

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

func (i8 Uint8) Bytes() []byte {
	return []byte{byte(i8)}
}

func (i8 Uint8) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i8), 10)), nil
}

func (i8 *Uint8) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

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

func (i16 Uint16) Bytes() []byte {
	return binary.BigEndian.AppendUint16(nil, uint16(i16))
}

func (i16 Uint16) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i16), 10)), nil
}

func (i16 *Uint16) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

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

func (i32 Uint32) Bytes() []byte {
	return binary.BigEndian.AppendUint32(nil, uint32(i32))
}

func (i32 Uint32) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i32), 10)), nil
}

func (i32 *Uint32) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

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

func (i64 Uint64) Bytes() []byte {
	return binary.BigEndian.AppendUint64(nil, uint64(i64))
}

func (i64 Uint64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i64), 10)), nil
}

func (i64 *Uint64) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

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
func (i Uint) Bytes() []byte {
	if runtime.GOARCH == "arm64" || runtime.GOARCH == "amd64" {
		return binary.BigEndian.AppendUint64(nil, uint64(i))
	} else {
		return binary.BigEndian.AppendUint32(nil, uint32(i))
	}
}

func (i Uint) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i), 10)), nil
}

func (i *Uint) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

	if c == "" {
		*i = 0
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)
	*i = Uint(tem)
	return e
}
