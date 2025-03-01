package conv

import (
	"encoding/binary"
	"errors"
	"math"
	"runtime"
	"time"
)

var ErrTooShort = errors.New("bytes.Buffer: too short")

// Bytes2Any 泛型数据类型，必须是引用类型
type b2a interface {
	*string | *bool | *[]byte | *time.Time | *int | *int8 | *int16 | *int32 | *int64 | *uint | *uint8 | *uint16 | *uint32 | *uint64 | *float32 | *float64 | any
}

// Bytes2Any r 必须是引用地址
func Bytes2Any[T b2a](b []byte, r T) error {
	var err error
	switch v := any(r).(type) {
	case *[]byte:
		*v = b
	case *string:
		*v = Bytes2String(b)
	case *uint64:
		if len(b) >= 8 {
			*v = binary.BigEndian.Uint64(b)
		} else {
			return ErrTooShort
		}
	case *uint32:
		if len(b) >= 4 {
			*v = binary.BigEndian.Uint32(b)
		} else {
			return ErrTooShort
		}
	case *uint16:
		if len(b) >= 2 {
			*v = binary.BigEndian.Uint16(b)
		} else {
			return ErrTooShort
		}
	case *uint8:
		if len(b) >= 1 {
			*v = b[0]
		} else {
			return ErrTooShort
		}
	case *int64:
		if len(b) >= 8 {
			*v = int64(binary.BigEndian.Uint64(b))
		} else {
			return ErrTooShort
		}
	case *int32:
		if len(b) >= 4 {
			*v = int32(binary.BigEndian.Uint32(b))
		} else {
			return ErrTooShort
		}
	case *int16:
		if len(b) >= 2 {
			*v = int16(binary.BigEndian.Uint16(b))
		} else {
			return ErrTooShort
		}
	case *int8:
		if len(b) >= 1 {
			*v = int8(b[0])
		} else {
			return ErrTooShort
		}
	case *float32:
		if len(b) >= 4 {
			*v = math.Float32frombits(binary.BigEndian.Uint32(b))
		} else {
			return ErrTooShort
		}
	case *float64:
		if len(b) >= 8 {
			*v = math.Float64frombits(binary.BigEndian.Uint64(b))
		} else {
			return ErrTooShort
		}
	case *int:
		if len(b) >= 8 {
			*v = int(binary.BigEndian.Uint64(b))
		} else if len(b) >= 4 {
			*v = int(binary.BigEndian.Uint32(b))
		} else {
			return ErrTooShort
		}
	case *uint:
		if len(b) >= 8 {
			*v = uint(binary.BigEndian.Uint64(b))
		} else if len(b) >= 4 {
			*v = uint(binary.BigEndian.Uint32(b))
		} else {
			return ErrTooShort
		}
	case *bool:
		if len(b) >= 1 {
			*v = b[0] != 0
		} else {
			return ErrTooShort
		}
	case *time.Time:
		if len(b) > 0 {
			*v, err = time.Parse(time.RFC3339Nano, Bytes2String(b))
			if err != nil {
				return err
			}
		} else {
			return ErrTooShort
		}
	default:
		return json.Unmarshal(b, r)
	}
	return nil
}

// Any2Bytes 泛型数据类型
type a2b interface {
	string | bool | []byte | time.Time | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | any
}

func Any2Bytes[T a2b](s T) ([]byte, error) {
	var buf []byte
	switch v := any(s).(type) {
	case string:
		// 直接将字符串转换为字节切片
		buf = String2Bytes(v)
	case []byte:
		// 直接返回字节切片
		buf = v
	case bool:
		if v {
			buf = []byte{1}
		} else {
			buf = []byte{0}
		}
	// 处理 time.Time 类型
	case time.Time:
		buf = []byte(v.Format(time.RFC3339Nano))
	case uint64:
		buf = make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v)
	case uint32:
		buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, v)
	case uint16:
		buf = make([]byte, 2)
		binary.BigEndian.PutUint16(buf, v)
	case uint8:
		buf = []byte{byte(v)}
	case int64:
		buf = make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(v))
	case int32:
		buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
	case int16:
		buf = make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(v))
	case int8:
		buf = []byte{byte(v)}
	case float32:
		buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, math.Float32bits(v))
	case float64:
		buf = make([]byte, 8)
		binary.BigEndian.PutUint64(buf, math.Float64bits(v))
	case uint:
		if runtime.GOARCH == "arm64" || runtime.GOARCH == "amd64" {
			buf = make([]byte, 8)
			binary.BigEndian.PutUint64(buf, uint64(v))
		} else {
			buf = make([]byte, 4)
			binary.BigEndian.PutUint32(buf, uint32(v))
		}
	case int:
		if runtime.GOARCH == "arm64" || runtime.GOARCH == "amd64" {
			buf = make([]byte, 8)
			binary.BigEndian.PutUint64(buf, uint64(v))
		} else {
			buf = make([]byte, 4)
			binary.BigEndian.PutUint32(buf, uint32(v))
		}
	default:
		return json.Marshal(v)
	}
	return buf, nil
}

// Bytes2String converts byte slice to string.
func Bytes2String(b []byte) string {
	// return *(*string)(unsafe.Pointer(&b))
	return string(b)
}

// String2Bytes converts string to byte slice.
func String2Bytes(s string) []byte {
	// return *(*[]byte)(unsafe.Pointer(
	// 	&struct {
	// 		string
	// 		Cap int
	// 	}{s, len(s)},
	// ))
	return []byte(s)
}
