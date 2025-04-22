package types

import (
	"bytes"
	"errors"
	"github.com/linbaozhong/gentity/pkg/gjson"
	"strconv"
	"time"
)

type Smap map[string]any

func NewSmap(size ...int) Smap {
	if len(size) > 0 {
		return make(Smap, size[0])
	}
	return make(Smap)
}

func (p Smap) MarshalJSON() ([]byte, error) {
	var _buf bytes.Buffer
	_buf.WriteString("{")
	for k, v := range p {
		_buf.WriteString(strconv.Quote(k) + `:` + Marshal(v) + `,`)
	}
	if _buf.Len() > 1 {
		_buf.Truncate(_buf.Len() - 1)
	}
	_buf.WriteString("}")
	return _buf.Bytes(), nil
}

// UnmarshalJSON 反序列化
func (p *Smap) UnmarshalJSON(data []byte) error {
	ok := gjson.ValidBytes(data)
	if !ok {
		return errors.New("invalid json")
	}

	_result := gjson.ParseBytes(data)
	_result.ForEach(func(key, value gjson.Result) bool {
		switch value.Type {
		case gjson.String:
			p.Set(key.Str, value.Str)
		case gjson.Number:
			p.Set(key.Str, value.Num)
		case gjson.True:
			p.Set(key.Str, true)
		case gjson.False:
			p.Set(key.Str, false)
		case gjson.Null:
			p.Set(key.Str, nil)
		case gjson.JSON:
			_p := NewSmap()
			e := Unmarshal(value, _p)
			if e != nil {
				return false
			}
			p.Set(key.Str, _p)
		}
		return true
	})
	return nil
}

func (p Smap) ConvertFrom(m map[string]any) Smap {
	for k, v := range m {
		if vv, ok := v.(map[string]any); ok {
			sm := Smap(vv)
			p[k] = sm
			sm.ConvertFrom(vv)
			continue
		}
		p[k] = v
	}
	return p
}

func (p Smap) Set(k string, v any) Smap {
	p[k] = v
	return p
}

func (p Smap) Get(k string) any {
	return p[k]
}

func (p Smap) Remove(k string) {
	delete(p, k)
}

// MinuteToString 日期分钟转字符串
func minuteToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}
