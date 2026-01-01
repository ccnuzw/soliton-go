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
./soliton-gen domain <EntityName> --fields "<field1>,<field2:type>,..."
./soliton-gen domain <EntityName> --fields "..." --force  # 强制覆盖
```

---

## 📋 字段类型参考

| 类型 | 格式 | Go 类型 | GORM 标签 |
|------|------|---------|-----------|
| 默认/string | `field` | `string` | `size:255` |
| text | `field:text` | `string` | `type:text` |
| int | `field:int` | `int` | `not null` |
| int64 | `field:int64` | `int64` | `not null` |
| float | `field:float` | `float64` | - |
| bool | `field:bool` | `bool` | `default:false` |
| uuid | `field:uuid` | `string` | `size:36;index` |
| time | `field:time` | `time.Time` | `autoCreateTime` |
| **enum** | `field:enum(a\|b\|c)` | `EntityField` | `size:50` |

---

## 🎯 使用示例

```bash
# User - 用户管理
./soliton-gen domain User --fields "username,email,password_hash,phone,role:enum(admin|seller|customer),status:enum(active|inactive|banned)"

# Product - 商品管理
./soliton-gen domain Product --fields "name,sku,price:int64,stock:int,status:enum(draft|active|inactive)"

# Order - 订单管理
./soliton-gen domain Order --fields "user_id:uuid,order_no,total_amount:int64,status:enum(pending|paid|shipped|completed|cancelled)"
```

---

## 🔄 修改字段的三种方式

### 方式 1: 手动修改（小改动推荐）

需要修改 **4 个文件**：

| 文件 | 修改内容 |
|------|----------|
| `domain/{name}/{name}.go` | Entity 结构体 + NewXxx() + Update() |
| `application/{name}/commands.go` | Command 结构体字段 |
| `application/{name}/dto.go` | Request/Response 结构体 |
| `interfaces/http/{name}_handler.go` | Handler 中 cmd 赋值 |

**示例：给 User 添加 `age:int` 字段**

```go
// 1. domain/user/user.go - 实体
type User struct {
    ...
    Age int `gorm:"default:0"`  // 新增
}

// 2. application/user/commands.go
type CreateUserCommand struct {
    ...
    Age int  // 新增
}

// 3. application/user/dto.go
type CreateUserRequest struct {
    ...
    Age int `json:"age"`  // 新增
}
```

---

### 方式 2: 删除后重新生成（大改动推荐）

```bash
# 1. 删除模块相关文件
rm -rf application/internal/domain/user
rm -rf application/internal/application/user
rm application/internal/infrastructure/persistence/user_repo.go
rm application/internal/interfaces/http/user_handler.go

# 2. 用新字段重新生成
./soliton-gen domain User --fields "username,email,age:int,status:enum(active|banned)"
```

---

### 方式 3: 使用 --force 强制覆盖（最简单）

```bash
./soliton-gen domain User --fields "username,email,age:int,status:enum(active|banned)" --force
```

**输出：**
```
🚀 Generating domain: User (force mode)

📦 Domain Layer
   [OVERWRITE] user.go
   [OVERWRITE] repository.go
   [OVERWRITE] events.go
...
```

> ⚠️ **警告**：`--force` 会覆盖所有文件，包括您手动修改的代码！

---

## 📁 生成文件清单

每次生成 **9 个文件**：

| 层 | 文件 | 说明 |
|---|------|------|
| Domain | `{name}.go` | Entity + ID + Enum 类型 |
| Domain | `repository.go` | Repository 接口 |
| Domain | `events.go` | 领域事件 |
| Infrastructure | `{name}_repo.go` | Repository GORM 实现 |
| Application | `commands.go` | Create/Update/Delete Handlers |
| Application | `queries.go` | Get/List Handlers |
| Application | `dto.go` | Request/Response DTOs |
| Application | `module.go` | Fx 依赖注入模块 |
| Interfaces | `{name}_handler.go` | HTTP CRUD Handler |

---

## 🔒 文件状态说明

| 状态 | 说明 |
|------|------|
| `[NEW]` | 新建文件 |
| `[SKIP]` | 文件已存在，跳过 |
| `[OVERWRITE]` | 使用 --force 覆盖 |
| `[ERROR]` | 生成失败 |

---

## ❓ FAQ

**Q: 如何添加/删除字段？**
A: 小改动手动修改 4 个文件；大改动用 `--force` 重新生成。

**Q: --force 会覆盖我的代码吗？**
A: 会！使用前请确保没有重要修改，或先备份。

**Q: 如何添加自定义 Repository 方法？**
A: 在 `repository.go` 添加声明，在 `{name}_repo.go` 实现。使用 --force 不会影响您新增的方法接口定义，但实现文件会被覆盖。
