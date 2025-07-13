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
	"reflect"
	"strconv"
	"time"
)

func Unmarshal(r gjson.Result, ptr any, args ...any) error {
	b := []byte(r.Raw)
	if j, ok := ptr.(json.Unmarshaler); ok {
		return j.UnmarshalJSON(b)
	}
	if len(args) > 0 {
		// 获取 ptr 的反射值
		ptrValue := reflect.ValueOf(ptr)
		// 检查 ptr 是否为指针
		if ptrValue.Kind() != reflect.Ptr {
			return fmt.Errorf("ptr must be a pointer")
		}
		// 获取指针指向的值
		elem := ptrValue.Elem()
		// 获取 args[0] 的反射值
		argValue := reflect.ValueOf(args[0])
		// 检查类型是否兼容
		if argValue.Type().AssignableTo(elem.Type()) {
			// 将 args[0] 的值赋给 ptr 指向的位置
			elem.Set(argValue)
		} else {
			return fmt.Errorf("type of args[0] is not assignable to the type pointed by ptr")
		}
		return nil
	}

	return json.Unmarshal(b, ptr)
}

func Marshal(s any) string {
	switch v := s.(type) {
	case string:
		return strconv.Quote(v)
	case []byte:
		return strconv.Quote(string(v))
	case time.Time:
		if v.IsZero() {
			return ""
		}
		return strconv.Quote(v.Format(time.DateTime))
	case time.Duration: //转为毫秒
		return strconv.Quote(strconv.FormatInt(v.Milliseconds(), 64))
	default:
		if s == nil {
			return "null"
		}
		if j, ok := s.(json.Marshaler); ok {
			b, e := j.MarshalJSON()
			if e != nil {
				return ""
			}
			return string(b)
		}
		switch reflect.Indirect(reflect.ValueOf(s)).Kind() {
		case reflect.Struct, reflect.Slice, reflect.Map:
			b, e := json.Marshal(s)
			if e != nil {
				return fmt.Sprintf("%+v", s)
			}

			return string(b)
		}
	}
	return conv.Any2String(s)
}
