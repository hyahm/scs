package script

func (s *Script) GetEnv() []string {
	env := make([]string, 0, len(s.Env))
	for k, v := range s.Env {
		env = append(env, k+"="+v)
	}
	return env
}
