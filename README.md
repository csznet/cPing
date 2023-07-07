cPing
----

编译
==

客户端编译

    go build client.go

服务端编译

    go build server.go

或者直接Make编译，进入cPing目录后  
  
    
    Make


使用
==

客户端运行时会自动注册到服务端，需要修改conf.json文件

    {
    "name": "湖南电信",
    "server": "http://192.168.88.9:7789",
    "client": "http://192.168.88.9:7788",
    "token": "31586"
    }

`server`为服务端地址，`client`为客户端地址，`token`为客户端密钥

docker使用
==

拉取镜像  

    docker pull csznet/cping:latest

启动服务端  

    docker run -d -p 7789:7789 --env mode=server csznet/cping:latest

启动客户端  

    docker run -d -p 7788:7788 --env mode=client -v $(pwd)/c.json:/app/conf.json csznet/cping:latest

其中`-v $(pwd)/c.json:/app/conf.json`代表将当前目录下的`c.json`作为客户端配置文件