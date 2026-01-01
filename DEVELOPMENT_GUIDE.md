# Soliton-Go 快速开发指南

## 🚀 使用生成器开发新项目

### 步骤 1: 安装生成器

```bash
cd tools/generator
go build -o soliton-gen .

# 可选：添加到 PATH
sudo mv soliton-gen /usr/local/bin/
```

---

### 步骤 2: 生成领域模块

```bash
# 从项目根目录或 tools/generator 执行
./soliton-gen domain User
./soliton-gen domain Product
./soliton-gen domain Order
```

**生成文件：**
```
domain/{name}/
├── {name}.go          # Entity + ID
├── repository.go      # Repository Interface
└── events.go          # 领域事件

application/{name}/
├── commands.go        # Create/Update/Delete Handlers
├── queries.go         # Get/List Handlers
├── dto.go             # Request/Response
└── module.go          # Fx 依赖注入模块

infrastructure/persistence/
└── {name}_repo.go     # Repository 实现 + 迁移函数

interfaces/http/
└── {name}_handler.go  # HTTP CRUD Handler
```

---

### 步骤 3: 扩展 Entity 字段

编辑 `domain/{name}/{name}.go`：

```go
type Order struct {
    ddd.BaseAggregateRoot
    ID        OrderID `gorm:"primaryKey"`
    Name      string  `gorm:"size:255"`
    // 添加业务字段
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

---

### 步骤 4: 更新 DTO

编辑 `application/{name}/dto.go`：

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
    CreatedAt string `json:"created_at"`
}
```

---

### 步骤 5: 配置 main.go

```go
package main

import (
    "go.uber.org/fx"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    
    orderapp "github.com/soliton-go/application/internal/application/order"
    productapp "github.com/soliton-go/application/internal/application/product"
    "github.com/soliton-go/application/internal/interfaces/http"
    "github.com/gin-gonic/gin"
)

func main() {
    fx.New(
        // 数据库
        fx.Provide(func() *gorm.DB {
            db, _ := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
            return db
        }),
        
        // 领域模块
        orderapp.Module,
        productapp.Module,
        
        // HTTP Handlers
        fx.Provide(http.NewOrderHandler),
        fx.Provide(http.NewProductHandler),
        
        // 启动服务
        fx.Invoke(func(db *gorm.DB, orderHandler *http.OrderHandler) {
            // 迁移表
            orderapp.RegisterMigration(db)
            productapp.RegisterMigration(db)
            
            // 注册路由
            r := gin.Default()
            orderHandler.RegisterRoutes(r)
            
            // 启动
            r.Run(":8080")
        }),
    ).Run()
}
```

---

### 步骤 6: 运行

```bash
# 编译
go build -o app ./cmd/main.go

# 运行
./app
```

**API 端点自动可用：**
- `POST /api/orders` - 创建
- `GET /api/orders` - 列表
- `GET /api/orders/:id` - 获取
- `PUT /api/orders/:id` - 更新
- `DELETE /api/orders/:id` - 删除

---

## 📁 项目结构参考

```
my-project/
├── cmd/
│   └── main.go                    # 入口
├── configs/
│   └── config.yaml                # 配置文件
├── internal/
│   ├── domain/                    # 领域层
│   │   ├── user/
│   │   ├── product/
│   │   └── order/
│   ├── application/               # 应用层
│   │   ├── user/
│   │   ├── product/
│   │   └── order/
│   ├── infrastructure/            # 基础设施层
│   │   └── persistence/
│   └── interfaces/                # 接口层
│       └── http/
└── go.mod
```

---

## ⚡ 开发流程总结

```
1. soliton-gen domain Xxx     # 生成模块骨架
2. 编辑 xxx.go                # 添加业务字段
3. 编辑 dto.go                # 更新请求/响应
4. 编辑 commands.go           # 调整业务逻辑
5. main.go 导入 xxxapp.Module # 注入依赖
6. go run ./cmd/main.go       # 启动服务
```

**无需编写的代码：**
- Repository 实现 ✅ 自动生成
- HTTP Handler ✅ 自动生成
- 依赖注入 ✅ 自动生成
- 数据库迁移 ✅ 自动生成
- 领域事件 ✅ 自动生成
