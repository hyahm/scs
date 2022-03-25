/*
 * @Author: your name
 * @Date: 2022-01-16 22:40:00
 * @LastEditTime: 2022-01-20 21:13:39
 * @LastEditors: your name
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /scs/probe/check.go
 */
package probe

type CheckPointer interface {
	Check()
	Update()
}

func test() {
}
