# scs

service control service or script  
类似supervisor,但是更高级，支持所有系统   
自带监控及通知     
服务控制脚本能否停止 最大程度防止脚本数据丢失   
码云地址: https://gitee.com/cander/scs

# 主要功能
1.  为了保护数据不会在执行中被人工手动中断丢失， 可以让服务在某段时间内才能停止（数据保证）
2.  监控硬件信息， 主要是磁盘， cpu， 内存,  scs服务  
3.  服务之间可以相互控制添加删除启动停止  
4.  报警功能api(适合代码内部问题的报警) 
5.  支持定时器功能执行命令或脚本
6.  客户端控制多台服务器 
7.  通过配置文件安装服务（简化安装）
8.  远程给某些脚本做更新和查看日志的权限（运维开发权限的桥梁）

[文档地址](https://scs.hyahm.com/#/)
[部分视频教程](https://www.bilibili.com/video/BV1bv411C7Qz/)

具体更新的内容请查看 [update.md](update.md)文件

QQ群:  346746477

