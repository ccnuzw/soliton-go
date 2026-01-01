# Soliton-Go 快速开发指南

## 1. 📋 环境要求

- **Go**: 1.22+
- **IDE**: VS Code 或 GoLand

---

## 2. 🚀 编译生成器

```bash
cd tools/generator
go build -o soliton-gen .
```

---

## 3. ⚡ 一键生成领域模块

### 基础用法
```bash
./soliton-gen domain User
```

### 🆕 使用 --fields 指定字段（推荐）

```bash
# User - 用户管理
./soliton-gen domain User --fields "username,email,password_hash,phone,nickname,role:enum(admin|seller|customer),status:enum(active|inactive|banned)"

# Product - 商品管理
./soliton-gen domain Product --fields "name,sku,description:text,price:int64,original_price:int64,stock:int,category_id:uuid,status:enum(draft|active|inactive)"

# Order - 订单管理
./soliton-gen domain Order --fields "user_id:uuid,order_no,total_amount:int64,status:enum(pending|paid|shipped|completed|cancelled),receiver_name,receiver_phone,receiver_address:text"
```

### 支持的字段类型

| 类型 | 写法 | Go 类型 | GORM 标签 |
|------|------|---------|-----------|
| string | `field` | `string` | `size:255` |
| text | `field:text` | `string` | `type:text` |
| int | `field:int` | `int` | `not null` |
| int64 | `field:int64` | `int64` | `not null` |
| float | `field:float` | `float64` | - |
| bool | `field:bool` | `bool` | `default:false` |
| uuid | `field:uuid` | `string` | `size:36;index` |
| time | `field:time` | `time.Time` | `autoCreateTime` |
| **enum** | `field:enum(a\|b\|c)` | 自定义枚举类型 | `size:50` |

---

## 4. 📁 生成文件清单

每个领域模块生成 **9 个文件**：

```
domain/{name}/
├── {name}.go          # Entity + ID + Enum类型
├── repository.go      # Repository 接口
└── events.go          # 领域事件 (Created/Updated/Deleted)

application/{name}/
├── commands.go        # Create/Update/Delete Handlers
├── queries.go         # Get/List Handlers
├── dto.go             # Request/Response DTOs
└── module.go          # Fx 依赖注入模块

infrastructure/persistence/
└── {name}_repo.go     # Repository 实现 + 数据库迁移

interfaces/http/
└── {name}_handler.go  # HTTP CRUD Handler
```

---

## 5. 🏗 配置 main.go

```go
import (
    userapp "github.com/soliton-go/application/internal/application/user"
    productapp "github.com/soliton-go/application/internal/application/product"
    orderapp "github.com/soliton-go/application/internal/application/order"
    "github.com/soliton-go/application/internal/interfaces/http"
)

func main() {
    fx.New(
        // 数据库
        fx.Provide(orm.NewGormDB),
        
        // 领域模块 (一行导入所有依赖)
        userapp.Module,
        productapp.Module,
        orderapp.Module,
        
        // HTTP Handlers
        fx.Provide(http.NewUserHandler),
        fx.Provide(http.NewProductHandler),
        fx.Provide(http.NewOrderHandler),
        
        // 启动
        fx.Invoke(func(db *gorm.DB, userH *http.UserHandler, productH *http.ProductHandler, orderH *http.OrderHandler) {
            // 自动建表
            userapp.RegisterMigration(db)
            productapp.RegisterMigration(db)
            orderapp.RegisterMigration(db)
            
            // 注册路由
            r := gin.Default()
            userH.RegisterRoutes(r)
            productH.RegisterRoutes(r)
            orderH.RegisterRoutes(r)
            
            r.Run(":8080")
        }),
    ).Run()
}
```

---

## 6. 🏃 运行

```bash
go run ./cmd/main.go
```

**自动可用的 API：**

| 模块 | 端点 |
|------|------|
| User | `/api/users` (POST/GET/PUT/DELETE) |
| Product | `/api/products` (POST/GET/PUT/DELETE) |
| Order | `/api/orders` (POST/GET/PUT/DELETE) |

---

## 7. ⚡ 开发流程总结

```
1. soliton-gen domain Xxx --fields "..."  # 一条命令生成完整模块
2. main.go 导入 xxxapp.Module            # 一行注入所有依赖
3. go run ./cmd/main.go                   # 启动服务
```

**自动生成，无需手写：**
- ✅ Entity + Enum 类型
- ✅ Repository 接口和实现
- ✅ Commands/Queries/DTOs
- ✅ HTTP Handler
- ✅ 依赖注入模块
- ✅ 数据库迁移
