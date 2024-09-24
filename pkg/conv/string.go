package conv

import (
	"strconv"
	"strings"
	"time"
)

type Stringer interface {
	String() string
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

// String2Int8Ptr 字符串转int8指针
func String2Int8Ptr(s string, def ...int64) *int8 {
	intValue := int8(String2Int64(s, def...))
	return &intValue
}

// String2IntPtr 字符串转int指针
func String2IntPtr(s string, def ...int64) *int {
	intValue := int(String2Int64(s, def...))
	return &intValue
}

// String2Int32Ptr 字符串转int32指针
func String2Int32Ptr(s string, def ...int64) *int32 {
	intValue := int32(String2Int64(s, def...))
	return &intValue
}

// String2Int64Ptr 字符串转int64指针
func String2Int64Ptr(s string, def ...int64) *int64 {
	intValue := String2Int64(s, def...)
	return &intValue
}

func String2Uint64(s string, def ...uint64) uint64 {
	if i, e := strconv.ParseUint(strings.TrimSpace(s), 10, 64); e == nil {
		return i
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Uint32 字符串转uint32
func String2Uint32(s string, def ...uint32) uint32 {
	if i, e := strconv.ParseUint(strings.TrimSpace(s), 10, 32); e == nil {
		return uint32(i)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Uint 字符串转uint
func String2Uint(s string, def ...uint) uint {
	if i, e := strconv.ParseUint(strings.TrimSpace(s), 10, 64); e == nil {
		return uint(i)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func String2Uint16(s string, def ...uint16) uint16 {
	if i, e := strconv.ParseUint(strings.TrimSpace(s), 10, 16); e == nil {
		return uint16(i)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Uint8 字符串转uint
func String2Uint8(s string, def ...uint8) uint8 {
	if i, e := strconv.ParseUint(strings.TrimSpace(s), 10, 8); e == nil {
		return uint8(i)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Int
func String2Int(s string, def ...int) int {
	if i, e := strconv.Atoi(strings.TrimSpace(s)); e == nil {
		return i
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Int8 字符串转int8
func String2Int8(s string, def ...int8) int8 {
	if b, e := strconv.ParseInt(strings.TrimSpace(s), 10, 8); e == nil {
		return int8(b)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Int16 字符串转int16
func String2Int16(s string, def ...int16) int16 {
	if b, e := strconv.ParseInt(strings.TrimSpace(s), 10, 16); e == nil {
		return int16(b)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Int32 字符串转int32
func String2Int32(s string, def ...int32) int32 {
	if b, e := strconv.ParseInt(strings.TrimSpace(s), 10, 32); e == nil {
		return int32(b)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Int64 字符串转int64
func String2Int64(s string, def ...int64) int64 {
	if b, e := strconv.ParseInt(strings.TrimSpace(s), 10, 64); e == nil {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Float32 字符串转float32
func String2Float32(s string, def ...float32) float32 {
	if b, e := strconv.ParseFloat(strings.TrimSpace(s), 32); e == nil {
		return float32(b)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Float64 字符串转float64
func String2Float64(s string, def ...float64) float64 {
	if b, e := strconv.ParseFloat(strings.TrimSpace(s), 64); e == nil {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// String2Time 如果转换失败,返回 def时间(如果存在)
func String2Time(s string, def ...time.Time) time.Time {
	if b, e := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); e == nil {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return time.Time{}
}
