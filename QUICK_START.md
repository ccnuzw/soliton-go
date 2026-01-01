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
./soliton-gen domain User --fields "username,email,password_hash,role:enum(admin|seller|customer),status:enum(active|inactive|banned)"

# Product - 商品管理
./soliton-gen domain Product --fields "name,sku,price:int64,stock:int,status:enum(draft|active|inactive)"

# Order - 订单管理
./soliton-gen domain Order --fields "user_id:uuid,order_no,total_amount:int64,status:enum(pending|paid|shipped|completed|cancelled)"
```

### 🔄 修改字段后重新生成

```bash
# 使用 --force 强制覆盖已存在的文件
./soliton-gen domain User --fields "username,email,age:int,status:enum(active|banned)" --force
```

---

## 4. 📋 支持的字段类型

| 类型 | 写法 | Go 类型 |
|------|------|---------|
| string | `field` | `string` |
| text | `field:text` | `string` |
| int | `field:int` | `int` |
| int64 | `field:int64` | `int64` |
| uuid | `field:uuid` | `string` |
| **enum** | `field:enum(a\|b\|c)` | 枚举类型 |

---

## 5. 📁 生成文件清单 (9个)

```
domain/{name}/
├── {name}.go          # Entity + Enum
├── repository.go      # Repository 接口
└── events.go          # 领域事件

application/{name}/
├── commands.go        # Create/Update/Delete
├── queries.go         # Get/List
├── dto.go             # Request/Response
└── module.go          # Fx 模块

infrastructure/persistence/
└── {name}_repo.go     # Repository 实现

interfaces/http/
└── {name}_handler.go  # HTTP Handler
```

---

## 6. 🏗 配置 main.go

```go
fx.New(
    fx.Provide(orm.NewGormDB),
    
    // 一行导入模块
    userapp.Module,
    productapp.Module,
    orderapp.Module,
    
    // HTTP Handlers
    fx.Provide(http.NewUserHandler),
    fx.Provide(http.NewProductHandler),
    fx.Provide(http.NewOrderHandler),
    
    fx.Invoke(func(db *gorm.DB, h1 *http.UserHandler, h2 *http.ProductHandler, h3 *http.OrderHandler) {
        // 自动建表
        userapp.RegisterMigration(db)
        productapp.RegisterMigration(db)
        orderapp.RegisterMigration(db)
        
        // 注册路由
        r := gin.Default()
        h1.RegisterRoutes(r)
        h2.RegisterRoutes(r)
        h3.RegisterRoutes(r)
        r.Run(":8080")
    }),
).Run()
```

---

## 7. 🏃 运行

```bash
go run ./cmd/main.go
```

**自动可用的 API：**

| 模块 | 端点 |
|------|------|
| User | `/api/users` |
| Product | `/api/products` |
| Order | `/api/orders` |

---

## 8. ⚡ 开发流程

```
1. soliton-gen domain Xxx --fields "..."  # 生成
2. main.go 导入 xxxapp.Module            # 注入
3. go run ./cmd/main.go                   # 启动

# 修改字段后
4. soliton-gen domain Xxx --fields "..." --force  # 重新生成
```
