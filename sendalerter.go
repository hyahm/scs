package scs

// 发送信息的接口
type SendAlerter interface {
	Send(body *Message, to ...string) error
}
