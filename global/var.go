package global

import "time"

// 配置文件重载
var CanReload int64

var _listen string
var _token string
var _logCount int
var _ignoreToken []string
var _monitored []string // 允许监控的服务器
var _disableTls bool
var _key string
var _pem string
var _continuityInterval time.Duration
var ProxyHeader string

const VERSION = "v3.7.0"
const FORMAT = "{{ .Ctime }} - [{{ .Level }}] - {{.Hostname}} - {{ .Msg }}"

var LogDir = "" // 日志目录
var CleanLog time.Duration

func SetToken(token string) {
	_token = token
}

func GetToken() string {
	return _token
}

func SetListen(listen string) {
	if listen == "" {
		listen = ":11111"
	}
	_listen = listen
}

func GetListen() string {
	return _listen
}

func SetLogCount(count int) {
	if count == 0 {
		count = 100
	}
	_logCount = count
}

func GetLogCount() int {
	return _logCount
}

func GetIgnoreToken() []string {
	return _ignoreToken
}

func SetIgnoreToken(it []string) {
	_ignoreToken = it
}

func GetMonitored() []string {
	return _monitored
}

func SetMonitored(it []string) {
	_monitored = it
}

func GetDisableTls() bool {
	return _disableTls
}

func SetDisableTls(tls bool) {
	_disableTls = tls
}

func SetKey(key string) {
	_key = key
}

func GetKey() string {
	return _key
}

func SetPem(pem string) {
	_pem = pem
}

func GetPem() string {
	return _pem
}

func SeContinuityInterval(ci time.Duration) {
	if ci == 0 {
		_continuityInterval = time.Hour * 1
	}
	_continuityInterval = ci
}

func GeContinuityInterval() time.Duration {
	return _continuityInterval
}
