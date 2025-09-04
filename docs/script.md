
# 脚本配置文件(以下配置都隶属于 `scripts`下)

> script 配置文件完整版配置
```yaml
    # 服务名 只能是字母下划线数字组成（唯一）
	name               string      
    # 执行命令的根目录     
	dir                string      
    # 执行的命令，不能为空，否则不生效    
	command            string    
    #  开启的副本数       
	replicate          int      
    #  如果异常退出会相差   
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



