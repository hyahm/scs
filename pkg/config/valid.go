package config

import "regexp"

// 检查配置文件name的值是否有效
func CheckScriptNameRule(name string) bool {
	reg, _ := regexp.Compile(`^\w+$`)
	return reg.MatchString(name)
}
