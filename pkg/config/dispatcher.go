package config

import (
	"crypto/md5"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/hyahm/scs/pkg/message"
)

func Md5(src string) string {
	m := md5.New()
	io.WriteString(m, src)
	return fmt.Sprintf("%x", m.Sum(nil))
}

type alertInfoMap struct {
	sync.RWMutex
	dispatcher map[string]AlertInfo
}

// 分发器, 每个报警器间隔多久发一次
// var dispatcherLock sync.RWMutex

// // 分发器  2个map的key 分别是  name key,  副本名和一个自定义的key
var dispatcher = alertInfoMap{
	dispatcher: make(map[string]AlertInfo),
}

// func init() {
// 	dispatcher = make(map[string]AlertInfo)
// 	dispatcherLock = sync.RWMutex{}
// }

func GetDispatcherList() map[string]AlertInfo {
	dispatcher.RLock()
	defer dispatcher.RUnlock()
	return dispatcher.dispatcher
}

func GetDispatcher(name string) (AlertInfo, bool) {
	dispatcher.RLock()
	defer dispatcher.RUnlock()
	val, ok := dispatcher.dispatcher[name]
	return val, ok
}

func SetDispatcher(name string, value AlertInfo) {
	dispatcher.RLock()
	defer dispatcher.RUnlock()
	dispatcher.dispatcher[name] = value
}

func DeleteDispatcher(name string) {
	dispatcher.Lock()
	defer dispatcher.Unlock()
	delete(dispatcher.dispatcher, name)
}

type RespAlert struct {
	Title              string   `json:"title"`
	Pname              string   `json:"pname"`
	Name               string   `json:"name"`
	Reason             string   `json:"reason"`
	ContinuityInterval int      `json:"continuityInterval"`
	To                 *AlertTo `json:"to"`
}

func (ra *RespAlert) SendAlert() {
	dispatcher.Lock()
	defer dispatcher.Unlock()
	// 异常的通知
	// 如果收到了报警
	val, ok := GetDispatcher(ra.Name)
	if !ok {
		// 如果是第一次， 那么初始化值并直接发送报警
		if ra.ContinuityInterval == 0 {
			ra.ContinuityInterval = 60 * 60
		}
		ai := AlertInfo{
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
		SetDispatcher(ra.Name, ai)
		AlertMessage(ai.AM, ai.To)
		return
	}
	// 否则的话， 检查上次报警的时间是否大于间隔时间
	if time.Since(val.AlertTime) > val.ContinuityInterval {
		AlertMessage(val.AM, val.To)
		val.AlertTime = time.Now()
		SetDispatcher(ra.Name, val)
	}

}

func CleanAlert() {
	// 删除超过10小时没发送信息的值， 每10分钟执行一次
	for {
		dispatcher.Lock()
		for name, di := range GetDispatcherList() {
			if time.Since(di.AlertTime) > time.Hour*10 {
				DeleteDispatcher(name)
			}

		}
		dispatcher.Unlock()
		time.Sleep(time.Minute * 10)
	}

}
