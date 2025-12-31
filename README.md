# Soliton-Go 分布式全栈开发框架

**Soliton-Go** 是一个基于 Go 语言构建的模块化、高性能后端开发脚手架。它借鉴了现代企业级框架（如 Spring Boot）的设计理念，专为构建**分布式系统**和**领域驱动设计 (DDD)** 应用而设计。

## 🚀 核心特性

*   **领域驱动设计 (DDD)**: 提供标准的 AggregateRoot, Entity, ValueObject 基类和分层架构规范，助您构建业务逻辑清晰的应用。
*   **模块化架构**: 采用 Go Workspace 管理，核心框架 (`framework`) 与业务应用 (`application`) 严格解耦。
*   **分布式能力**:
    *   **分布式锁**: 基于 Redis 的高可用分布式锁，支持自动续期 (Watchdog)，防止死锁。
    *   **事件驱动**: 内置 Watermill 事件总线，支持基于 Redis Stream 或 NATS 的可靠消息投递。
    *   **分布式事务**: 提供 **Saga 模式** 的轻量级编排器，支持补偿机制 (Compensating Transaction)，轻松处理跨服务一致性。
*   **CQRS 模式**: 内置 Command Bus 和 Query Bus，轻松实现读写分离架构。
*   **极致开发体验**:
    *   **依赖注入**: 全项目集成 **Uber Fx**，零手动组装，享受自动装配的便利。
    *   **泛型仓储**: 提供 `GormRepository[T, ID]`，一行代码获得 CRUD 能力。
    *   **SQL Mapper**: 提供类似 MyBatis 的 `SQLMapper[T]`，支持原生 SQL 查询与映射，解决复杂查询难题。
    *   **通用服务**: 提供 `BaseService[T, ID]`，封装标准业务逻辑，减少样板代码。
    *   **GraphQL 支持**: 开箱即用的 GraphQL 服务器 (基于 99designs/gqlgen)。
*   **代码生成器**: 内置 CLI 工具 `soliton-gen`，一键生成 **Domain + Infrastructure** 全栈代码，实现真正的快速开发。

## 🛠 技术栈

| 模块 | 技术选型 | 说明 |
| :--- | :--- | :--- |
| **语言** | Go 1.22+ | 需要支持 Generic 和 Workspace |
| **Web 框架** | [Gin](https://github.com/gin-gonic/gin) | 高性能 HTTP 路由 |
| **API 接口** | [gqlgen](https://github.com/99designs/gqlgen) | Schema-first GraphQL |
| **ORM** | [GORM](https://gorm.io/) | 强大的 ORM 库 |
| **依赖注入** | [Uber Fx](https://go.uber.org/fx) | 模块化依赖注入容器 |
| **配置管理** | [Viper](https://github.com/spf13/viper) | 支持多格式、热加载 |
| **日志** | [Zap](https://github.com/uber-go/zap) | 高性能日志库 |
| **消息队列** | [Watermill](https://watermill.io/) | 通用事件驱动库 (Redis/Kafka/NATS) |
| **分布式锁** | [Redislock](https://github.com/bsm/redislock) | Redis 分布式锁实现 |

---

## 📂 项目结构

```text
soliton-go/
├── go.work                 # Go Workspace 定义文件
├── framework/              # 核心框架层 (通用库，可独立发版)
│   ├── core/               # 基础模块 (Config, Log, DI)
│   ├── ddd/                # DDD 原语 (AggregateRoot, Entity)
│   ├── orm/                # GORM 封装与泛型 Repository
│   │   └── mapper.go       # [新增] SQLMapper (MyBatis 风格原生 SQL 支持)
│   ├── service/            # [新增] 通用业务服务层 (BaseService)
│   ├── transaction/        # [新增] 分布式事务 (Saga 模式实现)
│   ├── web/                # Gin Server & GraphQL 封装
│   ├── event/              # 事件总线 (Watermill 实现)
│   ├── lock/               # 分布式锁 (Redis 实现)
│   └── cqrs/               # CQRS 命令/查询总线
├── application/            # 业务应用层 (参考实现)
│   ├── cmd/                # 程序入口
│   │   └── server/         # API Server 入口
│   ├── configs/            # 配置文件
│   └── internal/           # 业务逻辑 (遵循 Clean Architecture)
│       ├── domain/         # [内层] 领域层 (Entities, Repository Interfaces)
│       ├── application/    # [中层] 应用层 (Use Cases, Command Handlers)
│       ├── infrastructure/ # [外层] 基础设施层 (Repo Impls, External APIs)
│       └── interfaces/     # [接入层] 接口层 (GraphQL Resolvers, HTTP Handlers)
└── tools/                  # 开发者工具
    └── generator/          # 代码生成器 CLI
```

---

## 🚦 快速开始

### 环境依赖
*   Go 1.22 或更高版本
*   Redis (用于测试分布式锁和事件总线，可选)
*   数据库 (PostgreSQL, MySQL 或 SQLite)

### 1. 启动示例应用
`application` 模块包含了一个完整的用户管理示例。

```bash
cd application

# 下载依赖 (利用 go workspace，会自动链接本地 framework)
go mod tidy

# 启动服务器
go run cmd/server/main.go
```
服务器默认监听 `:8080` 端口。

### 2. 使用代码生成器
`tools/generator` 提供了一个 CLI 工具来加速开发。

```bash
cd tools/generator

# 编译生成器
go build -o soliton-gen main.go

# 生成名为 "Order" 的领域对象
./soliton-gen domain Order
```

---

## 📖 核心模块开发指南

### 1. 定义领域实体 (Domain Layer)
在 `internal/domain` 中定义你的聚合根。继承 `ddd.BaseAggregateRoot` 以获得事件能力。

```go
// internal/domain/order/order.go
package order

import "github.com/soliton-go/framework/ddd"

type Order struct {
    ddd.BaseAggregateRoot
    ID     string `gorm:"primaryKey"`
    Amount float64
}

func NewOrder(id string, amount float64) *Order {
    o := &Order{ID: id, Amount: amount}
    // 记录领域事件
    o.AddDomainEvent(NewOrderCreatedEvent(id))
    return o
}
```

### 2. 定义仓储接口 (Domain Layer)
```go
// internal/domain/order/repository.go
type OrderRepository interface {
    orm.Repository[*Order, string] // 继承泛型接口，自动获得 CRUD
    FindByAmount(ctx context.Context, amount float64) ([]*Order, error)
}
```

### 3. 实现基础设施 (Infrastructure Layer)
在 `internal/infrastructure` 中实现仓储。

```go
// internal/infrastructure/persistence/order_repo.go
package persistence

type OrderRepo struct {
    *orm.GormRepository[*order.Order, string]
}

func NewOrderRepo(db *gorm.DB) order.OrderRepository {
    return &OrderRepo{
        GormRepository: orm.NewGormRepository[*order.Order, string](db),
    }
}
```

### 4. 编写应用逻辑 (Application Layer)
使用 CQRS 模式编写 Command Handler。

```go
// internal/application/order/command.go
type CreateOrderHandler struct {
    repo order.OrderRepository
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrderCommand) error {
    order := order.NewOrder(cmd.ID, cmd.Amount)
    return h.repo.Save(ctx, order) // Save 自动触发领域事件分发 (需配置中间件)
}
```

### 5. 并发控制 (分布式锁)
在需要并发保护的业务逻辑中注入 `lock.Locker`。

```go
func (s *Service) ProcessPayment(ctx context.Context, orderID string) error {
    // 获取分布式锁，TTL 10秒，带 Watchdog 自动续期
    lock, err := s.locker.Obtain(ctx, "lock:order:"+orderID, 10*time.Second)
    if err != nil {
        return err
    }
    defer lock.Release(ctx)

    // 安全的业务逻辑...
    return nil
}

### 6. 分布式事务 (Saga 模式)
使用 `transaction.SagaOrchestrator` 编排跨服务/跨步骤的原子操作。

```go
func (s *Service) CreateOrderWithSaga(ctx context.Context, orderID string) error {
    // 创建 Saga 编排器
    saga := transaction.NewSaga()

    // 添加步骤：扣减库存
    saga.AddStep(
        "DeductInventory",
        func(ctx context.Context) error { return s.inventoryClient.Deduct(orderID) }, // 正向操作
        func(ctx context.Context) error { return s.inventoryClient. add(orderID) },   // 补偿操作
    )

    // 添加步骤：保存订单
    saga.AddStep(
        "SaveOrder",
        func(ctx context.Context) error { return s.repo.Save(order) },
        func(ctx context.Context) error { return s.repo.Delete(order) },
    )

    // 执行事务 (失败自动回滚)
    return saga.Execute(ctx)
}
```

---

## 🏗 应用装配 (Dependency Injection)

我们使用 **Uber Fx** 来管理生命周期。在 `main.go` 中：

```go
func main() {
    fx.New(
        // 1. 提供框架模块
        fx.Provide(
            config.NewConfig,
            logger.NewLogger,
            orm.NewGormDB,
            lock.NewRedisLocker,
        ),
        // 2. 提供业务模块
        fx.Provide(
            persistence.NewOrderRepo,
            application.NewCreateOrderHandler,
        ),
        // 3. 启动入口
        fx.Invoke(func(server *web.Server) {
            go server.Run(":8080")
        }),
    ).Run()
}
```

---

## 🤝 贡献与支持
欢迎提交 Issue 和 PR 共同完善 Soliton-Go。
