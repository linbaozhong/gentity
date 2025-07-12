package types

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"strconv"
)

type BigInt uint64

func (i BigInt) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(i.String())), nil
}

func (i *BigInt) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)

	if c == "" {
		*i = BigInt(0)
		return nil
	}
	tem, e := strconv.ParseUint(c, 10, 64)
	*i = BigInt(tem)
	return e
}

// String
func (i BigInt) String() string {
	return strconv.FormatUint(i.Uint64(), 10)
}

// int64
func (i BigInt) Int64() int64 {
	return int64(i)
}

// Uint64
func (i BigInt) Uint64() uint64 {
	return uint64(i)
}

func (i BigInt) Uint() uint {
	return uint(i)
}

func (i BigInt) Bytes() []byte {
	return binary.BigEndian.AppendUint64(nil, uint64(i))
}

func (i *BigInt) FromBytes(b []byte) {
	*i = BigInt(binary.BigEndian.Uint64(b))
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

func (i BigInt) Value() (driver.Value, error) {
	return int64(i), nil
}
