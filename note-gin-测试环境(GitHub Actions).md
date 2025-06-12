## note-gin 测试环境配置指南

### 1. GitHub Actions 自动化部署流程

本项目使用 GitHub Actions 实现了自动化的测试、构建和部署流程，主要针对测试环境的持续集成与部署。

#### 1.1 触发条件
- 推送到 master 分支

- 创建新的版本标签 (格式: v*)

- 手动触发工作流（支持自定义部署信息和环境选择）

#### 1.2 工作流程设计思路

整个 CI/CD 流程分为三个主要阶段：

**测试阶段 (test)**

````yaml
test:
  runs-on: ubuntu-latest
  steps:
    - 检出代码
    - 设置 Go 环境 (版本 1.21)
````



**构建阶段 (build)**

````yaml
build:
  needs: test  # 依赖测试阶段成功完成
  runs-on: ubuntu-latest
  if: github.event_name != 'pull_request'  # 非PR时执行
````

1. 多平台构建：同时支持 linux/amd64 和 linux/arm64 架构
2. 镜像标签策略：
   - 分支名称标签
   - 提交哈希标签
   - 语义化版本标签（适用于发布标签）
   - latest 标签（仅用于默认分支）
3. 构建缓存：使用 GitHub Actions 缓存加速构建过程

**部署阶段 (deploy-testing)**

````yaml
deploy-testing:
  needs: build  # 依赖构建阶段成功完成
  runs-on: ubuntu-latest
  if: github.ref == 'refs/heads/master' && github.event_name == 'push'  # 仅在master分支推送时执行
  environment: testing  # 使用testing环境配置
````

部署阶段的核心思路：

1. **零停机部署**：先停止旧容器，再启动新容器

2. **环境隔离**：使用卷挂载实现配置和数据的隔离

3. **资源清理**：自动清理旧镜像，避免磁盘空间浪费

#### 1.3 环境变量与密钥管理

- 使用 GitHub Secrets 存储敏感信息（Docker Hub 凭证、服务器信息）
- 使用环境变量定义镜像名称和仓库地址   

----

### 2. Dockerfile 设计思路

本项目采用了多阶段构建的 Dockerfile 设计，优化了镜像大小和构建效率。

#### 2.1 主 Dockerfile 设计

````yaml
FROM golang:1.21-alpine AS builder
# 构建阶段
...

# 运行阶段
FROM alpine:latest
...
````

设计思路：

1. **多阶段构建**：

- 第一阶段（builder）：编译 Go 应用
- 第二阶段：仅复制编译好的二进制文件和必要配置

2. **优化构建速度**：

- 分层复制依赖文件（先复制 go.mod/go.sum）
- 利用 Docker 缓存机制

3. **减小镜像体积**：

- 使用 Alpine 作为基础镜像
- 编译时使用 -ldflags "-s -w" 压缩二进制文件
- 仅复制必要的文件到最终镜像

4. **配置文件处理**：

- 通过命令行参数 -c /app/config/file/BootLoader.yaml 指定配置文件路径

----

### 3. 测试环境配置说明

#### 3.1 配置文件结构
测试环境使用以下配置文件结构：

在linux虚拟机上创建文件夹（/note-gin/testing/config/），把config/file.example下的配置文件复制过去

````yaml
/note-gin/testing/config/
├── AppConfig.yaml
├── BootLoader.yaml
├── MySqlConfig.yaml
├── RedisConfig.yaml
└── ServerConfig.yaml
````

此外BootLoader.yaml需要更改为对应的路径（与代码上的有所不同），例如我的环境：

````yaml
AppPath: /app/config/file/AppConfig.yaml
ServerPath: /app/config/file/ServerConfig.yaml
MySqlPath: /app/config/file/MySqlConfig.yaml
RedisPath: /app/config/file/RedisConfig.yaml
````

app:工作目录

config/file：挂载自动创建的

Tips：例如MySqlConfig.yaml上的配置也可根据环境更改

#### 3.2 卷挂载策略

````yaml
-v /note-gin/testing/config:/app/config/file 
-v /note-gin/testing/data:/app/data 
-v /note-gin/testing/logs:/app/logs 
````

这种挂载策略实现了：

1. 配置隔离：测试环境使用独立的配置文件
2. 数据持久化：数据和日志存储在宿主机上，容器重建不会丢失
3. 环境变量注入：通过 -e SERVER_MODE=release 设置运行模式

----


### 4. 测试环境部署流程

1. 代码提交到 master 分支：触发 GitHub Actions 工作流
2. 自动化测试：确保代码质量
3. 构建 Docker 镜像：多平台构建并推送到 Docker Hub
4. 自动部署到测试服务器：
   - 拉取最新镜像
   - 停止并删除旧容器
   - 启动新容器，挂载配置和数据卷
5. 验证部署：访问 http://<测试服务器IP>:9000/ping 确认服务正常

### 5.手动部署测试环境
 ````yaml
 # 拉取最新镜像
 docker pull vino2snax/note-gin:master
 
 # 停止并删除旧容器
 docker stop note-gin-testing || true
 docker rm note-gin-testing || true
 
 # 启动新容器
 docker run -d \
   --name note-gin-testing \
   --restart unless-stopped \
   -p 9000:9000 \
   -e SERVER_MODE=release \
   -v /note-gin/testing/config:/app/config/file \
   -v /note-gin/testing/data:/app/data \
   -v /note-gin/testing/logs:/app/logs \
   vino2snax/note-gin:master
 ````

### 6.测试环境调试技巧

1. 查看容器日志：

````yaml
docker logs -f note-gin-testing
````

2. 进入容器内部：

````yaml
docker exec -it note-gin-testing /bin/sh
````

3. 修改测试环境配置： 直接编辑宿主机上的配置文件，无需重建容器：

````yaml
vim /note-gin/testing/config/ServerConfig.yaml
````

4. 修改后重启容器：

````yaml
docker restart note-gin-testing
````

----

通过以上配置和流程，可以实现 note-gin 项目测试环境的自动化部署和管理，确保测试环境与生产环境隔离，同时提供便捷的调试和验证能力。