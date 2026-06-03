# 智能招聘系统 - 数据库设计文档

---

## 一、数据库概述

本系统使用 **MySQL 8.0+** 作为数据库管理系统，数据库名为 `recruitment`，采用 InnoDB 存储引擎，字符集为 `utf8mb4`。

### ER 关系图

```
 users ──1:N──> positions ──1:N──> applications ──1:1──> resumes
   │                                    │
   └────────── 1:N ────────────────────>│ (applicant)
                                        │
 users ──1:N──> chat_history
```

### 实体关系说明

| 关联方向       | 关联字段                                         | 说明             |
| -------------- | ------------------------------------------------ | ---------------- |
| HR → 岗位     | `positions.hr_id` → `users.id`              | HR 发布岗位      |
| 候选人 → 投递 | `applications.candidate_id` → `users.id`    | 候选人投递岗位   |
| 岗位 → 投递   | `applications.position_id` → `positions.id` | 岗位关联投递记录 |
| 候选人 → 简历 | `resumes.candidate_id` → `users.id`         | 候选人上传简历   |
| HR → 对话     | `chat_history.hr_id` → `users.id`           | HR 与 AI 对话    |

---

## 二、数据表结构

> 所有表的 GORM Model 定义位于 `logic-grpc-service/internal/model/`，遵循以下规范：
>
> - 可为 NULL 的字段使用指针类型（`*string`、`*int64`）
> - 软删除使用 `gorm.DeletedAt` 类型，GORM 自动注入 `WHERE deleted_at IS NULL`
> - 必须显式定义 `TableName()` 方法
> - **禁止使用 `db.AutoMigrate()`**，所有表结构必须通过 SQL 脚本管理

### 2.1 用户表（users）

**用途**：存储 HR 和候选人的统一用户信息，通过 `role` 字段区分角色。

| 字段名            | 类型                   | 约束                                                            | 说明                                      |
| ----------------- | ---------------------- | --------------------------------------------------------------- | ----------------------------------------- |
| `id`            | BIGINT UNSIGNED        | PRIMARY KEY, AUTO_INCREMENT                                     | 用户唯一标识                              |
| `username`      | VARCHAR(64)            | NOT NULL                                                        | 登录账号                                  |
| `password_hash` | VARCHAR(255)           | NOT NULL                                                        | bcrypt 哈希密码                           |
| `role`          | ENUM('hr','candidate') | NOT NULL                                                        | 角色：hr（人力资源）/ candidate（候选人） |
| `real_name`     | VARCHAR(64)            | DEFAULT NULL                                                    | 真实姓名（候选人必填，HR 可为空）         |
| `phone`         | VARCHAR(20)            | DEFAULT NULL                                                    | 联系电话（候选人必填）                    |
| `education`     | VARCHAR(32)            | DEFAULT NULL                                                    | 最高学历（候选人必填）                    |
| `school`        | VARCHAR(128)           | DEFAULT NULL                                                    | 毕业院校（候选人必填）                    |
| `experience`    | TEXT                   | DEFAULT NULL                                                    | 工作/项目经历（候选人必填）               |
| `skills`        | VARCHAR(255)           | DEFAULT NULL                                                    | 核心技能标签（逗号分隔，候选人必填）      |
| `created_at`    | DATETIME               | NOT NULL, DEFAULT CURRENT_TIMESTAMP                             | 创建时间                                  |
| `updated_at`    | DATETIME               | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 更新时间                                  |
| `deleted_at`    | DATETIME               | DEFAULT NULL                                                    | 软删除时间                                |

**索引**：

- `uk_username` (username) - 登录唯一性校验

**建表语句**：

```sql
CREATE TABLE users (
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username      VARCHAR(64)  NOT NULL COMMENT '登录账号',
    password_hash VARCHAR(255) NOT NULL COMMENT 'bcrypt 哈希密码',
    role          ENUM('hr','candidate') NOT NULL COMMENT '角色：hr | candidate',
    real_name     VARCHAR(64)  DEFAULT NULL COMMENT '真实姓名',
    phone         VARCHAR(20)  DEFAULT NULL COMMENT '联系电话',
    education     VARCHAR(32)  DEFAULT NULL COMMENT '最高学历',
    school        VARCHAR(128) DEFAULT NULL COMMENT '毕业院校',
    experience    TEXT         DEFAULT NULL COMMENT '工作/项目经历',
    skills        VARCHAR(255) DEFAULT NULL COMMENT '核心技能标签（逗号分隔）',
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    DATETIME     DEFAULT NULL COMMENT '软删除时间',
    UNIQUE KEY uk_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表（HR+候选人）';
```

### GORM 模型定义

```go
type User struct {
    ID           int64          `gorm:"column:id;primaryKey"`
    Username     string         `gorm:"column:username"`
    PasswordHash string         `gorm:"column:password_hash"`
    Role         string         `gorm:"column:role"`
    RealName     *string        `gorm:"column:real_name"`
    Phone        *string        `gorm:"column:phone"`
    Education    *string        `gorm:"column:education"`
    School       *string        `gorm:"column:school"`
    Experience   *string        `gorm:"column:experience"`
    Skills       *string        `gorm:"column:skills"`
    CreatedAt    time.Time      `gorm:"column:created_at"`
    UpdatedAt    time.Time      `gorm:"column:updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (User) TableName() string { return "users" }
```

角色常量：`model.RoleHR = "hr"`、`model.RoleCandidate = "candidate"`

---

### 2.2 招聘岗位表（positions）

**用途**：存储 HR 发布的招聘岗位信息。

| 字段名           | 类型                        | 约束                                                            | 说明                                         |
| ---------------- | --------------------------- | --------------------------------------------------------------- | -------------------------------------------- |
| `id`           | BIGINT UNSIGNED             | PRIMARY KEY, AUTO_INCREMENT                                     | 岗位唯一标识                                 |
| `hr_id`        | BIGINT UNSIGNED             | NOT NULL                                                        | 发布者 HR 的用户 ID                          |
| `title`        | VARCHAR(128)                | NOT NULL                                                        | 岗位名称                                     |
| `description`  | TEXT                        | NOT NULL                                                        | 岗位职责描述                                 |
| `requirements` | TEXT                        | DEFAULT NULL                                                    | 任职要求                                     |
| `salary_min`   | INT UNSIGNED                | DEFAULT NULL                                                    | 最低薪资（单位：K）                          |
| `salary_max`   | INT UNSIGNED                | DEFAULT NULL                                                    | 最高薪资（单位：K）                          |
| `location`     | VARCHAR(128)                | DEFAULT NULL                                                    | 工作地点                                     |
| `status`       | ENUM('published','offline') | NOT NULL, DEFAULT 'published'                                   | 状态：published（发布中）/ offline（已下架） |
| `created_at`   | DATETIME                    | NOT NULL, DEFAULT CURRENT_TIMESTAMP                             | 创建时间                                     |
| `updated_at`   | DATETIME                    | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 更新时间                                     |
| `deleted_at`   | DATETIME                    | DEFAULT NULL                                                    | 软删除时间                                   |

**索引**：

- `idx_hr` (hr_id) - HR 查询自己发布的岗位
- `idx_status` (status) - 按状态筛选岗位
- `idx_title` (title) - 岗位名称搜索
- `idx_salary_range` (salary_min, salary_max) - 薪资区间筛选和排序

**建表语句**：

```sql
CREATE TABLE positions (
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    hr_id         BIGINT UNSIGNED NOT NULL COMMENT '发布者 HR 的用户 ID',
    title         VARCHAR(128) NOT NULL COMMENT '岗位名称',
    description   TEXT         NOT NULL COMMENT '岗位职责描述',
    requirements  TEXT         COMMENT '任职要求',
    salary_min    INT UNSIGNED COMMENT '最低薪资（单位：K）',
    salary_max    INT UNSIGNED COMMENT '最高薪资（单位：K）',
    location      VARCHAR(128) COMMENT '工作地点',
    status        ENUM('published','offline') NOT NULL DEFAULT 'published' COMMENT 'published=发布中 offline=已下架',
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    DATETIME     DEFAULT NULL COMMENT '软删除时间',
    INDEX idx_hr (hr_id),
    INDEX idx_status (status),
    INDEX idx_title (title),
    INDEX idx_salary_range (salary_min, salary_max)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='招聘岗位表';
```

### GORM 模型定义

```go
type Position struct {
    ID           int64          `gorm:"column:id;primaryKey"`
    HRID         int64          `gorm:"column:hr_id"`
    Title        string         `gorm:"column:title"`
    Description  string         `gorm:"column:description"`
    Requirements string         `gorm:"column:requirements"`
    SalaryMin    *int32         `gorm:"column:salary_min"`
    SalaryMax    *int32         `gorm:"column:salary_max"`
    Location     *string        `gorm:"column:location"`
    Status       string         `gorm:"column:status"`
    CreatedAt    time.Time      `gorm:"column:created_at"`
    UpdatedAt    time.Time      `gorm:"column:updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Position) TableName() string { return "positions" }
```

状态常量：`model.PositionStatusPublished = "published"`、`model.PositionStatusOffline = "offline"`

---

### 2.3 投递记录表（applications）

**用途**：存储候选人投递岗位的记录，支持撤回后重新投递。

| 字段名           | 类型                                             | 约束                                                            | 说明                |
| ---------------- | ------------------------------------------------ | --------------------------------------------------------------- | ------------------- |
| `id`           | BIGINT UNSIGNED                                  | PRIMARY KEY, AUTO_INCREMENT                                     | 投递记录唯一标识    |
| `position_id`  | BIGINT UNSIGNED                                  | NOT NULL                                                        | 投递的岗位 ID       |
| `candidate_id` | BIGINT UNSIGNED                                  | NOT NULL                                                        | 投递的候选人用户 ID |
| `resume_id`    | BIGINT UNSIGNED                                  | DEFAULT NULL                                                    | 关联的简历 ID       |
| `status`       | ENUM('pending','reviewed','rejected','accepted') | NOT NULL, DEFAULT 'pending'                                     | 投递状态            |
| `created_at`   | DATETIME                                         | NOT NULL, DEFAULT CURRENT_TIMESTAMP                             | 创建时间            |
| `updated_at`   | DATETIME                                         | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 更新时间            |
| `deleted_at`   | DATETIME                                         | DEFAULT NULL                                                    | 软删除时间          |

**状态说明**：

- `pending` - 待审核
- `reviewed` - 已查看
- `rejected` - 已拒绝
- `accepted` - 已录用

**索引**：

- `uk_position_candidate_deleted` (position_id, candidate_id, deleted_at) - 防重复投递（支持撤回后重新投递）
- `idx_candidate_deleted` (candidate_id, deleted_at) - 查询候选人的全部投递
- `idx_position_deleted` (position_id, deleted_at) - HR 查询某岗位的全部投递者
- `idx_status` (status) - 按投递状态筛选

**建表语句**：

```sql
CREATE TABLE applications (
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    position_id   BIGINT UNSIGNED NOT NULL COMMENT '岗位 ID',
    candidate_id  BIGINT UNSIGNED NOT NULL COMMENT '候选人用户 ID',
    resume_id     BIGINT UNSIGNED DEFAULT NULL COMMENT '关联的简历 ID',
    status        ENUM('pending','reviewed','rejected','accepted') NOT NULL DEFAULT 'pending' COMMENT '投递状态',
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    DATETIME     DEFAULT NULL COMMENT '软删除时间',
    UNIQUE KEY uk_position_candidate_deleted (position_id, candidate_id, deleted_at),
    INDEX idx_candidate_deleted (candidate_id, deleted_at),
    INDEX idx_position_deleted (position_id, deleted_at),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='投递记录表';
```

### GORM 模型定义

```go
type Application struct {
    ID          int64          `gorm:"column:id;primaryKey"`
    PositionID  int64          `gorm:"column:position_id"`
    CandidateID int64          `gorm:"column:candidate_id"`
    ResumeID    *int64         `gorm:"column:resume_id"`
    Status      string         `gorm:"column:status"`
    CreatedAt   time.Time      `gorm:"column:created_at"`
    UpdatedAt   time.Time      `gorm:"column:updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Application) TableName() string { return "applications" }
```

状态常量：`model.ApplicationStatusPending = "pending"`、`model.ApplicationStatusReviewed = "reviewed"`、`model.ApplicationStatusRejected = "rejected"`、`model.ApplicationStatusAccepted = "accepted"`

---

### 2.4 简历文件表（resumes）

**用途**：存储简历文件的 OSS 元信息，下载 URL 实时生成不存储。

| 字段名           | 类型                     | 约束                                | 说明                     |
| ---------------- | ------------------------ | ----------------------------------- | ------------------------ |
| `id`           | BIGINT UNSIGNED          | PRIMARY KEY, AUTO_INCREMENT         | 简历唯一标识             |
| `candidate_id` | BIGINT UNSIGNED          | NOT NULL                            | 所属候选人用户 ID        |
| `file_name`    | VARCHAR(255)             | NOT NULL                            | 原始文件名               |
| `file_type`    | ENUM('pdf','doc','docx') | NOT NULL                            | 文件格式                 |
| `file_size`    | BIGINT                   | NOT NULL                            | 文件大小（字节）         |
| `oss_key`      | VARCHAR(512)             | NOT NULL                            | OSS 对象 Key（唯一路径） |
| `uploaded_at`  | DATETIME                 | NOT NULL, DEFAULT CURRENT_TIMESTAMP | 上传时间                 |
| `deleted_at`   | DATETIME                 | DEFAULT NULL                        | 软删除时间               |

**支持的文件格式**：

- `pdf` - PDF 格式
- `doc` - Microsoft Word 97-2003 格式
- `docx` - Microsoft Word 2007+ 格式

**索引**：

- `idx_candidate` (candidate_id) - 查询候选人的简历列表

**建表语句**：

```sql
CREATE TABLE resumes (
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    candidate_id  BIGINT UNSIGNED NOT NULL COMMENT '所属候选人用户 ID',
    file_name     VARCHAR(255) NOT NULL COMMENT '原始文件名',
    file_type     ENUM('pdf','doc','docx') NOT NULL COMMENT '文件格式',
    file_size     BIGINT       NOT NULL COMMENT '文件大小（字节）',
    oss_key       VARCHAR(512) NOT NULL COMMENT 'OSS 对象 Key',
    uploaded_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    DATETIME     DEFAULT NULL COMMENT '软删除时间',
    INDEX idx_candidate (candidate_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='简历文件表（OSS 元数据）';
```

### GORM 模型定义

```go
type Resume struct {
    ID          int64          `gorm:"column:id;primaryKey"`
    CandidateID int64          `gorm:"column:candidate_id"`
    FileName    string         `gorm:"column:file_name"`
    FileType    string         `gorm:"column:file_type"`
    FileSize    int64          `gorm:"column:file_size"`
    OSSKey      string         `gorm:"column:oss_key"`
    UploadedAt  time.Time      `gorm:"column:uploaded_at"`
    DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Resume) TableName() string { return "resumes" }
```

文件类型常量：`model.ResumeFileTypePDF = "pdf"`、`model.ResumeFileTypeDOC = "doc"`、`model.ResumeFileTypeDOCX = "docx"`

OSS 对象键格式：`resumes/{candidateID}/{unixTimestamp}_{originalFileName}`（由 `oss.GenerateObjectKey()` 生成）

---

### 2.5 AI 对话历史表（chat_history）

**用途**：存储 HR 与 AI 助手的对话历史，支持多会话管理。

| 字段名         | 类型                     | 约束                                | 说明                      |
| -------------- | ------------------------ | ----------------------------------- | ------------------------- |
| `id`         | BIGINT UNSIGNED          | PRIMARY KEY, AUTO_INCREMENT         | 对话记录唯一标识          |
| `hr_id`      | BIGINT UNSIGNED          | NOT NULL                            | HR 用户 ID                |
| `session_id` | VARCHAR(64)              | NOT NULL                            | 会话 ID，用于分组管理对话 |
| `role`       | ENUM('user','assistant') | NOT NULL                            | 发言角色                  |
| `content`    | TEXT                     | NOT NULL                            | 对话内容                  |
| `created_at` | DATETIME                 | NOT NULL, DEFAULT CURRENT_TIMESTAMP | 创建时间                  |
| `deleted_at` | DATETIME                 | DEFAULT NULL                        | 软删除时间                |

**角色说明**：

- `user` - HR 用户
- `assistant` - AI 助手

**索引**：

- `idx_hr_session` (hr_id, session_id) - 按 HR + 会话分组查询对话
- `idx_session_created` (session_id, created_at) - 按会话+时间顺序加载对话历史
- `idx_hr_created` (hr_id, created_at) - 按 HR + 时间加载对话历史

**建表语句**：

```sql
CREATE TABLE chat_history (
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    hr_id         BIGINT UNSIGNED NOT NULL COMMENT 'HR 用户 ID',
    session_id    VARCHAR(64)   NOT NULL COMMENT '会话 ID',
    role          ENUM('user','assistant') NOT NULL COMMENT '发言角色',
    content       TEXT         NOT NULL COMMENT '对话内容',
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    DATETIME     DEFAULT NULL COMMENT '软删除时间',
    INDEX idx_hr_session (hr_id, session_id),
    INDEX idx_session_created (session_id, created_at),
    INDEX idx_hr_created (hr_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI对话历史表';
```

### GORM 模型定义

```go
type ChatHistory struct {
    ID        int64          `gorm:"column:id;primaryKey"`
    HRID      int64          `gorm:"column:hr_id"`
    SessionID string         `gorm:"column:session_id"`
    Role      string         `gorm:"column:role"`
    Content   string         `gorm:"column:content"`
    CreatedAt time.Time      `gorm:"column:created_at"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (ChatHistory) TableName() string { return "chat_history" }
```

角色常量：`model.ChatRoleUser = "user"`、`model.ChatRoleAssistant = "assistant"`

会话 ID 由时间戳纳秒生成：`fmt.Sprintf("%d", time.Now().UnixNano())`

---

## 三、索引汇总

| 表名         | 索引名                        | 字段                                  | 用途                  |
| ------------ | ----------------------------- | ------------------------------------- | --------------------- |
| users        | uk_username                   | username                              | 登录唯一性校验        |
| positions    | idx_hr                        | hr_id                                 | HR 查询自己发布的岗位 |
| positions    | idx_status                    | status                                | 按状态筛选岗位        |
| positions    | idx_title                     | title                                 | 岗位名称搜索          |
| positions    | idx_salary_range              | salary_min, salary_max                | 薪资区间筛选和排序    |
| applications | uk_position_candidate_deleted | position_id, candidate_id, deleted_at | 防重复投递            |
| applications | idx_candidate_deleted         | candidate_id, deleted_at              | 查询候选人的全部投递  |
| applications | idx_position_deleted          | position_id, deleted_at               | HR 查询某岗位的投递者 |
| applications | idx_status                    | status                                | 按投递状态筛选        |
| resumes      | idx_candidate                 | candidate_id                          | 查询候选人的简历列表  |
| chat_history | idx_hr_session                | hr_id, session_id                     | 按 HR + 会话分组查询  |
| chat_history | idx_session_created           | session_id, created_at                | 按会话+时间加载历史   |
| chat_history | idx_hr_created                | hr_id, created_at                     | 按 HR + 时间加载历史  |

---

## 四、表关系约束（代码层保证）

本系统不使用数据库外键约束，所有关联关系由代码层保证一致性。

| 关联方向       | 关联字段                                         | 代码校验点                                |
| -------------- | ------------------------------------------------ | ----------------------------------------- |
| HR → 岗位     | `positions.hr_id` → `users.id`              | 岗位写操作校验 `hr_id == token.user_id` |
| 候选人 → 投递 | `applications.candidate_id` → `users.id`    | 投递前校验候选人角色和资料完整性          |
| 岗位 → 投递   | `applications.position_id` → `positions.id` | 投递前校验岗位存在且 status='published'   |
| 候选人 → 简历 | `resumes.candidate_id` → `users.id`         | 简历上传/下载校验 candidate_id 归属       |
| HR → 对话     | `chat_history.hr_id` → `users.id`           | 对话历史按 hr_id 隔离查询                 |

---

## 五、初始化脚本

**创建数据库**：

```sql
CREATE DATABASE IF NOT EXISTS recruitment 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE recruitment;
```

**完整建表脚本**：请参考上述各表的建表语句，按顺序执行。

---

## 六、数据安全

### 软删除机制

所有表均采用软删除模式，通过 `deleted_at` 字段标记删除状态：

- 删除操作：设置 `deleted_at = CURRENT_TIMESTAMP`
- 查询操作：默认过滤 `deleted_at IS NULL` 的记录

### 敏感数据处理

| 字段              | 处理方式                |
| ----------------- | ----------------------- |
| `password_hash` | bcrypt 加密存储，不可逆 |
| `phone`         | 仅存储                  |

---

## 七、注意事项

1. **禁止使用 GORM AutoMigrate**：所有表结构变更必须通过 SQL 脚本管理
2. **索引维护**：根据业务查询需求定期评估索引有效性
3. **数据备份**：定期备份数据库，建议每日增量备份 + 每周全量备份
4. **性能优化**：对于大数据量的表（如 `applications`、`chat_history`），考虑分表或分区策略
