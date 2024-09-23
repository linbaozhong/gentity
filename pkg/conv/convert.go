// @Title 各种类型转换方法
// @Description
// @Author 蔺保仲 2020/04/20
// @Update 蔺保仲 2020/04/20
package conv

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

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

// Wrap 将float64转成精确的int64, 小数点向右移动retain位
func Wrap(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}

// Unwrap 将int64恢复成正常的float64，小数点向左移动retain位
func Unwrap(num int64, retain int) float64 {
	return float64(num) / math.Pow10(retain)
}

// WrapToFloat64 精准float64，小数点向右移动retain位
func WrapToFloat64(num float64, retain int) float64 {
	return num * math.Pow10(retain)
}

// UnwrapToInt64 精准int64，小数点向左移动retain位
func UnwrapToInt64(num int64, retain int) int64 {
	return int64(Unwrap(num, retain))
}

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
