package config

import (
	"os"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg/message"
)

type AlertInfo struct {
	AlertTime          time.Time
	Interval           int // 上次报警的时间
	Broken             bool
	Start              time.Time // 报警时间
	BrokenTime         time.Time
	AM                 *message.Message
	To                 *AlertTo
	ContinuityInterval time.Duration
}

func (ai *AlertInfo) BreakDown(title string) {
	ai.AM.Title = title
	if !ai.Broken {
		ai.Broken = true
		ai.AM.BrokenTime = time.Now().String()
		ai.Start = time.Now()
		ai.AlertTime = time.Now()
		AlertMessage(ai.AM, nil)
	} else {
		if time.Since(ai.AlertTime) >= ai.ContinuityInterval {
			ai.AlertTime = time.Now()
			AlertMessage(ai.AM, nil)
		}
	}
}

func (ai *AlertInfo) Recover(title string) {
	if ai.Broken {
		ai.AM.Title = title
		ai.AM.FixTime = time.Now().Local().String()
		AlertMessage(ai.AM, nil)
		ai.Broken = false
	}
}

// AlertMessage 将报警消息发送到所有已配置的报警通道。
func AlertMessage(msg *message.Message, at *AlertTo) {
	msg.HostName, _ = os.Hostname()
	alerts := alerter.Alerts
	if len(alerts) == 0 {
		golog.Warnf("无报警通道配置，仅打印日志: %s", msg.String())
		return
	}
	for name, alert := range alerts {
		if at == nil {
			go func(n string, a message.SendAlerter) {
				if err := a.Send(msg); err != nil {
					golog.Errorf("报警发送失败 [%s]: %v", n, err)
				}
			}(name, alert)
			continue
		}
		switch a := alert.(type) {
		case *AlertEmail:
			go func(a *AlertEmail) {
				if err := a.Send(msg, at.Email...); err != nil {
					golog.Errorf("报警发送失败 [email]: %v", err)
				}
			}(a)
		case *AlertRocket:
			go func(a *AlertRocket) {
				if err := a.Send(msg, at.Rocket...); err != nil {
					golog.Errorf("报警发送失败 [rocket]: %v", err)
				}
			}(a)
		case *AlertTelegram:
			go func(a *AlertTelegram) {
				if err := a.Send(msg); err != nil {
					golog.Errorf("报警发送失败 [telegram]: %v", err)
				}
			}(a)
		case *AlertWeiXin:
			go func(a *AlertWeiXin) {
				if err := a.Send(msg); err != nil {
					golog.Errorf("报警发送失败 [weixin]: %v", err)
				}
			}(a)
		case *AlertDingDing:
			go func(a *AlertDingDing) {
				if err := a.Send(msg); err != nil {
					golog.Errorf("报警发送失败 [dingding]: %v", err)
				}
			}(a)
		case *Callback:
			go func(a *Callback) {
				if err := a.Send(msg, at.Callback...); err != nil {
					golog.Errorf("报警发送失败 [callback]: %v", err)
				}
			}(a)
		}
	}
}
