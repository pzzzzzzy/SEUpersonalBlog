# 个人博客系统

一个功能完整的个人博客系统，包含用户注册登录、文章管理、评论互动等功能。

## 功能特性

- **用户系统**: 注册、登录、个人主页
- **文章管理**: 发布、编辑、删除、草稿保存
- **内容分类**: 自定义专栏、标签系统
- **评论互动**: 一级评论、回复、点赞
- **数据统计**: 阅读量、评论数统计
- **响应式设计**: 适配移动端和桌面端

## 技术栈

### 后端
- Go 1.21+
- Gin Web框架
- GORM ORM
- SQLite数据库
- JWT认证

### 前端
- HTML5 + CSS3 + JavaScript
- 原生DOM操作
- 响应式布局

## 快速开始

### 1. 克隆项目
```bash
cd personal-blog
```

### 2. 安装依赖
```bash
cd backend
go mod download
```

### 3. 启动服务
```bash
go run cmd/main.go
```

### 4. 访问应用
打开浏览器访问: http://localhost:8080

## API接口

### 认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

### 文章相关
- `GET /api/articles` - 获取文章列表
- `GET /api/articles/:id` - 获取文章详情
- `POST /api/articles` - 创建文章（需要认证）
- `PUT /api/articles/:id` - 更新文章（需要认证）
- `DELETE /api/articles/:id` - 删除文章（需要认证）

### 评论相关
- `GET /api/comments/article/:articleId` - 获取文章评论
- `POST /api/comments` - 发表评论（需要认证）
- `POST /api/comments/:id/like` - 点赞评论

## 项目结构

```
personal-blog/
├── backend/                 # Go后端代码
│   ├── cmd/                # 主程序入口
│   ├── config/             # 配置文件
│   ├── db/                 # 数据库初始化
│   ├── handlers/           # HTTP处理器
│   ├── middleware/         # 中间件
│   ├── models/             # 数据模型
│   ├── routes/             # 路由配置
│   └── utils/              # 工具函数
├── frontend/                # 前端代码
│   ├── static/             # 静态资源
│   │   ├── css/            # 样式文件
│   │   ├── js/             # JavaScript文件
│   │   └── images/         # 图片资源
│   └── templates/          # HTML模板
└── docs/                    # 文档
```

## 开发计划

- [x] 基础项目结构
- [x] 用户注册登录
- [x] 文章CRUD操作
- [x] 评论系统
- [x] 前端基础页面
- [ ] 富文本编辑器
- [ ] 文件上传功能
- [ ] 第三方登录
- [ ] 消息通知
- [ ] 文章搜索
- [ ] 数据可视化

## 贡献指南

欢迎提交Issue和Pull Request来改进这个项目。

## 许可证

MIT License