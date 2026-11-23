package conv

import (
	"strconv"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/exp/constraints"
)

type Stringer interface {
	String() string
}

func String2Ptr(s string) *string {
	return &s
}

// String2Int8Ptr 字符串转int8指针
func String2Int8Ptr(s string, def ...int8) *int8 {
	return Ptr(String2Int8(s, def...))
}

// String2IntPtr 字符串转int指针
func String2IntPtr(s string, def ...int) *int {
	return Ptr(String2Int(s, def...))
}

// String2Int32Ptr 字符串转int32指针
func String2Int32Ptr(s string, def ...int32) *int32 {
	return Ptr(String2Int32(s, def...))
}

// String2Int64Ptr 字符串转int64指针
func String2Int64Ptr(s string, def ...int64) *int64 {
	return Ptr(String2Int64(s, def...))
}

// /////////////////////

func String2Signed[T constraints.Signed](s string, def ...T) T {
	var zero T
	bits := int(unsafe.Sizeof(zero)) * 8
	if v, e := strconv.ParseInt(strings.TrimSpace(s), 0, bits); e == nil {
		return T(v)
	}
	if len(def) > 0 {
		return def[0]
	}
	return zero
}

func String2Unsigned[T constraints.Unsigned](s string, def ...T) T {
	var zero T
	bits := int(unsafe.Sizeof(zero)) * 8
	if v, e := strconv.ParseUint(strings.TrimSpace(s), 0, bits); e == nil {
		return T(v)
	}
	if len(def) > 0 {
		return def[0]
	}
	return zero
}

func String2Float[T constraints.Float](s string, def ...T) T {
	var zero T
	bits := int(unsafe.Sizeof(zero)) * 8
	if v, e := strconv.ParseFloat(strings.TrimSpace(s), bits); e == nil {
		return T(v)
	}
	if len(def) > 0 {
		return def[0]
	}
	return zero
}

// String2Uint 字符串转uint
func String2Uint(s string, def ...uint) uint {
	return String2Unsigned(s, def...)
}

func String2Uint64(s string, def ...uint64) uint64 {
	return String2Unsigned(s, def...)
}

// String2Uint32 字符串转uint32
func String2Uint32(s string, def ...uint32) uint32 {
	return String2Unsigned[uint32](s, def...)
}

func String2Uint16(s string, def ...uint16) uint16 {
	return String2Unsigned(s, def...)
}

// String2Uint8 字符串转uint
func String2Uint8(s string, def ...uint8) uint8 {
	return String2Unsigned(s, def...)
}

// String2Int
func String2Int(s string, def ...int) int {
	return String2Signed(s, def...)
}

// String2Int8 字符串转int8
func String2Int8(s string, def ...int8) int8 {
	return String2Signed(s, def...)
}

// String2Int16 字符串转int16
func String2Int16(s string, def ...int16) int16 {
	return String2Signed(s, def...)
}

// String2Int32 字符串转int32
func String2Int32(s string, def ...int32) int32 {
	return String2Signed(s, def...)
}

// String2Int64 字符串转int64
func String2Int64(s string, def ...int64) int64 {
	return String2Signed(s, def...)
}

// String2Float32 字符串转float32
func String2Float32(s string, def ...float32) float32 {
	return String2Float(s, def...)
}

// String2Float64 字符串转float64
func String2Float64(s string, def ...float64) float64 {
	return String2Float(s, def...)
}

// String2Time 如果转换失败,返回 def时间(如果存在)
func String2Time(s string, def ...time.Time) time.Time {
	l := len(time.DateTime)
	if len(s) > l {
		s = s[:l]
	}
	if pos := strings.Index(s, "T"); pos > 0 {
		s = strings.Replace(s, "T", " ", 1)
	}
	if b, e := time.ParseInLocation(time.DateTime, s, time.Local); e == nil {
		return b
	} else {
		if _e, ok := e.(*time.ParseError); ok {
			switch _e.LayoutElem {
			case "15", ":", "04", "05":
				b, e = time.Parse(time.DateOnly, s)
			case "2006", "-", "01", "02":
				b, e = time.Parse(time.TimeOnly, s)
			}
			if e == nil {
				return b
			}
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return time.Time{}
}

// String2Bool 字符串转bool
func String2Bool(s string, def ...bool) bool {
	if b, e := strconv.ParseBool(strings.TrimSpace(s)); e == nil {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}
