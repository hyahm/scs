package alert

import (
	"encoding/json"
	"scs/internal"
	"sync"
	"time"
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
	AlertTime  time.Time
	Interval   int // 上次报警的时间
	Broken     bool
	Start      time.Time // 报警时间
	BrokenTime time.Time
	AM         *Message
	To         *internal.AlertTo
}

type RespAlert struct {
	Title    string            `json:"title"`
	Pname    string            `json:"pname"`
	Name     string            `json:"name"`
	Reason   string            `json:"reason"`
	Broken   bool              `json:"broken"`
	Interval int               `json:"interval"`
	To       *internal.AlertTo `json:"to"`
}

func (ra *RespAlert) SendAlert() {
	dispatcherLock.Lock()
	defer dispatcherLock.Unlock()
	// 异常的通知
	if _, ok := dispatcher[ra.Pname]; !ok {
		dispatcher[ra.Pname] = make(map[string]*AlertInfo)
	}
	if _, ok := dispatcher[ra.Pname][ra.Name]; !ok && ra.Broken {
		dispatcher[ra.Pname][ra.Name] = &AlertInfo{
			AlertTime:  time.Now().Add(-time.Duration(ra.Interval) * time.Second * 10),
			Broken:     ra.Broken,
			Start:      time.Now(),
			BrokenTime: time.Now(),
			AM: &Message{
				Title:      ra.Title,
				Pname:      ra.Pname,
				Name:       ra.Name,
				Reason:     ra.Reason,
				BrokenTime: time.Now().String(),
			},
			To:       ra.To,
			Interval: ra.Interval,
		}
	} else {
		// 如果存在这个报警器
		if dispatcher[ra.Pname][ra.Name].Broken == ra.Broken {
			return
		}
		dispatcher[ra.Pname][ra.Name].Broken = ra.Broken
		dispatcher[ra.Pname][ra.Name].Interval = ra.Interval
	}
}

func SendNetAlert() {
	for {
		dispatcherLock.Lock()
		for pname := range dispatcher {
			for name, di := range dispatcher[pname] {
				if !di.Broken {
					// 如果恢复了， 发完报警后删除key
					di.AlertTime = time.Now()
					di.AM.Title += "(已恢复)"
					di.AM.FixTime = time.Now().String()
					AlertMessage(di.AM, di.To)
					di.AM.Reason = "问题已修复"
					delete(dispatcher[pname], name)
				} else {
					// 间隔时间内才发送报警
					if time.Since(di.AlertTime) >= time.Duration(di.Interval)*time.Second {
						di.AlertTime = time.Now()
						AlertMessage(di.AM, di.To)
					}
				}
			}
		}
		dispatcherLock.Unlock()
		time.Sleep(time.Second * 10)
	}

}
