package config

import "regexp"

func CheckScriptNameRule(name string) bool {
	reg, _ := regexp.Compile(`^\w+$`)
	return reg.MatchString(name)
}
