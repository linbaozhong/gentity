// Copyright © 2023 SnowIM. All rights reserved.
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

import "context"

const (
	VisitorKey = "__Visitor__"
)

type (
	Checker interface {
		Check() error
	}
	Initializer interface {
		Init()
	}
	Visiter interface {
		Visiter(ctx context.Context)
	}
)

// Validate 校验参数
// 注意：如果参数实现了Checker接口，会调用Check方法
func Validate(arg any) error {
	if checker, ok := arg.(Checker); ok {
		return checker.Check()
	}
	return nil
}

// Visit 访问参数
// 注意：如果参数实现了Visiter接口，会调用Visiter方法
func Visit(ctx context.Context, arg any) error {
	if vis, ok := arg.(Visiter); ok {
		vis.Visiter(ctx)
	}
	return nil
}
