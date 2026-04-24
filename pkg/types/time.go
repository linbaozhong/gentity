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
		// 尝试多种格式
		layouts := []string{
			"2006-01-02 15:04:05",  // MySQL datetime
			"2006-01-02T15:04:05Z", // ISO 8601
			"2006-01-02T15:04:05-07:00",
			"2006-01-02 15:04:05.000",
			"2006-01-02",
		}
		s := string(v)
		for _, layout := range layouts {
			if parsed, err := time.Parse(layout, s); err == nil {
				t.Time = parsed
				return nil
			}
		}
		return fmt.Errorf("cannot parse time: %s", s)
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

// IsZero 是否零值
func (t Time) IsZero() bool {
	return t.Time.IsZero()
}

// IsEmpty 是否空值或零值
func (t Time) IsEmpty() bool {
	return t.Time == NilTime || t.Time.IsZero()
}

func (t Time) Now() Time {
	return Time{time.Now()}
}

func Now() Time {
	return Time{time.Now()}
}

func (t Time) Unix() int64 {
	return t.Time.Unix()
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

// 支持自定义格式化
func (t Time) FormatStr(layout string) string {
	if t.IsNil() {
		return ""
	}
	return t.Time.Format(layout)
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
