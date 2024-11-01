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
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
	"github.com/hyahm/xmux"
)

// 日志心跳时间检测
var HEARTBEAT = time.Second * 10

func Log(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	line := xmux.Var(r)["line"]
	num, _ := strconv.Atoi(line)

	ws, err := xmux.UpgradeWebSocket(w, r)
	if err != nil {
		golog.Error(err)
		return
	}
	logfile := filepath.Join(global.LogDir, subname.Subname(name).String()+".log")
	f, err := os.Open(logfile)
	if err != nil {
		ws.SendMessage([]byte("file not found, yes, just without any print on, please wait"), xmux.TypeMsg)
		ws.Close()
		return
	}
	count := 0
	scan := bufio.NewScanner(f)
	for scan.Scan() {
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
	// ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	// 设置超时时间
	latest := time.Now()
	for {
		if time.Since(latest) > HEARTBEAT {
			err := ws.Ping([]byte("ping"))
			if err != nil {
				break
			}
			latest = time.Now()
		}
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				time.Sleep(time.Second)
				continue
			} else {
				golog.Error(err)
				break
			}

		}
		count--
		if count >= num {
			continue
		}
		err = ws.SendMessage(line, xmux.TypeMsg)
		if err != nil {
			golog.Error(err)
			return
		}
	}
	golog.Info("show log exit ", name)
}
