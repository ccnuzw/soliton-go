# Soliton-Go 快速上手指南

## 1. 📋 环境要求

*   **Go**: 版本 **1.22** 或更高
*   **Git**: 版本控制
*   **IDE**: 推荐 VS Code 或 GoLand

---

## 2. 🚀 一键生成领域模块

### 2.1 编译生成器
```bash
cd tools/generator
go build -o soliton-gen .
```

### 2.2 生成模块
```bash
./soliton-gen domain Order
```

**输出：**
```
🚀 Generating domain: Order

📦 Domain Layer
   [NEW] order.go
   [NEW] repository.go
   [NEW] events.go

🔧 Infrastructure Layer
   [NEW] order_repo.go

⚙️ Application Layer
   [NEW] commands.go
   [NEW] queries.go
   [NEW] dto.go

🌐 Interfaces Layer
   [NEW] order_handler.go

📌 Fx Module
   [NEW] module.go

✅ Domain generation complete!
```

---

## 3. 👨‍💻 扩展生成的代码

### 3.1 添加 Entity 字段
编辑 `domain/order/order.go`：

```go
type Order struct {
    ddd.BaseAggregateRoot
    ID        OrderID `gorm:"primaryKey"`
    Name      string  `gorm:"size:255"`
    // 新增字段
    Amount    int64   `gorm:"not null"`
    Status    string  `gorm:"size:50;default:'pending'"`
    CreatedAt time.Time
}

func NewOrder(id, name string, amount int64) *Order {
    e := &Order{
        ID:        OrderID(id),
        Name:      name,
        Amount:    amount,
        Status:    "pending",
        CreatedAt: time.Now(),
    }
    e.AddDomainEvent(NewOrderCreatedEvent(id))
    return e
}
```

### 3.2 更新 DTO
编辑 `application/order/dto.go`：

```go
type CreateOrderRequest struct {
    Name   string `json:"name" binding:"required"`
    Amount int64  `json:"amount" binding:"required,gt=0"`
}

type OrderResponse struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Amount    int64  `json:"amount"`
    Status    string `json:"status"`
}
```

### 3.3 更新 Commands
编辑 `application/order/commands.go`，同步新字段。

---

## 4. 🏗 配置 main.go

```go
package main

import (
    "go.uber.org/fx"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    
    orderapp "github.com/soliton-go/application/internal/application/order"
    "github.com/soliton-go/application/internal/interfaces/http"
)

func main() {
    fx.New(
        // 数据库
        fx.Provide(func() *gorm.DB {
            db, _ := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
            return db
        }),
        
        // 领域模块 (一行导入所有依赖)
        orderapp.Module,
        
        // HTTP Handler
        fx.Provide(http.NewOrderHandler),
        
        // 启动
        fx.Invoke(func(db *gorm.DB, handler *http.OrderHandler) {
            // 自动建表
            orderapp.RegisterMigration(db)
            
            // 注册路由
            r := gin.Default()
            handler.RegisterRoutes(r)
            r.Run(":8080")
        }),
    ).Run()
}
```

---

## 5. 🏃 运行

```bash
go run ./cmd/main.go
```

**自动可用的 API：**

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/orders` | 创建订单 |
| GET | `/api/orders` | 订单列表 |
| GET | `/api/orders/:id` | 获取订单 |
| PUT | `/api/orders/:id` | 更新订单 |
| DELETE | `/api/orders/:id` | 删除订单 |

---

## 6. ⚡ 开发流程总结

```
1. soliton-gen domain Xxx     # 生成模块骨架 (9个文件)
2. 编辑 xxx.go                # 添加业务字段
3. 编辑 dto.go                # 更新请求/响应
4. 编辑 commands.go           # 调整业务逻辑
5. main.go 导入 xxxapp.Module # 一行注入依赖
6. go run ./cmd/main.go       # 启动服务
```

**无需手写：**
- ✅ Repository 实现
- ✅ HTTP Handler
- ✅ 依赖注入配置
- ✅ 数据库迁移
- ✅ 领域事件注册

---

## 7. 🔥 高级功能

### 分布式锁
```go
lock, _ := locker.Obtain(ctx, "lock:order:"+id, 10*time.Second)
defer lock.Release(ctx)
// 安全的业务逻辑
```

### 领域事件订阅
```go
bus.Subscribe(ctx, "order.created", func(ctx context.Context, e ddd.DomainEvent) error {
    event := e.(*order.OrderCreatedEvent)
    // 发送通知...
    return nil
})
```

### Saga 分布式事务
```go
saga := transaction.NewSaga()
saga.AddStep("DeductStock", deductFunc, compensateFunc)
saga.AddStep("CreateOrder", createFunc, deleteFunc)
saga.Execute(ctx)  // 失败自动回滚
```

---

## 8. ❓ 常见问题

**Q: `go mod tidy` 报错？**
A: 检查 `go.work` 是否包含所有子模块。

**Q: 如何切换到 PostgreSQL？**
A: 修改 `configs/config.yaml`：
```yaml
database:
  driver: postgres
  postgres:
    host: localhost
    port: 5432
    user: postgres
    database: myapp
```

---

更多详情请参阅 [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md)
