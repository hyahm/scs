package config

type Logger struct {
	Path  string `yaml:"path"`
	Clear int    `json:"clear"`
}
