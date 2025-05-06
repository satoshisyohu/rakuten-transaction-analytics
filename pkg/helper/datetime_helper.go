package helper

import "time"

// StringToDate stringをdateに変換する
func StringToDate(layout, value string) time.Time {
	dateValue, _ := time.Parse(layout, value)
	return dateValue
}
