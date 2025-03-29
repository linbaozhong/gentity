package do

import (
	"encoding/json"
	"github.com/linbaozhong/gentity/pkg/types"
	"testing"
)

func TestBills_MarshalJSON(t *testing.T) {
	wallet := NewWallet()
	wallet.Utime = types.Now()

	b, e := json.Marshal(wallet)
	if e != nil {
		t.Error(e)
	}
	t.Log(string(b))
}
