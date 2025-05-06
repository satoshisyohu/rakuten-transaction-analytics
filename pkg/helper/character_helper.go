package helper

import "strconv"

// StringToInt64 stringをint64に変換する
func StringToInt64(s string) int64 {
	i, _ := strconv.Atoi(s)
	return int64(i)
}
