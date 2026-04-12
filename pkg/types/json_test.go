package types

import (
	"testing"
)

type data struct {
	Name map[string]string `json:"name"`
}

func TestName(t *testing.T) {
	var d data
	t.Log(Marshal(d.Name))

}
