package db

import (
	"testing"
)

func TestName(t *testing.T) {
	//c := NewCompany()
	//c.FullCorpName = "test"
	//t.Log(c)

	b := Company{}
	b.FullCorpName = "test1"

	a := b
	t.Log(a)
}
