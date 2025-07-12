# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖（使用阿里云 Go 模块代理加速）
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/ && go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 生产阶段
FROM alpine:latest

# 安装ca证书
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY config.yaml .

# 创建日志目录
RUN mkdir -p /var/log/app

EXPOSE 8080

CMD ["./main"]