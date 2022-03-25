package alert

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg/message"
)

type AlertTelegram struct {
	Server string  `yaml:"server,omitempty" json:"server,omitempty"`
	To     []int64 `yaml:"to,omitempty" json:"to,omitempty"`
}

var telegramFormat = "```\\nTitle: {{.Title }} \\nHostname: {{.HostName}} \\nAddr: {{.Addr}} \\n{{ if .Pname  }}pname:{{.Pname}} \\n{{end}}{{ if .Name }}name:{{.Name}} \\n{{end}}{{ if .DiskPath }}DiskPath:{{.DiskPath}}\\n{{end}}{{ if .UsePercent }}UsePercent:{{.UsePercent}}% \\n{{end}}{{ if .Use }}Use:{{.Use}}G \\n{{end}}{{ if .Total }}Total:{{.Total}}G\\n{{end}}{{ if .BrokenTime }}BrokenTime:{{.BrokenTime}} \\n{{end}}{{ if .FixTime }}FixTime:{{.FixTime}}\\n{{end}}{{ if .Reason }}Reason:{{.Reason}}\\n{{end}}{{ if .Top }}Top1: {{.Top}}\\n{{end}}\\n```"

func (telegram *AlertTelegram) Send(body *message.Message, to ...string) error {
	user := make(map[string]bool)
	for _, channel := range telegram.To {
		if _, ok := user[strconv.FormatInt(channel, 10)]; !ok {
			user[strconv.FormatInt(channel, 10)] = true
		}
	}
	for _, channel := range to {
		if _, ok := user[channel]; !ok {
			user[channel] = true
		}
	}
	text := body.FormatBody(telegramFormat)
	for k := range user {
		resp, err := http.Post(telegram.Server, "application/json",
			strings.NewReader(
				fmt.Sprintf(`{"chat_id": "%s","text":"%s"}`, k, text),
			),
		)

		if err != nil {
			golog.Error(err)
			continue
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			golog.Error(err)
			continue
		}
		if resp.StatusCode != 200 {
			golog.Error(string(b))
			continue
		}
		golog.Info(string(b))
	}
	return nil
}
