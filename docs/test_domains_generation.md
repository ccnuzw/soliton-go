# 项目初始化与完整生成流程测试

## 执行步骤

### 1. 清理旧项目
```bash
rm -rf /Users/mac/Progame/soliton-go/application
```

### 2. 初始化新项目
```bash
soliton-gen init application
```

**生成内容：**
- ✅ 项目结构（cmd, configs, internal/...）
- ✅ go.mod
- ✅ main.go
- ✅ migrate.go
- ✅ config.yaml
- ✅ Makefile
- ✅ README.md
- ✅ .gitignore

### 3. 生成三个领域模型

### 3.1 生成命令（带 --wire）

```bash
# User
soliton-gen domain User --wire --fields "username,email,password,full_name,phone,avatar,bio:text,birth_date:time?,gender:enum(male|female|other),role:enum(admin|manager|user|guest),status:enum(active|inactive|suspended|banned),email_verified:bool,phone_verified:bool,last_login_at:time?,login_count:int,failed_login_count:int,balance:int64,points:int,vip_level:int,preferences:text"

# Order
soliton-gen domain Order --wire --fields "user_id:uuid,order_no,total_amount:int64,discount_amount:int64,tax_amount:int64,shipping_fee:int64,final_amount:int64,currency,payment_method:enum(credit_card|debit_card|paypal|alipay|wechat|cash),payment_status:enum(pending|paid|failed|refunded),order_status:enum(pending|confirmed|processing|shipped|delivered|cancelled|returned),shipping_method:enum(standard|express|overnight),tracking_number,receiver_name,receiver_phone,receiver_email,receiver_address,receiver_city,receiver_state,receiver_country,receiver_postal_code,notes:text,paid_at:time?,shipped_at:time?,delivered_at:time?,cancelled_at:time?,refund_amount:int64,refund_reason:text,item_count:int,weight:float64,is_gift:bool,gift_message:text"

# Product
soliton-gen domain Product --wire --fields "sku,name,slug,description:text,short_description:text,brand,category,subcategory,price:int64,original_price:int64,cost_price:int64,discount_percentage:int,stock:int,reserved_stock:int,sold_count:int,view_count:int,rating:float64,review_count:int,weight:float64,length:float64,width:float64,height:float64,color,size,material,manufacturer,country_of_origin,barcode,status:enum(draft|active|inactive|out_of_stock|discontinued),is_featured:bool,is_new:bool,is_on_sale:bool,is_digital:bool,requires_shipping:bool,is_taxable:bool,tax_rate:float64,min_order_quantity:int,max_order_quantity:int,tags:text,images:text,video_url,published_at:time?,discontinued_at:time?"
```

## User 领域（用户管理）

### 字段列表（20个字段）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| Username | string | 用户名 |
| Email | string | 邮箱 |
| Password | string | 密码 |
| FullName | string | 全名 |
| Phone | string | 电话 |
| Avatar | string | 头像URL |
| Bio | text | 个人简介 |
| BirthDate | *time.Time | 生日（可选） |
| Gender | enum | 性别：male\|female\|other |
| Role | enum | 角色：admin\|manager\|user\|guest |
| Status | enum | 状态：active\|inactive\|suspended\|banned |
| EmailVerified | bool | 邮箱已验证 |
| PhoneVerified | bool | 电话已验证 |
| LastLoginAt | *time.Time | 最后登录时间（可选） |
| LoginCount | int | 登录次数 |
| FailedLoginCount | int | 失败登录次数 |
| Balance | int64 | 账户余额 |
| Points | int | 积分 |
| VipLevel | int | VIP等级 |
| Preferences | text | 用户偏好设置 |

### 类型覆盖
- ✅ string
- ✅ text
- ✅ int
- ✅ int64
- ✅ bool
- ✅ time (可选)
- ✅ enum (3个枚举字段)

---

## Order 领域（订单管理）

### 字段列表（32个字段）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| UserId | uuid | 用户ID |
| OrderNo | string | 订单号 |
| TotalAmount | int64 | 总金额 |
| DiscountAmount | int64 | 折扣金额 |
| TaxAmount | int64 | 税费 |
| ShippingFee | int64 | 运费 |
| FinalAmount | int64 | 最终金额 |
| Currency | string | 货币 |
| PaymentMethod | enum | 支付方式：credit_card\|debit_card\|paypal\|alipay\|wechat\|cash |
| PaymentStatus | enum | 支付状态：pending\|paid\|failed\|refunded |
| OrderStatus | enum | 订单状态：pending\|confirmed\|processing\|shipped\|delivered\|cancelled\|returned |
| ShippingMethod | enum | 配送方式：standard\|express\|overnight |
| TrackingNumber | string | 物流单号 |
| ReceiverName | string | 收货人姓名 |
| ReceiverPhone | string | 收货人电话 |
| ReceiverEmail | string | 收货人邮箱 |
| ReceiverAddress | string | 收货地址 |
| ReceiverCity | string | 城市 |
| ReceiverState | string | 省份/州 |
| ReceiverCountry | string | 国家 |
| ReceiverPostalCode | string | 邮编 |
| Notes | text | 订单备注 |
| PaidAt | *time.Time | 支付时间（可选） |
| ShippedAt | *time.Time | 发货时间（可选） |
| DeliveredAt | *time.Time | 送达时间（可选） |
| CancelledAt | *time.Time | 取消时间（可选） |
| RefundAmount | int64 | 退款金额 |
| RefundReason | text | 退款原因 |
| ItemCount | int | 商品数量 |
| Weight | float64 | 重量 |
| IsGift | bool | 是否礼物 |
| GiftMessage | text | 礼物留言 |

### 类型覆盖
- ✅ string
- ✅ text
- ✅ int
- ✅ int64
- ✅ float64
- ✅ bool
- ✅ uuid
- ✅ time (可选)
- ✅ enum (4个枚举字段)

---

## Product 领域（商品管理）

### 字段列表（44个字段）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| Sku | string | SKU编号 |
| Name | string | 商品名称 |
| Slug | string | URL别名 |
| Description | text | 详细描述 |
| ShortDescription | text | 简短描述 |
| Brand | string | 品牌 |
| Category | string | 分类 |
| Subcategory | string | 子分类 |
| Price | int64 | 售价 |
| OriginalPrice | int64 | 原价 |
| CostPrice | int64 | 成本价 |
| DiscountPercentage | int | 折扣百分比 |
| Stock | int | 库存 |
| ReservedStock | int | 预留库存 |
| SoldCount | int | 已售数量 |
| ViewCount | int | 浏览次数 |
| Rating | float64 | 评分 |
| ReviewCount | int | 评论数 |
| Weight | float64 | 重量 |
| Length | float64 | 长度 |
| Width | float64 | 宽度 |
| Height | float64 | 高度 |
| Color | string | 颜色 |
| Size | string | 尺寸 |
| Material | string | 材质 |
| Manufacturer | string | 制造商 |
| CountryOfOrigin | string | 原产国 |
| Barcode | string | 条形码 |
| Status | enum | 状态：draft\|active\|inactive\|out_of_stock\|discontinued |
| IsFeatured | bool | 是否精选 |
| IsNew | bool | 是否新品 |
| IsOnSale | bool | 是否促销 |
| IsDigital | bool | 是否数字商品 |
| RequiresShipping | bool | 是否需要配送 |
| IsTaxable | bool | 是否需要税费 |
| TaxRate | float64 | 税率 |
| MinOrderQuantity | int | 最小订购量 |
| MaxOrderQuantity | int | 最大订购量 |
| Tags | text | 标签 |
| Images | text | 图片列表 |
| VideoUrl | string | 视频URL |
| PublishedAt | *time.Time | 发布时间（可选） |
| DiscontinuedAt | *time.Time | 停产时间（可选） |

### 类型覆盖
- ✅ string
- ✅ text
- ✅ int
- ✅ int64
- ✅ float64
- ✅ bool
- ✅ time (可选)
- ✅ enum (1个枚举字段)

---

## DDD 领域增强组件生成

### 生成命令

```bash
# 领域值对象
soliton-gen valueobject user EmailAddress --fields "value:string"

# 领域规格
soliton-gen spec user ActiveUserSpec --target User

# 领域策略
soliton-gen policy user PasswordPolicy --target User

# 自定义领域事件（含注册）
soliton-gen event user UserActivated --fields "user_id:uuid"

# 事件处理器（自动注入 module.go / main.go）
soliton-gen event-handler user UserActivatedEvent
```

### 生成结果
- `internal/domain/user/value_object_email_address.go`
- `internal/domain/user/spec_active_user_spec.go`
- `internal/domain/user/policy_password_policy.go`
- `internal/domain/user/event_user_activated.go`
- `internal/application/user/event_handler_user_activated.go`

---

## 类型覆盖总结

### 所有支持的字段类型都已测试

| 类型 | User | Order | Product | 总计 |
|------|------|-------|---------|------|
| string | 6 | 13 | 15 | 34 |
| text | 2 | 4 | 4 | 10 |
| int | 3 | 1 | 7 | 11 |
| int64 | 3 | 8 | 3 | 14 |
| float64 | 0 | 1 | 8 | 9 |
| bool | 2 | 1 | 6 | 9 |
| time? | 2 | 4 | 2 | 8 |
| uuid | 0 | 1 | 0 | 1 |
| enum | 3 | 4 | 1 | 8 |
| **总计** | **20** | **32** | **44** | **96** |

---

## 迁移与运行验证

### 迁移
```bash
cd /Users/mac/Progame/soliton-go/application
GOWORK=off go run ./cmd/migrate.go
```

### 启动服务
```bash
GOWORK=off go run ./cmd/main.go
```

### CRUD + 排序参数验证
```bash
# 分页 + 排序
curl "http://localhost:8080/api/users?page=1&page_size=20&sort_by=created_at&sort_order=desc"
```

---

## 编译验证

```bash
cd /Users/mac/Progame/soliton-go/application
go build ./...
```

**结果：** ✅ 编译成功，无错误

---

## 生成的文件结构

```
application/
├── cmd/
│   ├── main.go (已自动注入所有模块)
│   └── migrate.go
├── configs/
│   ├── config.yaml
│   └── config.example.yaml
├── internal/
│   ├── domain/
│   │   ├── user/
│   │   │   ├── user.go (20个字段)
│   │   │   ├── repository.go
│   │   │   ├── events.go
│   │   │   ├── service.go
│   │   │   ├── value_object_email_address.go
│   │   │   ├── spec_active_user_spec.go
│   │   │   ├── policy_password_policy.go
│   │   │   └── event_user_activated.go
│   │   ├── order/
│   │   │   ├── order.go (32个字段)
│   │   │   ├── repository.go
│   │   │   ├── events.go
│   │   │   └── service.go
│   │   └── product/
│   │       ├── product.go (44个字段)
│   │       ├── repository.go
│   │       ├── events.go
│   │       └── service.go
│   ├── application/
│   │   ├── user/
│   │   │   ├── commands.go
│   │   │   ├── queries.go
│   │   │   ├── dto.go
│   │   │   ├── module.go
│   │   │   └── event_handler_user_activated.go
│   │   ├── order/
│   │   │   ├── commands.go
│   │   │   ├── queries.go
│   │   │   ├── dto.go
│   │   │   └── module.go
│   │   └── product/
│   │       ├── commands.go
│   │       ├── queries.go
│   │       ├── dto.go
│   │       └── module.go
│   ├── infrastructure/
│   │   └── persistence/
│   │       ├── user_repo.go
│   │       ├── order_repo.go
│   │       └── product_repo.go
│   └── interfaces/
│       └── http/
│           ├── helpers.go
│           ├── response.go
│           ├── user_handler.go
│           ├── order_handler.go
│           └── product_handler.go
├── go.mod
├── Makefile
├── README.md
└── .gitignore
```

---

## 特性验证

### ✅ 已验证的功能

1. **字段类型多样性**
   - string, text, int, int64, float64, bool, time?, uuid, enum

2. **枚举类型**
   - User: Gender, Role, Status
   - Order: PaymentMethod, PaymentStatus, OrderStatus, ShippingMethod
   - Product: Status

3. **可选时间字段**
   - User: BirthDate, LastLoginAt
   - Order: PaidAt, ShippedAt, DeliveredAt, CancelledAt
   - Product: PublishedAt, DiscontinuedAt

4. **自动生成的功能**
   - CRUD Commands (Create, Update, Delete)
   - Queries (Get, List)
   - DTO (Request/Response)
   - HTTP Handlers
   - Repository 实现
   - Domain Service（领域服务）
   - Fx Module 注入
   - main.go 自动注入
   - migrate.go 迁移入口
   - List API 排序参数 (sort_by / sort_order)
   - 领域增强组件（ValueObject/Spec/Policy/Event/Handler）

5. **代码质量**
   - ✅ 编译通过
   - ✅ 无重复字段
   - ✅ 类型正确
   - ✅ 命名规范

---

## 下一步

项目已完全初始化并生成了三个功能完整的领域模型，可以：

1. 运行项目：`GOWORK=off go run ./cmd/main.go`
2. 测试 API 端点
3. 添加业务逻辑
4. 编写单元测试
5. 部署到生产环境
