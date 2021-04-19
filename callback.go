package scs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hyahm/golog"
)

type Callback struct {
	Urls    []string      `yaml:"urls"`   // 请求url
	Method  string        `yaml:"method"` // 请求方式
	Headers http.Header   `yaml:"header"` // 请求头
	Timeout time.Duration `yaml:"timeout"`
}

func (c *Callback) Send(body *Message, to ...string) error {
	cli := &http.Client{
		Timeout: c.Timeout,
	}
	data, err := json.Marshal(body)
	if err != nil {
		golog.Error(err)
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
		golog.Info(c.Headers)
		req, err := http.NewRequest(c.Method, url, bytes.NewReader(data))
		if err != nil {
			golog.Error(err)
			continue
		}
		req.Header = c.Headers
		resp, err := cli.Do(req)
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
		golog.Info(string(b))
	}
	return nil
}
