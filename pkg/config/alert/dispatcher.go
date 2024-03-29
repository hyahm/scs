package alert

import (
	"crypto/md5"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/hyahm/scs/pkg/config/alert/to"
	"github.com/hyahm/scs/pkg/message"
)

func Md5(src string) string {
	m := md5.New()
	io.WriteString(m, src)
	return fmt.Sprintf("%x", m.Sum(nil))
}

// 分发器, 每个报警器间隔多久发一次
var dispatcherLock sync.RWMutex

// 分发器  2个map的key 分别是  name key,  副本名和一个自定义的key
var dispatcher map[string]*AlertInfo

func init() {
	dispatcher = make(map[string]*AlertInfo)
	dispatcherLock = sync.RWMutex{}
}

func GetDispatcher() map[string]*AlertInfo {
	return dispatcher
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
	// 如果收到了报警
	if _, ok := dispatcher[ra.Name]; !ok {
		// 如果是第一次， 那么初始化值并直接发送报警
		if ra.ContinuityInterval == 0 {
			ra.ContinuityInterval = 60 * 60
		}
		dispatcher[ra.Name] = &AlertInfo{
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
		AlertMessage(dispatcher[ra.Name].AM, dispatcher[ra.Name].To)

	} else {
		// 否则的话， 检查上次报警的时间是否大于间隔时间
		if time.Since(dispatcher[ra.Name].AlertTime) > dispatcher[ra.Name].ContinuityInterval {
			AlertMessage(dispatcher[ra.Name].AM, dispatcher[ra.Name].To)
			dispatcher[ra.Name].AlertTime = time.Now()
		}
	}
}

func CleanAlert() {
	// 删除超过10小时没发送信息的值， 每10分钟执行一次
	for {
		dispatcherLock.Lock()
		for name, di := range dispatcher {
			if time.Since(di.AlertTime) > time.Hour*10 {
				delete(dispatcher, name)
			}

		}
		dispatcherLock.Unlock()
		time.Sleep(time.Minute * 10)
	}

}
