package alert

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hyahm/golog"
)

type AlertWeiXin struct {
	Server string `yaml:"server"`
}

var weixinFormat = "```\\nTitle: {{.Title }} \\nHostname: {{.HostName}} \\nAddr: {{.Addr}} \\n{{ if .Pname  }}pname:{{.Pname}} \\n{{end}}{{ if .Name }}name:{{.Name}} \\n{{end}}{{ if .DiskPath }}DiskPath:{{.DiskPath}}\\n{{end}}{{ if .UsePercent }}UsePercent:{{.UsePercent}}% \\n{{end}}{{ if .Use }}Use:{{.Use}}G \\n{{end}}{{ if .Total }}Total:{{.Total}}G\\n{{end}}{{ if .BrokenTime }}BrokenTime:{{.BrokenTime}} \\n{{end}}{{ if .FixTime }}FixTime:{{.FixTime}}\\n{{end}}{{ if .Reason }}Reason:{{.Reason}}\\n{{end}}{{ if .Top }}Top1: {{.Top}}\\n{{end}}\\n```"

func (weixin *AlertWeiXin) Send(body *Message, to ...string) error {
	text := body.FormatBody(telegramFormat)
	resp, err := http.Post(weixin.Server, "application/json",
		strings.NewReader(
			fmt.Sprintf(`{"msgtype": "markdown", "markdown": {"content": "%s"}}`, text),
		),
	)
	defer resp.Body.Close()
	if err != nil {
		golog.Error(err)
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		golog.Error(err)
		return err
	}
	if resp.StatusCode != 200 {
		golog.Error(string(b))
		return errors.New(string(b))
	}
	golog.Error(err)
	golog.Info(string(b))
	return nil
}
