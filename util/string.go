package util

import (
	"math/rand"
	"strings"
	"unicode"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Camel2Snake(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := runes[i-1]
				// 当前大写字母前面是小写字母，或前面不是下划线并后面是小写字母时，插入下划线
				if unicode.IsLower(prev) || (i+1 < len(runes) && unicode.IsLower(runes[i+1])) {
					b.WriteByte('_')
				}
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func RandStringRunes(n int) string {
	ans := make([]rune, n, n)
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ许昊龙")
	for i, _ := range ans {
		ans[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(ans)
}
