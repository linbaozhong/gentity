package types

import (
	"github.com/linbaozhong/gentity/pkg/conv"
	"testing"
)

func TestTypes(t *testing.T) {
	var err error
	var f32 Float32 = -127.33
	f, err := conv.Any2Bytes(f32)
	//f := f32.Bytes()
	t.Log(f, err)

	var f322 Float32
	err = conv.Bytes2Any(f, &f322)
	//f322.FromBytes(f)
	t.Log(f322, err)

}
