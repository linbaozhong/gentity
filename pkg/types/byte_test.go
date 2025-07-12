package types

import (
	"github.com/linbaozhong/gentity/pkg/conv"
	"testing"
)

func TestTypes(t *testing.T) {
	var err error
	var f32 Money = -127
	f, err := conv.Any2Bytes(f32)
	t.Log(f, err)

	f32, err = conv.Bytes2Base[Money](f)
	t.Log(f32, err)

}
