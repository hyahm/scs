package scs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type AlertTelegram struct {
	Server   string  `yaml:"server"`
	Username string  `yaml:"username"`
	Password string  `yaml:"password"`
	To       []int64 `yaml:"to"`
}

var telegramFormat = "```\\nTitle: {{.Title }} \\nHostname: {{.HostName}} \\nAddr: {{.Addr}} \\n{{ if .Pname  }}pname:{{.Pname}} \\n{{end}}{{ if .Name }}name:{{.Name}} \\n{{end}}{{ if .DiskPath }}DiskPath:{{.DiskPath}}\\n{{end}}{{ if .UsePercent }}UsePercent:{{.UsePercent}}% \\n{{end}}{{ if .Use }}Use:{{.Use}}G \\n{{end}}{{ if .Total }}Total:{{.Total}}G\\n{{end}}{{ if .BrokenTime }}BrokenTime:{{.BrokenTime}} \\n{{end}}{{ if .FixTime }}FixTime:{{.FixTime}}\\n{{end}}{{ if .Reason }}Reason:{{.Reason}}\\n{{end}}{{ if .Top }}Top1: {{.Top}}\\n{{end}}\\n```"

func (telegram *AlertTelegram) Send(body *Message, to ...string) error {
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
				fmt.Sprintf(`{"chat_id": "%s", "username": "%s","password":"%s","text":"%s"}`, k, telegram.Username, telegram.Password,
					text),
			),
		)

		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			fmt.Println(err)
			continue
		}
		if resp.StatusCode != 200 {
			fmt.Println(string(b))
			continue
		}
		fmt.Println(err)
		fmt.Println(string(b))
	}
	return nil
}
