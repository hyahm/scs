package pkg

import (
	"math/rand"
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
