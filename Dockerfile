FROM golang:latest as builder

# 设置工作目录
WORKDIR /src

# 复制源代码到容器中
COPY . .

# 编译二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kv-raft .

FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从builder镜像中复制编译好的二进制文件
COPY --from=builder /src/kv-raft /app/kv-raft
COPY conf/app-docker-voter.yaml /app/app.yaml

EXPOSE 2315

# 设置执行权限
RUN chmod +x /app/kv-raft

# 运行二进制文件
CMD ["./kv-raft"]