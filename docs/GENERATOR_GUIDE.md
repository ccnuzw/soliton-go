# Soliton-Go 代码生成器使用指南

## 📦 安装

```bash
cd tools/generator
go build -o soliton-gen .
```

---

## 🚀 命令列表

| 命令 | 说明 |
|------|------|
| `init` | 初始化新项目 |
| `domain` | 生成领域模块 (Entity/Repo/Events/Handler) |
| `service` | 生成应用服务 (跨领域业务逻辑) |

---

## init - 初始化项目

```bash
./soliton-gen init my-project
./soliton-gen init my-project --module github.com/myorg/my-project
```

**生成内容：** `cmd/main.go`, `configs/`, `internal/`, `go.mod`, `Makefile`, `README.md`

---

## domain - 生成领域模块

```bash
./soliton-gen domain User
./soliton-gen domain User --fields "username,email,status:enum(active|inactive)"
./soliton-gen domain User --fields "..." --force  # 强制覆盖
./soliton-gen domain User --fields "..." --wire   # 自动接入 main.go
```

### --wire 自动接线
使用 `--wire` 标志时，生成器会自动修改 `main.go`：
- 取消注释所需 imports（gorm、module、handler）
- 取消注释 Module 注册
- 取消注释 Handler Provider
- 取消注释路由和迁移注册

> **注意**: `--wire` 仅在 main.go 保持 init 模板结构时生效。如已手动修改，请手动接线。

### 字段类型
| 类型 | 格式 | 示例 |
|------|------|------|
| string | `field` | `username` |
| int64 | `field:int64` | `price:int64` |
| text | `field:text` | `description:text` |
| uuid | `field:uuid` | `user_id:uuid` |
| time? | `field:time?` | `last_login_at:time?` |
| enum | `field:enum(a\|b)` | `status:enum(active\|banned)` |

### 生成文件 (9个)
- `domain/{name}/` - 实体 + Repository + Events
- `application/{name}/` - Commands + Queries + DTO + Module
- `infrastructure/persistence/{name}_repo.go`
- `interfaces/http/{name}_handler.go`

---

## service - 生成应用服务

用于生成跨领域的业务编排服务。

```bash
./soliton-gen service OrderService
./soliton-gen service OrderService --methods "CreateOrder,CancelOrder,GetUserOrders"
```

### 生成文件 (2个)
- `application/services/{name}_service.go` - 服务结构和方法
- `application/services/{name}_dto.go` - 请求/响应 DTO

📖 **详细文档**: [Service 应用服务使用指南](./docs/SERVICE_GUIDE.md)

---

## 🔄 修改已生成代码

| 场景 | 推荐方式 |
|------|----------|
| 小改动 | 手动编辑 |
| 大改动 | `--force` 重新生成 |

```bash
./soliton-gen domain User --fields "..." --force
```

---

## 🎯 完整开发流程

```bash
# 1. 初始化项目
./soliton-gen init my-shop && cd my-shop

# 2. 生成领域模块 (--wire 自动接入)
soliton-gen domain User --fields "username,email,role:enum(admin|customer)" --wire
soliton-gen domain Product --fields "name,price:int64,stock:int" --wire
soliton-gen domain Order --fields "user_id:uuid,total:int64,status:enum(pending|paid)" --wire

# 3. 生成跨领域服务
soliton-gen service OrderService --methods "CreateOrder,CancelOrder"

# 4. 运行（在 monorepo 中需要 GOWORK=off）
GOWORK=off go mod tidy && GOWORK=off go run ./cmd/main.go
```

> **Monorepo 提示**: 如果在包含 `go.work` 的 monorepo 中运行，请使用 `GOWORK=off` 前缀。
