package conv

import (
	"github.com/linbaozhong/gentity/pkg/types"
	"time"
	"unsafe"
)

// Bytes2Any r 必须是引用地址
func Bytes2Any(b []byte, r any) error {
	var err error
	switch r.(type) {
	case []byte:
		r = b
	case string:
		r = Bytes2String(b)
	case uint64:
		r = Any2Uint64(Bytes2String(b))
	case uint32:
		r = Any2Uint32(Bytes2String(b))
	case uint16:
		r = Any2Uint16(Bytes2String(b))
	case uint8:
		r = Any2Uint8(Bytes2String(b))
	case int64:
		r = Any2Int64(Bytes2String(b))
	case int32:
		r = Any2Int32(Bytes2String(b))
	case int16:
		r = Any2Int16(Bytes2String(b))
	case int8:
		r = Any2Int8(Bytes2String(b))
	case float32:
		r = Any2Float32(Bytes2String(b))
	case float64:
		r = Any2Float64(Bytes2String(b))
	case bool:
		r = Any2Bool(Bytes2String(b))
	case uint:
		r = Any2Uint(Bytes2String(b))
	case int:
		r = Any2Int(Bytes2String(b))
	case time.Time:
		r, err = time.Parse(time.DateTime, Bytes2String(b))
	default:
		return types.JSON.Unmarshal(b, r)
	}
	return err
}

// func Any2Bytes(s any) ([]byte, error) {
// 	buf := new(bytes.Buffer)
// 	switch v := s.(type) {
// 	case string:
// 		// 直接将字符串转换为字节切片
// 		return String2Bytes(v), nil
// 	case []byte:
// 		// 直接返回字节切片
// 		return v, nil
// 	default:
// 		// 对于其他固定大小的类型，使用 binary.Write
// 		err := binary.Write(buf, binary.BigEndian, v)
// 		if err != nil {
// 			return []byte{}, err
// 		}
// 		return buf.Bytes(), nil
// 	}
// }
//
// func Bytes2Any(b []byte, t any) error {
// 	bufReader := bytes.NewReader(b)
// 	return binary.Read(bufReader, binary.BigEndian, t)
// }

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
