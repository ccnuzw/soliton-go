# Soliton-Go 代码生成器使用指南

## 📦 安装

```bash
cd tools/generator
go build -o soliton-gen .
```

---

## 🚀 命令

### domain - 生成领域模块

生成完整的 DDD 领域模块，包含 9 个文件。

**基础用法：**
```bash
./soliton-gen domain <EntityName>
```

**带字段用法（推荐）：**
```bash
./soliton-gen domain <EntityName> --fields "<field1>,<field2:type>,..."
```

---

## 📋 字段类型参考

| 类型 | 格式 | Go 类型 | GORM 标签 |
|------|------|---------|-----------|
| 默认/string | `field` 或 `field:string` | `string` | `size:255` |
| text | `field:text` | `string` | `type:text` |
| int | `field:int` | `int` | `not null;default:0` |
| int64/long | `field:int64` | `int64` | `not null;default:0` |
| float/double | `field:float` | `float64` | `default:0` |
| bool | `field:bool` | `bool` | `default:false` |
| uuid/id | `field:uuid` | `string` | `size:36;index` |
| time/datetime | `field:time` | `time.Time` | `autoCreateTime` |
| **enum** | `field:enum(a\|b\|c)` | `EntityField` | `size:50;default:'a'` |

---

## 🎯 使用示例

### 用户管理

```bash
./soliton-gen domain User --fields "username,email,password_hash,phone,avatar,role:enum(admin|seller|customer),status:enum(active|inactive|banned)"
```

**生成：**
- `UserRole` 枚举: `admin`, `seller`, `customer`
- `UserStatus` 枚举: `active`, `inactive`, `banned`
- 7 个业务字段 + `CreatedAt`/`UpdatedAt`

### 商品管理

```bash
./soliton-gen domain Product --fields "name,sku,description:text,price:int64,original_price:int64,stock:int,category_id:uuid,brand_id:uuid,images:text,status:enum(draft|active|inactive|out_of_stock)"
```

### 订单管理

```bash
./soliton-gen domain Order --fields "user_id:uuid,order_no,total_amount:int64,discount_amount:int64,payable_amount:int64,status:enum(pending|paid|shipped|delivered|completed|cancelled),receiver_name,receiver_phone,receiver_address:text,tracking_no"
```

---

## 📁 生成文件清单

每次生成 **9 个文件**：

| 层 | 文件 | 说明 |
|---|------|------|
| Domain | `{name}.go` | Entity + ID + Enum 类型 |
| Domain | `repository.go` | Repository 接口 |
| Domain | `events.go` | 领域事件 (Created/Updated/Deleted) |
| Infrastructure | `{name}_repo.go` | Repository GORM 实现 + 迁移 |
| Application | `commands.go` | Create/Update/Delete Handlers |
| Application | `queries.go` | Get/List Handlers |
| Application | `dto.go` | Request/Response DTOs |
| Application | `module.go` | Fx 依赖注入模块 |
| Interfaces | `{name}_handler.go` | HTTP CRUD 5 接口 |

---

## 🔒 文件保护

生成器会**跳过已存在的文件**，不会覆盖您的修改。

```
[NEW] user.go        # 新建
[SKIP] user.go       # 已存在，跳过
```

---

## ⚡ 集成到项目

### 1. 导入模块

```go
import userapp "github.com/soliton-go/application/internal/application/user"

fx.New(
    userapp.Module,  // 一行导入所有依赖
)
```

### 2. 注册迁移

```go
userapp.RegisterMigration(db)  // 自动建表
```

### 3. 注册路由

```go
userHandler.RegisterRoutes(r)  // 注册 CRUD 路由
```

---

## ❓ FAQ

**Q: 如何添加自定义字段？**
A: 使用 `--fields` 参数或在生成后直接编辑 `{name}.go` 文件。

**Q: 如何添加自定义 Repository 方法？**
A: 在 `repository.go` 中添加方法声明，在 `{name}_repo.go` 中实现。

**Q: 如何防止文件被覆盖？**
A: 生成器自动跳过已存在的文件。
