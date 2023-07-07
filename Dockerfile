# 使用官方的 Ubuntu 基础镜像
FROM ubuntu:latest

# 安装必要的依赖
RUN apt-get update && \
    apt-get install -y net-tools iproute2 procps

# 将编译好的 server 和 client 二进制文件复制到容器中
COPY server /app/server
COPY client /app/client

# 设置工作目录
WORKDIR /app

# 设置暴露的端口
EXPOSE 7788 7789

# 设置容器启动时执行的命令，其中 mode 默认值为 "client"
CMD [ "/bin/bash", "-c", "if [ \"$mode\" = \"server\" ]; then ./server; else ./client; fi" ]
