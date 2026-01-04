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

## 3. 🎨 使用 Web GUI（推荐）

### 启动可视化界面

```bash
# 在项目目录或空目录中运行
./soliton-gen serve
```

访问 http://127.0.0.1:3000 即可使用图形界面：

- **初始化项目**：可视化配置项目信息
- **生成领域**：拖拽式字段编辑器，支持预览
- **生成服务**：可视化方法配置

**优势：**
- ✨ 无需记忆命令参数
- 👁️ 生成前预览代码
- 📖 详细的操作提示
- 🌐 中英双语界面

详细使用说明：[Web GUI 使用指南](./WEB_GUI_GUIDE.md)

---

## 4. ⚡ 使用命令行（传统方式）

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

### 🆕 启用软删除

```bash
# 使用 --soft-delete 添加 DeletedAt 字段
./soliton-gen domain User --fields "username,email" --soft-delete --wire
```

### 🔄 修改字段后重新生成

```bash
# 使用 --force 强制覆盖已存在的文件
./soliton-gen domain User --fields "username,email,age:int,status:enum(active|banned)" --force
```

---

## 5. 📋 支持的字段类型

| 类型 | 写法 | Go 类型 |
|------|---------|---------|
| string | `field` | `string` |
| text | `field:text` | `string` |
| int | `field:int` | `int` |
| int64 | `field:int64` | `int64` |
| uuid | `field:uuid` | `string` |
| **enum** | `field:enum(a\|b\|c)` | 枚举类型 |

---

## 6. 📁 生成文件清单 (9个)

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

## 7. 🏗 配置 main.go

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

---

## 8. 🏃 运行

```bash
GOWORK=off go run ./cmd/main.go
```

> **提示**: 生成的 `Makefile` 默认 `GOWORK=off`，需要时可 `GOWORK=on make run`。

**自动可用的 API：**

| 模块 | 端点 |
|------|------|
| User | `/api/users` |
| Product | `/api/products` |
| Order | `/api/orders` |

**🆕 分页查询：**
```bash
# 获取第1页，每页20条
curl "http://localhost:8080/api/users?page=1&page_size=20"
```

响应包含分页信息：
```json
{
  "items": [...],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "total_pages": 5
}
```

---

## 9. ⚡ 开发流程

```
# 方式1：使用 Web GUI（推荐）
1. soliton-gen serve                    # 启动 Web 界面
2. 访问 http://127.0.0.1:3000          # 可视化操作
3. GOWORK=off go run ./cmd/main.go      # 启动

# 方式2：使用命令行
1. soliton-gen domain Xxx --fields "..."  # 生成
2. main.go 导入 xxxapp.Module            # 注入
3. GOWORK=off go run ./cmd/main.go         # 启动

# 修改字段后
4. soliton-gen domain Xxx --fields "..." --force  # 重新生成
```
