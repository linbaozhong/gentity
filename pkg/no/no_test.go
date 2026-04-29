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

package no

import (
	"encoding/json"
	"github.com/linbaozhong/gentity/pkg/types"
	"testing"
)

func TestBigInt_MarshalJSON(t *testing.T) {
	// 生成10个id
	var ids = make([]types.BigInt, 0, 10)
	for i := 0; i < 10; i++ {
		ids = append(ids, GetMustId(0))
	}
	t.Logf("ids: %v", ids)
	// 	ids 转换为 json 字符串
	jsonStr, err := json.Marshal(ids)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	t.Logf("jsonStr: %s", string(jsonStr))
	// 	json 字符串转换为 ids
	var ids2 []types.BigInt
	err = json.Unmarshal(jsonStr, &ids2)
	if err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	t.Logf("ids2: %v", ids2)
}
