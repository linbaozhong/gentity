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
	"encoding/json"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/conv"
	"github.com/linbaozhong/gentity/pkg/gjson"
	"go/types"
	"time"
)

func Unmarshal(r gjson.Result, ptr any, args ...any) error {
	b := conv.String2Bytes(r.Raw)
	if j, ok := ptr.(json.Unmarshaler); ok {
		return j.UnmarshalJSON(b)
	}
	if len(args) > 0 {
		ptr = args[0]
		return nil
	}

	return json.Unmarshal(b, ptr)
}

func Marshal(s any) string {
	if j, ok := s.(json.Marshaler); ok {
		b, e := j.MarshalJSON()
		if e != nil {
			return ""
		}
		return conv.Bytes2String(b)
	}
	if ss, ok := s.(conv.Stringer); ok {
		return `"` + ss.String() + `"`
	}
	switch v := s.(type) {
	case string:
		return `"` + v + `"`
	case []byte:
		return `"` + conv.Bytes2String(v) + `"`
	case time.Time:
		if v.IsZero() {
			return ""
		}
		return v.Format(time.DateTime)
	case types.Slice, types.Struct, types.Map:
		b, e := json.Marshal(v)
		if e != nil {
			return fmt.Sprintf("%+v", v)
		}
		return `"` + conv.Bytes2String(b) + `"`
	default:
		if s == nil {
			return "null"
		}
	}
	return conv.Any2String(s)
}
