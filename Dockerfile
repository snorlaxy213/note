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
# -ldflags "-s -w" 压缩编译后的二进制文件 -s：去掉符号表和调试信息，减少二进制大小。 -w：不生成 DWARF 调试信息，进一步减小体积。
RUN go build -ldflags "-s -w" -o note-gin .

# 启动应用
CMD ["./note-gin", "-c", "/app/config/file/BootLoader.yaml"]