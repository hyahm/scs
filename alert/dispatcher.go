package alert

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/hyahm/scs/internal"
)

// 分发器, 每个报警器间隔多久发一次
var dispatcherLock sync.RWMutex
var dispatcher map[string]map[string]*AlertInfo

func init() {
	dispatcher = make(map[string]map[string]*AlertInfo)
	dispatcherLock = sync.RWMutex{}
}

func GetDispatcher() []byte {
	b, err := json.Marshal(dispatcher)
	if err != nil {
		return []byte(err.Error())
	}
	return b
}

type AlertInfo struct {
	AlertTime          time.Time
	Interval           int // 上次报警的时间
	Broken             bool
	Start              time.Time // 报警时间
	BrokenTime         time.Time
	AM                 *Message
	To                 *internal.AlertTo
	ContinuityInterval time.Duration
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

type RespAlert struct {
	Title              string            `json:"title"`
	Pname              string            `json:"pname"`
	Name               string            `json:"name"`
	Reason             string            `json:"reason"`
	Broken             bool              `json:"broken"`
	ContinuityInterval int               `json:"continuityInterval"`
	To                 *internal.AlertTo `json:"to"`
}

func (ra *RespAlert) SendAlert() {
	dispatcherLock.Lock()
	defer dispatcherLock.Unlock()
	// 异常的通知
	if _, ok := dispatcher[ra.Pname]; !ok {
		dispatcher[ra.Pname] = make(map[string]*AlertInfo)
	}
	// 如果收到了报警
	if _, ok := dispatcher[ra.Pname][ra.Name]; !ok {
		// 如果是第一次， 那么初始化值并直接发送报警
		if ra.ContinuityInterval == 0 {
			ra.ContinuityInterval = 60 * 60
		}
		dispatcher[ra.Pname][ra.Name] = &AlertInfo{
			AlertTime:  time.Now(),
			Start:      time.Now(),
			BrokenTime: time.Now(),
			Broken:     true,
			AM: &Message{
				Title:      ra.Title,
				Pname:      ra.Pname,
				Name:       ra.Name,
				Reason:     ra.Reason,
				BrokenTime: time.Now().String(),
			},
			To:                 ra.To,
			ContinuityInterval: time.Duration(ra.ContinuityInterval) * time.Second,
		}
		AlertMessage(dispatcher[ra.Pname][ra.Name].AM, dispatcher[ra.Pname][ra.Name].To)

	} else {
		// 否则的话， 检查上次报警的时间是否大于间隔时间
		if time.Since(dispatcher[ra.Pname][ra.Name].AlertTime) > dispatcher[ra.Pname][ra.Name].ContinuityInterval {
			AlertMessage(dispatcher[ra.Pname][ra.Name].AM, dispatcher[ra.Pname][ra.Name].To)
			dispatcher[ra.Pname][ra.Name].AlertTime = time.Now()
		}
	}
}

func SendNetAlert() {
	// 删除超过1天没发送信息的值
	for {
		dispatcherLock.Lock()
		for pname := range dispatcher {
			for name, di := range dispatcher[pname] {
				if time.Since(di.AlertTime) > time.Hour*10 {
					delete(dispatcher[pname], name)
				}
				// if !di.Broken {
				// 	// 如果恢复了， 发完报警后删除key
				// 	di.AlertTime = time.Now()
				// 	di.AM.Title += "(已恢复)"
				// 	di.AM.FixTime = time.Now().String()
				// 	AlertMessage(di.AM, di.To)
				// 	di.AM.Reason = "问题已修复"
				// 	delete(dispatcher[pname], name)
				// } else {
				// 	// 间隔时间内才发送报警
				// 	if time.Since(di.AlertTime) >= time.Duration(di.Interval)*time.Second {
				// 		di.AlertTime = time.Now()
				// 		AlertMessage(di.AM, di.To)
				// 	}
				// }
			}
		}
		dispatcherLock.Unlock()
		time.Sleep(time.Second * 10)
	}

}
