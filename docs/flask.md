## Flask部署 

先创建一个简单的例子
- python3.10

> 创建一个flask虚拟环境， 并激活
```powershell
PS E:\> python -m venv test
PS E:\> cd .\test\
PS E:\test> .\Scripts\activate
(test) PS E:\test>  // 进入虚拟环境后前面会出现一个虚拟目录名
(test) PS E:\test> ls
Directory: E:\test
Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
d-----         2022/3/25     22:42                Include
d-----         2022/3/25     22:42                Lib
d-----         2022/3/25     22:42                Scripts
-a----         2022/3/25     22:42             91 pyvenv.cfg
```
> 安装`Flask`先
```
pip install Flask
```
> 在test目录下面创建一个基础文件
```
from flask import Flask

app = Flask(__name__)

@app.route("/")
def hello_world():
    return "<p>Hello, World!</p>"
```