// @Title 各种类型转换方法
// @Description
// @Author 蔺保仲 2020/04/20
// @Update 蔺保仲 2020/04/20
package conv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/types"
	"math"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

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

func Interface2Time(s any, def ...time.Time) time.Time {
	if b, ok := s.(time.Time); ok {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return time.Time{}
}

func Interface2Int(s any, def ...int) int {
	switch v := s.(type) {
	case uint64:
		return int(v)
	case uint32:
		return int(v)
	case uint16:
		return int(v)
	case uint8:
		return int(v)
	case uint:
		return int(v)
	case int64:
		return int(v)
	case int32:
		return int(v)
	case int16:
		return int(v)
	case int8:
		return int(v)
	case int:
		return v
	case string:
		return String2Int(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Interface2Uint8(s any, def ...uint8) uint8 {
	switch v := s.(type) {
	case uint64:
		return uint8(v)
	case uint32:
		return uint8(v)
	case uint16:
		return uint8(v)
	case uint8:
		return v
	case uint:
		return uint8(v)
	case int64:
		return uint8(v)
	case int32:
		return uint8(v)
	case int16:
		return uint8(v)
	case int8:
		return uint8(v)
	case int:
		return uint8(v)
	case string:
		return String2Uint8(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}
func Interface2Uint(s any, def ...uint) uint {
	switch v := s.(type) {
	case uint64:
		return uint(v)
	case uint32:
		return uint(v)
	case uint16:
		return uint(v)
	case uint8:
		return uint(v)
	case uint:
		return v
	case int64:
		return uint(v)
	case int32:
		return uint(v)
	case int16:
		return uint(v)
	case int8:
		return uint(v)
	case int:
		return uint(v)
	case string:
		return String2Uint(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Interface2Uint32(s any, def ...uint32) uint32 {
	switch v := s.(type) {
	case uint64:
		return uint32(v)
	case uint32:
		return v
	case uint16:
		return uint32(v)
	case uint8:
		return uint32(v)
	case uint:
		return uint32(v)
	case int64:
		return uint32(v)
	case int32:
		return uint32(v)
	case int16:
		return uint32(v)
	case int8:
		return uint32(v)
	case int:
		return uint32(v)
	case string:
		return String2Uint32(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Interface2Int64(s any, def ...int64) int64 {
	if b, ok := s.(int64); ok {
		return b
	}
	switch v := s.(type) {
	case uint64:
		return int64(v)
	case uint32:
		return int64(v)
	case uint16:
		return int64(v)
	case uint8:
		return int64(v)
	case uint:
		return int64(v)
	case int32:
		return int64(v)
	case int16:
		return int64(v)
	case int8:
		return int64(v)
	case int:
		return int64(v)
	case string:
		return String2Int64(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Interface2Uint64(s any, def ...uint64) uint64 {
	switch v := s.(type) {
	case uint64:
		return v
	case uint32:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint:
		return uint64(v)
	case int64:
		return uint64(v)
	case int32:
		return uint64(v)
	case int16:
		return uint64(v)
	case int8:
		return uint64(v)
	case int:
		return uint64(v)
	case string:
		return String2Uint64(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Interface2Int8(s any, def ...int8) int8 {
	switch v := s.(type) {
	case uint64:
		return int8(v)
	case uint32:
		return int8(v)
	case uint16:
		return int8(v)
	case uint8:
		return int8(v)
	case uint:
		return int8(v)
	case int64:
		return int8(v)
	case int32:
		return int8(v)
	case int16:
		return int8(v)
	case int8:
		return v
	case int:
		return int8(v)
	case string:
		return String2Int8(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}
func Interface2Int16(s any, def ...int16) int16 {
	if b, ok := s.(int16); ok {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func Interface2String(s any) string {
	switch v := s.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		if s.(time.Time).IsZero() {
			return ""
		}
		return Interface2Time(s).Format(time.DateTime)
	case bool:
		b := s.(bool)
		if b {
			return "1"
		}
		return "0"
	default:
		return fmt.Sprintf("%+v", s)
	}
}

// Interface2Bytes
func Interface2Bytes(s any) ([]byte, error) {
	switch v := s.(type) {
	case []byte:
		return v, nil
	case string:
		return StringToBytes(v), nil
	case uint64, uint32, uint16, uint8, int64, int32, int16, int8, float32, float64, bool:
		return any2Bytes(v)
	case uint:
		return any2Bytes(uint64(v))
	case int:
		return any2Bytes(int64(v))
	case time.Time:
		if s.(time.Time).IsZero() {
			return []byte{}, nil
		}
		return StringToBytes(Interface2Time(s).Format(time.DateTime)), nil
	default:
		return types.JSON.Marshal(v)
	}
}

// Bytes2Interface r 必须是引用地址
func Bytes2Interface(b []byte, r any) error {
	var err error
	switch v := r.(type) {
	case []byte:
		r = b
	case string:
		r = BytesToString(b)
	case uint64, uint32, uint16, uint8, int64, int32, int16, int8, float32, float64, bool:
		return bytes2Any(b, r)
	case uint:
		r = uint64(v)
		return bytes2Any(b, r)
	case int:
		r = int64(v)
		return bytes2Any(b, r)
	case time.Time:
		r, err = time.Parse(time.DateTime, BytesToString(b))
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

// BytesToString converts byte slice to string.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts string to byte slice.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// Interface2StringSlice any转[]string
func Interface2StringSlice(s any, def ...[]string) []string {
	if b, ok := s.([]string); ok {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return []string{}
}

// Interface2IntSlice any转[]int
func Interface2IntSlice(s any, def ...[]int) []int {
	if b, ok := s.([]int); ok {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return []int{}
}

// Interface2StringMap any转map[string]any
func Interface2StringMap(s any, def ...map[string]any) map[string]any {
	if b, ok := s.(map[string]any); ok {
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return map[string]any{}
}

// IntToFloat64 IntToFloat64
func IntToFloat64(i int) float64 {
	intValueString := strconv.Itoa(i)
	value, err := strconv.ParseFloat(intValueString, 64)
	if err != nil {
		return 0
	}
	return value
}

func Uint8ToString(i uint8) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Int16ToString Int16ToString
func Int16ToString(i int16) string {
	valueString := strconv.FormatInt(int64(i), 10)
	return valueString
}

// Int32ToString Int32ToString
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

// UintToString
func UintToString(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Uint16ToString
func Uint16ToString(i uint16) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Uint32ToString
func Uint32ToString(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Uint64ToString
func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

// IntToString IntToString
func IntToString(i int) string {
	valueString := strconv.Itoa(i)
	return valueString
}

// Int64ToString Int64ToString
func Int64ToString(i int64) string {
	valueString := strconv.FormatInt(i, 10)
	return valueString
}

// Int8ToString Int8ToString
func Int8ToString(i int8) string {
	valueString := strconv.FormatInt(int64(i), 10)
	return valueString
}

// Float32ToString Float32ToString
func Float32ToString(i float32) string {
	valueString := strconv.FormatFloat(float64(i), 'f', -1, 32)
	return valueString
}

// Float64ToString Float64ToString
func Float64ToString(i float64) string {
	valueString := strconv.FormatFloat(i, 'f', -1, 64)
	return valueString
}

// Wrap 将float64转成精确的int64
func Wrap(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}

// Unwrap 将int64恢复成正常的float64
func Unwrap(num int64, retain int) float64 {
	return float64(num) / math.Pow10(retain)
}

// WrapToFloat64 精准float64
func WrapToFloat64(num float64, retain int) float64 {
	return num * math.Pow10(retain)
}

// UnwrapToInt64 精准int64
func UnwrapToInt64(num int64, retain int) int64 {
	return int64(Unwrap(num, retain))
}

// // 处理float64精度,保留n位小数
// func Round(f float64, n int) float64 {
//	n10 := math.Pow10(n)
//	return math.Trunc((f+0.5/n10)*n10) / n10
// }

// 处理float64精度,保留n位小数,可以是负数
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)

	if f > 0 {
		return math.Trunc((f+0.5/n10)*n10) / n10
	}
	return math.Trunc((f-0.5/n10)*n10) / n10
}

// Hex2Dec 十六进制转十进制
func Hex2Dec(val string) uint64 {
	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	return n
}

// YuanString2Fen 金额元字符串 转 分
func YuanString2Fen(s string) int64 {
	s += "00"
	pos := strings.IndexByte(s, '.')
	if pos == 0 { // 第0位是。
		return String2Int64(s[1:3])
	} else if pos > 0 {
		return String2Int64(s[:pos] + s[pos+1:pos+3])
	}
	return String2Int64(s)
}
