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

package core

// Checker is implemented by types that can validate themselves.
type Checker interface {
	Check() error
}

// Initializer is implemented by types that need initialization before use.
type Initializer interface {
	Init() error
}

// UnmarshalValueser is implemented by types that can unmarshal form values.
type UnmarshalValueser interface {
	UnmarshalValues(vals map[string][]string) error
}

// Common errors shared across frameworks.
var (
	Param_Invalid *Error // placeholder — set by framework adapter
	UnKnown       *Error // placeholder — set by framework adapter
)

// Error is re-exported from pkg/types for convenience.
// Actual type is types.Error; this alias allows core to reference it without
// circular imports. Framework adapters must set these vars at init time.
type Error struct {
	Code    int
	Message string
}
