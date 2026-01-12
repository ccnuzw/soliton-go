# Service 应用服务使用指南

## 📖 什么是 Service？

**Service（应用服务）** 是 DDD 分层架构中的**应用层组件**，用于编排跨多个领域的复杂业务逻辑。

### Domain vs Service 的区别

| 类型 | 关注点 | 示例 |
|------|--------|------|
| **Domain** | 单一领域的实体和规则 | User 的状态变更、Order 的价格计算 |
| **Service** | 跨多个领域的业务编排 | 下单流程（User + Product + Order） |

---

## 🎯 何时使用 Service？

**适合使用 Service 的场景：**
- ✅ 涉及多个 Domain 的协作（如下单需要 User + Product + Order）
- ✅ 需要事务协调（如支付成功后更新订单状态并发送通知）
- ✅ 复杂的用例流程（如退款需要校验订单状态、计算退款金额、更新库存）

**不适合使用 Service 的场景：**
- ❌ 单一实体的 CRUD 操作（用 Domain + Handler 即可）
- ❌ 简单的查询（用 Query Handler 即可）

---

## 🚀 快速生成

```bash
# 基础用法
./soliton-gen service OrderService

# 指定方法
./soliton-gen service OrderService --methods "CreateOrder,CancelOrder,GetUserOrders"

# 指定方法并添加备注
./soliton-gen service OrderService --methods "CreateOrder::创建订单,CancelOrder::取消订单"

# 添加服务备注
./soliton-gen service OrderService --methods "CreateOrder,CancelOrder" --remark "订单服务"

# 支付服务示例
./soliton-gen service PaymentService --methods "ProcessPayment,Refund,QueryStatus"
```

### 参数说明
| 参数 | 说明 | 示例 |
|------|------|------|
| `--methods` | 方法列表（支持备注） | `--methods "Create::创建,Update::更新"` |
| `--remark` | 服务备注 | `--remark "支付服务"` |
| `--force` | 强制覆盖文件 | `--force` |

---

## 🔍 智能服务类型检测

生成器会自动检测并区分两种类型的服务：

### 服务类型对照表

| 类型 | 判断条件 | 生成目录 | GUI 卡片颜色 |
|------|----------|----------|--------------|
| **领域服务** (Domain Service) | 存在同名 Domain（如 `domain/order/`） | `application/order/service.go` | 🟢 绿色边框 |
| **跨域服务** (Cross-domain Service) | 不存在同名 Domain | `application/payment/service.go` | 🟣 紫色边框 |

### CLI 智能提示

运行 `soliton-gen service` 时会输出检测结果：

```bash
$ ./soliton-gen service OrderService
📋 服务类型检测
━━━━━━━━━━━━━━━
✅ 类型：领域服务 (Domain Service)
📁 目标路径：application/order
📝 DTO：service_dto.go

正在生成 Service OrderService...
```

```bash
$ ./soliton-gen service PaymentService
📋 服务类型检测
━━━━━━━━━━━━━━━
ℹ️  类型：跨领域服务 (Cross-domain Service)
📁 目标路径：application/payment
📝 DTO：service_dto.go

正在生成 Service PaymentService...
```

### GUI 颜色标识

在 Web GUI 的"已生成服务"列表中，卡片会通过颜色区分：

- 🟢 **绿色左边框 + "领域" 徽章**：表示此服务有对应的 Domain
- 🟣 **紫色左边框 + "跨域" 徽章**：表示此服务是独立的跨域编排服务

### Service DTO 生成逻辑

| 场景 | 行为 |
|------|------|
| 领域服务/跨域服务 | 生成 `service_dto.go` |
| `service_dto.go` 已存在且未 `--force` | 跳过生成 |
| 使用 `--force` | 覆盖生成 |

---

## 📁 生成文件

```
application/{servicename}/
├── service.go    # 服务结构和方法
├── service_dto.go # 请求/响应 DTO
└── module.go     # Fx 模块注册
```

---

## 📝 完整示例

### 1. 生成服务

```bash
./soliton-gen service OrderService --methods "CreateOrder,CancelOrder"
```

### 2. 注入 Repositories

编辑 `order_service.go`：

```go
package services

import (
    "context"
    
    "your-project/internal/domain/user"
    "your-project/internal/domain/product"
    "your-project/internal/domain/order"
)

type OrderService struct {
    userRepo    user.UserRepository
    productRepo product.ProductRepository
    orderRepo   order.OrderRepository
}

func NewOrderService(
    userRepo user.UserRepository,
    productRepo product.ProductRepository,
    orderRepo order.OrderRepository,
) *OrderService {
    return &OrderService{
        userRepo:    userRepo,
        productRepo: productRepo,
        orderRepo:   orderRepo,
    }
}
```

### 3. 定义 DTO

编辑 `service_dto.go`：

```go
type CreateOrderServiceRequest struct {
    UserID     string           `json:"user_id"`
    Items      []OrderItemInput `json:"items"`
    Address    string           `json:"address"`
}

type OrderItemInput struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

type CreateOrderServiceResponse struct {
    OrderID     string `json:"order_id"`
    OrderNo     string `json:"order_no"`
    TotalAmount int64  `json:"total_amount"`
}
```

### 4. 实现业务逻辑

```go
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderServiceRequest) (*CreateOrderServiceResponse, error) {
    // 1. 验证用户
    user, err := s.userRepo.Find(ctx, user.UserID(req.UserID))
    if err != nil {
        return nil, errors.New("user not found")
    }
    if user.Status != user.UserStatusActive {
        return nil, errors.New("user is not active")
    }

    // 2. 检查库存并计算价格
    var totalAmount int64
    for _, item := range req.Items {
        product, err := s.productRepo.Find(ctx, product.ProductID(item.ProductID))
        if err != nil {
            return nil, errors.New("product not found")
        }
        if product.Stock < item.Quantity {
            return nil, errors.New("insufficient stock")
        }
        totalAmount += product.Price * int64(item.Quantity)
    }

    // 3. 创建订单
    orderID := uuid.New().String()
    orderNo := generateOrderNo()
    newOrder := order.NewOrder(orderID, orderNo, req.UserID, totalAmount)
    
    if err := s.orderRepo.Save(ctx, newOrder); err != nil {
        return nil, err
    }

    // 4. 扣减库存
    for _, item := range req.Items {
        s.productRepo.DeductStock(ctx, product.ProductID(item.ProductID), item.Quantity)
    }

    return &CreateOrderServiceResponse{
        OrderID:     orderID,
        OrderNo:     orderNo,
        TotalAmount: totalAmount,
    }, nil
}
```

### 5. 注册到 Fx

```go
// main.go
fx.New(
    userapp.Module,
    productapp.Module,
    orderapp.Module,
    
    // 注册 Service
    fx.Provide(services.NewOrderService),
    
    fx.Invoke(StartServer),
)
```

### 6. 创建 HTTP Handler（可选）

```go
// interfaces/http/order_service_handler.go
type OrderServiceHandler struct {
    service *services.OrderService
}

func (h *OrderServiceHandler) CreateOrder(c *gin.Context) {
    var req services.CreateOrderServiceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        BadRequest(c, err.Error())
        return
    }
    
    resp, err := h.service.CreateOrder(c.Request.Context(), req)
    if err != nil {
        InternalError(c, err.Error())
        return
    }
    
    Success(c, resp)
}
```

---

## 🏗 架构位置

```
┌─────────────────────────────────────────────┐
│              Interfaces Layer               │
│  (HTTP Handlers, GraphQL, gRPC)             │
└───────────────────┬─────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────┐
│            Application Layer                │
│  ┌─────────────────────────────────────┐    │
│  │         Services (跨领域)            │    │
│  │   OrderService, PaymentService      │    │
│  └─────────────────────────────────────┘    │
│  ┌─────────────────────────────────────┐    │
│  │      Commands / Queries (单领域)     │    │
│  │   CreateUserCommand, GetUserQuery   │    │
│  └─────────────────────────────────────┘    │
└───────────────────┬─────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────┐
│              Domain Layer                   │
│   User, Product, Order (Aggregates)         │
└─────────────────────────────────────────────┘
```

---

## 💡 最佳实践

1. **Service 不持有状态**：只通过 Repository 访问数据
2. **保持方法单一职责**：每个方法对应一个用例
3. **使用事务**：跨多个写操作时使用数据库事务
4. **发布领域事件**：在业务完成后发布事件通知其他系统
5. **错误处理**：返回业务错误而非技术错误
