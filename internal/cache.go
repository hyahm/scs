package internal

import (
	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/scs/pkg/message"
)

var msgCache chan message.Message

// 为了避免信息错乱 将有一个1000缓冲区来存放信息
func init() {
	msgCache = make(chan message.Message, 1000)
}

// PushAlert 将报警消息推入缓存通道，由消费者异步发送。
func PushAlert(msg message.Message) {
	select {
	case msgCache <- msg:
	default:
		golog.Warn("报警通道已满，丢弃消息: ", msg.String())
	}
}

// StartAlertConsumer 启动报警消费协程，从 msgCache 读取并发送。
func StartAlertConsumer() {
	go func() {
		for msg := range msgCache {
			config.AlertMessage(&msg, nil)
		}
	}()
}
