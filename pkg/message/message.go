package message

// 用来组装body
import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"text/template"

	"github.com/hyahm/golog"
)

var addr string

type Response struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
	Time    string `json:"time"`
}

// 定义 Data 结构体
type Data struct {
	IP string `json:"ip"`
}

func GetIp() {
	url := "https://cz88.net/api/cz88/ip/base?ip="
	resp, err := http.Get(url)
	if err != nil {
		addr = "127.0.0.1"
		golog.Error(err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		addr = "127.0.0.1"
		golog.Error(err)
		return
	}
	res := Response{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		addr = "127.0.0.1"
		golog.Error(err)
		return
	}

	// 将响应内容转换为字符串
	addr = res.Data.IP
}

type Message struct {
	Key        string  `json:"key,omitempty"`
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

func (msg *Message) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		golog.Error(err)
	}
	return string(b)
}

func (am *Message) FormatBody(format string) string {
	am.Addr = addr
	buf := bytes.NewBuffer([]byte(""))
	tmpl, err := template.New("test").Parse(format) //建立一个名字为test的模版"hello, {{.}}"
	if err != nil {
		golog.Error(err)
		return ""
	}

	err = tmpl.Execute(buf, am) //将str的值合成到tmpl模版的{{.}}中，并将合成得到的文本输入到os.Stdout,返回hello, world
	if err != nil {
		golog.Error(err)
		return ""
	}
	return buf.String()
}
