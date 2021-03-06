package alert

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/hyahm/scs/message"
	"github.com/hyahm/scs/to"
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

type RespAlert struct {
	Title              string      `json:"title"`
	Pname              string      `json:"pname"`
	Name               string      `json:"name"`
	Reason             string      `json:"reason"`
	ContinuityInterval int         `json:"continuityInterval"`
	To                 *to.AlertTo `json:"to"`
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
			AM: &message.Message{
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
	// 删除超过10小时没发送信息的值
	for {
		dispatcherLock.Lock()
		for pname := range dispatcher {
			for name, di := range dispatcher[pname] {
				if time.Since(di.AlertTime) > time.Hour*10 {
					delete(dispatcher[pname], name)
				}

			}
		}
		dispatcherLock.Unlock()
		time.Sleep(time.Minute * 10)
	}

}
