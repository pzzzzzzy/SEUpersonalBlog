# 使用官方Go镜像作为构建阶段
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制所有源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go

# 使用Alpine作为运行时基础镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/auth-service .

# 设置环境变量
ENV PORT=8081
ENV DB_HOST=mysql
ENV DB_USER=root
ENV DB_PASSWORD=password
ENV DB_NAME=blog
ENV CONSUL_ADDR=consul:8500
ENV SERVICE_NAME=auth-service
ENV SERVICE_ID=auth-service-1
ENV SERVICE_PORT=8081

# 暴露端口
EXPOSE 8081

# 运行应用
CMD ["./auth-service"]