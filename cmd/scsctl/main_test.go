package main

import (
	"sync"
	"testing"

	"github.com/hyahm/scs/client/cliconfig"
)

func BenchmarkMain(t *testing.B) {
	// 如果不是windows系统
	// 配置文件就放在 /etc/ 下面
	cliconfig.NewClientConfig()

	// command.Execute()
}

func TestMain(t *testing.T) {
	// 如果不是windows系统
	// 配置文件就放在 /etc/ 下面
	s := make([]int, 0)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {

		for i := 0; i < 30; i++ {
			s = append(s, i)
		}
		wg.Done()
	}()
	wg.Wait()
	t.Log(s)
	// command.Execute()
}
