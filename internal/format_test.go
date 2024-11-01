package internal

import "testing"

func TestFormat(t *testing.T) {
	m := make(map[string]string)
	m["aaa"] = "bbb"
	m["PARAMETER"] = "windows"
	t.Log()
	t.Log(Format("python test.py {{ .PARAMETER }}", m))
	t.Log(Format("{{ if .aaa }}{{.aaa}}{{else}}88888{{end}}", m))
	t.Log(Format(`{{ if eq .aaa "bbb" }}1111111111{{else}}2222222222222{{end}}`, m))
	t.Log(Format(`{{ if eq .OS "windows" }}main.exe{{ else }} main{{ end}}`, m))
}
