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

import "fmt"

type Bytes []byte

// ////////////////////////////
// Byte
func (b *Bytes) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*b = []byte{}
	case []byte:
		*b = v
	default:
		return fmt.Errorf("unsupported scan type for Byte: %T", src)
	}
	return nil
}

func (b Bytes) String() string {
	return string(b)
}
