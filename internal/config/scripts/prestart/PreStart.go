/*
 * @Author: your name
 * @Date: 2022-02-27 10:24:34
 * @LastEditTime: 2022-02-27 10:29:08
 * @LastEditors: Please set LastEditors
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /scs/lookpath/lookpath.go
 */
package prestart

type PreStart struct {
	// 判断路径是否存在
	Path string `yaml:"path,omitempty" json:"path,omitempty"`
	// 判断命令是否存在
	Command string `yaml:"command,omitempty" json:"command,omitempty"`
	// 执行命令
	ExecCommand string `yaml:"execCommand,omitempty" json:"execCommand,omitempty"`
	// 比较的分隔符， 比较版本类似 v1.4.6， 需要设置为.
	Separation string `yaml:"separation,omitempty" json:"separation,omitempty"`
	// 与 exec_command 配合对比
	// 	-eq           等于
	// -ne           不等于
	// -gt            大于
	// -lt            小于
	// -ge            大于等于
	// -le            小于等于
	EQ string `yaml:"eq,omitempty" json:"eq,omitempty"`
	NE string `yaml:"ne,omitempty" json:"ne,omitempty"`
	GT string `yaml:"gt,omitempty" json:"gt,omitempty"`
	LT string `yaml:"lt,omitempty" json:"lt,omitempty"`
	GE string `yaml:"ge,omitempty" json:"ge,omitempty"`
	LE string `yaml:"le,omitempty" json:"le,omitempty"`
	// 执行的命令
	Install string `yaml:"install,omitempty" json:"install,omitempty"`
	// 配置文件默认模板
	Template string `yaml:"template,omitempty" json:"template,omitempty"`
}

// 判断2个 preStart 的值是否相同
func EqualPreStart(p1, p2 []*PreStart) bool {
	if len(p1) != len(p2) {
		return false
	}
	for i := 0; i < len(p1); i++ {
		if !equalPreStart(p1[i], p2[i]) {
			return false
		}
	}
	return true
}

func equalPreStart(p1, p2 *PreStart) bool {

	return !(p1.Path != p2.Path ||
		p1.Command != p2.Command ||
		p1.ExecCommand != p2.ExecCommand ||
		p1.Separation != p2.Separation ||
		p1.EQ != p2.EQ ||
		p1.NE != p2.NE ||
		p1.GT != p2.GT ||
		p1.LT != p2.LT ||
		p1.GE != p2.GE ||
		p1.LE != p2.LE ||
		p1.Install != p2.Install ||
		p1.Template != p2.Template)
}
