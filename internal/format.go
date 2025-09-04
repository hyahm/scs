package internal

import (
	"bytes"
	"text/template"
	"time"

	"github.com/hyahm/golog"
)

// 格式化 text/template
func Format(format string, data interface{}) string {
	tpl, err := template.New(time.Now().String()).Parse(format)
	if err != nil {
		golog.Error(err)
		return ""
	}
	buf := bytes.NewBuffer(nil)
	err = tpl.Execute(buf, data)
	if err != nil {
		golog.Error(err)
		return ""
	}
	return buf.String()
}
