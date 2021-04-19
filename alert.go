package scs

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func (s *Server) successAlert() {
	// 启动成功后恢复的通知
	if !s.AI.Broken {
		return
	}
	for {
		select {
		// 每3秒一次操作
		case <-time.After(time.Second * 3):
			am := &Message{
				Title:      "service recover",
				Pname:      s.Script.Name,
				Name:       s.SubName,
				BrokenTime: s.AI.Start.String(),
				FixTime:    time.Now().String(),
			}
			AlertMessage(am, s.AT)
			s.AI.Broken = false
			return
		case <-s.Ctx.Done():
			return
		}
	}

}

type AlertTo struct {
	Email    []string `yaml:"email"`
	Rocket   []string `yaml:"rocket"`
	Telegram []string `yaml:"telegram"`
	WeiXin   []string `yaml:"weixin"`
	Callback []string `yaml:"callback"`
}

// 比较新的与之前的是否相等， 调用者必须是新的
func (at *AlertTo) IsEqual(pname string) bool {
	if at == nil && ss.Scripts[pname].AT == nil {
		return true
	}
	if (at == nil && ss.Scripts[pname].AT != nil) || (at != nil && ss.Scripts[pname].AT == nil) {
		return false
	}
	if !EqualStringArray(at.Email, ss.Scripts[pname].AT.Email) ||
		!EqualStringArray(at.Rocket, ss.Scripts[pname].AT.Rocket) ||
		!EqualStringArray(at.Telegram, ss.Scripts[pname].AT.Telegram) ||
		!EqualStringArray(at.WeiXin, ss.Scripts[pname].AT.WeiXin) {
		return false
	}
	return true

}

// 暂时只支持邮箱
type Alert struct {
	Email    *AlertEmail    `yaml:"email,omitempty"`
	Rocket   *AlertRocket   `yaml:"rocket,omitempty"`
	Telegram *AlertTelegram `yaml:"telegram,omitempty"`
	WeiXin   *AlertWeiXin   `yaml:"weixin,omitempty"`
	Callback *Callback      `yaml:"weixin,omitempty"`
}
type Alerter struct {
	Alert        *Alert
	Alerts       map[string]SendAlerter
	alertsLocker *sync.RWMutex
}

var alerter *Alerter // 保存报警器配置文件
func RunAlert(a *Alert) {

	if alerter == nil {
		alerter = &Alerter{
			Alert:        a,
			Alerts:       make(map[string]SendAlerter),
			alertsLocker: &sync.RWMutex{},
		}
	}
	alerter.Alert = a
	// 运行报警器
	// 启动goroutine
	InitAlert()
}

// var Alerts map[string]SendAlerter

// func init() {
// 	Alerts = make(map[string]SendAlerter)
// }

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
	if alerter.Alert.Callback != nil {
		alerter.Alerts["callback"] = alerter.Alert.Callback
	}
}

// func (alert *Alert) InitAlert() {
// 	// 报警配置转移到了 Alerts
// 	if alert.Email != nil {
// 		if alert.Email.Host != "" && alert.Email.UserName != "" &&
// 			alert.Email.Password != "" {
// 			if alert.Email.Port == 0 {
// 				alert.Email.Port = 465
// 			}
// 			Alerts["email"] = alert.Email

// 		}
// 	}
// 	if alert.Rocket != nil {
// 		if alert.Rocket.Server != "" && alert.Rocket.Username != "" &&
// 			alert.Rocket.Password != "" {
// 			Alerts["rocket"] = alert.Rocket
// 		}

// 	}

// 	if alert.Telegram != nil {
// 		if alert.Telegram.Server != "" && alert.Telegram.Username != "" &&
// 			alert.Telegram.Password != "" {
// 			Alerts["telegram"] = alert.Telegram
// 		}

// 	}
// 	if alert.WeiXin != nil {
// 		Alerts["weixin"] = alert.WeiXin
// 	}
// }

func AlertMessage(msg *Message, at *AlertTo) {
	for _, alert := range alerter.Alerts {
		al := alert
		msg.HostName, _ = os.Hostname()
		if at == nil {
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
				alertErr := al.Send(msg, at.Email...)
				if alertErr != nil {
					fmt.Println(alertErr)
				}

			}()
		case *AlertRocket:
			go func() {
				alertErr := al.Send(msg, at.Rocket...)
				if alertErr != nil {
					fmt.Println(alertErr)
				}

			}()
		case *AlertTelegram:
			go func() {
				alertErr := al.Send(msg, at.Telegram...)
				if alertErr != nil {
					fmt.Println(alertErr)
				}

			}()
		case *AlertWeiXin:
			go func() {
				alertErr := al.Send(msg, at.WeiXin...)
				if alertErr != nil {
					fmt.Println(alertErr)
				}

			}()
		case *Callback:
			go func() {
				alertErr := al.Send(msg, at.Callback...)
				if alertErr != nil {
					fmt.Println(alertErr)
				}

			}()
		}

	}
}
