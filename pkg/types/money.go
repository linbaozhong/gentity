package types

import (
	"database/sql/driver"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/conv"
	"strconv"
	"strings"
)

type Money int64

func (m Money) MarshalJSON() ([]byte, error) {
	yuan := strconv.FormatFloat((float64(m) / 100), 'f', -1, 64)
	return []byte(yuan), nil
}

func (m *Money) UnmarshalJSON(b []byte) error {
	c := bytes2String(b)
	fen, e := strconv.ParseFloat(c, 10)
	if e != nil {
		return e
	}
	i, e := strconv.ParseInt(strconv.FormatFloat(fen*100, 'f', 0, 64), 10, 64)
	if e == nil {
		*m = Money(i)
	}
	return e
}

func (m *Money) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*m = 0
		return nil
	case int64:
		*m = Money(v)
		return nil
	case int:
		*m = Money(v)
		return nil
	case int32:
		*m = Money(v)
		return nil
	case int16:
		*m = Money(v)
		return nil
	case int8:
		*m = Money(v)
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Money: %T", src)
	}
}
func (m Money) Value() (driver.Value, error) {
	return int64(m), nil
}

func (m Money) Int() int {
	return int(m)
}

func (m Money) Int64() int64 {
	return int64(m)
}

func (m Money) String() string {
	return strconv.FormatInt(int64(m), 10)
}

func (m Money) Bytes() []byte {
	return conv.Base2Bytes(m)
}

// Yuan 金额分精确到元
func (m Money) Yuan() float64 {
	return float64(m) / 100
}

// 金额分小写转中文大写金额
func (m Money) ToCNY() string {
	if m == 0 {
		return "零元整"
	}
	if m < 0 {
		var mm Money = m * -1
		return "负" + mm.ToCNY()
	}

	numstr := []rune(strconv.Itoa(int(m)))
	numlen := len(numstr)
	moneyUnit := []string{"仟", "佰", "拾", "亿", "仟", "佰", "拾", "万", "仟", "佰", "拾", "元", "角", "分"}
	unit := moneyUnit[len(moneyUnit)-numlen:]
	num := map[rune]string{48: "零", 49: "壹", 50: "贰", 51: "叁", 52: "肆", 53: "伍", 54: "陆", 55: "柒", 56: "捌", 57: "玖"}

	var hasZero bool
	var buf strings.Builder
	for i := 0; i < numlen; i++ {
		if numstr[i] == 48 {
			if strings.Index("亿万", unit[i]) > -1 {
				buf.WriteString(unit[i])
			}
			if hasZero {
				continue
			}
			hasZero = true
		} else {
			if hasZero {
				buf.WriteString("零")
			}
			hasZero = false
			buf.WriteString(num[numstr[i]] + unit[i])
		}
	}

	result := buf.String()
	if strings.HasSuffix(result, "元") || strings.HasSuffix(result, "角") {
		buf.WriteString("整")
	} else if strings.HasSuffix(result, "分") {
	} else {
		buf.WriteString("元整")
	}

	result = strings.Replace(buf.String(), "亿万", "亿", -1)
	return result
}
