[安装](install.md)  
[客户端scsctl使用](scsctl.md)  
[报警](alert.md)  
[服务添加删除接口](script.md)  
[硬件监控配置说明](hardware.md)

#  通过api添加删除挂载的服务(所有接口请带上token请求头)

可用于一个服务控制另外一个服务的作用， 请求后，配置文件自动更新配置文件  


> https://127.0.0.1:11111/script  POST <添加>
```
{
    "name": "addscript0",             
	"command" : "pwd", 
    "replicate": 10,         
	"dir": "/data/scs",
	"version": "v0.0.1",
    "disableAlert": true,
    "continuityInterval": 600, 
}
```

> script 配置文件完整版配置
```
    # 服务名 建议字母下划线数字组成（唯一）
	name               string      
    # 执行命令的根目录     
	dir                string      
    # 执行的命令，不能为空，否则不生效    
	command            string    
    #  开启的副本数       
	replicate          int      
    #  如果异常退出会相差    killTime 的时间退出     
	always             bool    
    # 警用脚本
    disable   bool  
    # 循环执行。 2次开始执行的间隔时间(s)
    loop  int
    # 禁止报警        
	disableAlert       bool   
    # 设置环境变量， 如果是PATH 就会追加， 其他的都是覆盖          
	env                map[string]string 
    #  报警的间隔
	continuityInterval time.Duration    json中单位是秒的整数，    yaml配置中是 1h的字符串 
    # 自定义端口,  会传入PORT环境变量中, 搭配    replicate 选项，会自增 
	port               int  
    # 报警收件人             
	alert:
      email: []string
      rocket: []string
      telegram: []string     
	version            string            服务的版本， 用户scsctl status 上的显示
```
```
 https://127.0.0.1:11111/delete/<pname>  POST <删除>  
```
这里的pname 就是配置文件中的name , scsctl status 中的panme



