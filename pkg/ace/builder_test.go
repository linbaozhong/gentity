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

package ace

import (
	"github.com/linbaozhong/gentity/example/model/define/table/tblcompany"
	"testing"
)

func TestConditions(t *testing.T) {
	cs := Conds(tblcompany.Id.Eq(1), tblcompany.State.In())
	for _, c := range cs.ToSlice() {
		s, v := c()
		switch val := v.(type) {
		case []any:
			t.Logf("[]any -- s:%s, v:%v", s, val)
		case error:
			t.Logf("error -- s:%s, v:%v", s, val)
		case [][]any:
			t.Logf("[][]any -- s:%s, v:%v", s, val)
		default:
			t.Logf("default -- s:%s, v:%v", s, val)
		}
	}
}
