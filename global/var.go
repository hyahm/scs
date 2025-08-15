package global

import "time"

// 配置文件重载
var CanReload int64

type ConfigStore struct {
	Listen             string
	Token              string
	LogCount           int
	ReadTime           time.Duration
	IgNoreToken        []string
	Monitored          []string
	EnableTLS          bool
	Key                string
	Cert               string
	ContinuityInterval time.Duration
	LogDir             string
	CleanLog           int
}

var CS ConfigStore
var ProxyHeader string

const VERSION = "v3.8.5"
const FORMAT = "{{ .Ctime }} - [{{ .Level }}] - {{.Hostname}} - {{ .Msg }}"
