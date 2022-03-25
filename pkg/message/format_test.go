package message

import (
	"bytes"
	"log"
	"testing"
	"text/template"
)

func TestFormat(t *testing.T) {
	a := `{{ .a }}is {{ .b}}`
	b := make(map[string]string)
	b["C"] = "aaa"
	b["a"] = "aaa"
	b["b"] = "bbb"
	tpl, err := template.New("test").Parse(a)
	if err != nil {
		log.Fatal(err)
	}
	buf := bytes.NewBuffer(nil)
	err = tpl.Execute(buf, b)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(buf.String())
}
