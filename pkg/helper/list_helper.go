package helper

import "strings"

func IsContain(val string, targets []string) bool {
	for _, s := range targets {
		if strings.Contains(val, s) {
			return true
		}
	}
	return false
}
