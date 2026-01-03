# Soliton-Go 分布式全栈开发框架

**Soliton-Go** 是一个基于 Go 语言构建的模块化、高性能后端开发脚手架，专为 **DDD (领域驱动设计)** 和 **分布式系统** 设计。

## 🚀 核心特性

- **一键生成可用代码**: `--fields` 参数直接生成带完整字段的领域模型
- **领域驱动设计**: AggregateRoot、Entity、ValueObject、Repository
- **分布式能力**: 分布式锁、事件驱动、Saga 分布式事务
- **CQRS 模式**: 内置 Command/Query 处理器
- **依赖注入**: 全项目集成 Uber Fx

## ⚡ 30 秒快速体验

```bash
# 1. 编译生成器
cd tools/generator && go build -o soliton-gen .

# 2. 创建新项目
./soliton-gen init my-project && cd my-project

# 3. 生成领域模块 (--wire 自动接入 main.go)
soliton-gen domain User --fields "username,email,status:enum(active|inactive)" --wire

# 4. 运行
GOWORK=off go mod tidy && GOWORK=off go run ./cmd/main.go
```

**生成结果：**
| 层 | 文件 |
|---|------|
| Domain | `user.go` (含 UserRole、UserStatus 枚举), `repository.go`, `events.go` |
| Application | `commands.go`, `queries.go`, `dto.go`, `module.go` |
| Infrastructure | `user_repo.go` |
| Interfaces | `user_handler.go` |

## 🛠 字段类型支持

| 类型 | 示例 | 生成结果 |
|------|------|----------|
| string | `username` | `Username string` |
| int64 | `price:int64` | `Price int64` |
| text | `desc:text` | `Desc string` (GORM: text) |
| uuid | `user_id:uuid` | `UserId string` (带索引) |
| time? | `login_at:time?` | `LoginAt *time.Time` (可选字段) |
| enum | `status:enum(a\|b\|c)` | 生成枚举类型和常量 |

## 🔌 命令列表

| 命令 | 说明 |
|------|------|
| `init <name>` | 初始化新项目（含 DDD 目录结构） |
| `domain <name>` | 生成领域模块（Entity/Repo/Handler 等） |
| `service <name>` | 生成应用服务（跨领域业务逻辑） |

### --wire 自动接线
```bash
# 支持多模块自动注入
soliton-gen domain User --fields "..." --wire
soliton-gen domain Product --fields "..." --wire
```
`--wire` 使用标记行追加模块，支持多模块无需手动接线。

### 其他参数
| 参数 | 说明 |
|------|------|
| `--table "xxx"` | 自定义数据库表名 |
| `--route "xxx"` | 自定义 API 路由基路径 |
| `--force` | 强制覆盖已存在文件 |

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
| [docs/DEVELOPMENT_GUIDE.md](./docs/DEVELOPMENT_GUIDE.md) | 开发指南 |
| [docs/GENERATOR_GUIDE.md](./docs/GENERATOR_GUIDE.md) | 生成器使用 |
| [docs/SERVICE_GUIDE.md](./docs/SERVICE_GUIDE.md) | Service 详解 |

---

## 🤝 贡献
欢迎提交 Issue 和 PR。
