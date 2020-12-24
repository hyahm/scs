package script

func (s *Script) GetEnv() []string {
	return s.cmd.Env
}
