# 开发包  只有go和python 版， 其他语言请参考上面的api自行封装  


## go版本， 本身就自带

`go get github.com/hyahm/scs`
```
package main

import (
	"fmt"
	"log"

	"github.com/hyahm/scs/client"
)

func main() {
	cli := client.NewClient()
	// 获取https://127.0.0.1:11111 的 脚本状态
	b, err := cli.Status()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

```
> 输出
```vim
{
        "data": [
                {
                        "name": "test_0",
                        "ppid": 0,
                        "status": "Stop",
                        "command": "python test.py",
                        "pname": "test",
                        "path": "F:\\scs",
                        "cannotStop": false,
                        "start": 0,
                        "version": "",
                        "Always": false,
                        "restartCount": 0
                }
        ],
        "code": 200
}
```

## python 版本

https://pypi.org/project/pyscs/
```
pip install pyscs
```