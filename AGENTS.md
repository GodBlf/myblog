# AGENTS.md

本文件用于指导在本仓库内工作的代码代理（Codex/Claude/Gemini 等）。

## 1) 项目概览

- 项目类型：Go Web 博客系统（Gin + GORM + MySQL + Redis + Zap + Viper）。
- 入口文件：`main.go`
- Go 模块：`myblog`
- 默认服务地址：`localhost:5678`
- 页面模板：`views/*.html`
- 静态资源：`views/js`、`views/css`、`views/img`

## 2) 目录职责

- `handler/`：HTTP 路由处理逻辑。
  - `handler/middleware/`：鉴权与指标中间件。
  - `handler/test/`：处理层测试。
- `database/`：数据模型与数据库访问（MySQL/Redis）。
  - `database/test/`：数据库相关测试。
- `util/`：通用工具（配置、日志、JWT、ORM 工具、字符串工具）。
  - `util/test/`：工具层测试。
- `config/`：YAML 配置文件（`mysql.yaml`、`redis.yaml`、`log.yaml`、`key.yaml`）。
- `global/`：全局路径与基础变量。
- `views/`：HTML 模板与前端静态资源。
- `doc/`：文档目录（若新增文档，优先放此处）。
- `log/`：运行日志输出目录。

补充（评论功能相关）：

- `database/comment.go`：公开评论的数据模型、查询、创建、删除与权限约束。
- `handler/comment.go`：公开评论 HTTP 接口（列表、创建、删除）。

## 3) 开发与运行

- 安装依赖：`go mod tidy`
- 启动服务：`go run .`
- 构建：`go build ./...`
- 运行测试：`go test ./...`

注意：当前测试中存在集成型用例（依赖 MySQL/Redis 或实际端口监听），若环境不完整可能失败。修改代码后优先运行与改动最相关的子包测试。

## 4) 代码风格与约定

- 保持现有代码风格，避免大规模无关重构。
- 新增处理器函数遵循现有命名风格（如 `NewXxx()` 返回 `gin.HandlerFunc`）。
- 错误处理优先明确返回 HTTP 状态码与文本信息。
- 日志使用 `zap`（`zap.L().Info/Error/Debug`），避免混用 `fmt.Println`。
- 配置读取保持通过 `util.CreateConfig(...)` 进行，不硬编码配置值。
- 数据库访问逻辑放在 `database/`，`handler/` 只做参数校验与流程编排。

## 5) 修改边界

- 不要在未被要求时修改：
  - `config/*.yaml` 中的环境配置与密钥；
  - 历史日志文件（`log/*.log*`）；
  - 与当前任务无关的路由、模板、测试。
- 不要引入新的大型依赖，除非任务明确要求且必要。
- 不要为了“看起来更整洁”而批量重命名公共 API。

## 6) 路由/页面相关变更清单

当任务涉及新页面或新接口时，按以下顺序检查：

1. 在 `handler/` 增加或修改对应 handler。
2. 在 `main.go` 注册路由（含鉴权中间件是否需要）。
3. 若返回 HTML，确认 `router.LoadHTMLFiles(...)` 已包含目标模板。
4. 若涉及静态资源，确认放置路径与 `router.Static(...)` 映射一致。
5. 补充最小必要测试（优先同目录 `test/`）。

评论相关路由约定（当前实现）：

- `GET /blog/public/:bid/comments`：获取公开文章评论列表。
- `POST /blog/public/:bid/comments`：创建评论（需 `auth_token`）。
- `DELETE /blog/public/:bid/comments/:cid`：删除评论（需 `auth_token`，仅作者本人可删）。

## 7) 数据层变更清单

- 任何与 Blog/User/Token 持久化相关改动，优先落在 `database/`。
- 保持函数语义单一：查询、创建、更新、发布/取消发布分离。
- 如涉及权限判断，业务入口在 `handler/`，数据一致性在 `database/`。

评论权限补充约束：

- 删除评论时，必须在 `database/` 做“评论归属用户”校验，不能仅依赖前端按钮控制。

## 8) 提交质量要求（代理执行）

- 改动应“最小且完整”：只改解决问题所需文件。
- 优先修复根因，避免表层补丁。
- 若无法在当前环境验证（例如缺 MySQL/Redis），在结果中明确说明未验证项与建议验证命令。
- 如新增行为，补充必要文档说明（可放 `doc/` 或相关代码邻近处）。
