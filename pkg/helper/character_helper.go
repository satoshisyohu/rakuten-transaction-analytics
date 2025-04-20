package helper

import "strconv"

func StringToInt64(s string) int64 {
	i, _ := strconv.Atoi(s)
	return int64(i)
}
