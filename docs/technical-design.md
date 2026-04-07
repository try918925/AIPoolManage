# 统一 AI 模型 API 网关 — 技术设计文档

> **版本**: v1.0  
> **日期**: 2026-04-03  
> **技术栈**: Go 1.22+ · Gin · PostgreSQL · Redis  

---

## 目录

1. [项目概述](#1-项目概述)
2. [系统架构](#2-系统架构)
3. [项目目录结构](#3-项目目录结构)
4. [数据库设计](#4-数据库设计)
5. [核心 API 接口设计](#5-核心-api-接口设计)
6. [认证与鉴权](#6-认证与鉴权)
7. [请求路由与代理](#7-请求路由与代理)
8. [负载均衡与故障转移](#8-负载均衡与故障转移)
9. [多厂商适配层](#9-多厂商适配层)
10. [流式 (SSE) 响应处理](#10-流式-sse-响应处理)
11. [限流与配额管理](#11-限流与配额管理)
12. [密钥安全](#12-密钥安全)
13. [错误处理](#13-错误处理)
14. [部署方案](#14-部署方案)
15. [后续扩展](#15-后续扩展)

---

## 1. 项目概述

### 1.1 背景

企业/团队在日常开发中需要同时使用多个 AI 厂商的模型（OpenAI、Anthropic Claude、阿里通义千问、百度文心、讯飞星火等）。直接分发各厂商原始 API Key 存在以下问题：

- 密钥泄露风险高，无法统一管控
- 各厂商 API 格式不同，前端/客户端需要分别适配
- 无法统一进行用量统计、限流与计费

### 1.2 目标

构建一个 **统一 AI 模型 API 网关**，实现：

| 能力 | 说明 |
|------|------|
| **多厂商管理** | 后台配置 OpenAI / Claude / 通义千问等厂商的 API Key |
| **模型管理** | 为每个厂商配置可用模型（如 gpt-4o、claude-sonnet-4-20250514、qwen-turbo） |
| **用户 Key 分发** | 用户在前台自助生成专属 API Key |
| **统一调用接口** | 用户使用自己的 Key，通过 **兼容 OpenAI 格式** 的统一 API 调用任意已配置模型 |
| **用量统计** | 记录每次调用的 token 消耗，支持按用户/模型维度统计 |
| **负载均衡与故障转移** | 同一模型可配置多个厂商通道，支持加权轮询、自动熔断与故障切换 |
| **安全管控** | 限流、配额、密钥加密存储 |

### 1.3 技术选型

| 组件 | 选型 | 理由 |
|------|------|------|
| 语言 | **Go 1.22+** | 高并发、低延迟、适合网关场景 |
| Web 框架 | **Gin** | 轻量高性能，中间件生态丰富 |
| ORM | **GORM** | Go 生态最流行 ORM，迁移方便 |
| 数据库 | **PostgreSQL 16** | JSONB 支持好，适合存储模型参数等半结构化数据 |
| 缓存 | **Redis 7** | 限流计数器、Key 缓存、会话管理 |
| 日志 | **Zap** | 高性能结构化日志 |
| 配置 | **Viper** | 支持多格式配置文件 + 环境变量 |
| 加密 | **AES-256-GCM** | 对称加密存储厂商 API Key |
| 容器化 | **Docker + docker-compose** | 一键部署 |

---

## 2. 系统架构

### 2.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                        客户端 / 前端                             │
│  (使用用户 API Key，兼容 OpenAI SDK 直接调用)                      │
└──────────────────────────┬──────────────────────────────────────┘
                           │  HTTPS
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Nginx / 负载均衡                            │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                    API Gateway (Go + Gin)                        │
│                                                                   │
│  ┌─────────┐  ┌──────────┐  ┌──────────┐  ┌────────────────┐   │
│  │ 认证中间件│  │ 限流中间件 │  │ 日志中间件 │  │ 请求路由 & 代理│   │
│  └─────────┘  └──────────┘  └──────────┘  └────────────────┘   │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │              负载均衡器 (Load Balancer)                       │ │
│  │  · 加权轮询 (Weighted Round-Robin)                           │ │
│  │  · 优先级故障转移 (Priority Failover)                         │ │
│  │  · 熔断器 (Circuit Breaker)                                  │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    厂商适配层 (Adapter)                       │ │
│  │  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐            │ │
│  │  │ OpenAI │  │ Claude │  │  Qwen  │  │ 更多... │            │ │
│  │  └────────┘  └────────┘  └────────┘  └────────┘            │ │
│  └─────────────────────────────────────────────────────────────┘ │
└──────┬──────────────┬───────────────┬───────────────────────────┘
       │              │               │
       ▼              ▼               ▼
┌──────────┐  ┌──────────────┐  ┌──────────────────────────────┐
│PostgreSQL│  │    Redis     │  │      AI 厂商 API (上游)       │
│ 持久存储  │  │ 缓存/限流/   │  │ ┌────────┐ ┌────────┐       │
│          │  │ 熔断状态     │  │ │OpenAI-1│ │OpenAI-2│ ...   │
└──────────┘  └──────────────┘  │ └────────┘ └────────┘       │
                                │ ┌────────┐ ┌────────┐       │
                                │ │ Claude │ │ Azure  │ ...   │
                                │ └────────┘ └────────┘       │
                                └──────────────────────────────┘
```

### 2.2 分层说明

| 层级 | 职责 |
|------|------|
| **API 层** (Handler) | 接收 HTTP 请求，参数校验，调用 Service |
| **Service 层** | 业务逻辑：鉴权、路由决策、请求转换、响应适配 |
| **Adapter 层** | 各厂商 API 协议适配，统一接口抽象 |
| **Repository 层** | 数据库 CRUD 操作 |
| **Middleware 层** | 认证、限流、日志、CORS 等横切关注点 |

---

## 3. 项目目录结构

```
awesomeProject/
├── cmd/
│   └── server/
│       └── main.go                 # 程序入口
├── internal/
│   ├── config/
│   │   └── config.go               # 配置加载 (Viper)
│   ├── middleware/
│   │   ├── auth.go                  # 认证中间件
│   │   ├── ratelimit.go             # 限流中间件
│   │   ├── logger.go                # 日志中间件
│   │   └── cors.go                  # CORS 中间件
│   ├── handler/
│   │   ├── admin_provider.go        # 管理端 - 厂商管理
│   │   ├── admin_model.go           # 管理端 - 模型管理
│   │   ├── admin_channel.go         # 管理端 - 通道健康管理
│   │   ├── user_apikey.go           # 用户端 - API Key 管理
│   │   ├── user_model.go            # 用户端 - 模型查询 & 通道路由
│   │   ├── user_usage.go            # 用户端 - 用量统计 & 明细
│   │   └── chat_completion.go       # 统一推理接口
│   ├── service/
│   │   ├── provider.go              # 厂商业务逻辑
│   │   ├── model.go                 # 模型业务逻辑
│   │   ├── apikey.go                # API Key 业务逻辑
│   │   ├── chat.go                  # 聊天推理业务逻辑
│   │   └── loadbalancer.go          # 负载均衡 & 故障转移
│   ├── adapter/
│   │   ├── adapter.go               # 适配器接口定义
│   │   ├── openai.go                # OpenAI 适配器
│   │   ├── claude.go                # Claude 适配器
│   │   ├── qwen.go                  # 通义千问适配器
│   │   └── registry.go              # 适配器注册表
│   ├── repository/
│   │   ├── provider.go              # 厂商数据访问
│   │   ├── model.go                 # 模型数据访问
│   │   ├── apikey.go                # API Key 数据访问
│   │   └── usage.go                 # 用量日志数据访问
│   ├── model/
│   │   └── entity.go                # 数据库实体定义
│   └── pkg/
│       ├── crypto/
│       │   └── aes.go               # AES-GCM 加解密工具
│       ├── hash/
│       │   └── sha256.go            # SHA-256 哈希工具
│       └── response/
│           └── response.go          # 统一响应格式
├── migrations/
│   ├── 001_create_providers.sql
│   ├── 002_create_models.sql
│   ├── 003_create_user_api_keys.sql
│   └── 004_create_usage_logs.sql
├── configs/
│   ├── config.yaml                  # 默认配置
│   └── config.prod.yaml             # 生产配置
├── docker-compose.yaml
├── Dockerfile
├── go.mod
├── go.sum
└── docs/
    └── technical-design.md          # 本文档
```

---

## 4. 数据库设计

### 4.1 ER 关系图

```
┌──────────────┐       ┌───────────────────┐
│  providers   │ 1───N │  provider_models   │
│──────────────│       │───────────────────│
│ id (PK)      │       │ id (PK)           │
│ name         │       │ provider_id (FK)  │
│ type         │       │ model_name        │
│ base_url     │       │ model_id          │
│ api_key_enc  │       │ weight            │
│ enabled      │       │ priority          │
│ created_at   │       │ enabled           │
│ updated_at   │       │ config (JSONB)    │
└──────────────┘       │ created_at        │
                       │ updated_at        │
                       └───────────────────┘

┌──────────────────┐       ┌──────────────────┐
│  user_api_keys   │ 1───N │   usage_logs     │
│──────────────────│       │──────────────────│
│ id (PK)          │       │ id (PK)          │
│ user_id          │       │ user_key_id (FK) │
│ key_hash         │       │ model_name       │
│ key_prefix       │       │ provider_id      │
│ name             │       │ prompt_tokens    │
│ enabled          │       │ completion_tokens│
│ rate_limit       │       │ total_tokens     │
│ quota_limit      │       │ latency_ms       │
│ quota_used       │       │ status           │
│ expires_at       │       │ created_at       │
│ created_at       │       └──────────────────┘
│ updated_at       │
└──────────────────┘

        (Redis 存储，非持久化表)
┌───────────────────────────┐
│   channel_health (Redis)  │
│───────────────────────────│
│ channel:{model_id}        │
│ status: healthy/unhealthy │
│ consecutive_failures      │
│ last_failure_at           │
│ recovery_at               │
└───────────────────────────┘
```

### 4.2 表结构详细定义

#### 4.2.1 `providers` — AI 厂商表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| `id` | `BIGSERIAL` | PK | 主键 |
| `name` | `VARCHAR(100)` | NOT NULL, UNIQUE | 厂商显示名称，如 "OpenAI" |
| `type` | `VARCHAR(50)` | NOT NULL | 厂商类型标识：`openai`, `claude`, `qwen`, `wenxin`, `spark` |
| `base_url` | `VARCHAR(500)` | NOT NULL | API 基础地址，如 `https://api.openai.com` |
| `api_key_encrypted` | `TEXT` | NOT NULL | AES-256-GCM 加密后的厂商 API Key |
| `org_id` | `VARCHAR(200)` | | 可选，如 OpenAI 的 Organization ID |
| `enabled` | `BOOLEAN` | DEFAULT true | 是否启用 |
| `config` | `JSONB` | | 扩展配置（如自定义 Header、超时时间等） |
| `created_at` | `TIMESTAMPTZ` | DEFAULT NOW() | 创建时间 |
| `updated_at` | `TIMESTAMPTZ` | DEFAULT NOW() | 更新时间 |

```sql
CREATE TABLE providers (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL UNIQUE,
    type          VARCHAR(50)  NOT NULL,
    base_url      VARCHAR(500) NOT NULL,
    api_key_encrypted TEXT NOT NULL,
    org_id        VARCHAR(200),
    enabled       BOOLEAN DEFAULT true,
    config        JSONB DEFAULT '{}',
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_providers_type ON providers(type);
CREATE INDEX idx_providers_enabled ON providers(enabled);
```

#### 4.2.2 `provider_models` — 厂商模型表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| `id` | `BIGSERIAL` | PK | 主键 |
| `provider_id` | `BIGINT` | FK → providers.id | 所属厂商 |
| `model_name` | `VARCHAR(100)` | NOT NULL | 对外暴露的模型名称（用户看到的） |
| `model_id` | `VARCHAR(100)` | NOT NULL | 上游厂商的实际模型 ID |
| `model_type` | `VARCHAR(20)` | DEFAULT 'chat' | 模型类型：`chat`, `embedding`, `image` |
| `enabled` | `BOOLEAN` | DEFAULT true | 是否启用 |
| `weight` | `INTEGER` | DEFAULT 1 | 负载均衡权重，值越大分配流量越多 |
| `priority` | `INTEGER` | DEFAULT 0 | 通道优先级，0=主通道，数值越大优先级越低（用于故障转移） |
| `max_context_tokens` | `INTEGER` | | 最大上下文长度 |
| `input_price` | `DECIMAL(10,6)` | | 输入价格 ($/1K tokens) |
| `output_price` | `DECIMAL(10,6)` | | 输出价格 ($/1K tokens) |
| `config` | `JSONB` | DEFAULT '{}' | 模型级扩展配置 |
| `created_at` | `TIMESTAMPTZ` | DEFAULT NOW() | |
| `updated_at` | `TIMESTAMPTZ` | DEFAULT NOW() | |

```sql
CREATE TABLE provider_models (
    id                 BIGSERIAL PRIMARY KEY,
    provider_id        BIGINT NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    model_name         VARCHAR(100) NOT NULL,
    model_id           VARCHAR(100) NOT NULL,
    model_type         VARCHAR(20)  DEFAULT 'chat',
    enabled            BOOLEAN DEFAULT true,
    weight             INTEGER DEFAULT 1,
    priority           INTEGER DEFAULT 0,
    max_context_tokens INTEGER,
    input_price        DECIMAL(10,6),
    output_price       DECIMAL(10,6),
    config             JSONB DEFAULT '{}',
    created_at         TIMESTAMPTZ DEFAULT NOW(),
    updated_at         TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(provider_id, model_name)
);

CREATE INDEX idx_provider_models_name ON provider_models(model_name);
CREATE INDEX idx_provider_models_enabled ON provider_models(enabled);
CREATE INDEX idx_provider_models_priority ON provider_models(model_name, priority, enabled);
```

#### 4.2.3 `user_api_keys` — 用户 API Key 表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| `id` | `BIGSERIAL` | PK | 主键 |
| `user_id` | `VARCHAR(100)` | NOT NULL | 用户标识 |
| `name` | `VARCHAR(100)` | | Key 别名，如 "开发测试用" |
| `key_hash` | `VARCHAR(64)` | NOT NULL, UNIQUE | SHA-256 哈希后的完整 Key |
| `key_prefix` | `VARCHAR(12)` | NOT NULL | Key 前缀用于展示，如 `sk-proj-a8Kx` |
| `enabled` | `BOOLEAN` | DEFAULT true | 是否启用 |
| `rate_limit` | `INTEGER` | DEFAULT 60 | 每分钟最大请求数 |
| `quota_limit` | `BIGINT` | DEFAULT 0 | Token 总配额，0=无限 |
| `quota_used` | `BIGINT` | DEFAULT 0 | 已使用 Token 数 |
| `allowed_models` | `TEXT[]` | | 允许访问的模型白名单，空=全部 |
| `expires_at` | `TIMESTAMPTZ` | | 过期时间，NULL=永不过期 |
| `last_used_at` | `TIMESTAMPTZ` | | 最后使用时间 |
| `created_at` | `TIMESTAMPTZ` | DEFAULT NOW() | |
| `updated_at` | `TIMESTAMPTZ` | DEFAULT NOW() | |

```sql
CREATE TABLE user_api_keys (
    id             BIGSERIAL PRIMARY KEY,
    user_id        VARCHAR(100) NOT NULL,
    name           VARCHAR(100),
    key_hash       VARCHAR(64)  NOT NULL UNIQUE,
    key_prefix     VARCHAR(12)  NOT NULL,
    enabled        BOOLEAN DEFAULT true,
    rate_limit     INTEGER DEFAULT 60,
    quota_limit    BIGINT  DEFAULT 0,
    quota_used     BIGINT  DEFAULT 0,
    allowed_models TEXT[],
    expires_at     TIMESTAMPTZ,
    last_used_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_user_api_keys_hash ON user_api_keys(key_hash);
CREATE INDEX idx_user_api_keys_user ON user_api_keys(user_id);
```

#### 4.2.4 `usage_logs` — 调用日志/用量表

```sql
CREATE TABLE usage_logs (
    id                BIGSERIAL PRIMARY KEY,
    user_key_id       BIGINT NOT NULL REFERENCES user_api_keys(id),
    user_id           VARCHAR(100) NOT NULL,
    provider_id       BIGINT NOT NULL REFERENCES providers(id),
    model_name        VARCHAR(100) NOT NULL,
    prompt_tokens     INTEGER DEFAULT 0,
    completion_tokens INTEGER DEFAULT 0,
    total_tokens      INTEGER DEFAULT 0,
    latency_ms        INTEGER,
    status            VARCHAR(20) NOT NULL DEFAULT 'success',  -- success / error
    error_message     TEXT,
    request_ip        VARCHAR(45),
    created_at        TIMESTAMPTZ DEFAULT NOW()
);

-- 按时间分区，便于归档和查询
CREATE INDEX idx_usage_logs_user ON usage_logs(user_id, created_at DESC);
CREATE INDEX idx_usage_logs_model ON usage_logs(model_name, created_at DESC);
CREATE INDEX idx_usage_logs_created ON usage_logs(created_at DESC);
```

---

## 5. 核心 API 接口设计

### 5.1 接口概览

| 模块 | 方法 | 路径 | 认证方式 | 说明 |
|------|------|------|----------|------|
| **管理端 - 厂商** | POST | `/admin/providers` | Admin JWT | 添加厂商 |
| | GET | `/admin/providers` | Admin JWT | 厂商列表 |
| | PUT | `/admin/providers/:id` | Admin JWT | 更新厂商 |
| | DELETE | `/admin/providers/:id` | Admin JWT | 删除厂商 |
| **管理端 - 模型** | POST | `/admin/providers/:id/models` | Admin JWT | 添加模型 |
| | GET | `/admin/providers/:id/models` | Admin JWT | 模型列表 |
| | PUT | `/admin/models/:id` | Admin JWT | 更新模型 |
| | DELETE | `/admin/models/:id` | Admin JWT | 删除模型 |
| **管理端 - 通道健康** | GET | `/admin/channels/health` | Admin JWT | 查看通道健康状态 |
| | POST | `/admin/channels/:id/reset` | Admin JWT | 重置通道健康状态 |
| **用户 - Key 管理** | POST | `/api/keys` | User Auth | 生成 API Key |
| | GET | `/api/keys` | User Auth | 查看我的 Key 列表 |
| | DELETE | `/api/keys/:id` | User Auth | 删除 Key |
| **用户 - 模型与路由** | GET | `/v1/models` | API Key | 获取可用模型列表 |
| | GET | `/v1/models/:model` | API Key | 模型详情（含可用通道） |
| | GET | `/v1/models/:model/channels` | API Key | 模型可用通道列表及健康状态 |
| **用户 - 用量查询** | GET | `/api/usage` | API Key | 我的用量统计 |
| | GET | `/api/usage/details` | API Key | 我的调用明细 |
| **统一推理** | POST | `/v1/chat/completions` | API Key | 聊天补全（兼容 OpenAI） |

---

### 5.2 管理端接口详细设计

#### POST `/admin/providers` — 添加 AI 厂商

**请求体：**
```json
{
    "name": "OpenAI",
    "type": "openai",
    "base_url": "https://api.openai.com",
    "api_key": "sk-xxxxxxxxxxxxxxxxxxxxxxxx",
    "org_id": "org-xxxxx",
    "config": {
        "timeout_seconds": 120,
        "max_retries": 3
    }
}
```

**响应 (201)：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "name": "OpenAI",
        "type": "openai",
        "base_url": "https://api.openai.com",
        "enabled": true,
        "created_at": "2026-04-03T10:00:00Z"
    }
}
```

> ⚠️ 注意：响应中 **不返回** `api_key`，密钥一经写入仅后端可读。

#### POST `/admin/providers/:id/models` — 为厂商添加模型

**请求体：**
```json
{
    "model_name": "gpt-4o",
    "model_id": "gpt-4o-2024-08-06",
    "model_type": "chat",
    "max_context_tokens": 128000,
    "input_price": 0.0025,
    "output_price": 0.01,
    "weight": 5,
    "priority": 0
}
```

**响应 (201)：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "provider_id": 1,
        "model_name": "gpt-4o",
        "model_id": "gpt-4o-2024-08-06",
        "model_type": "chat",
        "enabled": true,
        "max_context_tokens": 128000,
        "weight": 5,
        "priority": 0
    }
}
```

---

### 5.3 用户 Key 管理接口

#### POST `/api/keys` — 生成 API Key

**请求体：**
```json
{
    "name": "我的开发 Key",
    "rate_limit": 60,
    "quota_limit": 1000000,
    "allowed_models": ["gpt-4o", "claude-sonnet-4-20250514"],
    "expires_at": "2026-12-31T23:59:59Z"
}
```

**响应 (201)：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "name": "我的开发 Key",
        "key": "sk-proj-a8KxM7nP2qR5tY9w...",
        "key_prefix": "sk-proj-a8Kx",
        "rate_limit": 60,
        "quota_limit": 1000000,
        "expires_at": "2026-12-31T23:59:59Z",
        "created_at": "2026-04-03T10:00:00Z"
    }
}
```

> ⚠️ **`key` 字段仅在创建时返回一次**，后续无法再次获取完整 Key。用户需自行妥善保存。

#### GET `/api/keys` — 查看 Key 列表

**响应 (200)：**
```json
{
    "code": 0,
    "data": [
        {
            "id": 1,
            "name": "我的开发 Key",
            "key_prefix": "sk-proj-a8Kx",
            "enabled": true,
            "rate_limit": 60,
            "quota_limit": 1000000,
            "quota_used": 23456,
            "last_used_at": "2026-04-03T09:30:00Z",
            "expires_at": "2026-12-31T23:59:59Z",
            "created_at": "2026-04-03T10:00:00Z"
        }
    ]
}
```

---

### 5.4 统一推理接口（核心）

#### POST `/v1/chat/completions` — 聊天补全

完全兼容 OpenAI Chat Completions API 格式，用户可直接使用 OpenAI SDK。

**请求头：**
```
Authorization: Bearer sk-proj-a8KxM7nP2qR5tY9w...
Content-Type: application/json
X-Preferred-Provider: openai          # 可选：指定优先厂商类型
X-Channel-Id: 3                       # 可选：指定具体通道 ID
```

> 💡 **路由控制说明**：
> - 不传任何路由头：按系统默认的加权轮询 + 优先级策略自动选择通道
> - `X-Preferred-Provider`：优先选择指定厂商类型的通道（如 `openai`、`claude`），无可用通道时回退到其他厂商
> - `X-Channel-Id`：强制指定某个通道 ID（可通过 `GET /v1/models/:model/channels` 获取），仅该通道可用时生效
> - 以上扩展字段不影响 OpenAI SDK 兼容性，SDK 会忽略未知 Header

**请求体：**
```json
{
    "model": "gpt-4o",
    "messages": [
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "Hello!"}
    ],
    "temperature": 0.7,
    "max_tokens": 2048,
    "stream": false
}
```

**非流式响应 (200)：**

响应头中会附带实际路由信息：
```
X-Channel-Id: 1
X-Provider-Type: openai
X-Provider-Name: OpenAI-主账号
```

响应体（兼容 OpenAI 格式）：
```json
{
    "id": "chatcmpl-xxxxxxxxxxxxx",
    "object": "chat.completion",
    "created": 1712000000,
    "model": "gpt-4o",
    "choices": [
        {
            "index": 0,
            "message": {
                "role": "assistant",
                "content": "Hello! How can I help you today?"
            },
            "finish_reason": "stop"
        }
    ],
    "usage": {
        "prompt_tokens": 20,
        "completion_tokens": 10,
        "total_tokens": 30
    }
}
```

**流式响应 (`stream: true`)：**
```
data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":null}]}

data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}

data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":"!"},"finish_reason":null}]}

data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}

data: [DONE]
```

#### GET `/v1/models` — 获取可用模型列表

**响应 (200)：**
```json
{
    "object": "list",
    "data": [
        {
            "id": "gpt-4o",
            "object": "model",
            "owned_by": "openai",
            "created": 1712000000
        },
        {
            "id": "claude-sonnet-4-20250514",
            "object": "model",
            "owned_by": "anthropic",
            "created": 1712000000
        },
        {
            "id": "qwen-turbo",
            "object": "model",
            "owned_by": "alibaba",
            "created": 1712000000
        }
    ]
}
```

---

### 5.5 用户路由与用量 API

#### GET `/v1/models/:model` — 模型详情

查看某个模型的详细信息，包括可用通道数、厂商类型等。

**请求示例：** `GET /v1/models/gpt-4o`

**响应 (200)：**
```json
{
    "code": 0,
    "data": {
        "model_name": "gpt-4o",
        "model_type": "chat",
        "available_channels": 3,
        "providers": ["openai", "azure-openai"],
        "max_context_tokens": 128000,
        "status": "healthy"
    }
}
```

#### GET `/v1/models/:model/channels` — 模型可用通道列表

查看某个模型下所有对当前用户可用的通道及其健康状态。用户可根据返回的 `channel_id` 在调用时通过 `X-Channel-Id` 请求头指定通道。

**请求示例：** `GET /v1/models/gpt-4o/channels`

**响应 (200)：**
```json
{
    "code": 0,
    "data": {
        "model_name": "gpt-4o",
        "channels": [
            {
                "channel_id": 1,
                "provider_type": "openai",
                "provider_name": "OpenAI-主账号",
                "priority": 0,
                "status": "healthy",
                "avg_latency_ms": 850
            },
            {
                "channel_id": 2,
                "provider_type": "openai",
                "provider_name": "OpenAI-备用账号",
                "priority": 0,
                "status": "healthy",
                "avg_latency_ms": 920
            },
            {
                "channel_id": 3,
                "provider_type": "azure-openai",
                "provider_name": "Azure-东亚",
                "priority": 1,
                "status": "healthy",
                "avg_latency_ms": 1100
            }
        ]
    }
}
```

> 💡 **安全说明**：通道列表 **不会** 暴露厂商的原始 API Key 或 base_url，仅展示通道 ID、厂商类型、健康状态和平均延迟等脱敏信息。

#### GET `/api/usage` — 用量统计

查看当前 API Key 的用量汇总，支持按时间范围、模型维度筛选。

**请求参数 (Query String)：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `start_date` | `string` | 否 | 起始日期 `YYYY-MM-DD`，默认 30 天前 |
| `end_date` | `string` | 否 | 结束日期 `YYYY-MM-DD`，默认今天 |
| `model` | `string` | 否 | 按模型名过滤 |
| `group_by` | `string` | 否 | 分组维度：`day` / `model` / `day,model`，默认 `day` |

**请求示例：** `GET /api/usage?start_date=2026-04-01&end_date=2026-04-03&group_by=model`

**响应 (200)：**
```json
{
    "code": 0,
    "data": {
        "total_requests": 1256,
        "total_tokens": 523400,
        "total_prompt_tokens": 312000,
        "total_completion_tokens": 211400,
        "breakdown": [
            {
                "model": "gpt-4o",
                "requests": 800,
                "prompt_tokens": 200000,
                "completion_tokens": 140000,
                "total_tokens": 340000,
                "avg_latency_ms": 1200
            },
            {
                "model": "claude-sonnet-4-20250514",
                "requests": 300,
                "prompt_tokens": 80000,
                "completion_tokens": 50000,
                "total_tokens": 130000,
                "avg_latency_ms": 1800
            },
            {
                "model": "qwen-turbo",
                "requests": 156,
                "prompt_tokens": 32000,
                "completion_tokens": 21400,
                "total_tokens": 53400,
                "avg_latency_ms": 600
            }
        ]
    }
}
```

#### GET `/api/usage/details` — 调用明细

查看具体的调用记录明细，支持分页。

**请求参数 (Query String)：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `start_date` | `string` | 否 | 起始日期 |
| `end_date` | `string` | 否 | 结束日期 |
| `model` | `string` | 否 | 按模型名过滤 |
| `status` | `string` | 否 | 按状态过滤：`success` / `error` |
| `page` | `int` | 否 | 页码，默认 1 |
| `page_size` | `int` | 否 | 每页条数，默认 20，最大 100 |

**响应 (200)：**
```json
{
    "code": 0,
    "data": {
        "total": 1256,
        "page": 1,
        "page_size": 20,
        "items": [
            {
                "id": 10001,
                "model": "gpt-4o",
                "provider_type": "openai",
                "channel_id": 1,
                "prompt_tokens": 150,
                "completion_tokens": 85,
                "total_tokens": 235,
                "latency_ms": 1340,
                "status": "success",
                "created_at": "2026-04-03T10:30:00Z"
            },
            {
                "id": 10000,
                "model": "claude-sonnet-4-20250514",
                "provider_type": "claude",
                "channel_id": 4,
                "prompt_tokens": 200,
                "completion_tokens": 0,
                "total_tokens": 200,
                "latency_ms": 5000,
                "status": "error",
                "error_message": "upstream timeout",
                "created_at": "2026-04-03T10:28:00Z"
            }
        ]
    }
}
```

---

## 6. 认证与鉴权

### 6.1 双层认证体系

```
┌──────────────────────────────────────────────────┐
│                   请求进入                        │
└──────────────────────┬───────────────────────────┘
                       │
                       ▼
              ┌────────────────┐
              │ 路径以 /admin/ │──── 是 ──→ JWT 认证
              │    开头？      │           (Admin Token)
              └───────┬────────┘
                      │ 否
                      ▼
              ┌────────────────┐
              │ 路径以 /v1/   │──── 是 ──→ API Key 认证
              │ 或 /api/ 开头？│           (Bearer Token)
              └───────┬────────┘
                      │ 否
                      ▼
                  401 Unauthorized
```

### 6.2 Admin JWT 认证

```go
// middleware/auth.go - 伪代码
func AdminAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractBearerToken(c)
        claims, err := jwt.ParseAndVerify(token, secretKey)
        if err != nil {
            c.AbortWithStatusJSON(401, errorResponse("invalid admin token"))
            return
        }
        if claims.Role != "admin" {
            c.AbortWithStatusJSON(403, errorResponse("admin access required"))
            return
        }
        c.Set("admin_id", claims.AdminID)
        c.Next()
    }
}
```

### 6.3 User API Key 认证

```go
// middleware/auth.go - 伪代码
func APIKeyAuthMiddleware(repo repository.UserAPIKeyRepo, cache *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := extractBearerToken(c)  // "sk-proj-a8KxM7nP2qR5..."
        if apiKey == "" {
            c.AbortWithStatusJSON(401, errorResponse("missing api key"))
            return
        }

        // 1. SHA-256 哈希
        keyHash := sha256Hex(apiKey)

        // 2. 先查 Redis 缓存
        cachedKey, err := cache.Get(ctx, "apikey:"+keyHash).Result()
        if err == redis.Nil {
            // 3. 缓存未命中，查数据库
            keyRecord, err := repo.FindByHash(keyHash)
            if err != nil || !keyRecord.Enabled {
                c.AbortWithStatusJSON(401, errorResponse("invalid api key"))
                return
            }
            // 4. 写入缓存 (TTL 5 分钟)
            cache.Set(ctx, "apikey:"+keyHash, serialize(keyRecord), 5*time.Minute)
            cachedKey = keyRecord
        }

        // 5. 检查过期和配额
        if cachedKey.ExpiresAt != nil && time.Now().After(*cachedKey.ExpiresAt) {
            c.AbortWithStatusJSON(401, errorResponse("api key expired"))
            return
        }
        if cachedKey.QuotaLimit > 0 && cachedKey.QuotaUsed >= cachedKey.QuotaLimit {
            c.AbortWithStatusJSON(429, errorResponse("quota exceeded"))
            return
        }

        c.Set("user_key", cachedKey)
        c.Next()
    }
}
```

---

## 7. 请求路由与代理

### 7.1 路由决策流程

```
用户请求: POST /v1/chat/completions { "model": "gpt-4o", ... }
          Headers: X-Preferred-Provider: openai (可选)
                   X-Channel-Id: 3           (可选)
                    │
                    ▼
          ┌──────────────────┐
          │ 1. API Key 认证   │
          └────────┬─────────┘
                   │
                   ▼
          ┌──────────────────┐
          │ 2. 提取 model 字段│
          │    "gpt-4o"      │
          └────────┬─────────┘
                   │
                   ▼
          ┌──────────────────┐     未找到     ┌─────────────┐
          │ 3. 查询          │─────────────→ │ 400 模型不存在│
          │  provider_models │               └─────────────┘
          │  获取所有通道     │
          └────────┬─────────┘
                   │ 找到 N 条记录
                   ▼
          ┌──────────────────┐
          │ 4. 检查用户 Key   │
          │  是否有此模型权限  │──── 无权限 ──→ 403
          └────────┬─────────┘
                   │ 有权限
                   ▼
          ┌──────────────────────────────┐
          │ 5. 解析用户路由偏好           │
          │  · X-Channel-Id → 置顶指定   │
          │  · X-Preferred-Provider →    │
          │    偏好厂商排前              │
          │  · 无偏好 → 默认策略          │
          └────────┬─────────────────────┘
                   │
                   ▼
          ┌──────────────────────────┐
          │ 6. 负载均衡器选择通道      │
          │  · 过滤健康通道            │
          │  · 按 priority 分组       │
          │  · 同优先级加权轮询(weight)│
          └────────┬─────────────────┘
                   │
                   ▼
          ┌──────────────────┐
          │ 7. 获取 Provider │
          │  + 解密 API Key  │
          └────────┬─────────┘
                   │
                   ▼
          ┌──────────────────┐
          │ 8. 适配器转换请求 │
          │  统一格式 → 厂商  │
          └────────┬─────────┘
                   │
                   ▼
          ┌──────────────────┐         ┌─────────────────────┐
          │ 9. 代理转发到上游 │── 失败 →│10. 标记通道不健康     │
          └────────┬─────────┘         │    选择下一通道       │
                   │                   │    重试步骤 7-9       │
                   │ 成功              └─────────────────────┘
                   ▼
          ┌──────────────────┐
          │11. 适配器转换响应 │
          │  厂商 → 统一格式  │
          └────────┬─────────┘
                   │
                   ▼
          ┌──────────────────┐
          │12. 记录用量日志   │
          │  更新 quota_used  │
          │  更新通道健康状态  │
          └────────┬─────────┘
                   │
                   ▼
          ┌──────────────────────┐
          │13. 返回响应 + 路由信息│
          │  X-Channel-Id: 1     │
          │  X-Provider-Type: .. │
          └──────────────────────┘
```

### 7.2 核心代理逻辑

```go
// service/chat.go - 伪代码
func (s *ChatService) ChatCompletion(ctx context.Context, req *ChatRequest, userKey *UserAPIKey, routeHint *RouteHint) (*ChatResponse, error) {
    // RouteHint 来自请求头:
    //   X-Channel-Id       → routeHint.ChannelID
    //   X-Preferred-Provider → routeHint.PreferredProvider

    // 1. 查找模型对应的所有通道（可能有多个 Provider）
    channels, err := s.modelRepo.FindAllByName(req.Model)
    if err != nil || len(channels) == 0 {
        return nil, ErrModelNotFound
    }

    // 2. 检查模型权限
    if !userKey.HasModelAccess(req.Model) {
        return nil, ErrModelForbidden
    }

    // 3. 应用用户路由控制
    sortedChannels := s.applyRouteHint(channels, routeHint)

    // 4. 按优先级依次尝试（故障转移）
    var lastErr error
    var usedChannel *ProviderModel
    for _, channel := range sortedChannels {
        // 4a. 跳过不健康通道
        if !s.loadBalancer.IsHealthy(channel.ID) {
            continue
        }

        // 4b. 获取 Provider 并解密 API Key
        provider, err := s.providerRepo.FindByID(channel.ProviderID)
        if err != nil {
            continue
        }
        providerAPIKey, err := s.crypto.Decrypt(provider.APIKeyEncrypted)
        if err != nil {
            continue
        }

        // 4c. 获取对应适配器
        adapter, err := s.adapterRegistry.Get(provider.Type)
        if err != nil {
            continue
        }

        // 4d. 转换请求并代理
        upstreamReq, err := adapter.ConvertRequest(req, channel, providerAPIKey)
        if err != nil {
            continue
        }

        upstreamResp, err := s.httpClient.Do(upstreamReq)
        if err != nil {
            s.loadBalancer.RecordFailure(channel.ID)
            lastErr = err
            continue
        }

        // 4e. 检查上游 HTTP 状态码
        if upstreamResp.StatusCode >= 500 {
            s.loadBalancer.RecordFailure(channel.ID)
            lastErr = fmt.Errorf("upstream returned %d", upstreamResp.StatusCode)
            upstreamResp.Body.Close()
            continue
        }

        // 5. 转换响应
        resp, usage, err := adapter.ConvertResponse(upstreamResp)
        if err != nil {
            s.loadBalancer.RecordFailure(channel.ID)
            lastErr = err
            continue
        }

        // 6. 成功：记录健康 + 异步记录用量
        usedChannel = channel
        s.loadBalancer.RecordSuccess(channel.ID)
        go s.recordUsage(userKey, channel, provider, usage)

        // 7. 在响应中附带路由信息
        resp.RouteInfo = &RouteInfo{
            ChannelID:    channel.ID,
            ProviderType: provider.Type,
            ProviderName: provider.Name,
        }

        return resp, nil
    }

    // 所有通道均失败
    return nil, fmt.Errorf("all channels exhausted: %w", lastErr)
}

// applyRouteHint 根据用户路由偏好调整通道排序
func (s *ChatService) applyRouteHint(channels []*ProviderModel, hint *RouteHint) []*ProviderModel {
    if hint == nil {
        // 无路由偏好，走默认负载均衡
        return s.loadBalancer.SelectChannels(channels)
    }

    // 场景 1: 指定了具体通道 ID
    if hint.ChannelID > 0 {
        for _, ch := range channels {
            if ch.ID == hint.ChannelID {
                // 优先用指定通道，其余作为 fallback
                result := []*ProviderModel{ch}
                for _, other := range channels {
                    if other.ID != hint.ChannelID {
                        result = append(result, other)
                    }
                }
                return result
            }
        }
        // 指定的通道不存在，回退到默认
    }

    // 场景 2: 指定了偏好厂商类型
    if hint.PreferredProvider != "" {
        preferred := make([]*ProviderModel, 0)
        fallback := make([]*ProviderModel, 0)
        for _, ch := range channels {
            provider, _ := s.providerRepo.FindByID(ch.ProviderID)
            if provider != nil && provider.Type == hint.PreferredProvider {
                preferred = append(preferred, ch)
            } else {
                fallback = append(fallback, ch)
            }
        }
        // 偏好厂商内部仍走负载均衡
        sorted := s.loadBalancer.SelectChannels(preferred)
        sorted = append(sorted, s.loadBalancer.SelectChannels(fallback)...)
        return sorted
    }

    return s.loadBalancer.SelectChannels(channels)
}
```

---

## 8. 负载均衡与故障转移

### 8.1 设计目标

同一个对外模型名（如 `gpt-4o`）可以绑定到 **多个厂商通道**（多条 `provider_models` 记录），实现：

| 能力 | 说明 |
|------|------|
| **多 Key 轮询** | 同一厂商配置多个 API Key，绕开单 Key 速率/配额限制 |
| **跨厂商冗余** | 同一模型由不同厂商提供（如 OpenAI 官方 + Azure OpenAI），一方宕机自动切换 |
| **加权流量分配** | 按 `weight` 字段分配流量比例，大权重通道承担更多请求 |
| **优先级故障转移** | 按 `priority` 字段分级，主通道不可用时自动降级到备用通道 |
| **自动熔断与恢复** | 连续失败达到阈值自动熔断，定时探测恢复 |

### 8.2 通道配置示例

```
模型名: gpt-4o 对应 3 条 provider_models 记录:

┌─────────────────────────────────────────────────────────────┐
│ 通道 A: provider=OpenAI-主账号   priority=0  weight=5       │  ← 主通道, 高权重
│ 通道 B: provider=OpenAI-备用账号  priority=0  weight=3       │  ← 主通道, 低权重
│ 通道 C: provider=Azure-OpenAI    priority=1  weight=1       │  ← 备用通道
└─────────────────────────────────────────────────────────────┘

正常情况: A 和 B 按 5:3 的权重轮询
A 挂了:   B 独自承担主通道流量
A+B 都挂: 自动降级到备用通道 C
C 也挂:   返回 503 所有通道不可用
```

### 8.3 负载均衡策略 — 加权轮询 (Weighted Round-Robin)

```go
// service/loadbalancer.go

type LoadBalancer struct {
    mu       sync.Mutex
    rdb      *redis.Client
    counters map[int64]int64  // channelID → 当前计数器（内存）
}

// SelectChannels 返回按优先级排序、同优先级按权重排序的通道列表
func (lb *LoadBalancer) SelectChannels(channels []*ProviderModel) []*ProviderModel {
    // 1. 按 priority 升序分组
    groups := groupByPriority(channels)  // map[int][]*ProviderModel

    var result []*ProviderModel
    for _, priority := range sortedKeys(groups) {
        group := groups[priority]

        // 2. 过滤健康通道
        healthy := lb.filterHealthy(group)
        if len(healthy) == 0 {
            continue
        }

        // 3. 加权轮询排序
        sorted := lb.weightedRoundRobin(healthy)
        result = append(result, sorted...)
    }

    return result
}

// weightedRoundRobin 基于平滑加权轮询算法选择
func (lb *LoadBalancer) weightedRoundRobin(channels []*ProviderModel) []*ProviderModel {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    // 计算总权重
    totalWeight := 0
    for _, ch := range channels {
        totalWeight += ch.Weight
    }

    // 平滑加权轮询 (Smooth Weighted Round-Robin)
    type weightedChannel struct {
        channel       *ProviderModel
        currentWeight int
    }

    items := make([]weightedChannel, len(channels))
    for i, ch := range channels {
        // 从上一轮累加权重
        items[i] = weightedChannel{
            channel:       ch,
            currentWeight: int(lb.counters[ch.ID]) + ch.Weight,
        }
    }

    // 选出 currentWeight 最大的
    sort.Slice(items, func(i, j int) bool {
        return items[i].currentWeight > items[j].currentWeight
    })

    // 被选中的减去 totalWeight
    items[0].currentWeight -= totalWeight

    // 更新计数器
    for _, item := range items {
        lb.counters[item.channel.ID] = int64(item.currentWeight)
    }

    result := make([]*ProviderModel, len(items))
    for i, item := range items {
        result[i] = item.channel
    }
    return result
}
```

### 8.4 故障转移流程

```
请求进入 → 选择 priority=0 的通道组
                    │
          ┌─────────▼─────────┐
          │ 加权轮询选中通道 A   │
          └─────────┬─────────┘
                    │
              ┌─────▼─────┐
              │ 调用上游 API│
              └─────┬─────┘
                    │
              ┌─────▼─────┐     失败 (超时/5xx/连接拒绝)
              │  成功？     │──── 否 ──→ 记录失败次数
              └─────┬─────┘              │
                    │ 是                 ▼
                    │          ┌─────────────────┐
                    │          │ 连续失败 ≥ 阈值?  │
                    │          └────┬────────┬───┘
                    │               │ 否     │ 是
                    │               ▼        ▼
                    │          重试同优先级  标记为不健康
                    │          下一通道     (熔断)
                    │               │        │
                    │               ▼        ▼
                    │          ┌────────────────────┐
                    │          │ 同优先级还有通道？    │
                    │          └────┬──────────┬────┘
                    │               │ 是       │ 否
                    │               ▼          ▼
                    │          选择下一通道  降级到 priority=1
                    │                          │
                    ▼                          ▼
              返回响应               重复以上流程...
                                       │
                                       ▼
                                 所有通道耗尽 → 503
```

### 8.5 熔断器 (Circuit Breaker) — Redis 实现

```go
// service/loadbalancer.go

const (
    FailureThreshold  = 5              // 连续失败 N 次触发熔断
    RecoveryTimeout   = 30 * time.Second // 熔断后多久进入半开状态
    HalfOpenMaxProbes = 3              // 半开状态最多放行 N 个探测请求
)

// 通道状态机:
//   Closed (正常) → 连续失败≥阈值 → Open (熔断)
//   Open → 等待 RecoveryTimeout → HalfOpen (半开, 放行少量请求探测)
//   HalfOpen → 探测成功 → Closed
//   HalfOpen → 探测失败 → Open

// IsHealthy 判断通道是否可用
func (lb *LoadBalancer) IsHealthy(channelID int64) bool {
    key := fmt.Sprintf("channel:health:%d", channelID)
    data, err := lb.rdb.HGetAll(ctx, key).Result()
    if err != nil || len(data) == 0 {
        return true  // 无记录，视为健康
    }

    status := data["status"]
    switch status {
    case "closed":
        return true
    case "open":
        // 检查是否到达恢复时间
        recoveryAt, _ := time.Parse(time.RFC3339, data["recovery_at"])
        if time.Now().After(recoveryAt) {
            // 进入半开状态
            lb.rdb.HSet(ctx, key, "status", "half_open", "probe_count", 0)
            return true
        }
        return false
    case "half_open":
        // 半开状态限制探测请求数
        probeCount, _ := strconv.Atoi(data["probe_count"])
        if probeCount < HalfOpenMaxProbes {
            lb.rdb.HIncrBy(ctx, key, "probe_count", 1)
            return true
        }
        return false
    }
    return true
}

// RecordFailure 记录一次失败
func (lb *LoadBalancer) RecordFailure(channelID int64) {
    key := fmt.Sprintf("channel:health:%d", channelID)

    failures, _ := lb.rdb.HIncrBy(ctx, key, "consecutive_failures", 1).Result()
    lb.rdb.HSet(ctx, key, "last_failure_at", time.Now().Format(time.RFC3339))

    if failures >= FailureThreshold {
        // 触发熔断
        lb.rdb.HSet(ctx, key,
            "status", "open",
            "recovery_at", time.Now().Add(RecoveryTimeout).Format(time.RFC3339),
        )
        lb.rdb.Expire(ctx, key, 10*time.Minute)  // 兜底过期
    }
}

// RecordSuccess 记录一次成功
func (lb *LoadBalancer) RecordSuccess(channelID int64) {
    key := fmt.Sprintf("channel:health:%d", channelID)
    lb.rdb.HSet(ctx, key,
        "status", "closed",
        "consecutive_failures", 0,
    )
    lb.rdb.Expire(ctx, key, 5*time.Minute)
}
```

### 8.6 管理端通道权重配置示例

管理员为同一模型配置多个通道（复用已有的模型管理接口）：

```bash
# 通道 A: OpenAI 主账号, 主通道, 高权重
curl -X POST http://localhost:8080/admin/providers/1/models \
  -H "Authorization: Bearer <admin-jwt>" \
  -d '{
    "model_name": "gpt-4o",
    "model_id": "gpt-4o-2024-08-06",
    "weight": 5,
    "priority": 0
  }'

# 通道 B: OpenAI 备用账号, 主通道, 低权重
curl -X POST http://localhost:8080/admin/providers/2/models \
  -H "Authorization: Bearer <admin-jwt>" \
  -d '{
    "model_name": "gpt-4o",
    "model_id": "gpt-4o-2024-08-06",
    "weight": 3,
    "priority": 0
  }'

# 通道 C: Azure OpenAI, 备用通道
curl -X POST http://localhost:8080/admin/providers/3/models \
  -H "Authorization: Bearer <admin-jwt>" \
  -d '{
    "model_name": "gpt-4o",
    "model_id": "gpt-4o",
    "weight": 1,
    "priority": 1
  }'
```

### 8.7 通道健康监控接口

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| GET | `/admin/channels/health` | Admin JWT | 查看所有通道健康状态 |
| POST | `/admin/channels/:id/reset` | Admin JWT | 手动重置通道为健康状态 |

**GET `/admin/channels/health` 响应示例：**
```json
{
    "code": 0,
    "data": [
        {
            "channel_id": 1,
            "provider_name": "OpenAI-主账号",
            "model_name": "gpt-4o",
            "status": "closed",
            "consecutive_failures": 0,
            "weight": 5,
            "priority": 0
        },
        {
            "channel_id": 2,
            "provider_name": "OpenAI-备用账号",
            "model_name": "gpt-4o",
            "status": "open",
            "consecutive_failures": 5,
            "last_failure_at": "2026-04-03T10:30:00Z",
            "recovery_at": "2026-04-03T10:30:30Z",
            "weight": 3,
            "priority": 0
        }
    ]
}
```

---

## 9. 多厂商适配层

### 9.1 适配器接口定义

```go
// adapter/adapter.go
type Adapter interface {
    // ConvertRequest 将统一格式请求转为厂商特定格式
    ConvertRequest(req *ChatRequest, model *ProviderModel, apiKey string) (*http.Request, error)

    // ConvertResponse 将厂商响应转为统一格式
    ConvertResponse(resp *http.Response) (*ChatResponse, *Usage, error)

    // ConvertStreamResponse 将厂商流式响应转为统一 SSE 格式
    ConvertStreamResponse(resp *http.Response, writer http.ResponseWriter) (*Usage, error)

    // ProviderType 返回厂商类型标识
    ProviderType() string
}
```

### 9.2 各厂商适配差异对照

| 维度 | OpenAI | Claude (Anthropic) | 通义千问 (Qwen) |
|------|--------|-------------------|-----------------|
| **API 路径** | `/v1/chat/completions` | `/v1/messages` | `/api/v1/services/aigc/text-generation/generation` |
| **认证方式** | `Authorization: Bearer sk-xxx` | `x-api-key: sk-ant-xxx` + `anthropic-version: 2023-06-01` | `Authorization: Bearer sk-xxx` |
| **请求格式** | `{ messages, model, ... }` | `{ messages, model, max_tokens, system(单独字段) }` | `{ model, input: { messages }, parameters: {...} }` |
| **System 消息** | 在 messages 数组中 | 独立的 `system` 字段 | 在 messages 数组中（role=system） |
| **流式格式** | `data: {...}\n\n` | `event: content_block_delta\ndata: {...}\n\n` | `data: {...}\n\n` |
| **流式结束** | `data: [DONE]` | `event: message_stop` | `data: [DONE]` (兼容模式) |
| **Token 统计位置** | `response.usage` | `response.usage` | `response.usage` |

### 9.3 OpenAI 适配器示例

```go
// adapter/openai.go
type OpenAIAdapter struct{}

func (a *OpenAIAdapter) ProviderType() string { return "openai" }

func (a *OpenAIAdapter) ConvertRequest(req *ChatRequest, model *ProviderModel, apiKey string) (*http.Request, error) {
    // OpenAI 格式与统一格式一致，直接透传
    body := map[string]interface{}{
        "model":       model.ModelID,  // 使用厂商实际模型 ID
        "messages":    req.Messages,
        "temperature": req.Temperature,
        "max_tokens":  req.MaxTokens,
        "stream":      req.Stream,
    }
    
    jsonBody, _ := json.Marshal(body)
    httpReq, _ := http.NewRequest("POST",
        provider.BaseURL+"/v1/chat/completions",
        bytes.NewReader(jsonBody))
    
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+apiKey)
    
    return httpReq, nil
}
```

### 9.4 Claude 适配器示例

```go
// adapter/claude.go
type ClaudeAdapter struct{}

func (a *ClaudeAdapter) ConvertRequest(req *ChatRequest, model *ProviderModel, apiKey string) (*http.Request, error) {
    // 提取 system 消息（Claude 要求 system 为独立字段）
    var systemMsg string
    var messages []Message
    for _, msg := range req.Messages {
        if msg.Role == "system" {
            systemMsg = msg.Content
        } else {
            messages = append(messages, msg)
        }
    }

    body := map[string]interface{}{
        "model":      model.ModelID,
        "max_tokens": req.MaxTokens,
        "messages":   messages,
        "stream":     req.Stream,
    }
    if systemMsg != "" {
        body["system"] = systemMsg
    }
    
    jsonBody, _ := json.Marshal(body)
    httpReq, _ := http.NewRequest("POST",
        provider.BaseURL+"/v1/messages",
        bytes.NewReader(jsonBody))
    
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("x-api-key", apiKey)
    httpReq.Header.Set("anthropic-version", "2023-06-01")
    
    return httpReq, nil
}
```

### 9.5 适配器注册表

```go
// adapter/registry.go
type Registry struct {
    adapters map[string]Adapter
}

func NewRegistry() *Registry {
    r := &Registry{adapters: make(map[string]Adapter)}
    r.Register(&OpenAIAdapter{})
    r.Register(&ClaudeAdapter{})
    r.Register(&QwenAdapter{})
    return r
}

func (r *Registry) Register(a Adapter) {
    r.adapters[a.ProviderType()] = a
}

func (r *Registry) Get(providerType string) (Adapter, error) {
    a, ok := r.adapters[providerType]
    if !ok {
        return nil, fmt.Errorf("unsupported provider type: %s", providerType)
    }
    return a, nil
}
```

---

## 10. 流式 (SSE) 响应处理

### 10.1 处理流程

```
用户 ←─ SSE ←─ 网关 ←─ SSE ←─ 上游厂商
         │              │
         │     ┌────────┴──────────┐
         │     │ Adapter 实时转换    │
         │     │ 厂商 SSE → 统一 SSE │
         │     └───────────────────┘
         │
    统一为 OpenAI 格式 SSE
```

### 10.2 核心实现

```go
// handler/chat_completion.go - 流式处理伪代码
func (h *ChatHandler) handleStream(c *gin.Context, req *ChatRequest, userKey *UserAPIKey) {
    // 设置 SSE 响应头
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    c.Header("X-Accel-Buffering", "no")  // 禁用 Nginx 缓冲

    // 获取上游流式响应
    upstreamResp, _ := adapter.DoStreamRequest(req)
    defer upstreamResp.Body.Close()

    // 适配器逐块转换并写入
    flusher := c.Writer.(http.Flusher)
    usage, err := adapter.ConvertStreamResponse(upstreamResp, c.Writer)
    flusher.Flush()

    // 写入结束标记
    fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
    flusher.Flush()

    // 异步记录用量
    go h.service.RecordUsage(userKey, usage)
}
```

### 10.3 Claude SSE → 统一格式转换示例

```go
// Claude SSE 事件流:
//   event: message_start
//   event: content_block_start
//   event: content_block_delta  ← 内容在这里
//   event: content_block_stop
//   event: message_delta        ← usage 在这里
//   event: message_stop

// 转换为 OpenAI 格式:
//   data: {"choices":[{"delta":{"content":"xxx"}}]}
//   data: [DONE]
```

---

## 11. 限流与配额管理

### 11.1 限流策略

使用 **Redis 滑动窗口** 实现每分钟请求数限制：

```go
// middleware/ratelimit.go - 伪代码
func RateLimitMiddleware(rdb *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        userKey := c.MustGet("user_key").(*UserAPIKey)
        
        key := fmt.Sprintf("ratelimit:%d:%d", userKey.ID, time.Now().Unix()/60)
        
        count, _ := rdb.Incr(ctx, key).Result()
        if count == 1 {
            rdb.Expire(ctx, key, 2*time.Minute)
        }
        
        if count > int64(userKey.RateLimit) {
            c.Header("X-RateLimit-Limit", strconv.Itoa(userKey.RateLimit))
            c.Header("X-RateLimit-Remaining", "0")
            c.Header("Retry-After", "60")
            c.AbortWithStatusJSON(429, errorResponse("rate limit exceeded"))
            return
        }
        
        c.Header("X-RateLimit-Limit", strconv.Itoa(userKey.RateLimit))
        c.Header("X-RateLimit-Remaining", strconv.FormatInt(int64(userKey.RateLimit)-count, 10))
        c.Next()
    }
}
```

### 11.2 配额管理

```go
// service/quota.go
func (s *QuotaService) CheckAndDeduct(userKey *UserAPIKey, tokens int64) error {
    if userKey.QuotaLimit == 0 {
        return nil  // 无限配额
    }
    
    // 使用 Redis + DB 双写保证准确性
    // Redis: 实时扣减（高性能）
    newUsed, err := s.redis.IncrBy(ctx, 
        fmt.Sprintf("quota:%d", userKey.ID), tokens).Result()
    if err != nil {
        return err
    }
    
    if newUsed > userKey.QuotaLimit {
        // 回滚 Redis
        s.redis.DecrBy(ctx, fmt.Sprintf("quota:%d", userKey.ID), tokens)
        return ErrQuotaExceeded
    }
    
    // 异步同步到数据库
    go s.syncQuotaToDB(userKey.ID, newUsed)
    
    return nil
}
```

---

## 12. 密钥安全

### 12.1 厂商 API Key — AES-256-GCM 加密存储

```go
// pkg/crypto/aes.go
func Encrypt(plaintext string, masterKey []byte) (string, error) {
    block, err := aes.NewCipher(masterKey)  // 32 bytes = AES-256
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    io.ReadFull(rand.Reader, nonce)
    
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encoded string, masterKey []byte) (string, error) {
    ciphertext, _ := base64.StdEncoding.DecodeString(encoded)
    
    block, _ := aes.NewCipher(masterKey)
    gcm, _ := cipher.NewGCM(block)
    
    nonceSize := gcm.NonceSize()
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    return string(plaintext), err
}
```

**Master Key 管理方式**（优先级从高到低）：
1. ✅ 使用 HashiCorp Vault / AWS KMS 等密钥管理服务
2. ✅ 通过环境变量注入 `MASTER_ENCRYPTION_KEY`
3. ⚠️ 配置文件（仅开发环境）

### 12.2 用户 API Key — SHA-256 哈希存储

```go
// pkg/hash/sha256.go
func GenerateAPIKey() (fullKey string, prefix string, hash string) {
    // 1. 生成 32 字节随机数
    randomBytes := make([]byte, 32)
    rand.Read(randomBytes)
    
    // 2. 编码为可读字符串
    fullKey = "sk-proj-" + base32.StdEncoding.EncodeToString(randomBytes)[:48]
    
    // 3. 提取前缀用于展示
    prefix = fullKey[:12]
    
    // 4. SHA-256 哈希用于存储和查找
    h := sha256.Sum256([]byte(fullKey))
    hash = hex.EncodeToString(h[:])
    
    return fullKey, prefix, hash
}
```

**安全要点：**
- 用户 Key 全文 **仅在生成时返回一次**
- 数据库中只存储 **SHA-256 哈希** + **前缀**
- 每次认证时对用户传入的 Key 做 SHA-256 后与数据库比对

---

## 13. 错误处理

### 13.1 统一错误响应格式

兼容 OpenAI 错误格式，方便客户端 SDK 处理：

```json
{
    "error": {
        "message": "The model 'gpt-5' does not exist or you do not have access to it.",
        "type": "invalid_request_error",
        "param": "model",
        "code": "model_not_found"
    }
}
```

### 13.2 错误码对照表

| HTTP 状态码 | error.type | error.code | 说明 |
|-------------|-----------|------------|------|
| 400 | `invalid_request_error` | `invalid_request` | 请求参数错误 |
| 401 | `authentication_error` | `invalid_api_key` | API Key 无效或缺失 |
| 403 | `permission_error` | `model_not_allowed` | 无权访问该模型 |
| 404 | `invalid_request_error` | `model_not_found` | 模型不存在 |
| 429 | `rate_limit_error` | `rate_limit_exceeded` | 请求频率超限 |
| 429 | `rate_limit_error` | `quota_exceeded` | Token 配额用尽 |
| 500 | `api_error` | `internal_error` | 服务端内部错误 |
| 502 | `api_error` | `upstream_error` | 上游厂商 API 错误 |
| 503 | `api_error` | `provider_unavailable` | 厂商服务不可用 |

### 13.3 错误处理中间件

```go
// middleware/error.go
func ErrorHandlerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            var apiErr *APIError
            if errors.As(err, &apiErr) {
                c.JSON(apiErr.HTTPStatus, gin.H{
                    "error": gin.H{
                        "message": apiErr.Message,
                        "type":    apiErr.Type,
                        "code":    apiErr.Code,
                    },
                })
            } else {
                c.JSON(500, gin.H{
                    "error": gin.H{
                        "message": "internal server error",
                        "type":    "api_error",
                        "code":    "internal_error",
                    },
                })
            }
        }
    }
}
```

---

## 14. 部署方案

### 14.1 Docker Compose

```yaml
# docker-compose.yaml
version: '3.8'

services:
  api-gateway:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=gateway
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=ai_gateway
      - REDIS_ADDR=redis:6379
      - MASTER_ENCRYPTION_KEY=${MASTER_ENCRYPTION_KEY}
      - JWT_SECRET=${JWT_SECRET}
      - GIN_MODE=release
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: gateway
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ai_gateway
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U gateway"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes --requirepass ${REDIS_PASSWORD}
    volumes:
      - redisdata:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pgdata:
  redisdata:
```

### 14.2 Dockerfile

```dockerfile
# 多阶段构建
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /server .
COPY configs/ ./configs/
EXPOSE 8080
CMD ["./server"]
```

### 14.3 配置文件

```yaml
# configs/config.yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 120s   # 流式响应需要较长超时

database:
  host: ${DB_HOST:localhost}
  port: ${DB_PORT:5432}
  user: ${DB_USER:gateway}
  password: ${DB_PASSWORD}
  dbname: ${DB_NAME:ai_gateway}
  max_open_conns: 50
  max_idle_conns: 10

redis:
  addr: ${REDIS_ADDR:localhost:6379}
  password: ${REDIS_PASSWORD}
  db: 0
  pool_size: 100

security:
  master_key: ${MASTER_ENCRYPTION_KEY}    # AES-256 主密钥 (32 字节 hex)
  jwt_secret: ${JWT_SECRET}
  jwt_expire: 24h

proxy:
  timeout: 120s            # 上游请求超时
  max_idle_conns: 200      # HTTP 连接池
  idle_conn_timeout: 90s

loadbalancer:
  strategy: weighted_round_robin   # 负载均衡策略: weighted_round_robin / random
  circuit_breaker:
    failure_threshold: 5           # 连续失败 N 次触发熔断
    recovery_timeout: 30s          # 熔断后探测恢复间隔
    half_open_max_probes: 3        # 半开状态最大探测请求数
```

### 14.4 部署架构（生产环境）

```
                    ┌─────────────┐
                    │   DNS / CDN  │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │    Nginx     │
                    │  SSL 终止    │
                    │  负载均衡    │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
        ┌─────▼─────┐ ┌───▼───┐ ┌─────▼─────┐
        │ Gateway-1 │ │ GW-2  │ │  GW-3     │
        └─────┬─────┘ └───┬───┘ └─────┬─────┘
              │            │            │
        ┌─────▼────────────▼────────────▼─────┐
        │         PostgreSQL (主从)             │
        │         Redis Cluster                │
        └─────────────────────────────────────┘
```

---

## 15. 后续扩展

### V2 规划

| 特性 | 说明 |
|------|------|
| **多租户支持** | 不同团队/组织独立管理厂商配置和用户 Key |
| **Embedding 接口** | 支持 `/v1/embeddings` 统一接口 |
| **图像生成接口** | 支持 `/v1/images/generations` 统一接口 |
| **Webhook 通知** | 配额即将用尽、异常调用等事件通知 |
| **Web 管理面板** | 可视化管理厂商、模型、Key、用量统计 |
| **OpenAPI 文档** | Swagger/OpenAPI 3.0 自动生成文档 |
| **Prometheus 监控** | 暴露 `/metrics` 端点，接入 Grafana 仪表盘 |

---

## 附录 A：快速启动

```bash
# 1. 克隆项目
git clone <repo-url>
cd awesomeProject

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 填入数据库密码、加密密钥等

# 3. 启动服务
docker-compose up -d

# 4. 添加厂商 (Admin)
curl -X POST http://localhost:8080/admin/providers \
  -H "Authorization: Bearer <admin-jwt>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "OpenAI",
    "type": "openai",
    "base_url": "https://api.openai.com",
    "api_key": "sk-your-openai-key"
  }'

# 5. 添加模型
curl -X POST http://localhost:8080/admin/providers/1/models \
  -H "Authorization: Bearer <admin-jwt>" \
  -H "Content-Type: application/json" \
  -d '{
    "model_name": "gpt-4o",
    "model_id": "gpt-4o-2024-08-06",
    "model_type": "chat",
    "max_context_tokens": 128000
  }'

# 6. 用户生成 API Key
curl -X POST http://localhost:8080/api/keys \
  -H "Authorization: Bearer <user-jwt>" \
  -H "Content-Type: application/json" \
  -d '{"name": "My Dev Key"}'
# 响应中获取 key: "sk-proj-xxxx..."

# 7. 使用统一接口调用
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer sk-proj-xxxx..." \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# 8. 查看模型可用通道
curl http://localhost:8080/v1/models/gpt-4o/channels \
  -H "Authorization: Bearer sk-proj-xxxx..."

# 9. 指定厂商调用
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer sk-proj-xxxx..." \
  -H "Content-Type: application/json" \
  -H "X-Preferred-Provider: openai" \
  -d '{
    "model": "gpt-4o",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# 10. 查看用量统计
curl "http://localhost:8080/api/usage?start_date=2026-04-01&group_by=model" \
  -H "Authorization: Bearer sk-proj-xxxx..."
```

## 附录 B：使用 OpenAI Python SDK 调用

```python
from openai import OpenAI
import httpx

# 只需修改 base_url 和 api_key
client = OpenAI(
    base_url="http://localhost:8080/v1",
    api_key="sk-proj-xxxx..."  # 用户自己生成的 Key
)

# ========================================
# 基本调用 — 自动路由（加权轮询 + 故障转移）
# ========================================

response = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "你好！"}]
)
print(response.choices[0].message.content)

# 同一个客户端调用 Claude 模型
response = client.chat.completions.create(
    model="claude-sonnet-4-20250514",
    messages=[{"role": "user", "content": "你好！"}]
)
print(response.choices[0].message.content)

# 流式调用通义千问
stream = client.chat.completions.create(
    model="qwen-turbo",
    messages=[{"role": "user", "content": "讲个笑话"}],
    stream=True
)
for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="")

# ========================================
# 高级用法 — 路由控制 & 用量查询
# ========================================

# 查看可用模型通道 (直接调 REST API)
import requests

headers = {"Authorization": "Bearer sk-proj-xxxx..."}
base = "http://localhost:8080"

# 查看 gpt-4o 的可用通道
channels = requests.get(f"{base}/v1/models/gpt-4o/channels", headers=headers).json()
for ch in channels["data"]["channels"]:
    print(f"  通道 {ch['channel_id']}: {ch['provider_name']} ({ch['status']})")

# 指定通道调用（通过 extra_headers）
response = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "你好！"}],
    extra_headers={"X-Preferred-Provider": "openai"}  # 优先走 OpenAI 通道
)
print(response.choices[0].message.content)

# 指定具体通道 ID
response = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "你好！"}],
    extra_headers={"X-Channel-Id": "3"}  # 走 Azure 备用通道
)
print(response.choices[0].message.content)

# 查看我的用量统计
usage = requests.get(
    f"{base}/api/usage?start_date=2026-04-01&group_by=model",
    headers=headers
).json()
print(f"总请求数: {usage['data']['total_requests']}")
print(f"总 Token: {usage['data']['total_tokens']}")
for item in usage["data"]["breakdown"]:
    print(f"  {item['model']}: {item['requests']} 次, {item['total_tokens']} tokens")
```

---

> **文档结束** — 如有疑问或需要修改，请联系技术架构组。

