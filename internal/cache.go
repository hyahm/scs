package internal

import "github.com/hyahm/scs/pkg/message"

var msgCache chan message.Message

// 为了避免信息错乱 将有一个1000缓冲区来存放信息
func init() {
	msgCache = make(chan message.Message, 1000)
}
