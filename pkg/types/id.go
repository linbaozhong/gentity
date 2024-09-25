package types

import (
	"bytes"
	"strconv"
)

type ID int64

func (i ID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + i.String() + `"`), nil
}

func (i *ID) UnmarshalJSON(b []byte) error {
	c := string(bytes.Trim(b, "\""))

	if c == "" {
		*i = ID(0)
		return nil
	}
	tem, e := strconv.ParseInt(c, 10, 64)

	*i = ID(tem)
	return e
}

// String
func (i ID) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// Uint64
func (i ID) Int64() int64 {
	return int64(i)
}
