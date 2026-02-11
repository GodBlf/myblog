# MyBlog 项目文档

本文档基于当前仓库代码编写，覆盖项目结构、运行方式、核心接口与常见排查。

## 1. 项目简介

MyBlog 是一个基于 Go 的简易博客系统，后端采用：

- Gin：HTTP 路由与中间件
- GORM + MySQL：数据存储
- Redis：双 token（refresh token -> auth token）映射
- Zap：日志系统（文件 + 控制台）
- Viper：配置读取
- Prometheus：基础 HTTP 指标采集

默认访问地址：`http://localhost:5678`

## 2. 目录结构

```text
myblog/
├─ main.go                 # 程序入口、路由注册
├─ handler/                # HTTP 处理层
│  ├─ middleware/          # 鉴权与指标中间件
│  └─ test/                # handler 测试
├─ database/               # 数据访问层（MySQL/Redis）
│  └─ test/                # database 测试
├─ util/                   # 工具层（配置、日志、JWT 等）
│  └─ test/                # util 测试
├─ config/                 # 配置文件（mysql/redis/log/key）
├─ global/                 # 全局路径等基础变量
├─ views/                  # HTML 模板与静态资源
│  ├─ js/
│  ├─ css/
│  └─ img/
├─ docker/                 # Docker 运行时配置
├─ doc/                    # 文档目录
└─ log/                    # 日志输出目录
```

## 3. 环境要求

- Go：建议 `1.24+`（以 `Dockerfile` 为准）
- MySQL：`8.x`
- Redis：`7.x`

本地默认配置见：

- `config/mysql.yaml`
- `config/redis.yaml`
- `config/log.yaml`

## 4. 本地启动

### 4.1 安装依赖

```bash
go mod tidy
```

### 4.2 启动 MySQL / Redis

确保数据库与缓存可用，并与配置文件一致：

- MySQL：`localhost:3308`，库名 `myblog`
- Redis：`localhost:6379`

### 4.3 运行服务

```bash
go run .
```

服务启动后可访问：

- 首页：`http://localhost:5678/`
- 登录页：`http://localhost:5678/login`
- Prometheus 指标：`http://localhost:5678/metrics`

## 5. Docker 启动

项目内置 `docker-compose.yml`，可直接启动 MySQL、Redis、App。

```bash
docker compose up --build
```

说明：

- Compose 会把 `docker/config/mysql.yaml`、`docker/config/redis.yaml` 挂载到容器内 `config/`
- App 暴露端口：`5678`

## 6. 核心数据模型

程序启动后会自动迁移（`AutoMigrate`）以下表：

- `user`
  - `id`（主键）
  - `name`
  - `password`
- `blog`
  - `id`（主键）
  - `user_id`
  - `title`
  - `article`
  - `update_time`
- `public_blog`
  - `blog_id`（主键）
  - `user_id`
  - `publish_time`
- `blog_comment`
  - `id`（主键）
  - `blog_id`
  - `user_id`
  - `content`
  - `create_time`

## 7. 认证与鉴权

### 7.1 登录态机制

登录成功后：

1. 后端生成 `auth token`（JWT）
2. 随机生成 `refresh_token`，写入 Redis（7 天过期）
3. 将 `refresh_token` 通过 Cookie 返回给浏览器
4. 前端将 `auth token` 存入 `sessionStorage`

前端 token 续期流程：

- 若 `sessionStorage` 中无 `auth_token`，前端读取 Cookie 中 `refresh_token`
- 调用 `POST /token` 换取新的 `auth_token`

### 7.2 鉴权头

受保护接口要求请求头携带：

```text
auth_token: <JWT>
```

## 8. 路由与接口

### 8.1 页面与公开接口

- `GET /`：首页（`home.html`）
- `GET /login`：登录页（`login.html`）
- `POST /login/submit`：登录
- `POST /register/submit`：注册
- `POST /token`：通过 `refresh_token` 获取 `auth_token`
- `GET /blog/public`：公开博客列表页
- `GET /blog/public/:bid`：公开博客详情页
- `GET /blog/public/:bid/comments`：公开博客评论列表（JSON）
- `GET /blog/list/:uid`：指定用户博客列表页
- `GET /blog/:bid`：博客详情页
- `GET /blog/belong?bid=<id>&token=<jwt>`：判断博客是否归属当前用户
- `GET /metrics`：Prometheus 指标

### 8.2 受保护接口（需 `auth_token`）

- `POST /blog/create`
  - 参数：`title`、`article`
  - 返回：`{"bid": <new_id>}`
- `POST /blog/update`
  - 参数：`bid`、`title`、`article`
- `POST /blog/publish`
  - 参数：`bid`
- `POST /blog/unpublish`
  - 参数：`bid`
- `POST /blog/public/:bid/comments`
  - 参数：`content`
- `DELETE /blog/public/:bid/comments/:cid`
  - 说明：仅评论作者可删除自己的评论

## 8.3 公共评论功能说明

- 公共文章详情页（`/blog/public/:bid`）支持评论展示与发布。
- 评论展示包含：用户名、评论时间、评论内容。
- 前端每条评论都显示“删除”按钮；删除请求会携带 `auth_token`。
- 后端会做权限校验：只有评论作者本人可以删除，非本人删除返回 `403`。

> 前端在登录与注册时会将明文密码做 MD5 后再提交，因此后端接口目前要求 `pass` 长度为 32。

## 9. 配置说明

### 9.1 MySQL（`config/mysql.yaml`）

```yaml
myblog:
  host: localhost
  port: 3308
  user: root
  password: asd456
```

### 9.2 Redis（`config/redis.yaml`）

```yaml
addr: localhost:6379
password: asd456
db: 0
```

### 9.3 日志（`config/log.yaml`）

```yaml
level: debug
file: log/blog.log
```

日志按小时滚动，保留 7 天。

## 10. 开发与测试

### 10.1 构建

```bash
go build ./...
```

### 10.2 测试

```bash
go test ./...
```

注意：当前仓库中包含依赖 MySQL/Redis 或端口监听的集成型测试。建议优先运行与改动相关的子包测试。

## 11. 常见问题排查

### Q1：启动时报数据库连接失败

- 检查 `config/mysql.yaml` 与实际 MySQL 实例是否一致
- 确认数据库 `myblog` 已创建，端口映射正确（默认 `3308`）

### Q2：登录成功后调用写接口仍提示 `auth failed`

- 确认请求头是否带了 `auth_token`
- 确认 `auth_token` 未过期，或可通过 `POST /token` 重新获取

### Q3：公共博客列表为空

- 仅发布后的文章会出现在 `/blog/public`
- 先调用 `/blog/publish`，再刷新公共列表

## 12. 后续改进建议

- 将 JWT 密钥从代码常量迁移到配置文件/环境变量
- 登录密码流程升级为服务端加盐哈希（如 bcrypt/argon2）
- 为关键接口补充更系统的单元测试与集成测试
