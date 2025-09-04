package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

type From struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID                          int64  `json:"id"`
	Title                       string `json:"title"`
	Type                        string `json:"type"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
}

type Message struct {
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
	Date int64  `json:"date"`
	Text string `json:"text"`
}

type Result struct {
	UpdateId int64   `json:"update_id"`
	Message  Message `json:"message"`
}

type Updates struct {
	OK     bool     `json:"ok"`
	Result []Result `json:"result"`
}

var (
	listen  = flag.String("l", ":8080", "listen default :8080")
	token   = flag.String("t", "", "token")
	logpath = flag.String("o", "", "log file path")
	update  = flag.Bool("u", false, "show update message")
	// url     = flag.String("i", "", "bot send message api // https://api.telegram.org/bot<token>/sendMessage")
)

type Text struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func proxy(w http.ResponseWriter, r *http.Request) {
	text := xmux.GetInstance(r).Data.(*Text)
	if text.ChatID == "" {
		golog.Info("chat_id is empty")
		w.Write([]byte("chat_id is empty"))
		return
	}
	if !client(text.ChatID, text.Text) {
		golog.Info("send message failed")
		w.Write([]byte("send message failed"))
		return
	}
}

func main() {
	defer golog.Sync()
	golog.InitLogger(*logpath, 0, true, 7)
	flag.Parse()
	if *token == "" {
		log.Fatal("token can not be empty")
	}
	if *update {
		res, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", *token), "application/json", nil)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		msg := &Updates{}
		err = json.Unmarshal(b, msg)
		if err != nil {
			log.Fatal(err)
		}
		show, err := json.MarshalIndent(msg, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		golog.Info(string(show))
	}
	router := xmux.NewRouter()
	router.Post("/", proxy).BindJson(Text{})
	router.Run(*listen)
}

func client(id, msg string) bool {
	r, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", *token),
		"application/json",
		strings.NewReader(fmt.Sprintf(`{
                        "chat_id": "%s",
						"text":"%s",
						"parse_mode": "markdown"
                }`, id, msg)),
	)
	if err != nil {
		golog.Error(err)
		return false
	}
	defer r.Body.Close()
	return r.StatusCode == 200
}
