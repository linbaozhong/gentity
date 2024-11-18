package conv

import (
	"fmt"
	"github.com/linbaozhong/gentity/pkg/types"
	"strconv"
	"time"
)

func Any2Time(s any, def ...time.Time) time.Time {
	switch t := s.(type) {
	case time.Time:
		return t
	case string:
		return String2Time(t, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return time.Time{}
	}
}

func Any2Int(s any, def ...int) int {
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
	case float64:
		return int(v)
	case float32:
		return int(v)
	case string:
		return String2Int(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Any2Uint8(s any, def ...uint8) uint8 {
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
	case float64:
		return uint8(v)
	case float32:
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

func Any2Uint16(s any, def ...uint16) uint16 {
	switch v := s.(type) {
	case uint64:
		return uint16(v)
	case uint32:
		return uint16(v)
	case uint16:
		return v
	case uint8:
		return uint16(v)
	case uint:
		return uint16(v)
	case int64:
		return uint16(v)
	case int32:
		return uint16(v)
	case int16:
		return uint16(v)
	case int8:
		return uint16(v)
	case int:
		return uint16(v)
	case float64:
		return uint16(v)
	case float32:
		return uint16(v)
	case string:
		return String2Uint16(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}
func Any2Uint(s any, def ...uint) uint {
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
	case float64:
		return uint(v)
	case float32:
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

func Any2Uint32(s any, def ...uint32) uint32 {
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
	case float64:
		return uint32(v)
	case float32:
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

func Any2Int32(s any, def ...int32) int32 {
	switch v := s.(type) {
	case uint64:
		return int32(v)
	case uint32:
		return int32(v)
	case uint16:
		return int32(v)
	case uint8:
		return int32(v)
	case uint:
		return int32(v)
	case int64:
		return int32(v)
	case int32:
		return v
	case int16:
		return int32(v)
	case int8:
		return int32(v)
	case int:
		return int32(v)
	case float64:
		return int32(v)
	case float32:
		return int32(v)
	case string:
		return String2Int32(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Any2Int64(s any, def ...int64) int64 {
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
	case float64:
		return int64(v)
	case float32:
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

func Any2Uint64(s any, def ...uint64) uint64 {
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
	case float64:
		return uint64(v)
	case float32:
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

func Any2Int8(s any, def ...int8) int8 {
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
	case float64:
		return int8(v)
	case float32:
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
func Any2Int16(s any, def ...int16) int16 {
	switch v := s.(type) {
	case uint64:
		return int16(v)
	case uint32:
		return int16(v)
	case uint16:
		return int16(v)
	case uint8:
		return int16(v)
	case uint:
		return int16(v)
	case int64:
		return int16(v)
	case int32:
		return int16(v)
	case int16:
		return v
	case int8:
		return int16(v)
	case int:
		return int16(v)
	case float64:
		return int16(v)
	case float32:
		return int16(v)
	case string:
		return String2Int16(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Any2Float32(s any, def ...float32) float32 {
	switch v := s.(type) {
	case float64:
		return float32(v)
	case float32:
		return v
	case uint64:
		return float32(v)
	case uint32:
		return float32(v)
	case uint16:
		return float32(v)
	case uint8:
		return float32(v)
	case uint:
		return float32(v)
	case int64:
		return float32(v)
	case int32:
		return float32(v)
	case int16:
		return float32(v)
	case int8:
		return float32(v)
	case int:
		return float32(v)
	case string:
		return String2Float32(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Any2Float64(s any, def ...float64) float64 {
	switch v := s.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case uint64:
		return float64(v)
	case uint32:
		return float64(v)
	case uint16:
		return float64(v)
	case uint8:
		return float64(v)
	case uint:
		return float64(v)
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	case int16:
		return float64(v)
	case int8:
		return float64(v)
	case int:
		return float64(v)
	case string:
		return String2Float64(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func Any2String(s any) string {
	if ss, ok := s.(Stringer); ok {
		return ss.String()
	}
	switch v := s.(type) {
	case string:
		return v
	case []byte:
		return Bytes2String(v)
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
		return Any2Time(s).Format(time.DateTime)
	case bool:
		return strconv.FormatBool(v)
	default:
		b, e := types.JSON.Marshal(v)
		if e != nil {
			return fmt.Sprintf("%+v", v)
		}
		return Bytes2String(b)
	}
}

func Any2Bool(s any, def ...bool) bool {
	switch v := s.(type) {
	case bool:
		return v
	case int64, int32, int16, int8, int, uint, uint64, uint32, uint16, uint8, float32, float64:
		return v != 0
	case string:
		return String2Bool(v, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return false
	}
}

// Any2Bytes
func Any2Bytes(s any) ([]byte, error) {
	switch v := s.(type) {
	case []byte:
		return v, nil
	case string:
		return String2Bytes(v), nil
	case uint64, uint32, uint16, uint8, int64, int32, int16, int8, float32, float64, bool, int, uint:
		return String2Bytes(Any2String(v)), nil
	case time.Time:
		if s.(time.Time).IsZero() {
			return []byte{}, nil
		}
		return String2Bytes(Any2Time(s).Format(time.DateTime)), nil
	default:
		return types.JSON.Marshal(v)
	}
}
