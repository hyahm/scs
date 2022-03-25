package alert

import (
	"sync"

	"github.com/hyahm/scs/pkg/message"
)

// 暂时只支持邮箱
type Alert struct {
	Email    *AlertEmail    `yaml:"email,omitempty" json:"email,omitempty"`
	Rocket   *AlertRocket   `yaml:"rocket,omitempty" json:"rocket,omitempty"`
	Telegram *AlertTelegram `yaml:"telegram,omitempty" json:"telegram,omitempty"`
	WeiXin   *AlertWeiXin   `yaml:"weixin,omitempty" json:"weixin,omitempty"`
	Callback *Callback      `yaml:"callback,omitempty" json:"callback,omitempty"`
}
type Alerter struct {
	Alert        *Alert
	Alerts       map[string]message.SendAlerter
	alertsLocker *sync.RWMutex
}

func GetAlerts() map[string]message.SendAlerter {
	alerter.alertsLocker.RLock()
	defer alerter.alertsLocker.RUnlock()
	return alerter.Alerts
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

func InitAlert() {
	if alerter.Alert == nil {
		alerter.Alerts = make(map[string]message.SendAlerter)
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
	} else {
		delete(alerter.Alerts, "email")
	}
	if alerter.Alert.Rocket != nil {

		if alerter.Alert.Rocket.Server != "" && alerter.Alert.Rocket.Username != "" &&
			alerter.Alert.Rocket.Password != "" {
			alerter.Alerts["rocket"] = alerter.Alert.Rocket
		}

	} else {
		delete(alerter.Alerts, "rocket")
	}

	if alerter.Alert.Telegram != nil && alerter.Alert.Telegram.Server != "" {
		alerter.Alerts["telegram"] = alerter.Alert.Telegram

	} else {
		delete(alerter.Alerts, "telegram")
	}
	if alerter.Alert.WeiXin != nil && alerter.Alert.WeiXin.Server != "" {
		alerter.Alerts["weixin"] = alerter.Alert.WeiXin
	} else {
		delete(alerter.Alerts, "weixin")
	}
	if alerter.Alert.Callback != nil {
		alerter.Alerts["callback"] = alerter.Alert.Callback
	} else {
		delete(alerter.Alerts, "callback")
	}
}
