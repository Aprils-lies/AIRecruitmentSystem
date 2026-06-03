# 智能招聘系统 — 接口说明文档

---

## 一、接口概述

### 1.1 基础信息

| 项目             | 说明                                                |
| ---------------- | --------------------------------------------------- |
| 基础 URL（开发） | Vite 代理 `/api` → `http://localhost:8080`     |
| 基础 URL（生产） | nginx 反向代理 `/api` → `web-gin-service:8080` |
| 数据格式         | JSON                                                |
| 字符编码         | UTF-8                                               |

### 1.2 统一响应格式

所有接口返回统一结构：

```json
{
    "code": 0,
    "message": "操作成功",
    "data": { ... }
}
```

### 1.3 错误码对照

| code | HTTP 状态码 | 含义             |
| ---- | ----------- | ---------------- |
| 0    | 200         | 成功             |
| 400  | 200         | 参数错误         |
| 401  | 401         | 未授权，请先登录 |
| 403  | 403         | 无权访问         |
| 404  | 200         | 资源不存在       |
| 500  | 200         | 服务器内部错误   |

### 1.4 认证方式

所有需要登录的接口在 HTTP Header 中携带 JWT Token：

```
Authorization: Bearer <token>
```

Token 通过登录接口获取，前端在 Axios 拦截器中自动注入。Token 过期由后端验证（通过 gRPC `AuthService.VerifyToken`），过期后会返回 `401`，前端会清除 Token 并跳转登录页。

### 1.5 分页参数

所有分页列表接口统一使用以下 query 参数：

| 参数      | 类型    | 默认值 | 说明     |
| --------- | ------- | ------ | -------- |
| page      | integer | 1      | 页码     |
| page_size | integer | 10     | 每页条数 |

响应中统一包含：

```json
{
    "total": 100,
    "page": 1
}
```

---

## 二、公开接口（无需登录）

### 2.1 用户注册

```
POST /api/auth/register
```

**请求体：**

```json
{
    "username": "hr1",
    "password": "123456",
    "role": "hr"
}
```

| 参数     | 类型   | 必填 | 说明                          |
| -------- | ------ | ---- | ----------------------------- |
| username | string | 是   | 登录账号，唯一                |
| password | string | 是   | 登录密码                      |
| role     | string | 是   | 角色：`hr` 或 `candidate` |

**成功响应：**

```json
{
    "code": 0,
    "message": "注册成功",
    "data": {
        "user_id": 1
    }
}
```

### 2.2 用户登录

```
POST /api/auth/login
```

**请求体：**

```json
{
    "username": "hr1",
    "password": "123456"
}
```

**成功响应：**

```json
{
    "code": 0,
    "message": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIs...",
        "user_id": 1,
        "role": "hr",
        "username": "hr1"
    }
}
```

### 2.3 公开岗位列表

```
GET /api/positions?page=1&page_size=10&keyword=&location=
```

**Query 参数：**

| 参数      | 类型    | 必填 | 说明                       |
| --------- | ------- | ---- | -------------------------- |
| page      | integer | 否   | 默认 1                     |
| page_size | integer | 否   | 默认 10                    |
| keyword   | string  | 否   | 搜索关键词，匹配标题和描述 |
| location  | string  | 否   | 工作地点筛选               |

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "positions": [
            {
                "id": 1,
                "hr_id": 1,
                "title": "高级前端工程师",
                "description": "负责前端架构设计...",
                "requirements": "3年以上React经验...",
                "salary_min": 20,
                "salary_max": 40,
                "location": "杭州",
                "status": "published",
                "created_at": "2026-05-01T10:00:00Z",
                "updated_at": "2026-05-01T10:00:00Z"
            }
        ],
        "total": 100,
        "page": 1
    }
}
```

### 2.4 公开岗位详情

```
GET /api/positions/:id
```

**路径参数：**

| 参数 | 类型    | 说明    |
| ---- | ------- | ------- |
| id   | integer | 岗位 ID |

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "id": 1,
        "hr_id": 1,
        "title": "高级前端工程师",
        "description": "负责前端架构设计...",
        "requirements": "3年以上React经验...",
        "salary_min": 20,
        "salary_max": 40,
        "location": "杭州",
        "status": "published",
        "created_at": "2026-05-01T10:00:00Z",
        "updated_at": "2026-05-01T10:00:00Z"
    }
}
```

---

## 三、HR 端接口（需 JWT + hr 角色）

### 3.1 创建岗位

```
POST /api/hr/positions
```

**请求体：**

```json
{
    "title": "高级前端工程师",
    "description": "负责前端架构设计及核心模块开发",
    "requirements": "3年以上React经验",
    "salary_min": 20,
    "salary_max": 40,
    "location": "杭州"
}
```

| 参数         | 类型    | 必填 | 说明                |
| ------------ | ------- | ---- | ------------------- |
| title        | string  | 是   | 岗位名称            |
| description  | string  | 是   | 岗位职责描述        |
| requirements | string  | 否   | 任职要求            |
| salary_min   | integer | 否   | 最低薪资（单位：K） |
| salary_max   | integer | 否   | 最高薪资（单位：K） |
| location     | string  | 否   | 工作地点            |

**成功响应：**

```json
{
    "code": 0,
    "message": "创建成功",
    "data": {
        "position_id": 1
    }
}
```

### 3.2 编辑岗位

```
PUT /api/hr/positions/:id
```

**路径参数：**

| 参数 | 类型    | 说明            |
| ---- | ------- | --------------- |
| id   | integer | 要编辑的岗位 ID |

**请求体（所有字段可选，仅更新非空字段）：**

```json
{
    "title": "高级前端工程师（急聘）",
    "description": "更新后的描述",
    "requirements": "更新后的要求",
    "salary_min": 25,
    "salary_max": 45,
    "location": "上海"
}
```

**成功响应：**

```json
{
    "code": 0,
    "message": "更新成功",
    "data": null
}
```

> 仅能编辑本人创建的岗位，操作他人岗位返回 `403`。

### 3.3 下架岗位

```
POST /api/hr/positions/:id/offline
```

**路径参数：**

| 参数 | 类型    | 说明    |
| ---- | ------- | ------- |
| id   | integer | 岗位 ID |

**成功响应：**

```json
{
    "code": 0,
    "message": "下架成功",
    "data": null
}
```

### 3.4 上架岗位

```
POST /api/hr/positions/:id/online
```

**路径参数：**

| 参数 | 类型    | 说明    |
| ---- | ------- | ------- |
| id   | integer | 岗位 ID |

**成功响应：**

```json
{
    "code": 0,
    "message": "上架成功",
    "data": null
}
```

### 3.5 我的岗位列表

```
GET /api/hr/my-positions?page=1&page_size=10
```

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "positions": [
            {
                "id": 1,
                "hr_id": 1,
                "title": "高级前端工程师",
                "description": "...",
                "requirements": "...",
                "salary_min": 20,
                "salary_max": 40,
                "location": "杭州",
                "status": "published",
                "created_at": "2026-05-01T10:00:00Z",
                "updated_at": "2026-05-01T10:00:00Z"
            }
        ],
        "total": 5,
        "page": 1
    }
}
```

> 仅返回当前 HR 自己发布的岗位列表。

### 3.6 候选人列表

```
GET /api/hr/positions/:id/candidates?page=1&page_size=10
```

**路径参数：**

| 参数 | 类型    | 说明    |
| ---- | ------- | ------- |
| id   | integer | 岗位 ID |

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "candidates": [
            {
                "user_id": 2,
                "username": "candidate1",
                "real_name": "张三",
                "phone": "13888888000",
                "education": "本科",
                "school": "浙江大学",
                "skills": "React,TypeScript",
                "resume_id": 1,
                "resume_name": "张三_前端工程师_简历.pdf",
                "applied_at": "2026-05-16T09:00:00Z",
                "status": "pending"
            }
        ],
        "total": 5,
        "page": 1
    }
}
```

> `status` 枚举值：`pending`（待审核）、`reviewed`（已查看）、`rejected`（已拒绝）、`accepted`（已通过）

### 3.7 候选人详情

```
GET /api/hr/candidates/:id?application_id=1
```

**路径参数：**

| 参数 | 类型    | 说明          |
| ---- | ------- | ------------- |
| id   | integer | 候选人用户 ID |

**Query 参数：**

| 参数           | 类型    | 必填 | 说明                                                |
| -------------- | ------- | ---- | --------------------------------------------------- |
| application_id | integer | 是   | 投递记录 ID，用于校验候选人确实投递了当前 HR 的岗位 |

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "user_id": 2,
        "username": "candidate1",
        "real_name": "张三",
        "phone": "13800138000",
        "education": "本科",
        "school": "浙江大学",
        "experience": "3年前端开发经验...",
        "skills": "React,TypeScript",
        "resume_id": 1,
        "resume_name": "张三_前端工程师_简历.pdf",
        "resume_type": "pdf",
        "applied_at": "2026-05-16T09:00:00Z"
    }
}
```

### 3.8 更新投递状态

```
PUT /api/hr/applications/:id/status
```

**路径参数：**

| 参数 | 类型    | 说明        |
| ---- | ------- | ----------- |
| id   | integer | 投递记录 ID |

**请求体：**

```json
{
    "new_status": "reviewed"
}
```

| 参数       | 类型   | 必填 | 说明                                                     |
| ---------- | ------ | ---- | -------------------------------------------------------- |
| new_status | string | 是   | `pending` / `reviewed` / `rejected` / `accepted` |

**成功响应：**

```json
{
    "code": 0,
    "message": "更新成功",
    "data": null
}
```

### 3.9 获取简历下载签名 URL（HR）

```
GET /api/hr/resumes/:id/download-url
```

**路径参数：**

| 参数 | 类型    | 说明    |
| ---- | ------- | ------- |
| id   | integer | 简历 ID |

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "download_url": "https://xxx.oss-cn-beijing.aliyuncs.com/resumes/1/xxx.pdf?Expires=xxx&OSSAccessKeyId=xxx&Signature=xxx",
        "file_name": "张三_前端工程师_简历.pdf",
        "expire_sec": 3600
    }
}
```

> HR 只能下载投递了本人岗位的候选人的简历，否则返回 `403`。
>
> 前端获取 URL 后创建 `<a>` 元素并触发点击即可下载，不经过服务器代理。

### 3.10 AI 对话（非流式）

```
POST /api/hr/ai/chat
```

**请求体：**

```json
{
    "question": "我的岗位投递情况如何？",
    "session_id": ""
}
```

| 参数       | 类型   | 必填 | 说明                          |
| ---------- | ------ | ---- | ----------------------------- |
| question   | string | 是   | 用户提问内容                  |
| session_id | string | 否   | 会话 ID，为空时自动创建新会话 |

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "answer": "您共有5个在招岗位，总计收到23份投递，其中待审核8份、已查看5份、已拒绝6份、已通过4份。",
        "message_id": 2,
        "session_id": "1747828800000000000"
    }
}
```

### 3.11 AI 对话（流式 SSE）

```
POST /api/hr/ai/chat/stream
```

**请求体**同 3.10。

**响应格式：** `text/event-stream`（Server-Sent Events）

```
data: {"chunk":"您共","done":false,"message_id":0,"session_id":"1747828800000000000"}

data: {"chunk":"有5个","done":false,"message_id":0,"session_id":"1747828800000000000"}

data: {"chunk":"","done":true,"message_id":2,"session_id":"1747828800000000000"}
```

| 字段       | 类型    | 说明                                  |
| ---------- | ------- | ------------------------------------- |
| chunk      | string  | 本次推送的文本片段                    |
| done       | boolean | 是否已结束                            |
| message_id | number  | 流结束前为 0，结束后为入库后的真实 ID |
| session_id | string  | 当前会话 ID                           |

> 前端使用原生 `fetch` + `ReadableStream` 读取，逐片追加到页面，实现打字机效果。

### 3.12 获取对话历史

```
GET /api/hr/ai/history?session_id=xxx&limit=20&offset=0
```

**Query 参数：**

| 参数       | 类型    | 必填 | 说明              |
| ---------- | ------- | ---- | ----------------- |
| session_id | string  | 是   | 会话 ID           |
| limit      | integer | 否   | 返回条数，默认 20 |
| offset     | integer | 否   | 偏移量，默认 0    |

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "items": [
            {
                "id": 1,
                "session_id": "1747828800000000000",
                "role": "user",
                "content": "我的岗位投递情况如何？",
                "created_at": "2026-05-16T10:00:00Z"
            },
            {
                "id": 2,
                "session_id": "1747828800000000000",
                "role": "assistant",
                "content": "您共有5个在招岗位...",
                "created_at": "2026-05-16T10:00:05Z"
            }
        ],
        "total": 2
    }
}
```

### 3.13 获取会话列表

```
GET /api/hr/ai/sessions
```

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "sessions": [
            {
                "session_id": "1747828800000000000",
                "last_message": "您共有5个在招岗位...",
                "created_at": "2026-05-16T10:00:00Z"
            }
        ]
    }
}
```

### 3.14 统计摘要

```
GET /api/hr/ai/stats
```

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "summary": "您共有5个在招岗位，总计收到23份投递。待审核8份，已查看5份，已拒绝6份，已通过4份。"
    }
}
```

### 3.15 删除会话

```
DELETE /api/hr/ai/sessions/:session_id
```

**成功响应：**

```json
{
    "code": 0,
    "message": "删除成功",
    "data": null
}
```

### 3.16 删除消息

```
DELETE /api/hr/ai/messages/:message_id
```

**成功响应：**

```json
{
    "code": 0,
    "message": "删除成功",
    "data": null
}
```

---

## 四、候选人端接口（需 JWT + candidate 角色）

### 4.1 获取个人资料

```
GET /api/candidate/profile
```

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "user_id": 2,
        "username": "candidate1",
        "role": "candidate",
        "real_name": "张三",
        "phone": "13800138000",
        "education": "本科",
        "school": "浙江大学",
        "experience": "3年前端开发经验...",
        "skills": "React,TypeScript"
    }
}
```

### 4.2 更新个人资料

```
PUT /api/candidate/profile
```

**请求体（所有字段可选，仅更新非空字段）：**

```json
{
    "real_name": "张三",
    "phone": "13800138000",
    "education": "本科",
    "school": "浙江大学",
    "experience": "3年前端开发经验...",
    "skills": "React,TypeScript"
}
```

| 参数       | 类型   | 必填 | 说明                                         |
| ---------- | ------ | ---- | -------------------------------------------- |
| real_name  | string | 否   | 真实姓名                                     |
| phone      | string | 否   | 联系电话                                     |
| education  | string | 否   | 最高学历：高中/中专/大专/本科/硕士/博士/其他 |
| school     | string | 否   | 毕业院校                                     |
| experience | string | 否   | 工作/项目经历                                |
| skills     | string | 否   | 核心技能标签，逗号分隔                       |

**成功响应：**

```json
{
    "code": 0,
    "message": "更新成功",
    "data": null
}
```

### 4.3 投递岗位

```
POST /api/candidate/positions/:id/apply
```

**路径参数：**

| 参数 | 类型    | 说明    |
| ---- | ------- | ------- |
| id   | integer | 岗位 ID |

**请求体：**

```json
{
    "resume_id": 1
}
```

| 参数      | 类型    | 必填 | 说明    |
| --------- | ------- | ---- | ------- |
| resume_id | integer | 是   | 简历 ID |

**成功响应：**

```json
{
    "code": 0,
    "message": "投递成功",
    "data": {
        "application_id": 1
    }
}
```

> 后端按顺序校验：JWT 有效性 → 角色是否为 candidate → 个人资料完整性（6 个字段非空）→ 简历是否存在 → 是否重复投递。任一不通过则返回错误消息。

### 4.4 我的投递记录

```
GET /api/candidate/applications?page=1&page_size=10
```

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "applications": [
            {
                "application_id": 1,
                "position_id": 1,
                "position_title": "高级前端工程师",
                "status": "pending",
                "applied_at": "2026-05-16T09:00:00Z"
            }
        ],
        "total": 5,
        "page": 1
    }
}
```

### 4.5 撤回投递

```
DELETE /api/candidate/applications/:id
```

**路径参数：**

| 参数 | 类型    | 说明        |
| ---- | ------- | ----------- |
| id   | integer | 投递记录 ID |

> 仅能撤回 `pending`（待审核）状态的投递。已查看 / 已拒绝 / 已通过的投递不可撤回。

**成功响应：**

```json
{
    "code": 0,
    "message": "撤回成功",
    "data": null
}
```

### 4.6 获取简历上传签名 URL

```
GET /api/candidate/resumes/upload-url?file_name=xxx&content_type=xxx
```

**Query 参数：**

| 参数         | 类型   | 必填 | 说明                                   |
| ------------ | ------ | ---- | -------------------------------------- |
| file_name    | string | 是   | 原始文件名（含扩展名）                 |
| content_type | string | 是   | 文件 MIME 类型，如 `application/pdf` |

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "upload_url": "https://xxx.oss-cn-beijing.aliyuncs.com/resumes/2/xxx.pdf?Expires=xxx&OSSAccessKeyId=xxx&Signature=xxx",
        "oss_key": "resumes/2/1747828800_张三_前端工程师_简历.pdf",
        "expire_sec": 3600
    }
}
```

> 上传流程三步走：
>
> **步骤 1** — 前端调用此接口获取 OSS 签名 PUT URL。
>
> **步骤 2** — 前端直接用 `fetch` 或 `XMLHttpRequest` 以 `PUT` 方式将文件直传到 OSS 的 `upload_url`，设置 `Content-Type` 头。**文件不经过后端服务器。**
>
> **步骤 3** — 上传完成后调用 4.7 确认接口通知后端。

### 4.7 确认简历上传

```
POST /api/candidate/resumes/confirm
```

**请求体：**

```json
{
    "oss_key": "resumes/2/1747828800_张三_前端工程师_简历.pdf",
    "file_name": "张三_前端工程师_简历.pdf",
    "file_type": "pdf",
    "file_size": 204800
}
```

| 参数      | 类型    | 必填 | 说明                                   |
| --------- | ------- | ---- | -------------------------------------- |
| oss_key   | string  | 是   | 步骤 1 返回的 OSS 对象 Key             |
| file_name | string  | 是   | 原始文件名                             |
| file_type | string  | 是   | 文件格式：`pdf` / `doc` / `docx` |
| file_size | integer | 是   | 文件大小（字节）                       |

> 后端校验流程：OSS 文件存在性校验 → 扩展名校验 → **文件魔数校验**（读取文件头字节，验证真实格式）→ 写入 resumes 表。

**成功响应：**

```json
{
    "code": 0,
    "message": "上传成功",
    "data": {
        "resume_id": 1
    }
}
```

### 4.8 我的简历列表

```
GET /api/candidate/resumes
```

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "resumes": [
            {
                "id": 1,
                "file_name": "张三_前端工程师_简历.pdf",
                "file_type": "pdf",
                "file_size": 204800,
                "uploaded_at": "2026-05-15T14:30:00Z"
            }
        ]
    }
}
```

### 4.9 获取简历下载签名 URL（候选人）

```
GET /api/candidate/resumes/:id/download-url
```

**路径参数：**

| 参数 | 类型    | 说明    |
| ---- | ------- | ------- |
| id   | integer | 简历 ID |

**成功响应：**

```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "download_url": "https://xxx.oss-cn-beijing.aliyuncs.com/resumes/2/xxx.pdf?Expires=xxx&OSSAccessKeyId=xxx&Signature=xxx",
        "file_name": "张三_前端工程师_简历.pdf",
        "expire_sec": 3600
    }
}
```

> 候选人只能下载自己的简历，无权下载他人简历（返回 `403`）。

### 4.10 删除简历

```
DELETE /api/candidate/resumes/:id
```

**路径参数：**

| 参数 | 类型    | 说明    |
| ---- | ------- | ------- |
| id   | integer | 简历 ID |

> 删除操作会同步删除 OSS 上的文件，以及关联的投递记录。
> 仅能删除自己的简历。

**成功响应：**

```json
{
    "code": 0,
    "message": "删除成功",
    "data": null
}
```

---

## 五、健康检查

```
GET /health
```

**成功响应：**

```json
{
    "status": "ok"
}
```
