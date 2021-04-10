package alert

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/hyahm/scs/server/at"
)

var alerter *Alerter // 保存报警器配置文件

type Alert struct {
	Email    *AlertEmail    `yaml:"email"`
	Rocket   *AlertRocket   `yaml:"rocket"`
	Telegram *AlertTelegram `yaml:"telegram"`
	WeiXin   *AlertWeiXin   `yaml:"weixin"`
}

type Alerter struct {
	Alert        *Alert
	Alerts       map[string]SendAlerter
	alertsLocker *sync.RWMutex
	Cancel       context.CancelFunc
	Ctx          context.Context
}

func Run(a *Alert) {
	if alerter == nil {
		alerter = &Alerter{
			Alert:        a,
			Alerts:       make(map[string]SendAlerter),
			alertsLocker: &sync.RWMutex{},
		}
	} else {
		// 对比传过来的alert， 是否存在值得变化，
		// 如果都是一样的，不错任何处理， 否则删除之前的goroutine， 重新启动
	}
	// 运行报警器
	// 启动goroutine
	alerter.Ctx, alerter.Cancel = context.WithCancel(context.Background())
	InitAlert()
}

func InitAlert() {
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
}

func AlertMessage(msg *Message, to *at.AlertTo) {
	for _, alert := range alerter.Alerts {
		al := alert
		msg.HostName, _ = os.Hostname()
		if to == nil {
			alertErr := al.Send(msg)
			if alertErr != nil {
				fmt.Println(alertErr)
			}
			continue
		}
		switch al.(type) {
		// 目前只支持邮箱
		case *AlertEmail:
			go func() {
				alertErr := al.Send(msg, to.Email...)
				if alertErr != nil {
					fmt.Println(alertErr)
				}

			}()
		case *AlertRocket:
			go func() {
				alertErr := al.Send(msg, to.Rocket...)
				if alertErr != nil {
					fmt.Println(alertErr)
				}

			}()
		case *AlertTelegram:
			go func() {
				alertErr := al.Send(msg, to.Telegram...)
				if alertErr != nil {
					fmt.Println(alertErr)
				}

			}()
		case *AlertWeiXin:
			go func() {
				alertErr := al.Send(msg, to.WeiXin...)
				if alertErr != nil {
					fmt.Println(alertErr)
				}

			}()
		}

	}
}
