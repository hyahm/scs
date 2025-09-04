package pkg

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
	m2 := make(map[string]string, len(s2))
	for _, v := range s2 {
		m2[v] = ""
	}
	return CompareMap(m1, m2)
}
