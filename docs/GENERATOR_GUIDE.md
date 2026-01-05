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
| `domain list` | 🆕 列出所有领域模块 |
| `domain delete` | 🆕 删除领域模块及所有相关文件 |
| `service` | 生成应用服务 (跨领域业务逻辑) |
| `service list` | 🆕 列出所有应用服务 |
| `service delete` | 🆕 删除应用服务 |
| `valueobject` | 🆕 生成领域值对象 |
| `spec` | 🆕 生成领域规格（Specification） |
| `policy` | 🆕 生成领域策略（Policy） |
| `event` | 🆕 生成领域事件（含注册） |
| `event-handler` | 🆕 生成事件处理器（含模块注入） |
| `tidy` | 🆕 运行 go mod tidy 更新依赖 |
| `serve` | 启动 Web GUI |

---

## init - 初始化项目

```bash
./soliton-gen init my-project
./soliton-gen init my-project --module github.com/myorg/my-project
```

**生成内容：** `cmd/main.go`, `cmd/migrate/main.go`, `configs/`, `internal/`, `go.mod`, `Makefile`, `README.md`

> **提示**: 生成的 `configs/config.example.yaml` 默认支持 sqlite/postgres，如需 MySQL 请扩展 `framework/orm/db.go`。

初始化后建议先执行依赖整理：

```bash
GOWORK=off go mod tidy
```

---

## domain - 生成领域模块

```bash
./soliton-gen domain User

# 简单格式（无备注）
./soliton-gen domain User --fields "username,email,status:enum(active|inactive)"

# 完整格式（带备注）
./soliton-gen domain User --fields "username:string:用户名,email::邮箱,status:enum(active|inactive):状态"

./soliton-gen domain User --fields "..." --force  # 强制覆盖
./soliton-gen domain User --fields "..." --wire   # 自动接入 main.go
```

### --wire 自动接线
使用 `--wire` 标志时，生成器会自动修改 `main.go`：
- 插入模块 import 和 handler import
- 添加 Module 注册
- 添加 Handler Provider
- 添加路由和迁移注册
同时会更新 `cmd/migrate/main.go`：
- 注入对应的迁移调用

**多模块支持**: 模板使用标记行 (`// soliton-gen:xxx`)，支持追加多个模块：
```go
// soliton-gen:imports     <- 自动插入 import
// soliton-gen:modules     <- 自动插入模块
// soliton-gen:handlers    <- 自动插入 Handler
// soliton-gen:routes      <- 自动插入路由注册
// soliton-gen:providers   <- 自动插入 Provider（如事件总线）
```

### 全部参数
| 参数 | 说明 | 示例 |
|------|------|------|
| `--fields`, `-f` | 定义字段 | `--fields "name,age:int"` |
| `--wire` | 自动注入 main.go | `--wire` |
| `--force` | 强制覆盖文件 | `--force` |
| `--table` | 自定义表名 | `--table "custom_users"` |
| `--route` | 自定义路由 | `--route "members"` |
| `--soft-delete` | 🆕 启用软删除 | `--soft-delete` |

### 字段格式

**基本格式：** `name:type:comment`（type 和 comment 可选）

| 格式 | 示例 | 说明 |
|------|------|------|
| `name` | `username` | string 类型，无备注 |
| `name:type` | `price:int64` | 指定类型，无备注 |
| `name:type:comment` | `username:string:用户名` | 完整格式 |
| `name::comment` | `email::邮箱` | 默认 string 类型 + 备注 |
| `name:enum(...):comment` | `status:enum(a\|b):状态` | 枚举 + 备注 |

### 字段类型
| 类型 | 格式 | 示例 | 说明 |
|------|------|------|------|
| string | `field` | `username` | 默认类型 |
| text | `field:text` | `description:text` | GORM text 类型 |
| int | `field:int` | `count:int` | 32位整数 |
| int64 | `field:int64` | `price:int64` | 64位整数 |
| float64 | `field:float64` | `score:float64` | 浮点数 |
| decimal | `field:decimal` | `amount:decimal` | 精确小数 (10,2 精度) |
| bool | `field:bool` | `active:bool` | 布尔值 |
| time | `field:time` | `created_at:time` | 时间戳 |
| time? | `field:time?` | `login_at:time?` | 可选时间字段 |
| date | `field:date` | `birth:date` | 日期 (无时间部分) |
| date? | `field:date?` | `expire:date?` | 可选日期 |
| uuid | `field:uuid` | `user_id:uuid` | 带索引的 UUID |
| json | `field:json` | `meta:json` | JSON 对象 (需 gorm.io/datatypes) |
| jsonb | `field:jsonb` | `data:jsonb` | PostgreSQL JSONB |
| bytes | `field:bytes` | `avatar:bytes` | 二进制数据 |
| enum | `field:enum(a\|b)` | `status:enum(active\|banned)` | 生成枚举类型 |

### 字段备注效果

生成的代码会在行尾添加注释：
```go
type User struct {
    ddd.BaseAggregateRoot
    ID        UserID `gorm:"primaryKey"`
    Username  string `gorm:"size:255"` // 用户名
    Email     string `gorm:"size:255"` // 邮箱
    Status    UserStatus `gorm:"size:50"` // 状态
}
```

### 🆕 新功能

#### 分页查询
所有生成的 List API 自动支持分页：
```bash
GET /api/users?page=1&page_size=20
```

响应格式：
```json
{
  "items": [...],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "total_pages": 5
}
```

#### 排序参数
List API 支持排序参数：
```bash
GET /api/users?page=1&page_size=20&sort_by=created_at&sort_order=desc
```

#### 数据库迁移入口
初始化项目会生成 `cmd/migrate/main.go`，用于执行迁移：
```bash
GOWORK=off go run ./cmd/migrate
```
Makefile 已内置 `migrate` 目标：
```bash
make migrate
```

#### 软删除
使用 `--soft-delete` 标志启用软删除：
```bash
soliton-gen domain User --fields "username,email" --soft-delete
```

生成的实体会包含 `DeletedAt` 字段：
```go
type User struct {
    ...
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

删除操作会自动变为软删除，GORM 查询默认排除已删除记录。

#### 错误码常量
生成的 `response.go` 包含预定义错误码：
```go
const (
    CodeSuccess      = 0     // 成功
    CodeBadRequest   = 400   // 请求错误
    CodeUnauthorized = 401   // 未授权
    CodeNotFound     = 404   // 未找到
    CodeInternal     = 500   // 内部错误
    CodeValidation   = 1001  // 验证失败
    CodeDuplicate    = 1002  // 重复条目
    CodeConflict     = 1003  // 业务冲突
)
```

### 生成文件 (10个)
- `domain/{name}/` - 实体 + Repository + Events + Domain Service
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

### 🆕 智能服务类型检测

生成器会自动检测服务类型并输出提示：

| 类型 | 判断条件 | 生成目录 | CLI 图标 | GUI 卡片颜色 |
|------|----------|----------|----------|--------------|
| **领域服务** | 存在同名 Domain | `application/{name}/` | ✅ | 🟢 绿色边框 |
| **跨域服务** | 不存在同名 Domain | `application/{name}/` | ℹ️ | 🟣 紫色边框 |

**CLI 输出示例：**
```
📋 服务类型检测
━━━━━━━━━━━━━━━
✅ 类型：领域服务 (Domain Service)
📁 目标路径：application/order
📝 DTO：service_dto.go

正在生成 Service OrderService...
```

### 🆕 Service DTO 生成逻辑

| 场景 | 行为 |
|------|------|
| 领域服务/跨域服务 | 生成 `service_dto.go` |
| `service_dto.go` 已存在且未 `--force` | 跳过生成 |
| 使用 `--force` | 覆盖生成 |

### 生成文件 (3个)
- `application/{name}/service.go` - 服务结构和方法
- `application/{name}/service_dto.go` - 请求/响应 DTO
- `application/{name}/module.go` - Fx 模块注册

📖 **详细文档**: [Service 应用服务使用指南](./SERVICE_GUIDE.md)

---

## 🆕 valueobject - 生成领域值对象

```bash
./soliton-gen valueobject user EmailAddress --fields "value:string"
./soliton-gen valueobject order Money --fields "amount:decimal,currency:string"
```

**生成文件：** `internal/domain/<domain>/value_object_<name>.go`

---

## 🆕 spec - 生成领域规格（Specification）

```bash
./soliton-gen spec user ActiveUserSpec --target User
./soliton-gen spec order PaidOrderSpec --target Order
```

**生成文件：** `internal/domain/<domain>/spec_<name>.go`

---

## 🆕 policy - 生成领域策略（Policy）

```bash
./soliton-gen policy user PasswordPolicy --target User
./soliton-gen policy order RefundPolicy --target Order
```

**生成文件：** `internal/domain/<domain>/policy_<name>.go`

---

## 🆕 event - 生成领域事件（含注册）

```bash
./soliton-gen event user UserActivated --fields "user_id:uuid"
./soliton-gen event order OrderPaid --fields "order_id:uuid,amount:decimal" --topic "order.paid"
```

**生成文件：** `internal/domain/<domain>/event_<name>.go`（自动注册到事件注册表）

---

## 🆕 event-handler - 生成事件处理器

```bash
./soliton-gen event-handler user UserCreated
./soliton-gen event-handler order OrderPaid --topic "order.paid"
```

**生成文件：** `internal/application/<domain>/event_handler_<name>.go`

> 事件名可带或不带 `Event` 后缀，生成时统一规范为 `<Name>Event`。

**自动更新：**
- `internal/application/<domain>/module.go` 注入 `fx.Provide`/`fx.Invoke`
- `cmd/main.go` 注入事件总线 Provider（EventBus 接口）

---

## 🆕 domain list - 列出领域

```bash
./soliton-gen domain list
```

**输出示例：**
```
已检测到 3 个领域模型：

  • order
  • product
  • user
```

---

## 🆕 domain delete - 删除领域

删除领域模块及所有相关文件。

```bash
./soliton-gen domain delete User          # 交互式确认
./soliton-gen domain delete User --force  # 跳过确认
```

**删除内容：**
- `internal/domain/<name>/`
- `internal/application/<name>/`
- `internal/infrastructure/persistence/<name>_repo.go`
- `internal/interfaces/http/<name>_handler.go`
- `main.go` 中的注入代码

---

## 🆕 service list / delete

```bash
./soliton-gen service list                    # 列出所有服务
./soliton-gen service delete OrderService     # 删除服务
./soliton-gen service delete OrderService --force
```

---

## 🆕 tidy - 更新依赖

运行 `go mod tidy` 更新项目依赖。

```bash
./soliton-gen tidy
```

等价于：`GOWORK=off go mod tidy`

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
soliton-gen domain User --fields "username:string:用户名,email::邮箱,role:enum(admin|customer):角色" --wire
soliton-gen domain Product --fields "name:string:商品名,price:int64:价格,stock:int:库存" --wire
soliton-gen domain Order --fields "user_id:uuid:用户ID,total:int64:总额,status:enum(pending|paid):订单状态" --wire

# 3. 生成跨领域服务
soliton-gen service OrderService --methods "CreateOrder,CancelOrder"

# 4. 运行（在 monorepo 中需要 GOWORK=off）
GOWORK=off go mod tidy && GOWORK=off go run ./cmd/main.go
```

> **Monorepo 提示**: 如果在包含 `go.work` 的 monorepo 中运行，请使用 `GOWORK=off` 前缀。
> **Makefile 默认**: 生成的 `Makefile` 默认 `GOWORK=off`，需要时可 `GOWORK=on make run`。
