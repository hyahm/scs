package controller

func DelScript(pname string) error {
	err := RemoveScript(pname)
	if err != nil {
		return err
	}

	for i, s := range cfg.SC {
		if s.Name == pname {
			cfg.SC = append(cfg.SC[:i], cfg.SC[i+1:]...)
			break
		}
	}
	return cfg.WriteConfigFile()
}
