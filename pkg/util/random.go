package util

import (
	"math/rand"
	"time"
)

const (
	lowerChars  = "abcdefghijkmnpqrstuvwxyz"
	upperChars  = "ABCDEFGHIJKMNPQRSTUVWXYZ"
	numberChars = "23456789"
	allChars    = lowerChars + upperChars + numberChars
)

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// generateRandString 生成随机字符串
func generateRandString(chars string, l int) string {
	charsLen := len(chars)
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = chars[rnd.Intn(charsLen)]
	}
	return string(bytes)
}

// GetRandLowerString 生成随机小写字母字符串
func GetRandLowerString(l int) string {
	return generateRandString(lowerChars, l)
}

// GetRandUpperString 生成随机大写字母字符串
func GetRandUpperString(l int) string {
	return generateRandString(upperChars, l)
}

// GetRandString 生成随机字符串
func GetRandString(l int) string {
	return generateRandString(allChars, l)
}

// GetRandDigit 生成范围随机数字
func GetRandDigit(min, max int64) int64 {
	if min >= max || min == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
