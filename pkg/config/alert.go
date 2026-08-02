package config

import (
	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg/message"
)

// Alert 聚合所有报警通道的配置。
type Alert struct {
	Email    AlertEmail    `yaml:"email,omitempty" json:"email,omitempty"`
	Rocket   AlertRocket   `yaml:"rocket,omitempty" json:"rocket,omitempty"`
	Telegram AlertTelegram `yaml:"telegram,omitempty" json:"telegram,omitempty"`
	WeiXin   AlertWeiXin   `yaml:"weixin,omitempty" json:"weixin,omitempty"`
	Callback Callback      `yaml:"callback,omitempty" json:"callback,omitempty"`
	DingDing AlertDingDing `yaml:"dingding,omitempty" json:"dingding,omitempty"`
}

var alerter = &Alerter{
	Alerts: make(map[string]message.SendAlerter),
}

type Alerter struct {
	Alert  Alert
	Alerts map[string]message.SendAlerter
}

func GetAlerts() map[string]message.SendAlerter {
	return alerter.Alerts
}

// InitAlert 根据配置初始化各报警器实例，存入 alerter.Alerts。
func InitAlert(cfg Alert) {
	alerter.Alert = cfg
	alerter.Alerts = make(map[string]message.SendAlerter)
	initAlertChannels(cfg)
}

// ReloadAlert 热更新报警通道配置，关闭旧通道并初始化新通道。
func ReloadAlert(cfg Alert) {
	alerter.Alert = cfg
	alerter.Alerts = make(map[string]message.SendAlerter)
	initAlertChannels(cfg)
}

func initAlertChannels(cfg Alert) {
	if cfg.Email.Host != "" && cfg.Email.UserName != "" && cfg.Email.Password != "" {
		email := cfg.Email
		if email.Port == 0 {
			email.Port = 465
		}
		alerter.Alerts["email"] = &email
		golog.Info("报警通道初始化: email")
	}
	if cfg.Rocket.Server != "" && cfg.Rocket.Username != "" && cfg.Rocket.Password != "" {
		alerter.Alerts["rocket"] = &cfg.Rocket
		golog.Info("报警通道初始化: rocket")
	}
	if cfg.Telegram.Server != "" {
		alerter.Alerts["telegram"] = &cfg.Telegram
		golog.Info("报警通道初始化: telegram")
	}
	if cfg.WeiXin.Server != "" {
		alerter.Alerts["weixin"] = &cfg.WeiXin
		golog.Info("报警通道初始化: weixin")
	}
	if cfg.DingDing.Server != "" {
		alerter.Alerts["dingding"] = &cfg.DingDing
		golog.Info("报警通道初始化: dingding")
	}
	if cfg.Callback.Urls != nil && len(cfg.Callback.Urls) > 0 {
		alerter.Alerts["callback"] = &cfg.Callback
		golog.Info("报警通道初始化: callback")
	}
}
