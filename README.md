# Go Clean Architecture

## 项目简介
基于 Go 语言的整洁架构 REST API 示例，使用 Gin 框架 + MongoDB。实现了用户认证（注册/登录/JWT）和书签管理（CRUD）功能。按照 Clean Architecture 分层：delivery（HTTP handler）、usecase（业务逻辑）、repository（数据访问）。

## 快速启动

### Docker 启动（推荐）

```bash
# 克隆项目
git clone <GitHub 地址>
cd solo-zj-00056-20260414

# 启动所有服务
docker compose up -d

# 查看运行状态
docker compose ps
```

### 访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| 后端 API | http://localhost:8000 | Go/Gin API |
| MongoDB | localhost:27017 | 数据库 |

### 停止服务

```bash
docker compose down
```

## 项目结构
- `auth/` - 认证模块（delivery/usecase/repository）
- `bookmark/` - 书签模块（delivery/usecase/repository）
- `cmd/api/` - 应用入口
- `config/` - 配置
- `models/` - 数据模型
- `server/` - 服务器初始化

## 来源
- 原始来源: https://github.com/zhashkevych/go-clean-architecture
- GitHub（上传）: https://github.com/11DingKing/solo-zj-00056-20260414
