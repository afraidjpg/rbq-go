package util

import "strconv"

func IntToString[T int | int8 | int16 | int32 | int64](i T) string {
	return strconv.FormatInt(int64(i), 10)
}
