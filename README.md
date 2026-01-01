# Soliton-Go 分布式全栈开发框架

**Soliton-Go** 是一个基于 Go 语言构建的模块化、高性能后端开发脚手架。它借鉴了现代企业级框架（如 Spring Boot）的设计理念，专为构建**分布式系统**和**领域驱动设计 (DDD)** 应用而设计。

## 🚀 核心特性

*   **领域驱动设计 (DDD)**: 提供标准的 AggregateRoot, Entity, ValueObject 基类和分层架构规范。
*   **一键生成可用代码**: `soliton-gen domain Order` 自动生成 **Entity + Repository + Events + Commands + Queries + HTTP Handler + Fx Module**，开箱即用。
*   **分布式能力**:
    *   **分布式锁**: 基于 Redis 的高可用分布式锁，支持自动续期 (Watchdog)。
    *   **事件驱动**: 内置 Watermill 事件总线，支持事件类型注册和自动反序列化。
    *   **分布式事务**: 提供 **Saga 模式** 编排器，支持补偿机制。
*   **CQRS 模式**: 内置 Command Bus 和 Query Bus。
*   **极致开发体验**:
    *   **依赖注入**: 全项目集成 **Uber Fx**，一行代码导入整个模块依赖。
    *   **泛型仓储**: 提供 `GormRepository[T, ID]`，自动获得 CRUD 能力。
    *   **HTTP Handler 自动生成**: REST API 无需手写。

## ⚡ 30 秒快速体验

```bash
# 1. 编译生成器
cd tools/generator && go build -o soliton-gen .

# 2. 生成一个完整的领域模块
./soliton-gen domain Order

# 3. 查看生成的 9 个文件
ls ../../application/internal/domain/order/
ls ../../application/internal/application/order/
```

**生成内容：**
| 层 | 文件 | 说明 |
|---|------|------|
| Domain | `order.go`, `repository.go`, `events.go` | 实体、仓储接口、领域事件 |
| Application | `commands.go`, `queries.go`, `dto.go`, `module.go` | CQRS Handler + Fx Module |
| Infrastructure | `order_repo.go` | Repository 实现 + 数据库迁移 |
| Interfaces | `order_handler.go` | HTTP CRUD 5 接口 |

## 🛠 技术栈

| 模块 | 技术选型 | 说明 |
| :--- | :--- | :--- |
| **语言** | Go 1.22+ | 需要 Generic 和 Workspace |
| **Web 框架** | [Gin](https://gin-gonic.com/) | 高性能 HTTP 路由 |
| **ORM** | [GORM](https://gorm.io/) | 强大的 ORM 库 |
| **依赖注入** | [Uber Fx](https://go.uber.org/fx) | 模块化依赖注入 |
| **事件总线** | [Watermill](https://watermill.io/) | 通用事件驱动库 |
| **分布式锁** | [Redislock](https://github.com/bsm/redislock) | Redis 分布式锁 |

---

## 📂 项目结构

```
soliton-go/
├── go.work                 # Go Workspace 定义
├── framework/              # 核心框架层 (通用库)
│   ├── ddd/                # DDD 原语 (AggregateRoot, Entity, ValueObject)
│   ├── orm/                # GORM 封装与泛型 Repository
│   ├── event/              # 事件总线 + 事件类型注册表
│   ├── lock/               # 分布式锁
│   ├── cqrs/               # CQRS 命令/查询总线
│   ├── service/            # 通用业务服务层
│   └── transaction/        # Saga 分布式事务
├── application/            # 业务应用层 (参考实现)
│   ├── configs/            # 配置文件 (config.example.yaml)
│   └── internal/
│       ├── domain/         # 领域层 (Entities, Repository Interfaces)
│       ├── application/    # 应用层 (Commands, Queries, Fx Module)
│       ├── infrastructure/ # 基础设施层 (Repo Implementations)
│       └── interfaces/     # 接口层 (HTTP Handlers)
└── tools/generator/        # 代码生成器 CLI
```

---

## 🚦 快速开始

详见 [QUICK_START.md](./QUICK_START.md)

### 简要步骤

```bash
# 1. 编译生成器
cd tools/generator && go build -o soliton-gen .

# 2. 生成领域模块
./soliton-gen domain Order

# 3. 在 main.go 中导入
import orderapp ".../application/order"
fx.New(orderapp.Module, ...)

# 4. 启动服务
go run ./cmd/main.go
```

**自动可用的 API：**
- `POST /api/orders` - 创建
- `GET /api/orders` - 列表
- `GET /api/orders/:id` - 获取
- `PUT /api/orders/:id` - 更新
- `DELETE /api/orders/:id` - 删除

---

## 📖 文档索引

| 文档 | 说明 |
|------|------|
| [QUICK_START.md](./QUICK_START.md) | 详细快速上手指南 |
| [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md) | 开发流程与最佳实践 |
| [GENERATOR_GUIDE.md](./GENERATOR_GUIDE.md) | 代码生成器使用说明 |

---

## 🏗 应用装配 (Dependency Injection)

使用 **Uber Fx** 管理生命周期，一行代码导入模块：

```go
func main() {
    fx.New(
        // 数据库
        fx.Provide(orm.NewGormDB),
        
        // 领域模块 (一行导入所有依赖)
        orderapp.Module,
        productapp.Module,
        
        // HTTP Handlers
        fx.Provide(http.NewOrderHandler),
        
        // 启动
        fx.Invoke(func(h *http.OrderHandler, db *gorm.DB) {
            orderapp.RegisterMigration(db)  // 自动建表
            
            r := gin.Default()
            h.RegisterRoutes(r)
            r.Run(":8080")
        }),
    ).Run()
}
```

---

## 🤝 贡献与支持
欢迎提交 Issue 和 PR 共同完善 Soliton-Go。
