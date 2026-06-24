package conv

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/constraints"
	"reflect"
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

// Any2Numeric converts any numeric value or string to the specified numeric type T.
// Supports int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
// float32, float64, and string types.
// For pointer types, dereferences and converts recursively.
//
// Parameters:
//   - s: source value (any type)
//   - def: default value of type T returned when conversion fails
//
// Returns the converted value of type T, or def[0] if conversion fails.
func Any2Numeric[T constraints.Integer | constraints.Float](s any, def ...T) T {
	switch v := s.(type) {
	case int:
		return T(v)
	case int64:
		return T(v)
	case int32:
		return T(v)
	case int16:
		return T(v)
	case int8:
		return T(v)
	case uint:
		return T(v)
	case uint64:
		return T(v)
	case uint32:
		return T(v)
	case uint16:
		return T(v)
	case uint8:
		return T(v)
	case float64:
		return T(v)
	case float32:
		return T(v)
	case bool:
		if v {
			return T(1)
		}
		return T(0)
	case string:
		var zero T
		switch any(zero).(type) {
		case float64, float32:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return T(f)
			}
		default:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				fmt.Println(i)
				return T(i)
			} else {
				fmt.Println(err)
			}
		}
		if len(def) > 0 {
			return def[0]
		}
		return zero
	default:
		rv := reflect.ValueOf(s)
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				if len(def) > 0 {
					return def[0]
				}
				var zero T
				return zero
			}
			return Any2Numeric[T](rv.Elem().Interface(), def...)
		}
		if len(def) > 0 {
			return def[0]
		}
		var zero T
		return zero
	}
}

func Any2Int(s any, def ...int) int {
	return Any2Numeric(s, def...)
}

func Any2Uint8(s any, def ...uint8) uint8 {
	return Any2Numeric(s, def...)
}

func Any2Uint16(s any, def ...uint16) uint16 {
	return Any2Numeric(s, def...)
}
func Any2Uint(s any, def ...uint) uint {
	return Any2Numeric(s, def...)
}

func Any2Uint32(s any, def ...uint32) uint32 {
	return Any2Numeric(s, def...)
}

func Any2Int32(s any, def ...int32) int32 {
	return Any2Numeric(s, def...)
}

func Any2Int64(s any, def ...int64) int64 {
	return Any2Numeric(s, def...)
}

func Any2Uint64(s any, def ...uint64) uint64 {
	return Any2Numeric(s, def...)
}

func Any2Int8(s any, def ...int8) int8 {
	return Any2Numeric(s, def...)
}
func Any2Int16(s any, def ...int16) int16 {
	return Any2Numeric(s, def...)
}

func Any2Float32(s any, def ...float32) float32 {
	return Any2Numeric(s, def...)
}

func Any2Float64(s any, def ...float64) float64 {
	return Any2Numeric(s, def...)
}

// Any2String converts any value to string type T.
// Supports Stringer interface, numeric types, time.Time, bool, and other types.
// For pointer types, dereferences and converts recursively.
// For complex types, marshals to JSON string.
//
// Parameters:
//   - s: source value (any type)
//   - def: default value of type T returned when conversion fails
//
// Returns the converted value of type T, or def[0] if conversion fails.
func Any2String[T string](s any, def ...T) T {
	if ss, ok := s.(Stringer); ok {
		return T(ss.String())
	}

	var zero T

	switch v := s.(type) {
	case string:
		return T(v)
	case []byte:
		return T(v)
	case uint64:
		return T(strconv.FormatUint(v, 10))
	case uint32:
		return T(strconv.FormatUint(uint64(v), 10))
	case uint16:
		return T(strconv.FormatUint(uint64(v), 10))
	case uint8:
		return T(strconv.FormatUint(uint64(v), 10))
	case uint:
		return T(strconv.FormatUint(uint64(v), 10))
	case int64:
		return T(strconv.FormatInt(v, 10))
	case int32:
		return T(strconv.FormatInt(int64(v), 10))
	case int16:
		return T(strconv.FormatInt(int64(v), 10))
	case int8:
		return T(strconv.FormatInt(int64(v), 10))
	case int:
		return T(strconv.FormatInt(int64(v), 10))
	case float32:
		return T(strconv.FormatFloat(float64(v), 'f', -1, 64))
	case float64:
		return T(strconv.FormatFloat(v, 'f', -1, 64))
	case time.Time:
		if v.IsZero() {
			return zero
		}
		return T(Any2Time(s).Format(time.DateTime))
	case bool:
		return T(strconv.FormatBool(v))
	default:
		rv := reflect.ValueOf(s)
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				if len(def) > 0 {
					return def[0]
				}
				return zero
			}
			return Any2String[T](rv.Elem().Interface(), def...)
		}
		b, e := json.Marshal(v)
		if e != nil {
			if len(def) > 0 {
				return def[0]
			}
			return zero
		}
		return T(b)
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
