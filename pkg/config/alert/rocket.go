package alert

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg/message"
)

type AlertRocket struct {
	Server   string   `yaml:"server,omitempty" json:"server,omitempty"`
	Username string   `yaml:"username,omitempty" json:"username,omitempty"`
	Password string   `yaml:"password,omitempty" json:"password,omitempty"`
	To       []string `yaml:"to,omitempty" json:"to,omitempty"`
}

type Token struct {
	UserId    string `json:"userId"`
	AuthToken string `json:"authToken"`
}

var rocketFormat = "```\nTitle: {{.Title }} \nHostname: {{.HostName}} \nAddr: {{.Addr}} \n{{ if .Pname  }}Pname:{{.Pname}}\n{{end}}{{ if .Name }}Name:{{.Name}}\n{{end}}{{ if .DiskPath }}DiskPath:{{.DiskPath}}\n{{end}}{{ if .UsePercent }}UsePercent:{{.UsePercent}}% \n{{end}}{{ if .Use }}Use:{{.Use}}G \n{{end}}{{ if .Total }}Total:{{.Total}}G\n{{end}}{{ if .BrokenTime }}BrokenTime:{{.BrokenTime}} \n{{end}}{{ if .FixTime }}FixTime:{{.FixTime}}\n{{end}}{{ if .Reason }}Reason:{{.Reason}}\n {{end}}{{ if .Top }}Top1: {{.Top}}\n{{end}}```"

func (rocket *AlertRocket) Send(body *message.Message, to ...string) error {
	token, err := rocket.getToken()
	if err != nil {
		golog.Error(err)
		return err
	}
	user := make(map[string]bool)
	for _, channel := range rocket.To {
		if _, ok := user[channel]; !ok {
			user[channel] = true
		}
	}
	for _, channel := range to {
		if _, ok := user[channel]; !ok {
			user[channel] = true
		}
	}

	for k := range user {
		err = rocket.sendMsg(token.AuthToken, token.UserId, body.FormatBody(rocketFormat), k)
		if err != nil {
			golog.Error(err)
			continue
		}
	}
	return err
}

func (rocket *AlertRocket) getToken() (*Token, error) {
	login := strings.NewReader(fmt.Sprintf(`{
		"username": "%s",
		"password": "%s"
	}`, rocket.Username, rocket.Password))
	resp, err := http.Post(fmt.Sprintf("%s/api/v1/login", rocket.Server), "application/json", login)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(string(b))
	}
	rd := &struct {
		Status  string `json:"status"`
		Data    Token  `json:"data"`
		Message string `json:"message"`
		Error   int    `json:"error"`
	}{}
	err = json.Unmarshal(b, rd)
	if err != nil {
		return nil, err
	}
	if rd.Status == "error" {
		golog.Error(rd.Message)
		return nil, fmt.Errorf(rd.Message)
	}
	return &rd.Data, nil
}

func (rocket *AlertRocket) sendMsg(token, uid, body, to string) (err error) {
	msg := struct {
		Channel string `json:"channel"`
		Text    string `json:"text"`
	}{
		Channel: to,
		Text:    body,
	}
	bt, err := json.Marshal(msg)
	if err != nil {
		golog.Error(err)
		return
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/chat.postMessage", rocket.Server), bytes.NewReader(bt))
	if err != nil {
		golog.Error(err)
		return
	}

	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("X-User-Id", uid)
	cli := &http.Client{}
	r, err := cli.Do(req)
	if err != nil {
		golog.Error(err)
		return
	}
	respmsg, err := io.ReadAll(r.Body)
	if err != nil {
		golog.Error(err)
		return
	}
	defer r.Body.Close()
	a := &struct {
		Success bool `json:"success"`
		Ts      int  `json:"ts"`
	}{}

	err = json.Unmarshal(respmsg, a)
	if err != nil {
		golog.Error(err)
		return
	}
	return
}
