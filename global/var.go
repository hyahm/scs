package global

import "time"

var Listen string
var Token string
var LogCount int
var IgnoreToken []string
var Monitored []string // 允许监控的服务器
var DisableTls bool
var Key string
var Pem string
var ContinuityInterval time.Duration

const VERSION = "v3.1.1"
