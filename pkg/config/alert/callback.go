package alert

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg/message"
)

type Callback struct {
	Urls    []string      `yaml:"urls,omitempty" json:"urls,omitempty"`     // 请求url
	Headers http.Header   `yaml:"header,omitempty" json:"header,omitempty"` // 请求头
	Timeout time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
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
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
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
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			golog.Error(err)
			continue
		}
		golog.Error(string(b))
	}
	return nil
}
