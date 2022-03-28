package logger

import (
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/hyahm/golog"
	"gopkg.in/yaml.v2"
)

func TestLogging(t *testing.T) {
	defer golog.Sync()
	logging := golog.NewLog("log\\test_0.log", 10<<10, false)

	logging.Info("44444444")
}

func TestDir(t *testing.T) {
	a := ""
	b := "log"
	c := ".\\log"
	t.Log(filepath.Dir(a))
	t.Log(filepath.Dir(b))
	t.Log(filepath.Dir(c))
}

func TestDuration(t *testing.T) {
	defer golog.Sync()

	temps := []string{
		`time: 4m`,
		`time: 4h`,
		`time: 4s`,
	}
	for _, temp := range temps {
		prinTime(temp)
	}
}

func prinTime(temp string) {
	type du struct {
		Time time.Duration `yaml:"time"`
	}
	d := &du{}
	err := yaml.Unmarshal([]byte(temp), d)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(d.Time)
}
