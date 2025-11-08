# 个人博客微服务架构

这是一个基于微服务架构的个人博客系统，采用Golang开发，使用Consul进行服务注册发现，Traefik作为API网关，MySQL作为数据库。

## 技术栈

| 组件类型 | 工具选择 | 核心原因 | 复杂度 |
| :--- | :--- | :--- | :--- |
| 业务服务开发 | Golang (1.20) | 无需学习新语言，直接复用原代码逻辑 | 低 |
| 服务注册发现 | Consul | 单二进制文件启动，支持KV配置，不用额外搭配置中心 | 低 |
| API网关 | Traefik | 自动发现Consul中的服务，配置极简（YAML即可） | 低 |
| 数据库 | MySQL | 沿用原SQL逻辑，数据持久化稳定 | 极低 |
| 服务间通信 | HTTP (REST) | 不用定义proto，直接用原单体的接口格式 | 极低 |
| 部署方式 | Docker + Docker Compose | 一键启动所有组件，不用手动配置环境 | 低 |

## 目录结构

```
your-blog-microservices/  # 总目录
├── auth-service/         # 认证服务（独立项目）
│   ├── cmd/              # 入口文件
│   ├── config/           # 配置文件
│   ├── models/           # 数据模型
│   ├── handlers/         # 请求处理器
│   ├── middleware/       # 中间件
│   ├── utils/            # 工具函数
│   ├── db/               # 数据库连接
│   └── Dockerfile        # 服务镜像构建文件
├── article-service/      # 文章服务（独立项目）
│   ├── cmd/              # 入口文件
│   ├── config/           # 配置文件
│   ├── models/           # 数据模型
│   ├── handlers/         # 请求处理器
│   ├── middleware/       # 中间件
│   ├── utils/            # 工具函数
│   ├── db/               # 数据库连接
│   └── Dockerfile        # 服务镜像构建文件
├── comment-service/      # 评论服务（独立项目）
│   ├── cmd/              # 入口文件
│   ├── config/           # 配置文件
│   ├── models/           # 数据模型
│   ├── handlers/         # 请求处理器
│   ├── middleware/       # 中间件
│   ├── utils/            # 工具函数
│   ├── db/               # 数据库连接
│   └── Dockerfile        # 服务镜像构建文件
├── configs/              # 全局配置文件
│   ├── init.sql          # 数据库初始化脚本
│   ├── consul.json       # Consul配置
│   ├── traefik.yml       # Traefik静态配置
│   └── dynamic-config.yml # Traefik动态配置
├── docker-compose.yml    # 一键启动所有组件
└── README.md             # 启动说明文档
```

## 服务说明

### 认证服务 (auth-service)
- **端口**: 8081
- **主要功能**: 用户注册、登录、JWT令牌验证
- **API前缀**: `/api/auth`

### 文章服务 (article-service)
- **端口**: 8082
- **主要功能**: 文章的增删改查、点赞功能
- **API前缀**: `/api/articles`

### 评论服务 (comment-service)
- **端口**: 8083
- **主要功能**: 评论的增删改查
- **API前缀**: `/api/comments`, `/api/user/comments`

## 快速启动

### 前提条件
- 已安装Docker和Docker Compose

### 启动步骤

1. 进入项目目录
```bash
cd your-blog-microservices
```

2. 启动所有服务
```bash
docker-compose up -d
```

3. 验证服务是否正常启动
```bash
docker-compose ps
```

## 访问地址

- **API网关**: http://localhost
- **Consul管理界面**: http://localhost:8500
- **Traefik管理界面**: http://localhost:8080

## API接口说明

### 认证服务

#### 用户注册
```
POST /api/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

#### 用户登录
```
POST /api/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

### 文章服务

#### 创建文章
```
POST /api/articles
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "title": "测试文章",
  "content": "这是一篇测试文章的内容"
}
```

#### 获取文章列表
```
GET /api/articles
```

#### 获取文章详情
```
GET /api/articles/:id
```

#### 更新文章
```
PUT /api/articles/:id
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "title": "更新后的标题",
  "content": "更新后的内容"
}
```

#### 删除文章
```
DELETE /api/articles/:id
Authorization: Bearer <JWT_TOKEN>
```

#### 点赞/取消点赞文章
```
POST /api/articles/:id/like
Authorization: Bearer <JWT_TOKEN>
```

### 评论服务

#### 创建评论
```
POST /api/comments
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "article_id": "文章ID",
  "content": "这是一条评论"
}
```

#### 获取文章评论列表
```
GET /api/comments/article/:article_id
```

#### 删除评论
```
DELETE /api/comments/:id
Authorization: Bearer <JWT_TOKEN>
```

#### 获取用户评论列表
```
GET /api/user/comments
Authorization: Bearer <JWT_TOKEN>
```

## 环境变量说明

### 通用环境变量
- `PORT`: 服务监听端口
- `DB_HOST`: 数据库主机地址
- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名称
- `CONSUL_ADDR`: Consul服务地址
- `SERVICE_NAME`: 服务名称
- `SERVICE_ID`: 服务ID
- `SERVICE_PORT`: 服务端口

### 特定服务环境变量
- `AUTH_SERVICE`: 认证服务地址（文章和评论服务需要）

## 停止服务

```bash
docker-compose down
```

## 查看服务日志

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs auth-service
```

## 注意事项

1. 本项目使用的是开发环境配置，生产环境需要调整安全设置
2. 数据库密码在生产环境中应该使用更安全的方式管理
3. JWT密钥在生产环境中应该使用强密钥
4. 建议在生产环境中启用HTTPS
5. 服务间通信可以考虑使用加密通道

## 故障排除

1. **服务无法启动**：检查Docker和Docker Compose是否正常安装，端口是否被占用
2. **数据库连接失败**：检查环境变量配置，确保MySQL服务正常运行
3. **服务注册失败**：检查Consul服务是否正常运行，网络连接是否正常
4. **API网关访问失败**：检查Traefik配置，确保路由规则正确

## 扩展建议

1. 添加日志服务（如ELK）
2. 添加监控服务（如Prometheus + Grafana）
3. 添加消息队列（如Kafka、RabbitMQ）
4. 实现服务熔断和限流机制
5. 添加分布式缓存（如Redis）