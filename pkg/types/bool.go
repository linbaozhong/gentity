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
func (b Bool) IsNil() bool {
	return b == NilBool
}

// IsNilZero 是否空值或零值
func (b Bool) IsNilZero() bool {
	return b.IsNil() || b == 0
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
	return []byte{byte(b)}
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
