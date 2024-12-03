package types

import (
	"bytes"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/conv"
	"strconv"
)

type BigInt int64

func (i BigInt) MarshalJSON() ([]byte, error) {
	return conv.String2Bytes(`"` + strconv.FormatInt(int64(i), 10) + `"`), nil
}

func (i *BigInt) UnmarshalJSON(b []byte) error {
	c := string(bytes.Trim(b, "\""))

	if c == "" {
		*i = BigInt(0)
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)

	*i = BigInt(tem)
	return e
}

// String
func (i BigInt) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// int64
func (i BigInt) Int64() int64 {
	return int64(i)
}

// Uint64
func (i BigInt) Uint64() uint64 {
	return uint64(i)
}

func (i *BigInt) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*i = 0
		return nil
	case int64:
		*i = BigInt(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for BigInt: %T", src)
	}
}

func (i *BigInt) Uint() uint {
	return uint(*i)
}
