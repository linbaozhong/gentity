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

type Json []map[string]any

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

func (j Json) Value() (driver.Value, error) {
	if j == nil || len(j) == 0 {
		return nil, nil
	}
	return MarshalString(j), nil
}

func (j Json) MarshalJSON() ([]byte, error) {
	if j == nil || len(j) == 0 {
		return []byte("null"), nil
	}
	return Marshal(j), nil
}

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

func (j Json) String() string {
	if j == nil || len(j) == 0 {
		return ""
	}
	return MarshalString(j)
}

func (j Json) Bytes() []byte {
	if j == nil || len(j) == 0 {
		return nil
	}
	return Marshal(j)
}

func (j *Json) FromMaps(maps []map[string]any) {
	if maps == nil || len(maps) == 0 {
		*j = nil
		return
	}
	*j = maps
}

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

func (j Json) Len() int {
	return len(j)
}

func (j Json) IsNil() bool {
	return j == nil
}

func (j Json) IsZero() bool {
	return len(j) == 0
}

func (j *Json) IsEmpty() bool {
	return j == nil || len(*j) == 0
}

func (j Json) Get(index int) map[string]any {
	if index < 0 || index >= len(j) {
		return nil
	}
	return j[index]
}

func (j *Json) Append(m map[string]any) {
	if m == nil {
		return
	}
	*j = append(*j, m)
}

func (j *Json) Remove(index int) {
	if index < 0 || index >= len(*j) {
		return
	}
	*j = append((*j)[:index], (*j)[index+1:]...)
}
