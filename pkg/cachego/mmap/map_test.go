package mmap

import (
	"context"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	m := New(WithExpired(time.Second * 10))
	m.Save(context.Background(), "key", "value")
	v, err := m.Fetch(context.Background(), "key")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(1, v)
	time.Sleep(time.Second * 10)
	v, err = m.Fetch(context.Background(), "key")
	if err != nil {
		//t.Fatal(err)
	}
	t.Log(2, v)

	m.Save(context.Background(), "name", "linbaozhong")
	time.Sleep(time.Second * 5)
	v, err = m.Fetch(context.Background(), "name")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(3, v)
}
