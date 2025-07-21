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
	"github.com/linbaozhong/gentity/pkg/conv"
	"strconv"
	"time"
)

type (
	Time struct {
		time.Time
	}
)

// //////////////////////////////////
// Time
func (t *Time) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*t = Time{}
		return nil
	case time.Time:
		*t = Time{v}
		return nil
	case []byte:
		t2, err := time.Parse(time.DateTime, string(v))
		if err != nil {
			return err
		}
		*t = Time{t2}
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Time: %T", src)
	}
}
func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

// IsNil 是否空值，注意空值!=零值
func (t Time) IsNil() bool {
	return t.Time == NilTime
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

func (t Time) Bytes() []byte {
	return []byte(t.String())
}

func (t *Time) FromBytes(b []byte) {
	t2, _ := time.Parse(time.DateTime, string(b))
	*t = Time{t2}
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(t.String())), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

	if c == "" {
		*t = Time{}
		return nil
	}
	*t = Time{conv.String2Time(c)}
	return nil
}
