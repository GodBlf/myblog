# MyBlog API 文档

本文档描述当前项目后端接口（基于 `main.go` 与 `handler/*` 实现）。

## 1. 基础信息

- Base URL：`http://localhost:5678`
- 数据格式：
  - 表单接口：`application/x-www-form-urlencoded` 或 `multipart/form-data`
  - 返回多为 JSON 或纯文本（以实际接口为准）
- 鉴权方式：请求头 `auth_token: <JWT>`

## 2. 认证相关

### 2.1 登录

- 方法：`POST`
- 路径：`/login/submit`
- 参数（表单）：
  - `user`：用户名
  - `pass`：32 位 MD5 字符串

成功响应（200）：

```json
{
  "code": 0,
  "msg": "success",
  "uid": 1,
  "token": "<jwt>"
}
```

失败响应示例：

- 400：`{"code":1,"msg":"must indicate user name"...}`
- 400：`{"code":2,"msg":"invalid password"...}`
- 403：`{"code":3,"msg":"user not exist"...}`
- 403：`{"code":4,"msg":"incorrect password"...}`
- 500：`{"code":5,"msg":"generate jwtToken failed"...}`

### 2.2 注册

- 方法：`POST`
- 路径：`/register/submit`
- 参数（表单）：
  - `user`：用户名
  - `pass`：32 位 MD5 字符串

成功响应（200）：

```json
{
  "code": 0,
  "msg": "success"
}
```

失败响应示例：

- 400：`{"code":1,"msg":"must indicate user name"}`
- 400：`{"code":2,"msg":"invalid password"}`
- 409：`{"code":3,"msg":"user already exist"}`
- 500：`{"code":4,"msg":"create user failed"}`

### 2.3 刷新 token

- 方法：`POST`
- 路径：`/token`
- 参数（表单）：
  - `refresh_token`：登录时服务端写入 Cookie 的值

成功响应（200）：

```text
<jwt字符串>
```

## 3. 博客查询（公开）

### 3.1 获取某用户博客列表页

- 方法：`GET`
- 路径：`/blog/list/:uid`
- 说明：返回 HTML 页面 `blog_list.html`

### 3.2 获取博客详情页

- 方法：`GET`
- 路径：`/blog/:bid`
- 说明：返回 HTML 页面 `blog.html`

错误：

- 400：`invalid blog id`
- 404：`blog not exist`

### 3.3 获取公开博客列表页

- 方法：`GET`
- 路径：`/blog/public`
- 说明：返回 HTML 页面 `public_blog_list.html`

### 3.4 获取公开博客详情页

- 方法：`GET`
- 路径：`/blog/public/:bid`
- 说明：返回 HTML 页面 `blog_public.html`

错误：

- 400：`invalid blog id`
- 404：`public blog not exist`

### 3.5 判断博客归属

- 方法：`GET`
- 路径：`/blog/belong`
- Query 参数：
  - `bid`：博客 ID
  - `token`：JWT

成功响应（200）：

- `true`：当前 token 用户是博客作者
- `false`：不是作者

错误：

- 400：`invalid blog id`
- 400：`blog id not exists`

### 3.6 获取公开博客评论列表

- 方法：`GET`
- 路径：`/blog/public/:bid/comments`
- 说明：返回公开博客评论 JSON 数组，按评论时间倒序（新到旧）

成功响应（200）示例：

```json
[
  {
    "id": 12,
    "user_name": "alice",
    "content": "写得很好",
    "create_time": "2026-02-11 22:10:00"
  }
]
```

错误：

- 400：`invalid blog id`
- 404：`public blog not exist`

## 4. 博客写操作（需鉴权）

> 以下接口都需要请求头：`auth_token: <JWT>`。

### 4.1 新建博客

- 方法：`POST`
- 路径：`/blog/create`
- 参数（表单）：
  - `title`
  - `article`

成功响应（200）：

```json
{
  "bid": 123
}
```

失败响应：

- 400：`invalid parameter`
- 403：`auth failed`
- 500：`create blog failed`

### 4.2 更新博客

- 方法：`POST`
- 路径：`/blog/update`
- 参数（表单）：
  - `bid`（>0）
  - `title`
  - `article`

成功响应（200）：

```text
update blog success
```

失败响应：

- 400：`invalid parameter`
- 400：`blog not exist`
- 403：`auth failed`
- 403：`no permission to update`
- 500：`update blog failed`

### 4.3 发布博客

- 方法：`POST`
- 路径：`/blog/publish`
- 参数（表单）：
  - `bid`（>0）

成功响应（200）：

```text
publish blog success
```

失败响应：

- 400：`invalid parameter`
- 400：`blog not exist`
- 403：`auth failed`
- 403：`no permission to publish`
- 500：`publish blog failed`

### 4.4 取消发布博客

- 方法：`POST`
- 路径：`/blog/unpublish`
- 参数（表单）：
  - `bid`（>0）

成功响应（200）：

```text
unpublish blog success
```

失败响应：

- 400：`invalid parameter`
- 400：`blog not exist`
- 403：`auth failed`
- 403：`no permission to unpublish`
- 500：`unpublish blog failed`

### 4.5 发表公开博客评论

- 方法：`POST`
- 路径：`/blog/public/:bid/comments`
- 参数（表单或 JSON）：
  - `content`（必填，去空白后 1~1000 字）

成功响应（200）示例：

```json
{
  "id": 13,
  "user_name": "bob",
  "content": "支持一下",
  "create_time": "2026-02-11 22:15:00"
}
```

失败响应：

- 400：`invalid blog id`
- 400：`invalid parameter`
- 403：`auth failed`
- 404：`public blog not exist`
- 500：`create comment failed`

### 4.6 删除公开博客评论（仅评论作者）

- 方法：`DELETE`
- 路径：`/blog/public/:bid/comments/:cid`
- 说明：仅评论创建者可删除自己的评论

成功响应（200）：

```text
delete comment success
```

失败响应：

- 400：`invalid blog id`
- 400：`invalid comment id`
- 400：`invalid parameter`
- 403：`auth failed`
- 403：`no permission to delete comment`
- 404：`public blog not exist`
- 404：`comment not exist`
- 500：`delete comment failed`

## 5. 系统接口

### 5.1 Prometheus 指标

- 方法：`GET`
- 路径：`/metrics`
- 说明：返回 Prometheus 文本格式指标

## 6. cURL 示例

### 6.1 注册

```bash
curl -X POST "http://localhost:5678/register/submit" \
  -d "user=test_user" \
  -d "pass=25d55ad283aa400af464c76d713c07ad"
```

### 6.2 登录

```bash
curl -i -X POST "http://localhost:5678/login/submit" \
  -d "user=test_user" \
  -d "pass=25d55ad283aa400af464c76d713c07ad"
```

说明：

- 响应 JSON 内 `token` 即 `auth_token`
- 响应头 `Set-Cookie` 中包含 `refresh_token`

### 6.3 创建博客

```bash
curl -X POST "http://localhost:5678/blog/create" \
  -H "auth_token: <your_jwt>" \
  -d "title=Hello" \
  -d "article=My first post"
```

### 6.4 发布博客

```bash
curl -X POST "http://localhost:5678/blog/publish" \
  -H "auth_token: <your_jwt>" \
  -d "bid=1"
```

### 6.5 通过 refresh token 获取 auth token

```bash
curl -X POST "http://localhost:5678/token" \
  -d "refresh_token=<your_refresh_token>"
```

### 6.6 发表评论

```bash
curl -X POST "http://localhost:5678/blog/public/1/comments" \
  -H "auth_token: <your_jwt>" \
  -d "content=这篇文章很赞"
```

### 6.7 删除自己的评论

```bash
curl -X DELETE "http://localhost:5678/blog/public/1/comments/13" \
  -H "auth_token: <your_jwt>"
```
