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
```

### 字段类型
| 类型 | 格式 | 示例 |
|------|------|------|
| string | `field` | `username` |
| int64 | `field:int64` | `price:int64` |
| text | `field:text` | `description:text` |
| uuid | `field:uuid` | `user_id:uuid` |
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
./soliton-gen service PaymentService --methods "ProcessPayment,Refund"
```

### 生成文件 (2个)
- `application/services/{name}_service.go` - 服务结构和方法
- `application/services/{name}_dto.go` - 请求/响应 DTO

### 使用场景
- **下单服务**: 涉及 User + Product + Order 多个领域
- **支付服务**: 涉及 Order + Payment + Wallet 多个领域
- **库存服务**: 涉及 Product + Inventory + Warehouse 多个领域

### 示例输出

```go
// OrderService handles cross-domain business logic.
type OrderService struct {
    userRepo  user.UserRepository
    orderRepo order.OrderRepository
    productRepo product.ProductRepository
}

// CreateOrder implements the CreateOrder use case.
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error) {
    // 1. 验证用户
    // 2. 检查库存
    // 3. 创建订单
    // 4. 扣减库存
    // 5. 发布事件
}
```

---

## 🔄 修改已生成代码

| 场景 | 推荐方式 |
|------|----------|
| 小改动 (1-2字段) | 手动编辑 |
| 大改动 | `--force` 重新生成 |

```bash
./soliton-gen domain User --fields "..." --force
```

---

## 🎯 完整开发流程

```bash
# 1. 初始化项目
./soliton-gen init my-shop && cd my-shop

# 2. 生成领域模块
soliton-gen domain User --fields "username,email,role:enum(admin|customer)"
soliton-gen domain Product --fields "name,price:int64,stock:int"
soliton-gen domain Order --fields "user_id:uuid,total:int64,status:enum(pending|paid)"

# 3. 生成跨领域服务
soliton-gen service OrderService --methods "CreateOrder,CancelOrder"

# 4. 更新 main.go (取消注释导入)
# 5. 运行
go mod tidy && go run ./cmd/main.go
```
