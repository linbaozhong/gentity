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
	"fmt"
	"reflect"
	"strconv"
)

type Bool int8

// //////////////////////////////////
// Bool
func (b *Bool) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*b = -1
		return nil
	case bool:
		if v {
			*b = 1
		} else {
			*b = 0
		}
		return nil
	case int64:
		if v > 0 {
			*b = 1
		} else if v == 0 {
			*b = 0
		} else {
			*b = -1
		}
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Bool: %T", src)
	}
}
func (b Bool) Value() (driver.Value, error) {
	return b > 0, nil
}

// IsNil 是否空值，注意空值!=零值
func (b *Bool) IsNil() bool {
	return b == nil
}

// IsZero 是否零值
func (b Bool) IsZero() bool {
	return b == 0
}

// IsEmpty 是否空值或零值
func (b *Bool) IsEmpty() bool {
	return b == nil || *b == 0
}

func (b Bool) Bool() bool {
	return b > 0
}

func (b Bool) String() string {
	if b > 0 {
		return "true"
	} else if b == 0 {
		return "false"
	} else {
		return "null"
	}
}

func (b Bool) Bytes() []byte {
	return []byte(b.String())
}

func (b Bool) Ptr() *bool {
	if b < 0 {
		return nil
	}
	tem := b > 0
	return &tem
}

// 3. 添加 Set 方法，方便链式创建
func (b Bool) SetTrue() Bool {
	b = 1
	return b
}
func (b Bool) SetFalse() Bool { b = 0; return b }
func (b Bool) SetNil() Bool   { b = -1; return b }

// 4. 添加构造函数
func NewBool(v bool) Bool {
	if v {
		return 1
	}
	return 0
}
func NewBoolFrom(v interface{}) (Bool, error) {
	switch v := v.(type) {
	case nil:
		return -1, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case int, int8, int16, int32, int64:
		if reflect.ValueOf(v).Int() > 0 {
			return 1, nil
		} else if reflect.ValueOf(v).Int() == 0 {
			return 0, nil
		}
		return -1, nil
	case uint, uint8, uint16, uint32, uint64:
		if reflect.ValueOf(v).Uint() > 0 {
			return 1, nil
		}
		return 0, nil
	case float32, float64:
		if reflect.ValueOf(v).Float() > 0 {
			return 1, nil
		} else if reflect.ValueOf(v).Float() == 0 {
			return 0, nil
		}
		return -1, nil
	case string:
		if v == "" || v == "null" {
			return -1, nil
		}
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return -1, fmt.Errorf("cannot parse bool from string: %w", err)
		}
		if parsed {
			return 1, nil
		}
		return 0, nil
	case []byte:
		return NewBoolFrom(string(v))
	default:
		return -1, fmt.Errorf("unsupported type for Bool: %T", v)
	}
}

func (b *Bool) FromBytes(buf []byte) {
	if len(buf) == 0 {
		*b = -1
		return
	}
	*b = Bool(buf[0])
}

func (b Bool) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *Bool) UnmarshalJSON(bs []byte) error {
	c := bytes2String(bs)
	if c == "" || c == "null" {
		*b = -1
		return nil
	}
	tem, e := strconv.ParseBool(c)
	if tem {
		*b = 1
	} else {
		*b = 0
	}
	return e
}
