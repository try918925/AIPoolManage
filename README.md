# AI Gateway — 统一 AI 模型 API 网关

统一管理多家 AI 厂商（OpenAI / Claude / 通义千问），提供兼容 OpenAI 格式的 API 代理，内置负载均衡、熔断、限流、用量统计和 Web 管理界面。

## 功能特性

- **多厂商管理** — 集中管理 OpenAI、Anthropic Claude、阿里通义千问的 API Key（AES-256-GCM 加密存储）
- **统一 API** — 兼容 OpenAI `/v1/chat/completions` 格式，支持流式（SSE）响应
- **负载均衡** — 加权轮询（Weighted Round-Robin），按优先级分组调度
- **熔断保护** — Circuit Breaker，自动摘除故障通道并探测恢复
- **速率限制** — 基于 Redis 的请求限流
- **用户系统** — 注册/登录、JWT 认证、角色区分（admin / user）
- **API Key 分发** — 用户自助生成 Key，支持配额、速率限制、过期时间
- **用量统计** — 按用户/模型记录 Token 消耗、延迟、调用状态
- **路由提示** — 支持 `X-Channel-Id` / `X-Preferred-Provider` 指定通道
- **Web UI** — 内置管理后台（`/`）和用户中心（`/portal`），多主题可切换

## 技术栈

| 组件 | 选型 |
|------|------|
| 语言 | Go 1.25 |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7 |
| 认证 | JWT |
| 日志 | Zap |
| 配置 | Viper |

## 快速开始

### Docker Compose（推荐）

```bash
docker-compose up -d
```

服务启动后访问：
- 管理后台：http://localhost:8080/
- 用户中心：http://localhost:8080/portal
- 健康检查：http://localhost:8080/health

### 本地开发

**前置条件：** Go 1.25+、PostgreSQL、Redis

1. 克隆项目并安装依赖：

```bash
git clone <repo-url>
cd AIPoolManage
go mod download
```

2. 复制并修改配置文件：

```bash
cp configs/config.yaml.example configs/config.yaml
# 编辑 configs/config.yaml 填写数据库、Redis、密钥等信息
```

3. 初始化数据库（执行 `migrations/` 下的 SQL）

4. 启动服务：

```bash
go run ./cmd/server
```

## 配置说明

配置文件路径：`configs/config.yaml`

```yaml
server:
  port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: ai_gateway
  max_open_conns: 25
  max_idle_conns: 5

redis:
  addr: localhost:6379
  password: ""
  db: 0

security:
  master_key: "32字节AES密钥"
  jwt_secret: "jwt签名密钥"
  jwt_expire: 24h

proxy:
  timeout: 120s
  max_idle_conns: 100

loadbalancer:
  strategy: weighted_round_robin
  circuit_breaker:
    failure_threshold: 5
    recovery_timeout: 30s
```

## API 概览

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/user/register` | 用户注册 |
| POST | `/user/login` | 用户登录 |
| POST | `/admin/login` | 管理员登录 |

### OpenAI 兼容接口（API Key 认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/v1/models` | 模型列表 |
| POST | `/v1/chat/completions` | Chat Completion（支持流式） |

### 管理接口（Admin JWT）

| 方法 | 路径 | 说明 |
|------|------|------|
| CRUD | `/admin/providers` | 厂商管理 |
| CRUD | `/admin/providers/:id/models` | 模型通道管理 |
| CRUD | `/admin/users` | 用户管理 |
| GET | `/admin/channels/health` | 通道健康状态 |

### 用户接口（User JWT）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET/DELETE | `/user/api/keys` | API Key 管理 |
| GET | `/user/api/usage` | 用量统计 |
| GET | `/user/api/models` | 可用模型 |

## 调用示例

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

```python
from openai import OpenAI

client = OpenAI(
    api_key="YOUR_API_KEY",
    base_url="http://localhost:8080/v1"
)

response = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "Hello!"}]
)
print(response.choices[0].message.content)
```

## 项目结构

```
├── cmd/server/          # 服务入口
├── configs/             # 配置文件
├── docs/                # 设计文档
├── internal/
│   ├── adapter/         # AI 厂商适配器（OpenAI/Claude/Qwen）
│   ├── config/          # 配置加载
│   ├── handler/         # HTTP 处理器
│   ├── middleware/       # 中间件（认证/CORS/日志/限流）
│   ├── model/           # 数据模型
│   ├── pkg/             # 工具包（加密/哈希/响应）
│   ├── repository/      # 数据访问层
│   └── service/         # 业务逻辑（负载均衡/聊天/用户等）
└── migrations/          # 数据库迁移
```

## License

MIT
