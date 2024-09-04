package types

import (
	"time"
)

import "github.com/json-iterator/go"

var (
	JSON = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Smap map[string]interface{}

func NewSmap(size ...int) Smap {
	if len(size) > 0 {
		return make(Smap, size[0])
	}
	return make(Smap)
}

func (p Smap) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	for k, v := range p {
		if vv, ok := v.(ID); ok {
			m[k] = vv.String()
			continue
		}
		if vv, ok := v.(Money); ok {
			m[k] = vv.Yuan()
			continue
		}
		if vv, ok := v.([]byte); ok {
			m[k] = string(vv)
			continue
		}
		if vv, ok := v.(time.Time); ok {
			m[k] = minuteToString(vv)
			continue
		}
		m[k] = v
	}
	return JSON.Marshal(m)
}

func (p Smap) ConvertFrom(m map[string]interface{}) Smap {
	for k, v := range m {
		if vv, ok := v.(map[string]interface{}); ok {
			sm := Smap(vv)
			p[k] = sm
			sm.ConvertFrom(vv)
			continue
		}
		p[k] = v
	}
	return p
}

func (p Smap) Set(k string, v interface{}) Smap {
	p[k] = v
	return p
}

func (p Smap) Get(k string) interface{} {
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