# Note-Gin 笔记系统
## 项目简介
Note-Gin 是一个基于 Gin 框架开发的在线笔记系统，提供文章、文件夹管理、图书管理等功能。系统采用 Go 语言开发后端 API，支持 Markdown 格式的笔记编辑与存储，并提供了完整的文件分类、标签管理功能。

## 主要功能模块
### 1. 文章管理
- 创建、编辑、删除文章
- Markdown 格式支持
- 文章标签管理
- 文章临时保存与恢复
- 回收站功能
- 批量操作支持
### 2. 文件夹管理
- 多级文件夹结构
- 文件夹导航
- 子文件夹与文章管理
### 3. 图书管理
- 图书信息记录（标题、作者、封面）
- 阅读状态管理（在读、读完、想读）
### 4. 七牛云存储集成
- 图片上传与管理
- 资源外链支持
### 5. 数据存储
- MySQL 数据库存储
- Redis 缓存支持
- GORM 数据库迁移
## 技术栈
- 后端框架：Gin
- 数据库：MySQL
- 缓存：Redis
- ORM：GORM
- 对象存储：七牛云
- 容器化：Docker
- CI/CD：GitHub Actions
## 项目结构
````tex
├── config/             # 配置相关
├── controller/         # 控制器
├── docker/             # Docker配置
├── middleware/         # 中间件
├── models/             # 数据模型
├── pkg/                # 工具包
│   ├── HttpCode/       # HTTP状态码
│   ├── QiniuClient/    # 七牛云客户端
│   ├── RedisClient/    # Redis客户端
│   ├── logging/        # 日志
│   └── utils/          # 工具函数
├── router/             # 路由
├── service/            # 服务层
└── view/               # 视图模型
````



## 贡献指南

- Fork 本仓库
- 创建特性分支 ( git checkout -b feature/amazing-feature )
- 提交更改 ( git commit -m 'Add some amazing feature' )
- 推送到分支 ( git push origin feature/amazing-feature )
- 创建 Pull Request



## 许可证

本项目采用 Apache License 许可证 - 详情见 LICENSE 文件