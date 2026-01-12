# Soliton-Go 分布式全栈开发框架

**Soliton-Go** 是一个基于 Go 语言构建的模块化、高性能后端开发脚手架，专为 **DDD (领域驱动设计)** 和 **分布式系统** 设计。

## 🚀 核心特性

- **一键生成可用代码**: `--fields` 参数直接生成带完整字段的领域模型
- **领域驱动设计**: AggregateRoot、Entity、ValueObject、Specification、Policy、Repository
- **分布式能力**: 分布式锁、事件驱动、Saga 分布式事务
- **CQRS 模式**: 内置 Command/Query 处理器
- **依赖注入**: 全项目集成 Uber Fx
- **迁移入口**: 自动生成 `cmd/migrate/main.go`，支持一键建表
- **默认可用配置**: 未提供 `config.yaml` 也可启动（默认 sqlite + log.level=info）

## ⚡ 30 秒快速体验

```bash
# 1. 编译生成器
cd tools/generator && go build -o soliton-gen .

# 2. 创建新项目
./soliton-gen init my-project && cd my-project

# 3. 生成领域模块 (--wire 自动接入 main.go)
soliton-gen domain User --fields "username,email,status:enum(active|inactive)" --wire

# 4. 生成应用服务（可选，支持方法备注）
soliton-gen service OrderService --methods "CreateOrder::创建订单,CancelOrder::取消订单"

# 5. 运行
GOWORK=off go mod tidy && GOWORK=off go run ./cmd/main.go
```

> **数据库驱动**: 默认支持 sqlite/postgres，如需 MySQL 请扩展 `framework/orm/db.go`。

**生成结果：**
| 层 | 文件 |
|---|------|
| Domain | `user.go` (含 UserRole、UserStatus 枚举), `repository.go`, `events.go`, `service.go` |
| Application | `commands.go`, `queries.go`, `dto.go`, `module.go` |
| Infrastructure | `user_repo.go` |
| Interfaces | `user_handler.go` |

## 🛠 字段类型支持

| 类型 | 示例 | 生成结果 |
|------|------|----------|
| string | `username` | `Username string` |
| text | `desc:text` | `Desc string` (GORM: text) |
| int | `count:int` | `Count int` |
| int64 | `price:int64` | `Price int64` |
| float64 | `score:float64` | `Score float64` |
| decimal | `amount:decimal` | `Amount float64` (GORM: decimal(10,2)) |
| bool | `active:bool` | `Active bool` |
| time | `created_at:time` | `CreatedAt time.Time` |
| time? | `login_at:time?` | `LoginAt *time.Time` (可选字段) |
| date | `birth:date` | `Birth time.Time` (GORM: date) |
| date? | `expire:date?` | `Expire *time.Time` (可选日期) |
| uuid | `user_id:uuid` | `UserId string` (带索引) |
| json | `meta:json` | `Meta datatypes.JSON` |
| jsonb | `data:jsonb` | `Data datatypes.JSON` (PostgreSQL) |
| bytes | `avatar:bytes` | `Avatar []byte` |
| enum | `status:enum(a\|b\|c)` | 生成枚举类型和常量 |

> [!WARNING]
> **已知限制 / Known Limitations**
> - **Domain 命令参数不全**: 缺少 `--api-only`（仅生成 API）、`--no-crud`（不生成 CRUD）、`--no-events`（不生成事件）等选项

## 🔌 命令列表

| 命令 | 说明 |
|------|------|
| `init <name>` | 初始化新项目（含 DDD 目录结构） |
| `domain <name>` | 生成领域模块（Entity/Repo/Handler 等） |
| `service <name>` | 生成应用服务（跨领域业务逻辑） |
| `valueobject <domain> <name>` | 生成领域值对象 |
| `spec <domain> <name>` | 生成领域规格（Specification） |
| `policy <domain> <name>` | 生成领域策略（Policy） |
| `event <domain> <name>` | 生成领域事件（含注册） |
| `event-handler <domain> <event>` | 生成事件处理器并注入 |
| `serve` | 🆕 启动 Web GUI（可视化代码生成器） |

### 🎨 Web GUI - 可视化代码生成

```bash
# 启动 Web 界面
soliton-gen serve

# 自定义端口
soliton-gen serve --port 8080
```

**功能特性：**
- ✨ 可视化字段编辑器，支持拖拽
- 👁️ 生成前预览代码
- 🔌 自动注入模块到 main.go
- 🆕 方法备注（服务方法用途说明，可回显）
- 📖 详细的操作提示和使用指南
- 🌐 中英双语界面

访问 http://127.0.0.1:3000 即可使用图形界面进行项目初始化、领域生成和服务生成。

详细使用说明请查看：[Web GUI 使用指南](./tools/generator/WEB_GUI_GUIDE.md)

### --wire 自动接线
```bash
# 支持多模块自动注入
soliton-gen domain User --fields "..." --wire
soliton-gen domain Product --fields "..." --wire
```
`--wire` 使用标记行追加模块，支持多模块无需手动接线。

### Domain 命令参数
| 参数 | 说明 |
|------|------|
| `--fields "..."` | 指定字段列表 |
| `--table "xxx"` | 自定义数据库表名 |
| `--route "xxx"` | 自定义 API 路由基路径 |
| `--soft-delete` | 🆕 启用软删除 (`DeletedAt` 字段) |
| `--force` | 强制覆盖已存在文件 |
| `--wire` | 自动接入 main.go |

### Domain 子命令
| 子命令 | 说明 |
|--------|------|
| `domain list` | 列出项目中所有已生成的领域模块 |
| `domain delete <name>` | 删除领域模块及其所有相关文件 |

> [!NOTE]
> 子命令支持不全：目前缺少 `domain show <name>`（查看领域详情）、`domain add-field`（添加字段）等子命令

## 🆕 新增功能

### 分页查询
生成的 List API 自动支持分页：
```bash
curl "http://localhost:8080/api/users?page=1&page_size=20"
```
返回结果：
```json
{
  "code": 0,
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

### 排序参数
List API 支持排序参数：
```bash
curl "http://localhost:8080/api/users?page=1&page_size=20&sort_by=created_at&sort_order=desc"
```

### 数据库迁移入口
生成项目包含 `cmd/migrate/main.go`，可单独执行迁移：
```bash
GOWORK=off go run ./cmd/migrate
```

### 软删除
```bash
soliton-gen domain User --fields "username,email" --soft-delete
```
自动添加 `DeletedAt gorm.DeletedAt` 字段，删除操作变为软删除。

### 错误码常量
生成的 `response.go` 包含预定义错误码：
```go
const (
    CodeSuccess      = 0     // 成功
    CodeBadRequest   = 400   // 请求错误
    CodeValidation   = 1001  // 验证失败
    CodeDuplicate    = 1002  // 重复条目
)
```

---

## 📂 项目结构

```
soliton-go/
├── framework/              # 核心框架
│   ├── ddd/                # DDD 原语
│   ├── orm/                # GORM 泛型 Repository
│   ├── event/              # 事件总线
│   └── lock/               # 分布式锁
├── application/            # 业务应用
│   └── internal/
│       ├── domain/         # 领域层
│       ├── application/    # 应用层
│       ├── infrastructure/ # 基础设施层
│       └── interfaces/     # 接口层
└── tools/generator/        # 代码生成器
```

---

## 📖 文档

| 文档 | 说明 |
|------|------|
| [QUICK_START.md](./QUICK_START.md) | 快速上手 |
| [tools/generator/WEB_GUI_GUIDE.md](./tools/generator/WEB_GUI_GUIDE.md) | 🆕 Web GUI 使用指南 |
| [docs/DEVELOPMENT_GUIDE.md](./docs/DEVELOPMENT_GUIDE.md) | 开发指南 |
| [docs/GENERATOR_GUIDE.md](./docs/GENERATOR_GUIDE.md) | 生成器使用 |
| [docs/SERVICE_GUIDE.md](./docs/SERVICE_GUIDE.md) | Service 详解 |

---

## 🤝 贡献
欢迎提交 Issue 和 PR。
