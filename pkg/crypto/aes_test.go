package crypto

import (
	"strings"
	"testing"
)

func TestCrypto(t *testing.T) {
	_tk, err := AesDecrypt("b0f22a96e94311d2d85721afe1962a74e5e14cbe7b6aff92b05fa1571e60856397abc299eef8cd7695d35490e3061788adeaf86091fb5cb22850eccbc1249c3e91bf4f73b824a535de5bd132b4fc2f5917c0b5258fc987", "snow19921115love")
	if err != nil {
		t.Error(err)
	}
	t.Logf("tk:%s", _tk)
	//token, err := AesEncrypt(_tk, "snow19921115love")
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Logf("token:%s", token)

	pos := strings.Index(_tk, "_")
	if pos < 1 {
		return
	}
	tk := _tk[:pos]
	key := _tk[pos+1:]
	t.Logf("tk:%s,key:%s", tk, key)
}
