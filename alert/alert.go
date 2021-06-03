package alert

import (
	"sync"

	"github.com/hyahm/scs/message"
)

// 暂时只支持邮箱
type Alert struct {
	Email    *AlertEmail    `yaml:"email,omitempty"`
	Rocket   *AlertRocket   `yaml:"rocket,omitempty"`
	Telegram *AlertTelegram `yaml:"telegram,omitempty"`
	WeiXin   *AlertWeiXin   `yaml:"weixin,omitempty"`
	Callback *Callback      `yaml:"callback,omitempty"`
}
type Alerter struct {
	Alert        *Alert
	Alerts       map[string]message.SendAlerter
	alertsLocker *sync.RWMutex
}

var alerter *Alerter // 保存报警器配置文件
func RunAlert(a *Alert) {

	if alerter == nil {
		alerter = &Alerter{
			Alerts:       make(map[string]message.SendAlerter),
			alertsLocker: &sync.RWMutex{},
		}
	}
	alerter.Alert = a
	// 运行报警器
	// 启动goroutine
	InitAlert()
}

func HaveAlert() bool {
	return len(alerter.Alerts) > 0
}

// var Alerts map[string]SendAlerter

// func init() {
// 	Alerts = make(map[string]SendAlerter)
// }

func InitAlert() {
	if alerter.Alert == nil {
		return
	}
	// 报警配置转移到了 Alerts
	if alerter.Alert.Email != nil {
		if alerter.Alert.Email.Host != "" && alerter.Alert.Email.UserName != "" &&
			alerter.Alert.Email.Password != "" {
			if alerter.Alert.Email.Port == 0 {
				alerter.Alert.Email.Port = 465
			}
			alerter.Alerts["email"] = alerter.Alert.Email

		}
	}
	if alerter.Alert.Rocket != nil {

		if alerter.Alert.Rocket.Server != "" && alerter.Alert.Rocket.Username != "" &&
			alerter.Alert.Rocket.Password != "" {
			alerter.Alerts["rocket"] = alerter.Alert.Rocket
		}

	}

	if alerter.Alert.Telegram != nil {
		if alerter.Alert.Telegram.Server != "" && alerter.Alert.Telegram.Username != "" &&
			alerter.Alert.Telegram.Password != "" {
			alerter.Alerts["telegram"] = alerter.Alert.Telegram
		}

	}
	if alerter.Alert.WeiXin != nil {
		alerter.Alerts["weixin"] = alerter.Alert.WeiXin
	}
	if alerter.Alert.Callback != nil {
		alerter.Alerts["callback"] = alerter.Alert.Callback
	}
}
