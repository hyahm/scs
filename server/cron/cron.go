package cron

type Cron struct {
	// 开始执行的时间戳
	Start string `yaml:"start,omitempty" json:"start,omitempty"`
	// 间隔的时间， 如果IsMonth 为true， loop 单位为月， 否则为秒
	IsMonth bool `yaml:"isMonth,omitempty" json:"isMonth,omitempty"`
	Loop    int  `yaml:"loop,omitempty" json:"loop,omitempty"`
}
