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
	"strconv"
)

type String string

// ////////////////////////////
// String
func (s *String) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*s = ""
	case []byte:
		*s = String(v)
	case string:
		*s = String(v)
	default:
		return fmt.Errorf("unsupported scan type for String: %T", src)
	}
	return nil
}

// IsNil 是否空值，注意空值!=零值
func (s String) IsNil() bool {
	return s == NilString
}

// IsNilZero 是否空值或零值
func (s String) IsNilZero() bool {
	return s.IsNil() || s == ""
}

func (s String) String() string {
	return string(s)
}

func (s String) Bytes() []byte {
	return []byte(s)
}

func (s String) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(s))), nil
}

func (s *String) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)
	*s = String(c)
	return nil
}
