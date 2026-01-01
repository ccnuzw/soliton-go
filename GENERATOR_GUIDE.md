# Soliton-Go 代码生成器使用指南

## 📦 安装

```bash
cd tools/generator
go build -o soliton-gen .
```

---

## 🚀 命令

### init - 初始化新项目

```bash
./soliton-gen init <project-name>
./soliton-gen init <project-name> --module github.com/myorg/my-project
```

**生成内容：**
```
my-project/
├── cmd/main.go              # 入口 (Fx + Gin)
├── configs/
│   ├── config.yaml          # 配置文件
│   └── config.example.yaml  # 配置示例
├── internal/
│   ├── domain/              # 领域层
│   ├── application/         # 应用层
│   ├── infrastructure/      # 基础设施层
│   └── interfaces/http/     # HTTP 层
├── go.mod                   # Go 模块
├── Makefile                 # 常用命令
├── README.md                # 项目说明
└── .gitignore               # Git 忽略
```

---

### domain - 生成领域模块

```bash
./soliton-gen domain <EntityName>
./soliton-gen domain <EntityName> --fields "<field1>,<field2:type>,..."
./soliton-gen domain <EntityName> --fields "..." --force  # 强制覆盖
```

---

## 📋 字段类型参考

| 类型 | 格式 | Go 类型 |
|------|------|---------|
| string | `field` | `string` |
| text | `field:text` | `string` |
| int | `field:int` | `int` |
| int64 | `field:int64` | `int64` |
| uuid | `field:uuid` | `string` |
| **enum** | `field:enum(a\|b\|c)` | 枚举类型 |

---

## 🎯 完整示例

```bash
# 1. 初始化项目
./soliton-gen init my-shop

# 2. 进入项目
cd my-shop

# 3. 生成领域模块
../soliton-gen domain User --fields "username,email,role:enum(admin|customer),status:enum(active|banned)"
../soliton-gen domain Product --fields "name,price:int64,stock:int,status:enum(draft|active)"
../soliton-gen domain Order --fields "user_id:uuid,total:int64,status:enum(pending|paid|shipped)"

# 4. 更新 main.go (取消注释导入)

# 5. 运行
go mod tidy
go run ./cmd/main.go
```

---

## 🔄 修改字段的三种方式

### 方式 1: 手动修改（小改动）
编辑 4 个文件：`{name}.go`, `commands.go`, `dto.go`, `{name}_handler.go`

### 方式 2: 删除后重新生成（大改动）
```bash
rm -rf internal/domain/user internal/application/user ...
./soliton-gen domain User --fields "..."
```

### 方式 3: --force 强制覆盖（最简单）
```bash
./soliton-gen domain User --fields "..." --force
```

---

## 📁 domain 生成文件 (9个)

| 层 | 文件 | 说明 |
|---|------|------|
| Domain | `{name}.go` | Entity + Enum |
| Domain | `repository.go` | Repository 接口 |
| Domain | `events.go` | 领域事件 |
| Infrastructure | `{name}_repo.go` | GORM 实现 |
| Application | `commands.go` | 命令处理器 |
| Application | `queries.go` | 查询处理器 |
| Application | `dto.go` | DTO |
| Application | `module.go` | Fx 模块 |
| Interfaces | `{name}_handler.go` | HTTP Handler |

---

## 🔒 状态说明

| 状态 | 说明 |
|------|------|
| `[NEW]` | 新建 |
| `[SKIP]` | 已存在，跳过 |
| `[OVERWRITE]` | --force 覆盖 |
| `[DIR]` | 创建目录 |
