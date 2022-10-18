FROM golang:alpine

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	   GOPROXY="https://goproxy.cn,direct"

# 移动到工作目录：/root/alliance/server
WORKDIR /root/alliance/server

# 将代码复制到容器中
COPY . .

# 代码编译成二进制可执行文件  可执行文件名为 app
RUN go build -o app .


# 移动到用于存放生成的二进制文件的 /root/alliance/server/bin 目录
WORKDIR /root/alliance/server/bin

# 将二进制文件从 /root/alliance/server 目录复制到这里
RUN cp /root/alliance/server/app .

# 启动容器时运行的命令
CMD ["/root/alliance/server/bin/app"]
