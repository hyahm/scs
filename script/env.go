package script

import "github.com/hyahm/golog"

func (s *Script) GetEnv() []string {
	env := make([]string, 0, len(s.Env))
	for k, v := range s.Env {
		env = append(env, k+"="+v)
	}
	golog.Info(env)
	return env
}
