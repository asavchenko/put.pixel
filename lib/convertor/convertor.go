package convertor

import (
	"assa.com/put.pixel/lib/log"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ToTime converts unknown value to time.Time
func ToTime(i interface{}) (time.Time, error) {
	switch v := i.(type) {
	case string:
		if t, err := time.Parse("2006-01-02 15:04:05.999 -0700 MST", v); err != nil {
			if t, err := time.Parse("20060102150405", v); err != nil {
				if t, err := time.Parse("02012006150405", v); err != nil {
					if t, err := time.Parse("2006-01-02 15:04:05", v); err != nil {
						if t, err := time.Parse("2006-01", v); err != nil {
							if t, err := time.Parse("2006-1", v); err != nil {
								if t, err := time.Parse("2006-01-02", v); err != nil {
									if t, err := time.Parse("2006-01-02T15:04:05.999Z07:00", v); err != nil {
										i := ToInt64IgnoreErrors(v)
										if i > 1000000000000 {
											i = i / 1000
										}
										if i > 10000 { // trying timestamp
											return time.Unix(i, 0).UTC(), nil
										}
										return time.Time{}, err
									} else {
										return t.UTC(), nil
									}
								} else {
									return t.UTC(), nil
								}
							} else {
								return t.UTC(), nil
							}
						} else {
							return t.UTC(), nil
						}
					} else {
						return t.UTC(), nil
					}
				} else {
					return t.UTC(), nil
				}
			} else {
				return t.UTC(), nil
			}
		} else {
			return t.UTC(), nil
		}

	case time.Time:
		return v.UTC(), nil
	case int:
		return time.Unix(int64(v), 0).UTC(), nil
	case int8:
		return time.Unix(int64(v), 0).UTC(), nil
	case int16:
		return time.Unix(int64(v), 0).UTC(), nil
	case int32:
		return time.Unix(int64(v), 0).UTC(), nil
	case int64:
		return time.Unix(v, 0).UTC(), nil
	case nil:
		return time.Time{}, nil
	default:
		return time.Time{}, fmt.Errorf("on conversion to time.Time got unexpected value %#v", i)
	}
}

func ToTimeInLocation(i interface{}, loc *time.Location) (time.Time, error) {
	switch v := i.(type) {

	case string:
		if t, err := time.ParseInLocation("2006-01-02 15:04:05.999 -0700 MST", v, loc); err != nil {
			if t, err := time.ParseInLocation("20060102150405", v, loc); err != nil {
				if t, err := time.ParseInLocation("02012006150405", v, loc); err != nil {
					if t, err := time.ParseInLocation("2006-01-02 15:04:05", v, loc); err != nil {
						if t, err := time.ParseInLocation("2006-01", v, loc); err != nil {
							if t, err := time.ParseInLocation("2006-1", v, loc); err != nil {
								if t, err := time.ParseInLocation("2006-01-02", v, loc); err != nil {
									if t, err := time.ParseInLocation("2006-01-02T15:04:05.999Z07:00", v, loc); err != nil {
										i := ToInt64IgnoreErrors(v)
										if i > 1000000000000 {
											i = i / 1000
										}
										if i > 10000 { // trying timestamp
											return time.Unix(i, 0).UTC(), nil
										}
										return time.Time{}, err
									} else {
										return t.UTC(), nil
									}
								} else {
									return t.UTC(), nil
								}
							} else {
								return t.UTC(), nil
							}
						} else {
							return t.UTC(), nil
						}
					} else {
						return t.UTC(), nil
					}
				} else {
					return t.UTC(), nil
				}
			} else {
				return t.UTC(), nil
			}
		} else {
			return t.UTC(), nil
		}

	case time.Time:
		return v.UTC(), nil
	case int:
		return time.Unix(int64(v), 0).UTC(), nil
	case int8:
		return time.Unix(int64(v), 0).UTC(), nil
	case int16:
		return time.Unix(int64(v), 0).UTC(), nil
	case int32:
		return time.Unix(int64(v), 0).UTC(), nil
	case int64:
		return time.Unix(v, 0).UTC(), nil
	default:
		return time.Time{}, fmt.Errorf("on conversion to time.Time got unexpected value %#v", i)
	}
}

func ToTimeInLocationIgnoreErrors(i interface{}, loc *time.Location) time.Time {
	v, err := ToTimeInLocation(i, loc)
	if err != nil {
		log.Info(err)
	}

	return v
}

func ToTimeIgnoreErrors(i interface{}) time.Time {
	v, err := ToTime(i)
	if err != nil {
		log.Info(err)
	}

	return v
}

// ToInt converts a unknown value to int
func ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	}
	result, err := ToFloat64(value)
	if err != nil {
		return 0, fmt.Errorf("on conversion to int got unexpected value %#v", value)
	}

	return int(result), nil
}

func ToIntIgnoreErrors(value interface{}) int {
	v, _ := ToInt(value)

	return v
}

func ToSliceByteIgnoreErrors(i interface{}) []byte {
	v, _ := ToSliceByte(i)

	return v
}

func ToSliceByte(i interface{}) ([]byte, error) {
	a := make([]byte, 0)

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []byte:
		return v, nil
	case []string:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []int:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
	case []int8:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []int16:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []int32:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []int64:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []float32:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []float64:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []uint:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []uint16:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []uint32:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []uint64:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []bool:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case nil:
		return a, nil
	case []json.Number:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []uintptr:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []time.Duration:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []time.Month:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case []time.Weekday:
		for _, u := range v {
			i, _ := ToByte(u)
			a = append(a, i)
		}
		return a, nil
	case interface{}:
		i, _ := ToByte(v)
		return []byte{i}, nil
	default:
		return a, fmt.Errorf("on conversion to []byte got unexpected value %#v, %T", i, i)
	}

	return make([]byte, 0), nil
}

func ToUint32IgnoreErrors(i interface{}) uint32 {
	v, _ := ToUint32(i)

	return v
}

func ToUint32(i interface{}) (uint32, error) {
	v, err := ToInt64(i)

	return uint32(v), err
}

func ToUint64IgnoreErrors(i interface{}) uint64 {
	v, _ := ToUint64(i)

	return v
}

func ToByteIgnoreErrors(value interface{}) byte {
	v, _ := ToByte(value)

	return v
}

func ToByte(i interface{}) (byte, error) {
	result, err := ToFloat64(i)
	if err != nil {
		return 0, fmt.Errorf("on conversion to byte got unexpected value %#v, %T", i, i)
	}

	return byte(result), nil
}

// ToFloat64 converts a unknown value to float64
func ToFloat64(i interface{}) (float64, error) {
	switch s := i.(type) {
	case float64:
		return s, nil
	case float32:
		return float64(s), nil
	case int64:
		return float64(s), nil
	case int32:
		return float64(s), nil
	case int16:
		return float64(s), nil
	case int8:
		return float64(s), nil
	case int:
		return float64(s), nil
	case uint:
		return float64(s), nil
	case uint8:
		return float64(s), nil
	case uint16:
		return float64(s), nil
	case uint32:
		return float64(s), nil
	case uint64:
		return float64(s), nil
	case uintptr:
		return float64(s), nil
	case time.Duration:
		return float64(s), nil
	case time.Month:
		return float64(s), nil
	case time.Weekday:
		return float64(s), nil
	case nil:
		return float64(0), nil
	case time.Time:
		return float64(s.Unix()), nil
	case bool:
		if s {
			return float64(1), nil
		}
		return float64(0), nil
	case json.Number:
		if result, err := strconv.ParseFloat(string(s), 64); err != nil {
			return float64(0), fmt.Errorf("on conversion to float64 got unexpected value %#v", i)

		} else {
			return result, nil
		}
	case string:
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			return v, nil
		}
		return float64(0), fmt.Errorf("on conversion to float64 got unexpected value %#v", i)
	default:
		return float64(0), fmt.Errorf("on conversion to float64 got unexpected value %#v, %T", i, i)
	}
}

func ToFloat64IgnoreErrors(i interface{}) float64 {
	v, _ := ToFloat64(i)

	return v
}

// ToFloat32 converts a unknown value to float32
func ToFloat32(i interface{}) (float32, error) {
	result, err := ToFloat64(i)
	if err != nil {
		return 0, fmt.Errorf("on conversion to float32 got unexpected value %#v", i)
	}

	return float32(result), nil
}

func ToFloat32IgnoreErrors(i interface{}) float32 {
	v, _ := ToFloat32(i)

	return v
}

// ToInt64 converts a unknown value to int64
func ToInt64(value interface{}) (int64, error) {
	result, err := ToFloat64(value)
	if err != nil {
		return 0, fmt.Errorf("on conversion to int64 got unexpected value %#v, %T", value, value)
	}

	return int64(result), nil
}

// ToUint64 converts a unknown value to uint64
func ToUint64(value interface{}) (uint64, error) {
	result, err := ToFloat64(value)
	if err != nil {
		return 0, fmt.Errorf("on conversion to uint64 got unexpected value %#v, %T", value, value)
	}

	return uint64(result), nil
}

func ToInt64IgnoreErrors(value interface{}) int64 {
	v, _ := ToInt64(value)

	return v
}

// ToBool casts an empty interface to a bool.
func ToBool(i interface{}) (bool, error) {
	switch b := i.(type) {
	case bool:
		return b, nil
	case nil:
		return false, nil
	case int, int8, int16, int32, int64:
		if b != 0 {
			return true, nil
		}
		return false, nil
	case float64:
		if b != float64(0) {
			return true, nil
		}
		return false, nil
	case float32:
		if b != float32(0) {
			return true, nil
		}

		return false, nil
	case string:
		return strconv.ParseBool(i.(string))
	default:
		return false, fmt.Errorf("on conversion to bool got unexpected value %#v, %T", i, i)
	}
}

func ToBoolIgnoreErrors(i interface{}) bool {
	v, _ := ToBool(i)

	return v
}

func ToError(value interface{}) error {
	if value == nil {
		return nil
	}
	switch t := value.(type) {
	case error:
		return t
	}

	return fmt.Errorf(ToStr(value)+"%#v %T", value, value)
}

func ToTimeDuration(value interface{}) (time.Duration, error) {
	v, err := ToInt64(value)

	return time.Duration(v), err
}

func ToTimeDurationIgnoreErrors(value interface{}) time.Duration {
	v, _ := ToTimeDuration(value)

	return v
}

func ToString(i interface{}) string {
	return ToStr(i)
}

func ToEscapedString(i interface{}) string {
	s := ToStr(i)
	q := ""
	for _, ch := range s {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			q = q + string(ch)
		} else {
			q = q + `\` + string(ch)
		}
	}
	return q
}

func ToAlphaNumericString(i interface{}) string {
	s := ToStr(i)
	q := ""
	for _, ch := range s {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			q = q + string(ch)
		}
	}
	return q
}

func ToAlphaNumericWithExceptionsString(i interface{}, exceptions []string) string {
	s := ToStr(i)
	q := ""
	for _, ch := range s {
		found := false
		for _, ex := range exceptions {
			if string(ch) == ex {
				found = true
				break
			}
		}
		if found {
			q = q + string(ch)
			continue
		}
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			q = q + string(ch)
		}
	}
	return q
}

func IsStr(value interface{}) bool {
	_, isString := value.(string)

	return isString
}

func IsNumerical(value interface{}) bool {
	var digitCheck = regexp.MustCompile(`^[0-9]+$`)

	return digitCheck.MatchString(ToString(value))
}

// ToStr converts unknown value to a string
func ToStr(value interface{}) string {
	switch s := value.(type) {
	case string:
		return s
	case bool:
		return strconv.FormatBool(s)
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64)
	case int64:
		return strconv.FormatInt(s, 10)
	case int:
		return strconv.FormatInt(int64(s), 10)
	case []byte:
		return string(s)
	case []string:
		return strings.Join(s, ",")
	case []interface{}:
		return strings.Join(ToSliceStrIgnoreErrors(s), ",")
	case json.Number:
		return string(s)
	case nil:
		return ""
	case fmt.Stringer:
		return s.String()
	case error:
		return s.Error()
	default:
		return fmt.Sprint(value)
	}
}

// ToMapStrStr converts unknown value to map string
func ToMapStrStr(i interface{}) (map[string]string, error) {
	result := make(map[string]string, 0)
	s, err := ToMapStrInterface(i)
	if err != nil {
		return result, fmt.Errorf("on conversion to map[string]string got unexpected value %#v error: %s", i, err.Error())
	}

	for k, val := range s {
		result[k] = ToStr(val)
	}

	return result, nil
}

func ToMapStrStrIgnoreErrors(i interface{}) map[string]string {
	v, _ := ToMapStrStr(i)

	return v
}

// ToMapStrFloat32 converts a unknown value a map of float32
func ToMapStrFloat32(i interface{}) (map[string]float32, error) {
	result := make(map[string]float32, 0)
	s, err := ToMapStrInterface(i)
	if err != nil {
		return result, fmt.Errorf("on conversion to map[string]float32 got unexpected value %#v error: %s", i, err.Error())
	}

	for k, val := range s {
		f, _ := ToFloat32(val)
		result[k] = f
	}

	return result, nil
}

// ToMapStrInt64 converts a unknown value a map of int64
func ToMapStrInt64(i interface{}) (map[string]int64, error) {
	result := make(map[string]int64, 0)
	s, err := ToMapStrInterface(i)
	if err != nil {
		return result, fmt.Errorf("on conversion to map[string]int64 got unexpected value %#v error: %s", i, err.Error())
	}

	for k, val := range s {
		f, _ := ToInt64(val)
		result[k] = f
	}

	return result, nil
}

// ToStringMapInt converts a unknown value a map of int
func ToMapStrInt(i interface{}) (map[string]int, error) {
	result := make(map[string]int, 0)
	s, err := ToMapStrInterface(i)
	if err != nil {
		return result, fmt.Errorf("on conversion to map[string]int64 got unexpected value %#v error: %s", i, err.Error())
	}

	for k, val := range s {
		result[k] = ToIntIgnoreErrors(val)
	}

	return result, nil
}

func ToMapStrFloat64IgnoreErrors(i interface{}) map[string]float64 {
	v, _ := ToMapStrFloat64(i)
	return v
}

// ToMapStrFloat64 to map string float64
func ToMapStrFloat64(i interface{}) (map[string]float64, error) {
	result := make(map[string]float64, 0)
	s, err := ToMapStrInterface(i)
	if err != nil {
		return result, fmt.Errorf("on conversion to map[string]float64 got unexpected value %#v error: %s", i, err.Error())
	}

	for k, val := range s {
		f, _ := ToFloat64(val)
		result[k] = f
	}

	return result, nil
}

// ToMapStrBool to map string bool
func ToMapStrBool(i interface{}) (map[string]bool, error) {
	result := make(map[string]bool, 0)
	s, err := ToMapStrInterface(i)
	if err != nil {
		return result, fmt.Errorf("on conversion to map[string]bool got unexpected value %#v error: %s", i, err.Error())
	}

	for k, val := range s {
		b, _ := ToBool(val)
		result[k] = b
	}

	return result, nil
}

func ToMapStrBoolIgnoreErrors(i interface{}) map[string]bool {
	result, _ := ToMapStrBool(i)

	return result
}

// ToSliceStr casts an empty interface to a []string.
func ToSliceStr(i interface{}) ([]string, error) {
	a := make([]string, 0)

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			s := ToStr(u)
			a = append(a, s)
		}
		return a, nil
	case []string:
		return v, nil
	case string:
		return strings.Fields(v), nil
	case interface{}:
		str := ToStr(v)
		return []string{str}, nil
	case nil:
		return make([]string, 0), nil
	default:
		return a, fmt.Errorf("on conversion to []string got unexpected value %#v, %T", i, i)
	}
}

// ToSliceStrIgnoreErrors casts an empty interface to []string ignoring conversion errors
func ToSliceStrIgnoreErrors(i interface{}) []string {
	res, _ := ToSliceStr(i)

	return res
}

// ToSliceIntIgnoreErrors casts an empty interface to []int ignoring conversion errors
func ToSliceIntIgnoreErrors(i interface{}) []int {
	res, _ := ToSliceInt(i)

	return res
}

// ToSliceInt64 casts an empty interface to a []int64.
func ToSliceInt64(i interface{}) ([]int64, error) {
	a := make([]int64, 0)
	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []string:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []int:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []int8:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []int16:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []int32:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []int64:
		return v, nil
	case []float32:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []float64:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []uint:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []uint8:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []uint16:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []uint32:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []uint64:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []bool:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case nil:
		return a, nil
	case []json.Number:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []uintptr:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []time.Duration:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []time.Month:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case []time.Weekday:
		for _, u := range v {
			i, _ := ToInt64(u)
			a = append(a, i)
		}
		return a, nil
	case interface{}:
		i, _ := ToInt64(v)
		return []int64{i}, nil
	default:
		return a, fmt.Errorf("on conversion to []int64 got unexpected value %#v, %T", i, i)
	}
}

// ToSliceInt casts an empty interface to a []int.
func ToSliceInt(i interface{}) ([]int, error) {
	a := make([]int, 0)

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []string:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []int:
		return v, nil
	case []int8:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []int16:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []int32:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []int64:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []float32:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []float64:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []uint:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []uint8:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []uint16:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []uint32:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []uint64:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []bool:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case nil:
		return a, nil
	case []json.Number:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []uintptr:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []time.Duration:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []time.Month:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case []time.Weekday:
		for _, u := range v {
			i, _ := ToInt(u)
			a = append(a, i)
		}
		return a, nil
	case interface{}:
		i, _ := ToInt(v)
		return []int{i}, nil
	default:
		return a, fmt.Errorf("on conversion to []int got unexpected value %#v, %T", i, i)
	}
}

// ToMapStrSliceStr casts an empty interface to a map[string][]string.
func ToMapStrSliceStr(i interface{}) (map[string][]string, error) {
	var m = map[string][]string{}
	switch v := i.(type) {
	case map[string][]string:
		return v, nil
	case map[string][]interface{}:
		for k, val := range v {
			key := ToStr(k)
			m[key], _ = ToSliceStr(val)
		}
		return m, nil
	case map[string]string:
		for k, val := range v {
			key := ToStr(k)
			m[key] = []string{val}
		}
	case map[string]interface{}:
		for key, val := range v {
			switch vt := val.(type) {
			case []interface{}:
				m[key], _ = ToSliceStr(vt)
			case []string:
				m[key] = vt
			default:
				s := ToStr(val)
				m[key] = []string{s}
			}
		}
		return m, nil
	case map[interface{}][]string:
		for k, val := range v {
			key := ToStr(k)
			m[key], _ = ToSliceStr(val)
		}
		return m, nil
	case map[interface{}]string:
		for k, val := range v {
			key := ToStr(k)
			m[key], _ = ToSliceStr(val)
		}
		return m, nil
	case map[interface{}][]interface{}:
		for k, val := range v {
			key := ToStr(k)
			m[key], _ = ToSliceStr(val)
		}
		return m, nil
	case map[interface{}]interface{}:
		for k, val := range v {
			key := ToStr(k)

			value, err := ToSliceStr(val)
			if err != nil {
				return m, fmt.Errorf("on conversion to map[string][]string got unexpected value %#v", i)
			}
			m[key] = value
		}
	default:
		return m, fmt.Errorf("on conversion to map[string][]string got unexpected value %#v, %T", i, i)
	}
	return m, nil
}

func ToMapStrSliceStrIgnoreErrors(i interface{}) map[string][]string {
	v, _ := ToMapStrSliceStr(i)

	return v
}

// ToMapStrInterface casts an interface to a map[string]interface{}
func ToMapStrInterface(i interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{}, 0)
	switch v := i.(type) {
	case map[string][]string:
		for k, val := range v {
			m[k] = val
		}
		return m, nil
	case map[string]map[string][]string:
		for k, val := range v {
			m[k] = val
		}
		return m, nil
	case map[string]map[string]bool:
		for k, val := range v {
			m[k] = val
		}
		return m, nil
	case map[string]map[string]map[string]bool:
		for k, val := range v {
			m[k] = val
		}
		return m, nil
	case nil:
		return m, nil
	case map[string]map[string]interface{}:
		for k, val := range v {
			m[k] = val
		}
		return m, nil
	case map[string]map[string]float64:
		for k, val := range v {
			m[k] = val
		}
		return m, nil
	case map[string][]map[string]interface{}:
		for k, val := range v {
			m[k] = val
		}
		return m, nil
	case map[string]interface{}:
		if v == nil {
			return m, nil
		}
		return v, nil
	case map[interface{}]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]int16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]int8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[string]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[int]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[int8]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int8]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int16]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[int32]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int32]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[uint]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[uint8]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint8]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[uint16]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint16]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[uint32]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint32]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[uint64]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uint64]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[uintptr]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[uintptr]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[int64]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[int64]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[float64]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float64]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[float32]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[float32]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[json.Number]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[json.Number]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[time.Time]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[time.Time]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil

	case map[bool]int32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]int:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]int64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]float32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]float64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]json.Number:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]string:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]interface{}:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]uint:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]uint8:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]uint16:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]uint32:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]uint64:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]uintptr:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]time.Time:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]time.Duration:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]time.Month:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]time.Weekday:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	case map[bool]bool:
		for k, val := range v {
			m[ToStr(k)] = val
		}
		return m, nil
	default:
		return m, fmt.Errorf("on conversion to map[string]interface{} got unexpected value %#v, %T", i, i)
	}
}

func ToMapStrInterfaceIgnoreErrors(i interface{}) map[string]interface{} {
	v, err := ToMapStrInterface(i)
	if err != nil {
		log.Info(err)
	}
	if v == nil {
		log.Info(v, err, i)
		return make(map[string]interface{}, 0)
	}

	return v
}

func ToMapStrMapStrInterface(i interface{}) (map[string]map[string]interface{}, error) {
	s, err := ToMapStrInterface(i)
	if err != nil {
		return make(map[string]map[string]interface{}, 0), err
	}
	r := make(map[string]map[string]interface{}, 0)
	for k, v := range s {
		val, err := ToMapStrInterface(v)
		if err != nil {
			return make(map[string]map[string]interface{}, 0), err
		}
		r[k] = val
	}

	return r, nil
}

func ToMapStrMapStrInterfaceIgnoreErrors(i interface{}) map[string]map[string]interface{} {
	s, _ := ToMapStrMapStrInterface(i)

	return s
}
