package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hyahm/scs/message"
)

type Callback struct {
	Urls    []string      `yaml:"urls"`   // 请求url
	Method  string        `yaml:"method"` // 请求方式
	Headers http.Header   `yaml:"header"` // 请求头
	Timeout time.Duration `yaml:"timeout"`
}

func (c *Callback) Send(body *message.Message, to ...string) error {
	cli := &http.Client{
		Timeout: c.Timeout,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	m := make(map[string]struct{})
	for _, url := range c.Urls {
		m[url] = struct{}{}
	}
	for _, url := range to {
		m[url] = struct{}{}
	}
	for url := range m {
		req, err := http.NewRequest(c.Method, url, bytes.NewReader(data))
		if err != nil {
			fmt.Println(err)

			continue
		}
		req.Header = c.Headers
		resp, err := cli.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(string(b))
	}
	return nil
}
