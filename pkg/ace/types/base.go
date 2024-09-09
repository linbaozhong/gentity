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

const (
	Inner_Join JoinType = " INNER"
	Left_Join  JoinType = " LEFT"
	Right_Join JoinType = " RIGHT"

	Operator_and = " AND "
	Operator_or  = " OR "

	MaxLimit uint = 1000
	PageSize uint = 20
)

var (
	ErrCreateEmpty        = fmt.Errorf("No data is created")
	ErrBeanEmpty          = fmt.Errorf("bean=nil 或者 len(beans)=0 或者 len(beans)>100")
	ErrNotFound           = fmt.Errorf("not found")
	ErrSetterEmpty        = fmt.Errorf("setter=nil 或者 len(setter)=0")
	ErrBeansEmpty         = fmt.Errorf("beans=nil 或者 len(beans)=0")
	ErrArgsNotMatch       = fmt.Errorf("args not match")
	ErrPrimaryKeyNotMatch = fmt.Errorf("primary key not match")
)

type (
	JoinType string

	AceString  string
	AceInt64   int64
	AceInt32   int32
	AceInt16   int16
	AceInt8    int8
	AceInt     int
	AceByte    byte
	AceFloat64 float64
	AceFloat32 float32
	AceBool    bool
	AceTime    time.Time
)

func (s *AceString) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*s = ""
		return nil
	case string:
		*s = AceString(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceString: %T", src)
	}
}

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

func (b *AceBool) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*b = false
		return nil
	case bool:
		*b = AceBool(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceBool: %T", src)
	}
}

func (t *AceTime) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*t = AceTime(time.Time{})
		return nil
	case time.Time:
		*t = AceTime(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for AceTime: %T", src)
	}
}
