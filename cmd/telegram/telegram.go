package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

var (
	listen   = flag.String("l", "", "listen default :8080")
	username = flag.String("u", "", "username")
	password = flag.String("p", "", "password")
	url      = flag.String("i", "", "bot send message api // https://api.telegram.org/bot<token>/sendMessage")
)

type Text struct {
	ChatID   string `json:"chat_id"`
	Text     string `json:"text"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

func proxy(w http.ResponseWriter, r *http.Request) {
	text := &Text{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(text)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if text.UserName != *username || text.Password != *password {
		w.WriteHeader(http.StatusProxyAuthRequired)
		return
	}
	golog.Info(text)
	if text.ChatID == "" {
		w.Write([]byte("chat_id is empty"))
		return
	}
	if !client(text.ChatID, text.Text) {
		w.Write([]byte("send message failed"))
		return
	}
	return
}

func main() {
	flag.Parse()
	router := xmux.NewRouter()
	router.Post("/", proxy)
	router.Run(*listen)
}

func client(id, msg string) bool {
	r, err := http.Post(
		*url,
		"application/json",
		strings.NewReader(fmt.Sprintf(`{
                        "chat_id": "%s",
						"text":"%s",
						"parse_mode": "markdown"
                }`, id, msg)),
	)
	defer r.Body.Close()
	if err != nil {
		golog.Error(err)
		return false
	}
	return r.StatusCode == 200
}
