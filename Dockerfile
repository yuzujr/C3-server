# Dockerfile
# ——— 构建阶段 ———
FROM golang:1.20-alpine AS builder

# 安装 git 用于 go mod
RUN apk add --no-cache git

WORKDIR /build

# 拷贝 go.mod、go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 拷贝项目源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o c3-server cmd/server/main.go

# ——— 运行阶段 ———
FROM alpine:3.18

# 安装 curl 用于健康检查
RUN apk add --no-cache curl

# 新建服务用户，避免使用 root
RUN addgroup -g 1001 c3group \
 && adduser -u 1001 -G c3group -s /sbin/nologin -D c3user

WORKDIR /app

# 拷贝可执行文件和 .env
COPY --from=builder /build/c3-server ./
COPY --from=builder /build/.env ./

# 创建 uploads 和 logs 目录，设置权限
RUN mkdir -p uploads logs \
 && chown -R c3user:c3group /app

USER c3user

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:3000/ || exit 1

# 启动命令
CMD ["./c3-server"]
