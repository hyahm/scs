package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// 随机生成 30 - 50之间的随机数
	s := `1234567890-=qwertyuiop[]asdfghjkl;zxcvbn#m,.!@%^&*()_+QWERTYUIOP{}ASDFGHJKL:|ZXCVBNM<>?`
	out := ""
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(20)
	for i := 0; i < n+30; i++ {
		time.Sleep(time.Nanosecond * 1)
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(87)
		out += s[r : r+1]
	}
	fmt.Print(out)
}
