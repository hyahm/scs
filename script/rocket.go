package script

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type AlertRocket struct {
	Server   string   `yaml:"server"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	To       []string `yaml:"to"`
}

type Token struct {
	UserId    string `json:"userId"`
	AuthToken string `json:"authToken"`
}

var rocketFormat = "```\nTitle: {{.Title }} \nHostname: {{.HostName}} \nAddr: {{.Addr}} \n{{ if .Pname  }}Pname:{{.Pname}}\n{{end}}{{ if .Name }}Name:{{.Name}}\n{{end}}{{ if .DiskPath }}DiskPath:{{.DiskPath}}\n{{end}}{{ if .UsePercent }}UsePercent:{{.UsePercent}}% \n{{end}}{{ if .Use }}Use:{{.Use}}G \n{{end}}{{ if .Total }}Total:{{.Total}}G\n{{end}}{{ if .BrokenTime }}BrokenTime:{{.BrokenTime}} \n{{end}}{{ if .FixTime }}FixTime:{{.FixTime}}\n{{end}}{{ if .Reason }}Reason:{{.Reason}}\n {{end}}{{ if .Top }}Top1: {{.Top}}\n{{end}}```"

func (rocket *AlertRocket) Send(body *Message, to ...string) error {

	token, err := rocket.getToken()
	if err != nil {
		fmt.Println(err)
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
			fmt.Println(err)
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
	b, err := ioutil.ReadAll(resp.Body)
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
		fmt.Println(rd.Message)
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
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/chat.postMessage", rocket.Server), bytes.NewReader(bt))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("X-User-Id", uid)
	cli := &http.Client{}
	r, err := cli.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	respmsg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Body.Close()
	a := &struct {
		Success bool `json:"success"`
		Ts      int  `json:"ts"`
	}{}

	err = json.Unmarshal(respmsg, a)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
