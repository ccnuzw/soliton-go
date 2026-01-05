# Soliton-Go 开发指南

## 🚀 使用生成器快速开发

### 步骤 1: 编译生成器

```bash
cd tools/generator && go build -o soliton-gen .
```

### 步骤 2: 一键生成领域模块

```bash
# 使用 --fields 直接生成带完整字段的模块
./soliton-gen domain User --fields "username,email,password_hash,role:enum(admin|user),status:enum(active|banned)"
./soliton-gen domain Product --fields "name,price:int64,stock:int,status:enum(draft|active)"
./soliton-gen domain Order --fields "user_id:uuid,total:int64,status:enum(pending|paid|shipped)"
```

### 步骤 2.1: 生成 DDD 领域增强组件

```bash
# 领域值对象
./soliton-gen valueobject user EmailAddress --fields "value:string"

# 领域规格
./soliton-gen spec user ActiveUserSpec --target User

# 领域策略
./soliton-gen policy user PasswordPolicy --target User

# 自定义领域事件
./soliton-gen event user UserActivated --fields "user_id:uuid"

# 事件处理器（自动注入 module.go / main.go）
./soliton-gen event-handler user UserActivated
```

### 步骤 2.2: 安装依赖

```bash
GOWORK=off go mod tidy
```

### 步骤 2.3: 执行数据库迁移

```bash
GOWORK=off go run ./cmd/migrate
```

### 步骤 3: 配置 main.go

```go
fx.New(
    fx.Provide(orm.NewGormDB),
    
    // 一行导入模块（含 Repository、Commands、Queries）
    userapp.Module,
    productapp.Module,
    orderapp.Module,
    
    // HTTP Handlers
    fx.Provide(http.NewUserHandler),
    fx.Provide(http.NewProductHandler),
    fx.Provide(http.NewOrderHandler),
    
    fx.Invoke(func(db *gorm.DB, h1 *http.UserHandler, h2 *http.ProductHandler, h3 *http.OrderHandler) error {
        // 自动建表
        if err := userapp.RegisterMigration(db); err != nil {
            return err
        }
        if err := productapp.RegisterMigration(db); err != nil {
            return err
        }
        if err := orderapp.RegisterMigration(db); err != nil {
            return err
        }
        
        // 注册路由
        r := gin.Default()
        h1.RegisterRoutes(r)
        h2.RegisterRoutes(r)
        h3.RegisterRoutes(r)
        return r.Run(":8080")
    }),
).Run()
```

### 步骤 4: 启动

```bash
GOWORK=off go run ./cmd/main.go
```

---

## 📁 项目结构

```
my-project/
├── cmd/main.go                           # 入口
├── cmd/migrate/main.go                   # 迁移入口
├── configs/config.yaml                   # 配置
└── internal/
    ├── domain/                           # 领域层（生成）
    │   ├── user/
    │   ├── product/
    │   └── order/
    ├── application/                      # 应用层（生成）
    │   ├── user/
    │   ├── product/
    │   └── order/
    ├── infrastructure/persistence/       # 持久层（生成）
    └── interfaces/http/                  # 接口层（生成）
```

---

## ⚡ 开发效率

| 传统开发 | 使用生成器 |
|---------|-----------|
| 手写 Entity + Enum | ✅ `--fields` 自动生成 |
| 手写 Repository | ✅ 自动生成 |
| 手写 Commands/Queries | ✅ 自动生成 |
| 手写 DTO | ✅ 自动生成 |
| 手写 HTTP Handler | ✅ 自动生成 |
| 手写依赖注入 | ✅ 自动生成 |
| 手写数据库迁移 | ✅ 自动生成 |

---

## 🔥 高级功能

### 分布式锁

```go
lock, _ := locker.Obtain(ctx, "lock:order:"+id, 10*time.Second)
defer lock.Release(ctx)
```

### 领域事件

```go
bus.Subscribe(ctx, "order.created", func(ctx context.Context, e ddd.DomainEvent) error {
    // 处理事件
    return nil
})
```

### Saga 分布式事务

```go
saga := transaction.NewSaga()
saga.AddStep("DeductStock", deductFunc, compensateFunc)
saga.AddStep("CreateOrder", createFunc, deleteFunc)
saga.Execute(ctx)
```
