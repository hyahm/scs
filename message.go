package scs

// 用来组装body
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

var Addr string

func GetIp() {
	r, err := http.Get("http://ip.hyahm.com")
	if err != nil {
		return
	}
	defer r.Body.Close()
	b, _ := ioutil.ReadAll(r.Body)
	Addr = string(b)
}

type Message struct {
	Title      string  `json:"Title,omitempty"`
	HostName   string  `json:"HostName,omitempty"`
	Pname      string  `json:"Pname,omitempty"`
	Name       string  `json:"Name,omitempty"`
	Addr       string  `json:"Addr,omitempty"`
	DiskPath   string  `json:"DiskPath,omitempty"`
	Use        uint64  `json:"Use,omitempty"`
	UsePercent float64 `json:"UsePercent,omitempty"`
	Total      uint64  `json:"Total,omitempty"`
	BrokenTime string  `json:"BrokenTime,omitempty"`
	FixTime    string  `json:"FixTime,omitempty"`
	Reason     string  `json:"Reason,omitempty"`
	Top        string  `json:"Top,omitempty"`
}

func (am *Message) FormatBody(format string) string {
	am.Addr = Addr
	buf := bytes.NewBuffer([]byte(""))
	tmpl, err := template.New("test").Parse(format) //建立一个名字为test的模版"hello, {{.}}"
	if err != nil {
		fmt.Println(err)
		return ""
	}

	err = tmpl.Execute(buf, am) //将str的值合成到tmpl模版的{{.}}中，并将合成得到的文本输入到os.Stdout,返回hello, world
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return buf.String()
}
