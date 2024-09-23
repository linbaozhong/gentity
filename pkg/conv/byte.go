package conv

import (
	"bytes"
	"encoding/binary"
	"github.com/linbaozhong/gentity/pkg/types"
	"time"
	"unsafe"
)

// Bytes2Interface r 必须是引用地址
func Bytes2Interface(b []byte, r any) error {
	var err error
	switch r.(type) {
	case []byte:
		r = b
	case string:
		r = Bytes2String(b)
	case uint64:
		r = Interface2Uint64(Bytes2String(b))
	case uint32:
		r = Interface2Uint32(Bytes2String(b))
	case uint16:
		r = Interface2Uint16(Bytes2String(b))
	case uint8:
		r = Interface2Uint8(Bytes2String(b))
	case int64:
		r = Interface2Int64(Bytes2String(b))
	case int32:
		r = Interface2Int32(Bytes2String(b))
	case int16:
		r = Interface2Int16(Bytes2String(b))
	case int8:
		r = Interface2Int8(Bytes2String(b))
	case float32:
		r = Interface2Float32(Bytes2String(b))
	case float64:
		r = Interface2Float64(Bytes2String(b))
	case bool:
		r = Interface2Bool(Bytes2String(b))
	case uint:
		r = Interface2Uint(Bytes2String(b))
	case int:
		r = Interface2Int(Bytes2String(b))
	case time.Time:
		r, err = time.Parse(time.DateTime, Bytes2String(b))
	default:
		return types.JSON.Unmarshal(b, r)
	}
	return err
}
func any2Bytes(s any) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, s)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
func bytes2Any(b []byte, t any) error {
	bufReader := bytes.NewReader(b)
	return binary.Read(bufReader, binary.BigEndian, t)
}

// Bytes2String converts byte slice to string.
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Bytes converts string to byte slice.
func String2Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
