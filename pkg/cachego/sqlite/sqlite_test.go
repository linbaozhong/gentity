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

package sqlite

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/conv"
	"testing"
)

func TestCache(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cc := New(ctx,
		WithName("a"),
		WithPrefix("abc"),
	)

	// e := cc.Save(ctx, "bbb", "456", time.Second*10)
	// if e != nil {
	// 	t.Fatal(e)
	// }
	v, e := cc.Fetch(ctx, "bbb")
	if e != nil {
		t.Fatal(e)
	}
	t.Log(conv.Bytes2String(v))
}
