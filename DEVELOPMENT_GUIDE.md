# 🚀 Soliton-Go 快速开发指南

## 一、环境准备

确保你的环境满足以下要求：
- **Go 1.22+**
- **Redis**（可选，用于分布式锁和事件总线）
- **数据库**（PostgreSQL、MySQL 或 SQLite）

---

## 二、使用代码生成器 (soliton-gen)

### 1. 编译生成器

```bash
cd tools/generator
go build -o soliton-gen main.go
```

### 2. 生成领域模块

假设你要创建一个 **订单 (Order)** 模块：

```bash
./soliton-gen domain Order
```

**自动生成的文件：**

| 文件路径 | 说明 |
|---------|------|
| `application/internal/domain/Order/Order.go` | 聚合根实体 |
| `application/internal/domain/Order/repository.go` | 仓储接口 |
| `application/internal/domain/Order/mapper.go` | SQL Mapper 接口 |
| `application/internal/infrastructure/persistence/Order_repo.go` | 仓储实现 |
| `application/internal/infrastructure/persistence/Order_mapper.go` | Mapper 实现 |

---

## 三、完整开发流程示例

以创建 **商品 (Product)** 模块为例：

### 步骤 1：生成基础代码

```bash
cd tools/generator
./soliton-gen domain Product
```

### 步骤 2：扩展实体字段

编辑 `application/internal/domain/Product/Product.go`：

```go
package Product

import "github.com/soliton-go/framework/ddd"

type ProductID string

func (id ProductID) String() string {
    return string(id)
}

type Product struct {
    ddd.BaseAggregateRoot
    ID          ProductID `gorm:"primaryKey"`
    Name        string    `gorm:"not null"`           // 添加字段
    Description string                                 // 添加字段
    Price       float64   `gorm:"not null"`           // 添加字段
    Stock       int       `gorm:"default:0"`          // 添加字段
}

func NewProduct(id, name, description string, price float64, stock int) *Product {
    return &Product{
        ID:          ProductID(id),
        Name:        name,
        Description: description,
        Price:       price,
        Stock:       stock,
    }
}

func (e *Product) GetID() ddd.ID {
    return e.ID
}
```

### 步骤 3：添加自定义仓储方法

编辑 `application/internal/domain/Product/repository.go`：

```go
package Product

import (
    "context"
    "github.com/soliton-go/framework/orm"
)

type ProductRepository interface {
    orm.Repository[*Product, ProductID]
    // 添加自定义方法
    FindByName(ctx context.Context, name string) ([]*Product, error)
    FindInStock(ctx context.Context) ([]*Product, error)
}
```

### 步骤 4：实现自定义方法

编辑 `application/internal/infrastructure/persistence/Product_repo.go`：

```go
package persistence

import (
    "context"
    "github.com/soliton-go/application/internal/domain/Product"
    "github.com/soliton-go/framework/orm"
    "gorm.io/gorm"
)

type ProductRepoImpl struct {
    *orm.GormRepository[*Product.Product, Product.ProductID]
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) Product.ProductRepository {
    return &ProductRepoImpl{
        GormRepository: orm.NewGormRepository[*Product.Product, Product.ProductID](db),
        db:             db,
    }
}

// 实现自定义方法
func (r *ProductRepoImpl) FindByName(ctx context.Context, name string) ([]*Product.Product, error) {
    var products []*Product.Product
    err := r.db.WithContext(ctx).Where("name LIKE ?", "%"+name+"%").Find(&products).Error
    return products, err
}

func (r *ProductRepoImpl) FindInStock(ctx context.Context) ([]*Product.Product, error) {
    var products []*Product.Product
    err := r.db.WithContext(ctx).Where("stock > 0").Find(&products).Error
    return products, err
}
```

### 步骤 5：创建应用层 Command Handler

创建 `application/internal/application/product/commands.go`：

```go
package productapp

import (
    "context"
    "github.com/google/uuid"
    "github.com/soliton-go/application/internal/domain/Product"
)

// Command 定义
type CreateProductCommand struct {
    Name        string
    Description string
    Price       float64
    Stock       int
}

// Handler 定义
type CreateProductHandler struct {
    repo Product.ProductRepository
}

func NewCreateProductHandler(repo Product.ProductRepository) *CreateProductHandler {
    return &CreateProductHandler{repo: repo}
}

func (h *CreateProductHandler) Handle(ctx context.Context, cmd CreateProductCommand) (*Product.Product, error) {
    p := Product.NewProduct(
        uuid.New().String(),
        cmd.Name,
        cmd.Description,
        cmd.Price,
        cmd.Stock,
    )
    if err := h.repo.Save(ctx, p); err != nil {
        return nil, err
    }
    return p, nil
}
```

### 步骤 6：注册到依赖注入容器

编辑 `application/cmd/server/main.go`：

```go
package main

import (
    "context"

    "github.com/soliton-go/framework/core/config"
    "github.com/soliton-go/framework/core/logger"
    "github.com/soliton-go/framework/orm"
    "github.com/soliton-go/framework/web"
    "go.uber.org/fx"

    userapp "github.com/soliton-go/application/internal/application/user"
    productapp "github.com/soliton-go/application/internal/application/product"  // 新增
    "github.com/soliton-go/application/internal/infrastructure/persistence"
)

func main() {
    fx.New(
        fx.Provide(
            // Framework modules
            config.NewConfig,
            logger.NewLogger,
            orm.NewGormDB,
            web.NewServer,

            // User module
            persistence.NewUserRepository,
            userapp.NewCreateUserHandler,

            // Product module (新增)
            persistence.NewProductRepository,
            productapp.NewCreateProductHandler,
        ),
        fx.Invoke(func(lc fx.Lifecycle, server *web.Server) {
            lc.Append(fx.Hook{
                OnStart: func(ctx context.Context) error {
                    go server.Run(":8080")
                    return nil
                },
                OnStop: func(ctx context.Context) error {
                    return nil
                },
            })
        }),
    ).Run()
}
```

---

## 四、使用高级特性

### 1. 使用分布式锁

```go
func (s *ProductService) UpdateStock(ctx context.Context, productID string, delta int) error {
    // 获取锁
    lock, err := s.locker.Obtain(ctx, "lock:product:"+productID, 10*time.Second)
    if err != nil {
        return fmt.Errorf("获取锁失败: %w", err)
    }
    defer lock.Release(ctx)

    // 安全更新库存
    product, _ := s.repo.Find(ctx, Product.ProductID(productID))
    product.Stock += delta
    return s.repo.Save(ctx, product)
}
```

### 2. 使用 Saga 分布式事务

```go
func (s *OrderService) CreateOrder(ctx context.Context, cmd CreateOrderCommand) error {
    saga := transaction.NewSaga()

    // 步骤1: 扣减库存
    saga.AddStep(
        "DeductStock",
        func(ctx context.Context) error {
            return s.productRepo.DeductStock(cmd.ProductID, cmd.Quantity)
        },
        func(ctx context.Context) error {
            return s.productRepo.AddStock(cmd.ProductID, cmd.Quantity) // 补偿
        },
    )

    // 步骤2: 创建订单
    saga.AddStep(
        "CreateOrder",
        func(ctx context.Context) error {
            order := order.NewOrder(cmd.UserID, cmd.ProductID, cmd.Quantity)
            return s.orderRepo.Save(ctx, order)
        },
        func(ctx context.Context) error {
            return s.orderRepo.Delete(ctx, order.ID) // 补偿
        },
    )

    return saga.Execute(ctx)
}
```

### 3. 使用 SQL Mapper 执行复杂查询

```go
func (r *ProductRepoImpl) FindTopSelling(ctx context.Context, limit int) ([]*Product.Product, error) {
    mapper := orm.NewGormMapper[Product.Product](r.db)
    return mapper.SelectList(ctx, `
        SELECT p.* FROM products p
        JOIN order_items oi ON p.id = oi.product_id
        GROUP BY p.id
        ORDER BY SUM(oi.quantity) DESC
        LIMIT ?
    `, limit)
}
```

### 4. 使用 CQRS 命令总线

```go
// 注册 Handler
bus := cqrs.NewCommandBus()
bus.Register(CreateProductCommand{}, func(ctx context.Context, cmd CreateProductCommand) error {
    return handler.Handle(ctx, cmd)
})

// 分发命令
err := bus.Dispatch(ctx, CreateProductCommand{
    Name:  "iPhone 15",
    Price: 7999.00,
    Stock: 100,
})
```

### 5. 发布领域事件

```go
// 在聚合根中添加事件
func NewProduct(id, name string, price float64) *Product {
    p := &Product{ID: ProductID(id), Name: name, Price: price}
    p.AddDomainEvent(ProductCreatedEvent{ProductID: id, Name: name})
    return p
}

// 订阅事件
eventBus.Subscribe(ctx, "ProductCreatedEvent", func(ctx context.Context, event ddd.DomainEvent) error {
    // 发送通知、更新搜索索引等
    return nil
})
```

---

## 五、运行项目

```bash
cd application
go run cmd/server/main.go
```

服务将在 **http://localhost:8080** 启动，GraphQL Playground 位于 **http://localhost:8080/query/playground**。

---

## 📋 开发流程速查表

| 步骤 | 命令/操作 |
|------|----------|
| 1. 生成领域代码 | `./soliton-gen domain <EntityName>` |
| 2. 扩展实体字段 | 编辑 `domain/<Name>/<Name>.go` |
| 3. 添加仓储方法 | 编辑 `domain/<Name>/repository.go` |
| 4. 实现仓储方法 | 编辑 `infrastructure/persistence/<Name>_repo.go` |
| 5. 创建 Command/Handler | 创建 `application/<name>/commands.go` |
| 6. 注册 DI | 编辑 `cmd/server/main.go` |
| 7. 运行测试 | `go run cmd/server/main.go` |

---

## 🔗 相关文档

- [README.md](./README.md) - 项目概述
- [QUICK_START.md](./QUICK_START.md) - 快速上手指南
- [GENERATOR_GUIDE.md](./GENERATOR_GUIDE.md) - 代码生成器详细指南
