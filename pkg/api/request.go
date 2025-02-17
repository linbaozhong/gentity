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

package api

import (
	"encoding/json"
	"errors"
	"io"
)

type UnmarshalValueser interface {
	UnmarshalValues(vals map[string][]string) error
}

// Content-Type为application/json的请求
func ReadJSON(ctx Context, ptr any) error {
	// 从请求体中读取 JSON 数据
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	// 关闭请求体
	defer ctx.Request().Body.Close()

	if len(body) == 0 {
		return errors.New("请求体为空")
	}
	// 解析 JSON 数据到结构体中
	if x, ok := ptr.(json.Unmarshaler); ok {
		return x.UnmarshalJSON(body)
	}
	return json.Unmarshal(body, ptr)
}

// Content-Type为application/x-www-form-urlencoded的请求
func ReadForm(ctx Context, ptr any) error {
	if x, ok := ptr.(UnmarshalValueser); ok {
		values := ctx.FormValues()

		if len(values) == 0 {
			return nil
		}
		return x.UnmarshalValues(values)
	}
	return ctx.ReadForm(ptr)
}

func ReadQuery(ctx Context, ptr any) error {
	if x, ok := ptr.(UnmarshalValueser); ok {
		values := ctx.Request().URL.Query()

		if len(values) == 0 {
			return nil
		}
		return x.UnmarshalValues(values)
	}
	return ctx.ReadQuery(ptr)
}
