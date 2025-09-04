/*
 * @Author: your name
 * @Date: 2021-04-25 19:08:58
 * @LastEditTime: 2021-04-25 20:31:07
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /scs/email.go
 */
package alert

import (
	"github.com/hyahm/scs/pkg/message"
	"gopkg.in/gomail.v2"
)

var GlobalEmail *AlertEmail

type AlertEmail struct {
	Host     string   `yaml:"host,omitempty" json:"host,omitempty"`
	NickName string   `yaml:"nickname,omitempty" json:"nickname,omitempty"`
	Port     int      `yaml:"port,omitempty" json:"port,omitempty"`
	UserName string   `yaml:"username,omitempty" json:"username,omitempty"`
	Password string   `yaml:"password,omitempty" json:"password,omitempty"`
	To       []string `yaml:"to,omitempty" json:"to,omitempty"`
}

var emailFormat = `<html><head>{{.Title}}</head><h5>hostname: {{.HostName}}</h5><h6>Addr: {{.Addr}}</h6>{{ if .Pname }}<div>pname:{{.Pname}} </div>{{end}}{{ if .Name }}<div>name:{{.Name}} </div>{{end}}{{ if .DiskPath }}<div>DiskPath:{{.DiskPath}} </div>{{end}}{{ if .UsePercent }}<div>UsePercent:{{.UsePercent}}% </div>{{end}}{{ if .Use }}<div>Use:{{.Use}}G </div>{{end}}{{ if .Total }}<div>Total:{{.Total}}G </div>{{end}}{{ if .BrokenTime }}<div>BrokenTime:{{.BrokenTime}} </div>{{end}}{{ if .FixTime }}<div>FixTime:{{.FixTime}} </div>{{end}}{{ if .Reason }}<div>Reason:{{.Reason}} </div>{{end}}{{ if .Reason }}<div>Top1:{{.Top}} </div>{{end}}</html>`

// SendEmail body支持html格式字符串
func (ae *AlertEmail) Send(body *message.Message, to ...string) error {
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
	d := gomail.NewDialer(ae.Host, ae.Port, ae.UserName, ae.Password)
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// 发送
	return d.DialAndSend(m)
}
