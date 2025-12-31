# Soliton-Go 超详细快速上手指南

欢迎来到 **Soliton-Go** 的世界！本指南将手把手带您从零开始，搭建开发环境、运行示例应用，并亲手开发一个新的业务模块。

---

## 1. 📋 环境准备

在开始之前，请确保您的开发环境满足以下要求：

*   **Go**: 版本 **1.22** 或更高（必需）。
    *   检查命令: `go version`
*   **Git**: 用于版本控制（必需）。
*   **Docker** (推荐): 用于快速启动 Redis 和 数据库。
*   **IDE**: 推荐使用 **VS Code** (配合 Go 插件) 或 **GoLand**。

---

## 2. 🚀 初始化项目

### 2.1 获取代码
假设您已经获取了本项目代码，进入项目根目录：
```bash
cd soliton-go
```

### 2.2 理解工作区 (Workspace)
本项目使用 **Go Workspace** 模式，这意味着你不需要 `go get` 远程仓库就可以直接引用本地模块。
查看 `go.work` 文件，确认包含以下目录：
```go
use (
    ./application
    ./framework
    ./tools/generator
)
```

### 2.3 整理依赖
在根目录下执行以下命令，确保所有子模块的依赖都已就绪：
```bash
# 进入 application 目录整理依赖
cd application
go mod tidy

# 进入 framework 目录整理依赖
cd ../framework
go mod tidy

# 回到根目录
cd ..
```

---

## 3. 🏃 运行示例应用 (User Service)

`application` 目录是一个完整的参考实现，包含了一个用户管理模块。

### 3.1 启动服务器
```bash
cd application
go run cmd/server/main.go
```
*   **成功标志**: 当你看到类似 `[GIN-debug] Listening and serving HTTP on :8080` 的日志时，说明启动成功。
*   **注意**: 默认配置下，应用使用 **SQLite 内存数据库**，无需安装额外的数据库服务。

### 3.2 验证 API (GraphQL Playgound)
打开浏览器访问: [http://localhost:8080/query/playground](http://localhost:8080/query/playground)

在 playground 中输入以下 GraphQL 查询（注意：实际接口需要根据生成的 Schema 定义）：
> *注意：目前示例代码仅展示了后端逻辑结构，GraphQL Schema 还需要您根据业务定义 `.graphql` 文件并使用 gqlgen 生成。*

作为替代，我们检查 HTTP 服务是否存活：
```bash
curl http://localhost:8080/health
# (如果未配置健康检查接口，可能会返回 404，这是正常的，取决于 gin 路由配置)
```

---

## 4. 🛠 使用代码生成器 (Soliton-Gen)

我们将编译并使用 CLI 工具来减少重复劳动。

### 4.1 编译工具
```bash
cd tools/generator
go build -o soliton-gen main.go
```
现在你会看到一个 `soliton-gen` 可执行文件。

### 4.2 使用工具生成代码
假设我们要开发一个 **"Product" (商品)** 模块。

```bash
# 生成 Product 领域对象
./soliton-gen domain Product
```
*输出预览:*
```text
Created internal/domain/Product.go
Created internal/domain/Product_repo.go
```
*(注：当前版本生成器通过 Stdout 模拟输出，后续可扩展为生成真实文件)*

---

## 5. 👨‍💻 实战：手动开发 "Product" 模块

代码生成器固然好用，但手动实现一遍能帮您深入理解 DDD 架构。我们将构建一个简单的商品创建功能。

### 5.1 领域层 (Domain Layer) - 不依赖任何外部框架
创建文件: `application/internal/domain/product/entity.go`

```go
package product

import "github.com/soliton-go/framework/ddd"

// 聚合根
type Product struct {
    ddd.BaseAggregateRoot
    ID    string `gorm:"primaryKey"`
    Name  string
    Price float64
}

func NewProduct(id string, name string, price float64) *Product {
    return &Product{
        ID:    id,
        Name:  name,
        Price: price,
    }
}

func (p *Product) GetID() ddd.ID {
    return ddd.ID(p.ID) // 简单转换
}
```

创建文件: `application/internal/domain/product/repository.go`

```go
package product

import (
    "context"
    "github.com/soliton-go/framework/orm"
)

// 定义接口
type ProductRepository interface {
    orm.Repository[*Product, string] // 继承通用 CRUD 接口
}
```

### 5.2 基础设施层 (Infrastructure Layer) - 实现数据持久化
创建文件: `application/internal/infrastructure/persistence/product_repo.go`

```go
package persistence

import (
    "github.com/soliton-go/application/internal/domain/product"
    "github.com/soliton-go/framework/orm"
    "gorm.io/gorm"
)

type ProductRepoImpl struct {
    *orm.GormRepository[*product.Product, string] // 复用 GORM 实现
}

func NewProductRepo(db *gorm.DB) product.ProductRepository {
    return &ProductRepoImpl{
        GormRepository: orm.NewGormRepository[*product.Product, string](db),
    }
}
```

### 5.3 应用层 (Application Layer) - 编排业务逻辑 (CQRS)
创建文件: `application/internal/application/product/commands.go`

```go
package productapp

import (
    "context"
    "github.com/soliton-go/application/internal/domain/product"
)

type CreateProductCommand struct {
    ID    string
    Name  string
    Price float64
}

type CreateProductHandler struct {
    repo product.ProductRepository
}

func NewCreateProductHandler(repo product.ProductRepository) *CreateProductHandler {
    return &CreateProductHandler{repo: repo}
}

func (h *CreateProductHandler) Handle(ctx context.Context, cmd CreateProductCommand) error {
    p := product.NewProduct(cmd.ID, cmd.Name, cmd.Price)
    return h.repo.Save(ctx, p)
}
```

### 5.4 接入层 (Configuration) - 依赖注入装配
最后，打开 `application/cmd/server/main.go`，将新模块注册到 DI 容器中。

```go
func main() {
    fx.New(
        fx.Provide(
            // ... 原有配置 ...
            
            // 注册 Product 模块
            persistence.NewProductRepo,
            productapp.NewCreateProductHandler,
        ),
        // ...
    ).Run()
}
```

---

## 6. 🔥 高级功能指南

### 6.1 使用分布式锁 (Redis Lock)
在集群环境下，确保同一时刻只有一个线程操作同一资源。

**前置条件**: 确保配置文件中已配置 Redis 地址。

```go
type InventoryService struct {
    locker lock.Locker
}

func (s *InventoryService) DeductStock(ctx context.Context, productID string) error {
    // 1. 获取锁 (锁名: "lock:stock:{id}", 过期时间: 5秒)
    // 框架会自动启动 Watchdog 为业务执行时间超过 5 秒的任务续期
    mutex, err := s.locker.Obtain(ctx, "lock:stock:"+productID, 5*time.Second)
    if err != nil {
        return fmt.Errorf("系统繁忙，请稍后再试") // 获取锁失败
    }
    // 2. 确保释放锁
    defer mutex.Release(ctx)

    // 3. 执行业务逻辑 (扣减库存)
    // ...
    return nil
}
```

### 6.2 发布与订阅领域事件 (Event Bus)
解耦业务逻辑，例如：创建订单后 -> 发送邮件。

**发布事件 (在聚合根中)**:
```go
func NewOrder(id string) *Order {
    o := &Order{ID: id}
    o.AddDomainEvent(OrderCreatedEvent{OrderID: id}) // 暂存事件
    return o
}
```

**持久化时分发 (在 Repository 中)**:
`GormRepository` 的 `Save` 方法可以扩展逻辑，在事务提交成功后，从聚合根中 `PullDomainEvents()` 并通过 `EventBus.Publish()` 发送出去。

---

## 7. ❓ 常见问题

**Q: 为什么 `go mod tidy` 报错？**
A: 请检查是否在项目根目录正确配置了 `go.work`。如果未配置 workspace，需要在 `application/go.mod` 中使用 `replace` 指向本地 framework 目录。

**Q: 如何切换数据库到 MySQL/PostgreSQL？**
A: 修改 `application/configs/config.yaml` (需创建) 或环境变量：
```yaml
database:
  driver: postgres
  dsn: "host=localhost user=user password=pass dbname=soliton port=5432 sslmode=disable"
```

祝您使用愉快！如有问题，请查阅源码或提交 Issue。
