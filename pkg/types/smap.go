package types

import (
	"encoding/json"
	"github.com/linbaozhong/gentity/pkg/conv"
	"time"
)

// import "github.com/json-iterator/go"

var (
// json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Smap map[string]any

func NewSmap(size ...int) Smap {
	if len(size) > 0 {
		return make(Smap, size[0])
	}
	return make(Smap)
}

func (p Smap) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	for k, v := range p {
		if vv, ok := v.([]byte); ok {
			m[k] = conv.Bytes2String(vv)
			continue
		}
		if vv, ok := v.(time.Time); ok {
			m[k] = minuteToString(vv)
			continue
		}
		m[k] = v
	}
	return json.Marshal(m)
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
