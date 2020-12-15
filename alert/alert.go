package alert

import (
	"os"
	"scs/internal"

	"github.com/hyahm/golog"
)

// 暂时只支持邮箱
type Alert struct {
	Email    *AlertEmail    `yaml:"email"`
	Rocket   *AlertRocket   `yaml:"rocket"`
	Telegram *AlertTelegram `yaml:"telegram"`
}

type SendAlerter interface {
	Send(body *Message, to ...string) error
}

var Alerts map[string]SendAlerter

func AlertMessage(msg *Message, at *internal.AlertTo) {
	for _, alert := range Alerts {
		al := alert
		msg.HostName, _ = os.Hostname()
		if at == nil {
			alertErr := al.Send(msg)
			if alertErr != nil {
				golog.Error(alertErr)
			}
			continue
		}
		switch al.(type) {
		// 目前只支持邮箱
		case *AlertEmail:
			go func() {
				alertErr := al.Send(msg, at.Email...)
				if alertErr != nil {
					golog.Error(alertErr)
				}

			}()
		case *AlertRocket:
			go func() {
				alertErr := al.Send(msg, at.Rocket...)
				if alertErr != nil {
					golog.Error(alertErr)
				}

			}()
		case *AlertTelegram:
			go func() {
				alertErr := al.Send(msg, at.Telegram...)
				if alertErr != nil {
					golog.Error(alertErr)
				}

			}()
		}

	}
}
