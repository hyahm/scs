

# 客户端
1.1.0 版开始， 配置文件必须要配置值， 不然什么也不会出来
```yaml
nodes:
  localhost: 
    url: https://127.0.0.1:11111
    token:
group:
  local:
    - localhost
```

```

scsctl status 
scsctl status pname
scsctl status pname name
scsctl start 
scsctl start pname
scsctl start pname name
scsctl restart --all
scsctl restart pname 
scsctl restart pname name
scsctl kill --all
scsctl kill pname 
scsctl kill pname name
scsctl stop --all
scsctl stop pname 
scsctl stop pname name
scsctl update --all
scsctl update pname 
scsctl update pname name
scsctl remove --all
scsctl remove pname 
scsctl remove pname name
scsctl enable pname
scsctl disable pname
scsctl log  name[:update|log|lookPath] # 不区分大小写
# 加载配置文件
scsctl config reload
```

