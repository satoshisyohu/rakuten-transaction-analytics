package helper

import "time"

func StringToDate(layout, value string) time.Time {
	dateValue, _ := time.Parse(layout, value)
	return dateValue
}
