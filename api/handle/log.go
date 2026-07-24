package handle

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/xmux"
)

// 日志心跳时间检测
var HEARTBEAT = time.Second * 10

// Log 通过 WebSocket 实时推送脚本日志（tailf）。
// 路由 /log/{name}/{int:line}：name 为副本名（subname），line 为显示最后 N 行（0 表示全部）。
func Log(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	line := r.FormValue("line")
	num, _ := strconv.Atoi(line)
	ws, err := xmux.UpgradeWebSocket(w, r)
	if err != nil {
		golog.Error(err)
		return
	}
	logfile := filepath.Join(internal.GetLogPath(), name+".log")
	f, err := os.Open(logfile)
	if err != nil {
		ws.SendMessage([]byte("file not found: "+logfile), xmux.TypeMsg)
		return
	}
	// 先统计总行数，用于跳过前面部分、只发送最后 num 行
	count := 0
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		if scan.Err() != nil {
			break
		}
		count++
	}
	f.Close()
	f, err = os.Open(logfile)
	if err != nil {
		golog.Error(err)
		return
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	latest := time.Now()
	for {
		if time.Since(latest) > HEARTBEAT {
			if err := ws.Ping([]byte("ping")); err != nil {
				break
			}
			latest = time.Now()
		}
		data, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				time.Sleep(time.Second)
				continue
			}
			golog.Error(err)
			break
		}
		count--
		// num>0 时只发送最后 num 行；num<=0 时发送全部
		if num > 0 && count >= num {
			continue
		}
		if err := ws.SendMessage(data, xmux.TypeMsg); err != nil {
			golog.Error(err)
			return
		}
	}
	golog.Info("show log exit ", name)
}
