package util

import "strconv"

// BoolToNumberString 将 bool 转换为数字字符串
func BoolToNumberString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func IntToString[T int | int8 | int16 | int32 | int64](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

func StringToInt[T int | int8 | int16 | int32 | int64](s string) T {
	i, _ := strconv.ParseInt(s, 10, 64)
	return T(i)
}
