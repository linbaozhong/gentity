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
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cc := New[[]byte](ctx, "a",
		WithPrefix("abc"),
	)

	e := cc.Save(ctx, "bbb", []byte("456"), time.Second*10)
	if e != nil {
		t.Fatal(e)
	}
	v, e := cc.Fetch(ctx, "bbb")
	if e != nil {
		t.Fatal(e)
	}
	t.Log(string(v))
}

func TestFetchMultiOrderAndPadding(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cc := New[[]byte](ctx, "b", WithPrefix("pfx"))

	if e := cc.Save(ctx, "x", []byte("1")); e != nil {
		t.Fatal(e)
	}
	if e := cc.Save(ctx, "z", []byte("3")); e != nil {
		t.Fatal(e)
	}

	// x 命中、y 缺失、z 命中：结果应与请求顺序一致，缺失项为零值
	vals, e := cc.FetchMulti(ctx, "x", "y", "z")
	if e != nil {
		t.Fatal(e)
	}
	if len(vals) != 3 {
		t.Fatalf("expected len 3, got %d", len(vals))
	}
	if string(vals[0]) != "1" {
		t.Errorf("vals[0] = %q, want 1", vals[0])
	}
	if vals[1] != nil {
		t.Errorf("vals[1] = %q, want nil (miss)", vals[1])
	}
	if string(vals[2]) != "3" {
		t.Errorf("vals[2] = %q, want 3", vals[2])
	}

	// 空请求不应 panic，且返回空切片
	empty, e := cc.FetchMulti(ctx)
	if e != nil {
		t.Fatal(e)
	}
	if len(empty) != 0 {
		t.Errorf("expected empty, got %v", empty)
	}
}

func TestPrefixDeleteBoundary(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cc := New[string](ctx, "c", WithPrefix("pfx"))

	mustSave := func(k, v string) {
		if e := cc.Save(ctx, k, v); e != nil {
			t.Fatal(e)
		}
	}
	mustSave("user:1", "a")
	mustSave("user:2", "b")
	mustSave("username", "c") // 相邻但不同命名空间，不应被误删
	mustSave("other:1", "d")

	if e := cc.PrefixDelete(ctx, "user"); e != nil {
		t.Fatal(e)
	}

	if _, e := cc.Fetch(ctx, "user:1"); e == nil {
		t.Error("user:1 should have been deleted")
	}
	if _, e := cc.Fetch(ctx, "user:2"); e == nil {
		t.Error("user:2 should have been deleted")
	}
	if v, _ := cc.Fetch(ctx, "username"); v != "c" {
		t.Errorf("username should be kept, got %q", v)
	}
	if v, _ := cc.Fetch(ctx, "other:1"); v != "d" {
		t.Errorf("other:1 should be kept, got %q", v)
	}
}
