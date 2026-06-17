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
	"bytes"
	"database/sql/driver"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/gjson"
)

// Json 表示 JSON 对象数组类型，用于处理数据库中的 JSON 数据类型
// 底层类型为 []map[string]any，支持存储多个键值对对象
// 实现了 sql.Scanner、driver.Valuer、json.Marshaler、json.Unmarshaler 接口
type Json []map[string]any

// Scan 实现 sql.Scanner 接口，从数据库读取 JSON 数据并解析为 Json 类型
//
// 参数:
//   - src: 数据库返回的原始数据，支持 nil、[]byte、string 类型
//
// 返回值:
//   - error: 解析错误，当数据类型不支持或 JSON 格式错误时返回
func (j *Json) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*j = nil
		return nil
	case []byte:
		if len(v) == 0 {
			*j = nil
			return nil
		}
		var data []map[string]any
		if err := Unmarshal(gjson.ParseBytes(v), &data); err != nil {
			return fmt.Errorf("failed to unmarshal Json: %w", err)
		}
		*j = data
		return nil
	case string:
		if v == "" {
			*j = nil
			return nil
		}
		var data []map[string]any
		if err := Unmarshal(gjson.Parse(v), &data); err != nil {
			return fmt.Errorf("failed to unmarshal Json: %w", err)
		}
		*j = data
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Json: %T", src)
	}
}

// Value 实现 driver.Valuer 接口，将 Json 类型转换为数据库可存储的值
//
// 返回值:
//   - driver.Value: JSON 字符串，空值返回 nil
//   - error: 序列化错误
func (j Json) Value() (driver.Value, error) {
	if j == nil || len(j) == 0 {
		return nil, nil
	}
	return MarshalString(j), nil
}

// MarshalJSON 实现 json.Marshaler 接口，将 Json 类型序列化为 JSON 字节数组
//
// 返回值:
//   - []byte: JSON 格式的字节数组，空值返回 "null"
//   - error: 序列化错误
func (j Json) MarshalJSON() ([]byte, error) {
	if j == nil || len(j) == 0 {
		return []byte("null"), nil
	}
	return Marshal(j), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口，从 JSON 字节数组反序列化为 Json 类型
//
// 参数:
//   - b: JSON 格式的字节数组
//
// 返回值:
//   - error: 反序列化错误，当 JSON 格式无效时返回
func (j *Json) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*j = nil
		return nil
	}
	var data []map[string]any
	if err := Unmarshal(gjson.ParseBytes(b), &data); err != nil {
		return fmt.Errorf("failed to unmarshal Json: %w", err)
	}
	*j = data
	return nil
}

// String 将 Json 类型转换为 JSON 字符串表示
//
// 返回值:
//   - string: JSON 格式的字符串，空值返回空字符串
func (j Json) String() string {
	if j == nil || len(j) == 0 {
		return ""
	}
	return MarshalString(j)
}

// Bytes 将 Json 类型转换为 JSON 字节数组
//
// 返回值:
//   - []byte: JSON 格式的字节数组，空值返回 nil
func (j Json) Bytes() []byte {
	if j == nil || len(j) == 0 {
		return nil
	}
	return Marshal(j)
}

// FromMaps 从 map 切片初始化 Json 类型
//
// 参数:
//   - maps: map 切片数据，nil 或空切片会设置为 nil
func (j *Json) FromMaps(maps []map[string]any) {
	if maps == nil || len(maps) == 0 {
		*j = nil
		return
	}
	*j = maps
}

// FromJSON 从 JSON 字符串解析并初始化 Json 类型
//
// 参数:
//   - jsonStr: JSON 格式的字符串
//
// 返回值:
//   - error: 解析错误，当 JSON 格式无效时返回
func (j *Json) FromJSON(jsonStr string) error {
	if jsonStr == "" {
		*j = nil
		return nil
	}
	var data []map[string]any
	if err := Unmarshal(gjson.Parse(jsonStr), &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	*j = data
	return nil
}

// ToMaps 将 Json 类型转换为 map 切片的深拷贝
//
// 返回值:
//   - []map[string]any: map 切片副本，空值返回 nil
func (j Json) ToMaps() []map[string]any {
	if j == nil || len(j) == 0 {
		return nil
	}
	result := make([]map[string]any, len(j))
	for i, m := range j {
		result[i] = make(map[string]any)
		for k, v := range m {
			result[i][k] = v
		}
	}
	return result
}

// Len 返回 Json 数组的长度
//
// 返回值:
//   - int: 数组中 map 的数量
func (j Json) Len() int {
	return len(j)
}

// IsNil 检查 Json 是否为 nil
//
// 返回值:
//   - bool: 为 nil 时返回 true
func (j Json) IsNil() bool {
	return j == nil
}

// IsZero 检查 Json 是否为零值（长度为 0）
//
// 返回值:
//   - bool: 长度为 0 时返回 true
func (j Json) IsZero() bool {
	return len(j) == 0
}

// IsEmpty 检查 Json 是否为空（nil 或长度为零）
//
// 返回值:
//   - bool: 为空时返回 true
func (j *Json) IsEmpty() bool {
	return j == nil || len(*j) == 0
}

// Get 获取指定索引位置的 map 对象
//
// 参数:
//   - index: 数组索引，从 0 开始
//
// 返回值:
//   - map[string]any: 指定位置的 map，索引越界返回 nil
func (j Json) Get(index int) map[string]any {
	if index < 0 || index >= len(j) {
		return nil
	}
	return j[index]
}

// Append 向 Json 数组追加一个 map 对象
//
// 参数:
//   - m: 要追加的 map 对象，nil 会被忽略
func (j *Json) Append(m map[string]any) {
	if m == nil {
		return
	}
	*j = append(*j, m)
}

// Remove 删除指定索引位置的 map 对象
//
// 参数:
//   - index: 要删除的元素索引，索引越界时不做任何操作
func (j *Json) Remove(index int) {
	if index < 0 || index >= len(*j) {
		return
	}
	*j = append((*j)[:index], (*j)[index+1:]...)
}
