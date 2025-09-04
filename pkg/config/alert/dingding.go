package alert

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg/message"
)

type AlertDingDing struct {
	Server string `yaml:"server,omitempty" json:"server,omitempty"`
}

var dingdingFormat = "```\\nTitle: {{.Title }} \\nHostname: {{.HostName}} \\nAddr: {{.Addr}} \\n{{ if .Pname  }}pname:{{.Pname}} \\n{{end}}{{ if .Name }}name:{{.Name}} \\n{{end}}{{ if .DiskPath }}DiskPath:{{.DiskPath}}\\n{{end}}{{ if .UsePercent }}UsePercent:{{.UsePercent}}% \\n{{end}}{{ if .Use }}Use:{{.Use}}G \\n{{end}}{{ if .Total }}Total:{{.Total}}G\\n{{end}}{{ if .BrokenTime }}BrokenTime:{{.BrokenTime}} \\n{{end}}{{ if .FixTime }}FixTime:{{.FixTime}}\\n{{end}}{{ if .Reason }}Reason:{{.Reason}}\\n{{end}}{{ if .Top }}Top1: {{.Top}}\\n{{end}}\\n```"

type DingTalkMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	At struct {
		IsAtAll bool `json:"isAtAll"`
	} `json:"at"`
}

func (dingding *AlertDingDing) Send(body *message.Message, to ...string) error {
	text := body.FormatBody(dingdingFormat)
	// golog.Info("send dingding")
	// resp, err := http.Post(dingding.Server, "application/json;charset=utf-8",
	// 	strings.NewReader(
	// 		fmt.Sprintf(`{"msgtype": "markdown", "markdown": {"title": "alert" ,"text": "%s"}}`, text),
	// 	),
	// )
	// if err != nil {
	// 	golog.Error(err)
	// 	return err
	// }
	// defer resp.Body.Close()

	message := fmt.Sprintf(`{
  "msgtype": "markdown",
  "markdown": {
    "title": "报警标题",
    "text": "%s"
  }

}`, text)
	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", dingding.Server, strings.NewReader(message))
	if err != nil {
		golog.Error(err)
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		golog.Error(err)
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		golog.Error(err)
		return err
	}

	type DingDingResp struct {
		Errcode int `json:"errcode"`
	}
	dd := &DingDingResp{}
	err = json.Unmarshal(b, dd)
	if err != nil {
		golog.Error(string(b))
		return err
	}
	if dd.Errcode != 0 {
		golog.Error(string(b))
		return errors.New(string(b))
	}
	return nil
}
