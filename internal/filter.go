package internal

import "github.com/hyahm/scs/pkg/message"

// 过滤器， 是一个map， 根据key来确定是否是重复的信息

var filter map[string]message.Message

// cpu mem disk 为 cpu内存磁盘对应的key
// 脚本报错退出的key为 name， 即副本名
// 也可以自定义key， 代码里面写
// 只有在需要报警的时候才会发送报警， 不会自动循环发送

func init() {
	filter = make(map[string]message.Message)
}
