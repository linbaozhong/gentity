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

package iface

import "ganji/pkg/types"

type Checker interface {
	Check() error
}

func Validate(arg any) error {
	if checker, ok := arg.(Checker); ok {
		if err := checker.Check(); err != nil {
			return err
		}
	}
	return nil
}

type Registerer interface {
	Register(visitor *types.UserClaims) error
}

func Register(arg any, fn func() *types.UserClaims) error {
	if registerer, ok := arg.(Registerer); ok {
		return registerer.Register(fn())
	}
	return nil
}
