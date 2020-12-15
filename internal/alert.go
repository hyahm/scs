package internal

// 报警相关配置
type AlertTo struct {
	Email    []string `yaml:"email"`
	Rocket   []string `yaml:"rocket"`
	Telegram []string `yaml:"telegram"`
}
