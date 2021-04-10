package alert

import (
	"gopkg.in/gomail.v2"
)

var GlobalEmail *AlertEmail

type AlertEmail struct {
	Host     string   `yaml:"host"`
	NickName string   `yaml:"nickname"`
	Port     int      `yaml:"port"`
	UserName string   `yaml:"username"`
	Password string   `yaml:"password"`
	To       []string `yaml:"to"`
}

var emailFormat = `<html><head>{{.Title}}</head><h5>hostname: {{.HostName}}</h5><h6>Addr: {{.Addr}}</h6>{{ if .Pname }}<div>pname:{{.Pname}} </div>{{end}}{{ if .Name }}<div>name:{{.Name}} </div>{{end}}{{ if .DiskPath }}<div>DiskPath:{{.DiskPath}} </div>{{end}}{{ if .UsePercent }}<div>UsePercent:{{.UsePercent}}% </div>{{end}}{{ if .Use }}<div>Use:{{.Use}}G </div>{{end}}{{ if .Total }}<div>Total:{{.Total}}G </div>{{end}}{{ if .BrokenTime }}<div>BrokenTime:{{.BrokenTime}} </div>{{end}}{{ if .FixTime }}<div>FixTime:{{.FixTime}} </div>{{end}}{{ if .Reason }}<div>Reason:{{.Reason}} </div>{{end}}{{ if .Reason }}<div>Top1:{{.Top}} </div>{{end}}</html>`

// SendEmail body支持html格式字符串
func (ae *AlertEmail) Send(body *Message, to ...string) error {
	// 主题
	m := gomail.NewMessage()
	receive := make([]string, 0, len(ae.To)+len(to))
	// 去重收件人
	duplicate := make(map[string]bool)
	for _, v := range ae.To {
		if _, ok := duplicate[v]; !ok {
			duplicate[v] = false
			receive = append(receive, v)
		}
	}
	for _, v := range to {
		if _, ok := duplicate[v]; !ok {
			duplicate[v] = false
			receive = append(receive, v)
		}
	}
	// 收件人可以有多个，故用此方式
	m.SetHeader("To", receive...)
	//抄送列表
	m.SetHeader("Cc", receive...)
	// 发件人
	// 第三个参数为发件人别名，如"李大锤"，可以为空（此时则为邮箱名称）
	m.SetAddressHeader("From", ae.UserName, ae.NickName)
	m.SetHeader("Subject", body.Title)
	// 正文
	m.SetBody("text/html", body.FormatBody(emailFormat))
	d := gomail.NewPlainDialer(ae.Host, ae.Port, ae.UserName, ae.Password)
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// 发送
	return d.DialAndSend(m)
}
