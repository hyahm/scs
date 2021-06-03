package alert

import (
	"fmt"
	"os"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/message"
	"github.com/hyahm/scs/to"
)

type AlertInfo struct {
	AlertTime          time.Time
	Interval           int // 上次报警的时间
	Broken             bool
	Start              time.Time // 报警时间
	BrokenTime         time.Time
	AM                 *message.Message
	To                 *to.AlertTo
	ContinuityInterval time.Duration
}

var cache []*AlertInfo

func init() {
	cache = make([]*AlertInfo, 4)
}

func (ai *AlertInfo) BreakDown(title string) {
	ai.AM.Title = title
	if !ai.Broken {
		// 第一次发送
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
		// ai.AM.BrokenTime = ai.Start.String()
		ai.AM.FixTime = time.Now().Local().String()
		AlertMessage(ai.AM, nil)
		ai.Broken = false
	}
}

func AlertMessage(msg *message.Message, at *to.AlertTo) {
	for _, alert := range alerter.Alerts {
		al := alert
		msg.HostName, _ = os.Hostname()
		if at == nil {
			alertErr := al.Send(msg)
			if alertErr != nil {
				golog.Error(alertErr)
				continue
			}
			golog.Info("send seccess")
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
