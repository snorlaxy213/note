# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置国内镜像源加速
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

# 安装必要的包
RUN apk --no-cache add tzdata ca-certificates git

# 设置工作目录
WORKDIR /app

# 设置Go代理和环境变量
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 首先复制go.mod和go.sum文件
COPY go.mod .
COPY go.sum .

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -ldflags "-s -w" -o note-gin .

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata wget

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件和配置
COPY --from=builder /app/note-gin .
COPY --from=builder /app/config/file.example ./config/file

# 启动应用
CMD ["./note-gin", "-c", "config/file/BootLoader.yaml"]