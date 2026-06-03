# 智能招聘系统

基于 **gRPC 两层微服务架构**的智能招聘平台，提供 HR 管理端和候选人用户端双前端，集成 OSS 文件直传、Eino AI 智能对话功能。

---

## 技术栈

### 后端

| 服务                         | 框架                      | 端口      | 职责                                                   |
| ---------------------------- | ------------------------- | --------- | ------------------------------------------------------ |
| **web-gin-service**    | Gin + gRPC Client         | `:8080` | HTTP 网关：接收请求、JWT 鉴权、参数校验、gRPC 转发     |
| **logic-grpc-service** | gRPC Server + GORM + Eino | `:8081` | 核心业务：鉴权、岗位 CRUD、投递管理、OSS 签名、AI 对话 |

### 前端

| 工程                    | 技术栈                                     | 端口（dev） | 说明         |
| ----------------------- | ------------------------------------------ | ----------- | ------------ |
| **hr-frontend**   | React 19 + TypeScript + Vite + TailwindCSS | `:3000`   | HR 管理端    |
| **user-frontend** | React 19 + TypeScript + Vite + TailwindCSS | `:3001`   | 候选人用户端 |

### 数据与中间件

| 组件         | 用途                                                   |
| ------------ | ------------------------------------------------------ |
| MySQL 8.0    | 业务数据存储（用户、岗位、投递、简历元数据、对话历史） |
| 阿里云 OSS   | 简历文件私有存储（签名 URL 直传，文件不经过服务器）    |
| DeepSeek API | AI 大模型（通过 Eino 框架调用）                        |

---

## 环境要求

| 工具                   | 版本建议 | 用途                   |
| ---------------------- | -------- | ---------------------- |
| Go                     | 1.21+    | 后端编译运行           |
| Node.js                | 18 LTS+  | 前端开发               |
| MySQL                  | 8.0+     | 数据库                 |
| protoc + protoc-gen-go | 最新     | 如需重新生成 gRPC 代码 |

---

## 快速启动

### 1. 创建数据库

```bash
mysql -u root -p < db/init.sql
```

### 2. 配置密钥

复制 `logic-grpc-service/.env.example` 为 `.env`，填入真实密钥：

```bash
# 数据库
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=recruitment

# JWT 密钥（自行修改）
JWT_SECRET=your-jwt-secret-key-change-me-in-production
JWT_EXPIRE=7200

# 阿里云 OSS（需自行注册）
OSS_ENDPOINT=https://oss-cn-beijing.aliyuncs.com
OSS_ACCESS_KEY_ID=your-access-key-id
OSS_ACCESS_KEY_SECRET=your-access-key-secret
OSS_BUCKET_NAME=your-private-bucket

# DeepSeek AI（需自行注册获取 API Key）
AI_PROVIDER=deepseek
AI_MODEL=deepseek-v4-flash
AI_API_KEY=sk-your-deepseek-api-key
AI_ENDPOINT=https://api.deepseek.com
AI_TEMPERATURE=0.7
```

`web-gin-service/.env` 只需配置：

```bash
SERVER_PORT=:8080
GRPC_ADDR=localhost:8081
JWT_SECRET=your-jwt-secret-key-change-me-in-production
JWT_EXPIRE=7200
```

> `JWT_SECRET` 必须与 logic-grpc-service 一致，否则 JWT 验证会失败。

### 3. 启动 gRPC 业务服务

```bash
cd logic-grpc-service
go run main.go
```

预期输出：

```
正在启动智能招聘系统 gRPC 服务...
配置初始化成功
数据库连接成功
OSS客户端初始化成功
AI客户端初始化成功
正在监听端口 :8081
所有服务已注册，准备启动 gRPC 服务器...
```

### 4. 启动 HTTP 网关

```bash
cd web-gin-service
go run main.go
```

预期输出：

```
正在启动智能招聘系统 HTTP 网关服务...
配置初始化成功
gRPC客户端初始化成功
正在监听端口 :8080
HTTP网关启动成功
```

### 5. 启动前端（需要新终端窗口）

```bash
# HR 管理端
cd hr-frontend
npm install
npm run dev

# 候选人用户端
cd user-frontend
npm install
npm run dev
```

### 6. 访问系统

| 入口         | 地址                             |
| ------------ | -------------------------------- |
| HR 管理端    | `http://localhost:3000`        |
| 候选人用户端 | `http://localhost:3001`        |
| 健康检查     | `http://localhost:8080/health` |

---

## 服务架构

```
┌─────────────────────────────────────────────────────┐
│                    用户浏览器                          │
│  ┌─────────────────┐  ┌─────────────────────┐       │
│  │  hr-frontend    │  │  user-frontend      │       │
│  │  (React :3000)  │  │  (React :3001)      │       │
│  └────────┬────────┘  └──────────┬──────────┘       │
│           │                      │                   │
│           │   /api/* (Vite proxy / nginx)            │
│           ▼                      ▼                   │
│  ┌──────────────────────────────────────┐            │
│  │      web-gin-service (:8080)         │            │
│  │  Gin HTTP 网关（JWT / 参数校验）       │            │
│  └──────────────┬───────────────────────┘            │
│                 │ gRPC                                │
│                 ▼                                     │
│  ┌──────────────────────────────────────┐            │
│  │    logic-grpc-service (:8081)        │            │
│  │  gRPC 核心业务服务                    │            │
│  │  鉴权 | 岗位 | 投递 | OSS签名 | Eino AI│            │
│  └──┬──────────┬──────────┬─────────────┘            │
│     │          │          │                           │
│     ▼          ▼          ▼                           │
│  ┌──────┐ ┌──────────┐ ┌──────────┐                  │
│  │MySQL │ │阿里云OSS │ │DeepSeek │ 外部依赖           │
│  │:3306 │ │(私有Bucket)│ │ 大模型  │                  │
│  └──────┘ └──────────┘ └──────────┘                  │
└─────────────────────────────────────────────────────┘
```

### 架构要点

- **gRPC 两层分离**：Web 服务仅做 HTTP 接入，所有业务逻辑通过 gRPC 调用 Logic 服务，禁止本地函数直连
- **无 CORS 问题**：开发期 Vite proxy 代理 `/api`，生产期 nginx 反向代理，浏览器始终请求同源
- **文件直传 OSS**：简历上传/下载均通过 OSS 签名 URL 直接操作，文件不经过应用服务器

---

## OSS 配置说明

### 1. 注册阿里云 OSS

前往 [阿里云 OSS 控制台](https://oss.console.aliyun.com/) 创建 Bucket。

### 2. 创建私有 Bucket

| 配置项     | 值                                     |
| ---------- | -------------------------------------- |
| 地域       | 任意（如北京 `oss-cn-beijing`）      |
| 存储类型   | 标准存储                               |
| 读写权限   | **私有**（关闭公共读、匿名访问） |
| 服务端加密 | 可选                                   |

### 3. 获取访问密钥

在 RAM 访问控制中创建 AccessKey，授予 OSS 管理权限。

### 4. 写入配置

将 `Endpoint`、`AccessKeyID`、`AccessKeySecret`、`BucketName` 写入 `logic-grpc-service/.env`。

### 5. 上传/下载流程

```
上传三步走：
 ① GET /api/candidate/resumes/upload-url  → 获取得签名 PUT URL
 ② 前端 PUT 文件直传到 OSS（不经服务器）
 ③ POST /api/candidate/resumes/confirm    → 后端确认入库

下载两步走：
 ① GET /api/*/resumes/:id/download-url    → 获取签名 GET URL
 ② 前端用 URL 创建 <a> 标签触发下载
```

---

## AI 对话说明

### 技术选型

| 组件       | 说明                                                           |
| ---------- | -------------------------------------------------------------- |
| 框架       | [Eino](https://github.com/cloudwego/eino)（字节跳动开源 AI 框架） |
| LLM        | DeepSeek（通过 OpenAI 兼容端点接入）                           |
| Agent 模式 | ReAct（思考-行动-观察循环）                                    |

### 内置工具

| 工具                        | 功能         | 查询的数据                             |
| --------------------------- | ------------ | -------------------------------------- |
| `query_application_stats` | 投递总统计   | 按 HR 汇总所有岗位的投递状态分布       |
| `query_position_stats`    | 单岗位统计   | 指定岗位的投递数、学历分布             |
| `query_candidates`        | 候选人筛选   | 按岗位名称、学历、技能、经验多条件筛选 |
| `query_position_hotness`  | 岗位热度排行 | 按投递数量降序排列岗位热度             |

### 对话持久化

每一条 HR 提问和 AI 回复会**成对自动写入** `chat_history` 表，绑定 HR 账号和会话 ID。页面刷新后自动加载历史会话，延续上下文。

---

## 数据库

5 张业务表：

| 表名             | 说明                                          |
| ---------------- | --------------------------------------------- |
| `users`        | 用户表（HR + 候选人统一存储，role 字段区分）  |
| `positions`    | 招聘岗位表                                    |
| `applications` | 投递记录表                                    |
| `resumes`      | 简历文件表（仅存 OSS 元数据，文件实际在 OSS） |
| `chat_history` | AI 对话历史表                                 |

详细表结构见 [db.md](db.md)。

---

## 项目亮点

1. **gRPC 两层微服务架构** — Web 网关与 Logic 业务服务严格分离，通过 gRPC 通信，守住微服务分层考点
2. **OSS 签名直传** — 简历文件前端直传阿里云 OSS，全程不经过服务器，后端仅管理签名 URL 和元数据
3. **文件魔数校验** — 后端确认上传时从 OSS 读取文件头字节，验证真实格式，防止改后缀名绕过
4. **Eino ReAct Agent** — 基于字节跳动 Eino 框架，4 个 SQL 工具查询 MySQL 真实业务数据，不做向量、不做 RAG
5. **SSE 流式 AI 对话** — 实时打字机效果，多会话管理，对话持久化入库
6. **双端分离** — HR 管理端和候选人用户端独立部署，角色互斥登录
7. **完整权限校验** — 投递四层校验链（登录 → 身份 → 资料完整 → 简历存在），岗位归属校验，简历下载权限校验
