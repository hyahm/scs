package pkg

import (
	"math/rand"
	"regexp"
	"time"
)

func RandomToken() string {
	s := `1234567890-=qwertyuiop[]asdfghjkl;zxcvbn#m,.!@%^&*()_+QWERTYUIOP{}ASDFGHJKL:|ZXCVBNM<>?`
	out := ""
	r := rand.New(rand.NewSource(time.Now().Unix()))
	n := r.Intn(20)
	for i := 0; i < n+30; i++ {
		r := rand.Intn(len(s))
		out += s[r : r+1]
	}
	return out
}

func IsNameWithSuffix(s string) bool {
	// 正则：以 _ 开头，后面跟着至少一个数字，并且到字符串末尾
	pattern := `_\d+$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}
