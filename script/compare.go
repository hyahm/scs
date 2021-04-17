package script

func CompareMap(m1, m2 map[string]string) bool {
	// 对比2个map 是否相等
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m2 {
		if mv, ok := m1[k]; ok {
			return v == mv
		} else {
			return false
		}
	}
	return true
}

func CompareSlice(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	m1 := make(map[string]string, len(s1))
	for _, v := range s1 {
		m1[v] = ""
	}
	m2 := make(map[string]string, len(s1))
	for _, v := range s1 {
		m2[v] = ""
	}
	return CompareMap(m1, m2)
}

func CompareScript(s1, s2 *Script) bool {
	if s1 == nil && s2 != nil || s1 != nil && s2 == nil {
		return false
	}
	if s1 == nil && s2 == nil {
		return true
	}
	// 这些有一个不同的。 那么就需要重启所有底下的server
	if s1.Name != s2.Name ||
		s1.Dir != s2.Dir ||
		s1.Command != s2.Command ||
		s1.Always != s2.Always ||
		!CompareMap(s1.Env, s2.Env) ||
		s1.ContinuityInterval != s2.ContinuityInterval ||
		!CompareAT(s1.AT, s2.AT) ||
		s1.DisableAlert != s2.DisableAlert ||
		s1.Disable != s2.Disable ||
		s1.Version != s2.Version ||
		!CompareCron(s1.Cron, s2.Cron) ||
		s1.Port != s2.Port {
		return false
	}

	return true
}

func CompareAT(a1, a2 *AlertTo) bool {
	if a1 == nil && a2 != nil || a1 != nil && a2 == nil {
		return false
	}
	if a1 == nil && a2 == nil {
		return true
	}
	if !CompareSlice(a1.Email, a2.Email) ||
		!CompareSlice(a1.Rocket, a2.Rocket) ||
		!CompareSlice(a1.Telegram, a2.Telegram) ||
		!CompareSlice(a1.WeiXin, a2.WeiXin) {
		return false
	}
	return true
}

func CompareCron(c1, c2 *Cron) bool {
	if c1 == nil && c2 != nil || c1 != nil && c2 == nil {
		return false
	}
	if c1 == nil && c2 == nil {
		return true
	}
	if c1.Start != c2.Start ||
		c1.IsMonth != c2.IsMonth ||
		c1.Loop != c2.Loop {
		return false
	}
	return true
}
