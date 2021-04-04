package alert

import (
	"os"

	"github.com/hyahm/golog"
)

// 报警相关配置
type AlertTo struct {
	Email    []string `yaml:"email"`
	Rocket   []string `yaml:"rocket"`
	Telegram []string `yaml:"telegram"`
	WeiXin   []string `yaml:"weixin"`
}

// 暂时只支持邮箱
type Alert struct {
	Email    *AlertEmail    `yaml:"email"`
	Rocket   *AlertRocket   `yaml:"rocket"`
	Telegram *AlertTelegram `yaml:"telegram"`
	WeiXin   *AlertWeiXin   `yaml:"weixin"`
}

var Alerts map[string]SendAlerter

func init() {
	Alerts = make(map[string]SendAlerter)
}

func (alert *Alert) InitAlert() {
	// 报警配置转移到了 Alerts
	if alert.Email != nil {
		if alert.Email.Host != "" && alert.Email.UserName != "" &&
			alert.Email.Password != "" {
			if alert.Email.Port == 0 {
				alert.Email.Port = 465
			}
			Alerts["email"] = alert.Email

		}
	}
	if alert.Rocket != nil {
		if alert.Rocket.Server != "" && alert.Rocket.Username != "" &&
			alert.Rocket.Password != "" {
			Alerts["rocket"] = alert.Rocket
		}

	}

	if alert.Telegram != nil {
		if alert.Telegram.Server != "" && alert.Telegram.Username != "" &&
			alert.Telegram.Password != "" {
			Alerts["telegram"] = alert.Telegram
		}

	}
	if alert.WeiXin != nil {
		Alerts["weixin"] = alert.WeiXin
	}
}

func AlertMessage(msg *Message, at *AlertTo) {
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
		case *AlertWeiXin:
			go func() {
				alertErr := al.Send(msg, at.WeiXin...)
				if alertErr != nil {
					golog.Error(alertErr)
				}

			}()
		}

	}
}
