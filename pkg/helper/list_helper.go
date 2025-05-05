package helper

import (
	"slices"
	"strings"
)

// IsContain 対象の文字列がリストに含まれているか判定する(部分一致)
// 下記でコメントアウトしている記載方法とは同義であり下記のほうがよりモダンな記載方法である
func IsContain(val string, targets []string) bool {
	return slices.ContainsFunc(targets, func(s string) bool {
		return strings.Contains(val, s)
	})
}
