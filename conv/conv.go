// Package conv provides safe type conversion utilities.
// Similar to spf13/cast but with explicit error returns for unsafe conversions.
package conv

import (
	"fmt"
	"strconv"
)

// ToString converts any value to its string representation.
// Unlike fmt.Sprint, it handles nil explicitly and omits braces for maps/slices.
func ToString(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case bool:
		return strconv.FormatBool(x)
	case int:
		return strconv.Itoa(x)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case fmt.Stringer:
		return x.String()
	case error:
		return x.Error()
	default:
		return fmt.Sprint(v)
	}
}

// ToBool converts v to bool.
func ToBool(v any) (bool, error) {
	switch x := v.(type) {
	case bool:
		return x, nil
	case string:
		return strconv.ParseBool(x)
	case int, int8, int16, int32, int64:
		n, _ := ToInt64(v)
		return n != 0, nil
	case uint, uint8, uint16, uint32, uint64:
		u, _ := ToUint64(v)
		return u != 0, nil
	case float32, float64:
		f, _ := ToFloat64(v)
		return f != 0, nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

// MustBool converts v to bool, panicking on error.
func MustBool(v any) bool {
	b, err := ToBool(v)
	if err != nil {
		panic(err)
	}
	return b
}

// ToInt converts v to int.
func ToInt(v any) (int, error) {
	n, err := ToInt64(v)
	return int(n), err
}

// ToInt64 converts v to int64.
func ToInt64(v any) (int64, error) {
	switch x := v.(type) {
	case int:
		return int64(x), nil
	case int8:
		return int64(x), nil
	case int16:
		return int64(x), nil
	case int32:
		return int64(x), nil
	case int64:
		return x, nil
	case uint:
		return int64(x), nil
	case uint8:
		return int64(x), nil
	case uint16:
		return int64(x), nil
	case uint32:
		return int64(x), nil
	case uint64:
		return int64(x), nil
	case float32:
		return int64(x), nil
	case float64:
		return int64(x), nil
	case bool:
		if x {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.ParseInt(x, 0, 64)
	case []byte:
		return strconv.ParseInt(string(x), 0, 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

// MustInt64 converts v to int64, panicking on error.
func MustInt64(v any) int64 {
	n, err := ToInt64(v)
	if err != nil {
		panic(err)
	}
	return n
}

// ToUint64 converts v to uint64.
func ToUint64(v any) (uint64, error) {
	switch x := v.(type) {
	case uint:
		return uint64(x), nil
	case uint8:
		return uint64(x), nil
	case uint16:
		return uint64(x), nil
	case uint32:
		return uint64(x), nil
	case uint64:
		return x, nil
	case int, int8, int16, int32, int64:
		n, err := ToInt64(v)
		if err != nil {
			return 0, err
		}
		if n < 0 {
			return 0, fmt.Errorf("cannot convert negative int %d to uint64", n)
		}
		return uint64(n), nil
	case float32:
		return uint64(x), nil
	case float64:
		return uint64(x), nil
	case bool:
		if x {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.ParseUint(x, 0, 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to uint64", v)
	}
}

// ToFloat64 converts v to float64.
func ToFloat64(v any) (float64, error) {
	switch x := v.(type) {
	case float32:
		return float64(x), nil
	case float64:
		return x, nil
	case int, int8, int16, int32, int64:
		n, err := ToInt64(v)
		return float64(n), err
	case uint, uint8, uint16, uint32, uint64:
		u, err := ToUint64(v)
		return float64(u), err
	case bool:
		if x {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.ParseFloat(x, 64)
	case []byte:
		return strconv.ParseFloat(string(x), 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

// MustFloat64 converts v to float64, panicking on error.
func MustFloat64(v any) float64 {
	f, err := ToFloat64(v)
	if err != nil {
		panic(err)
	}
	return f
}

// ToStringSlice converts a slice of any to []string.
func ToStringSlice(v []any) []string {
	result := make([]string, len(v))
	for i, x := range v {
		result[i] = ToString(x)
	}
	return result
}

// ToIntSlice converts a slice of any to []int64, returning an error if any element fails.
func ToIntSlice(v []any) ([]int64, error) {
	result := make([]int64, len(v))
	for i, x := range v {
		n, err := ToInt64(x)
		if err != nil {
			return nil, fmt.Errorf("element %d: %w", i, err)
		}
		result[i] = n
	}
	return result, nil
}
