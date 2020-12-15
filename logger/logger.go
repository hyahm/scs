package logger

type Logger struct {
	Path string `yaml:"path"`
	Size int64  `yaml:"size"`
	Day  bool   `yaml:"day"`
}
