package utils

import (
	"fmt"
	"nuvem/engine/logger"
	"strconv"
)

func GetInterfaceUint32(str string, qurey map[string]interface{}) uint32 {
	if qurey == nil {
		return 0
	}
	value, err := qurey[str]
	if !err {
		return 0
	}
	return uint32(InterfaceToInt64(value))
}

func ForceUint32(str string, qurey map[string]interface{}) uint32 {
	if qurey == nil {
		return 0
	}
	value, err := qurey[str]
	if !err {
		return 0
	}
	return uint32(InterfaceToInt64(value))
}

func ForceInt(str string, qurey map[string]interface{}) int {
	if qurey == nil {
		return 0
	}
	value, err := qurey[str]
	if !err {
		return 0
	}
	return int(InterfaceToInt64(value))
}

func GetInterfaceString(str string, qurey map[string]interface{}) string {
	if qurey == nil {
		return "0"
	}
	value, err := qurey[str]
	if !err {
		return "0"
	}

	return InterfaceToString(value)
}

func GetString(str string, qurey map[string]string) string {
	if qurey == nil {
		return "0"
	}
	value, err := qurey[str]
	if !err {
		return "0"
	}

	return value
}

func ToUint32(str string) uint32 {
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return uint32(value)
}

func ToInt(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return value
}

func ToString(value uint32) string {
	return strconv.FormatUint(uint64(value), 10)
}

func InterfaceToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return value.(string)
	case int:
		return strconv.Itoa(value.(int))
	case int64:
		return strconv.FormatInt(value.(int64), 10)
	case uint32:
		return fmt.Sprintf("%d", value.(uint32))
	case float64:
		return strconv.FormatFloat(value.(float64), 'E', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'E', -1, 32)
	default:
		logger.Error(fmt.Sprintf("unknow type is %T", v))
	}
	return ""
}

func InterfaceToInt64(value interface{}) int64 {
	switch v := value.(type) {
	case string:
		val, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return 0
		}
		return val
	case int:
		return int64(value.(int))
	case int64:
		return int64(value.(int64))
	case uint32:
		return int64(value.(uint32))
	case float64:
		return int64(value.(float64))
	case float32:
		return int64(value.(float32))
	default:
		logger.Error(fmt.Sprintf("unknow type is %T", v))
	}
	return 0
}

func CharCodeAt(s string, n int) rune {
	i := 0
	for _, r := range s {
		if i == n {
			return r
		}
		i++
	}
	return 0
}

func FromCharCode(r []rune, n int) string {
	i := 0
	for _, s := range r {
		if i == n {
			return string(s)
		}
		i++
	}
	return ""
}

//return position
func IndexOf(slicearray []int, val int) int {
	ret := -1
	for i, v := range slicearray {
		if v == val {
			return i
		}
	}
	return ret
}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func DumpSocketData(data interface{}) {
	switch v := data.(type) {
	case string:
		logger.Debug("is string", v)
	case map[string]interface{}:
		logger.Debug("is json", v)
	case []interface{}:
		logger.Debug("is array", v)
	default:
		logger.Debug(fmt.Sprintf("type is %T", v))
	}
}
