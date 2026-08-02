package config

// var Cfg = &Config{}

type Repo struct {
	Url        []string `yaml:"url,omitempty" json:"url"`
	Derivative string   `yaml:"derivative,omitempty" json:"derivative"`
}

type Config struct {
	Listen string `yaml:"listen,omitempty"`
	Token  string `yaml:"token,omitempty"`
	// ProxyHeader string `yaml:"proxyHeader,omitempty"`
	Key       string `yaml:"key,omitempty"`
	Cert      string `yaml:"cert,omitempty"`
	EnableTLS bool   `yaml:"enableTLS,omitempty"`
	Debug     bool   `yaml:"debug,omitempty"`
	// Packet      bool     `yaml:"packet,omitempty"`
	Log         Logger   `yaml:"log,omitempty"`
	IgnoreToken []string `yaml:"ignoreToken,omitempty"`
	// ReadTimeout time.Duration `yaml:"readTimeout,omitempty"`
	// Repo        *Repo          `yaml:"repo,omitempty"`
	Alert   Alert    `yaml:"alert,omitempty"`
	Probe   Probe    `yaml:"probe,omitempty"`
	Repo    *Repo    `yaml:"repo,omitempty"`
	Scripts []Script `yaml:"scripts,omitempty"`
}

// 保存的配置文件路径

// 读文件
