# Soliton-Gen 代码生成器

Soliton-Gen 是 Soliton-Go 框架的代码生成工具，支持命令行和 Web GUI 两种使用方式。

## 🚀 快速开始

### 编译

```bash
cd tools/generator
go build -o soliton-gen .
```

### 使用 Web GUI（推荐）

**新项目：**
```bash
# 1. 在空目录启动 GUI
mkdir my-workspace && cd my-workspace
./soliton-gen serve

# 2. 访问 http://127.0.0.1:3000
# 3. 点击"初始化项目"，在 GUI 中创建项目
# 4. 项目创建后，cd 到新项目目录继续开发
```

**现有项目：**
```bash
# 在项目根目录（包含 go.mod）启动
cd /path/to/your/project
./soliton-gen serve

# 访问 http://127.0.0.1:3000 使用可视化界面
```

**自定义端口：**
```bash
./soliton-gen serve --port 8080 --host 0.0.0.0
```

### 使用命令行

```bash
# 初始化项目
./soliton-gen init my-project

# 生成领域模块（简单格式，无备注）
./soliton-gen domain User --fields "username,email,status:enum(active|inactive)" --wire

# 生成领域模块（完整格式，带备注）
./soliton-gen domain User --fields "username:string:用户名,email::邮箱,status:enum(active|inactive):账户状态" --wire

# 生成应用服务（带方法备注）
./soliton-gen service Order --methods "CreateOrder::创建订单,ProcessPayment::发起支付" --wire
```

## 📖 文档

- [Web GUI 使用指南](./WEB_GUI_GUIDE.md) - 详细的 Web 界面使用说明
- [项目 README](../../README.md) - 框架总体说明
- [快速开始](../../QUICK_START.md) - 快速上手指南

## 🎨 Web GUI 功能

### 1. Dashboard（首页）
- 项目状态检测
- 快速导航
- 使用指南

### 2. Init Project（初始化项目）
- 可视化配置项目信息
- 三步向导流程（配置 → 预览 → 完成）
- 详细的字段说明和提示

### 3. Domain Editor（领域编辑器）
- 可视化字段编辑器
- 🆕 **字段备注**：为每个字段添加注释，自动生成行尾注释
- 🆕 **领域备注**：为领域添加说明，便于识别与检索
- 🆕 **字段排序**：通过 ↑↓ 按钮调整字段顺序
- 支持多种字段类型（string、int、enum 等）
- 枚举值可视化配置
- 软删除选项
- 自动注入到 main.go
- 🆕 **自动更新依赖**：生成后自动运行 go mod tidy
- 代码预览功能

### 4. Service Editor（服务编辑器）
- 可视化方法配置
- 🆕 **方法备注**：为每个方法添加用途说明，便于理解与回显
- 🆕 **服务备注**：为服务添加说明，卡片列表可回显
- 默认方法生成
- 代码预览功能

### 5. DDD Enhancer（领域增强）
- 以中文说明为主，专业术语保留英文
- 支持 Value Object / Specification / Policy / Event / Handler 可视化生成
- 支持 Event + Handler 组合生成
- 🆕 支持已有组件加载、重命名与删除
- 🆕 支持 Diff 对比预览与批量导入/导出字段

### 6. Migration Center（迁移中心）
- 详细迁移日志（SYSTEM / TIDY / MIGRATE）
- 支持自动 tidy 与执行前确认
- 支持历史记录、复制与下载日志

## 🔌 命令列表

| 命令 | 说明 | 示例 |
|------|------|------|
| `init <name>` | 初始化新项目 | `soliton-gen init my-project` |
| `domain <name>` | 生成领域模块 | `soliton-gen domain User --fields "username,email"` |
| `domain list` | 🆕 列出所有领域 | `soliton-gen domain list` |
| `domain delete <name>` | 🆕 删除领域模块 | `soliton-gen domain delete User` |
| `service <name>` | 生成应用服务 | `soliton-gen service Order --methods "CreateOrder"` |
| `service list` | 🆕 列出所有服务 | `soliton-gen service list` |
| `service delete <name>` | 🆕 删除应用服务 | `soliton-gen service delete OrderService` |
| `valueobject <domain> <name>` | 生成领域值对象 | `soliton-gen valueobject user EmailAddress` |
| `spec <domain> <name>` | 生成领域规格 | `soliton-gen spec user ActiveUserSpec` |
| `policy <domain> <name>` | 生成领域策略 | `soliton-gen policy user PasswordPolicy` |
| `event <domain> <name>` | 生成领域事件 | `soliton-gen event user UserActivated` |
| `event-handler <domain> <event>` | 生成事件处理器 | `soliton-gen event-handler user UserActivated` |
| `tidy` | 🆕 更新依赖 | `soliton-gen tidy` |
| `serve` | 启动 Web GUI | `soliton-gen serve --port 3000` |

### Domain 命令参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `--fields` | 字段列表 | 见下方字段格式 |
| `--remark` | 领域备注 | `--remark "用户领域"` |
| `--table` | 自定义表名 | `--table "sys_users"` |
| `--route` | 自定义路由 | `--route "/v1/users"` |
| `--soft-delete` | 启用软删除 | `--soft-delete` |
| `--wire` | 自动注入到 main.go | `--wire` |
| `--force` | 强制覆盖/跳过确认 | `--force` |

#### 字段格式

**基本格式：** `name:type:comment`（type 和 comment 可选）

| 格式 | 示例 | 说明 |
|------|------|------|
| `name` | `username` | string 类型，无备注 |
| `name:type` | `price:int64` | 指定类型，无备注 |
| `name:type:comment` | `username:string:用户名` | 完整格式 |
| `name::comment` | `email::邮箱` | 默认 string 类型 + 备注 |
| `name:enum(...):comment` | `status:enum(a\|b):状态` | 枚举 + 备注 |

### Service 命令参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `--methods` | 方法列表（支持备注） | `--methods "Create::创建,Update::更新"` |
| `--remark` | 服务备注 | `--remark "支付服务"` |
| `--force` | 强制覆盖/跳过确认 | `--force` |

### Serve 命令参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--port` | 端口号 | `3000` |
| `--host` | 主机地址 | `127.0.0.1` |

## 📂 项目结构

```
tools/generator/
├── cmd/                    # CLI 命令
│   ├── init.go            # 初始化命令
│   ├── domain.go          # 领域生成命令
│   ├── service.go         # 服务生成命令
│   ├── serve.go           # Web GUI 命令
│   └── layout.go          # 布局工具
├── core/                   # 核心逻辑
│   ├── types.go           # 类型定义
│   ├── layout.go          # 项目布局
│   ├── helpers.go         # 工具函数
│   ├── project.go         # 项目初始化
│   ├── domain.go          # 领域生成
│   ├── service.go         # 服务生成
│   ├── templates_project.go   # 项目模板
│   ├── templates_domain.go    # 领域模板
│   └── templates_service.go   # 服务模板
├── server/                 # Web 服务器
│   ├── server.go          # 服务器主文件
│   ├── handlers/          # API 处理器
│   │   ├── project.go     # 项目 API
│   │   ├── domain.go      # 领域 API
│   │   └── service.go     # 服务 API
│   └── static/            # 前端静态文件（嵌入）
├── web/                    # Vue 前端
│   ├── src/
│   │   ├── views/         # 页面组件
│   │   ├── api.ts         # API 客户端
│   │   └── router.ts      # 路由配置
│   └── vite.config.ts     # Vite 配置
├── main.go                 # 入口文件
├── go.mod                  # Go 模块
└── WEB_GUI_GUIDE.md       # Web GUI 使用指南
```

## 🛠 开发

### 前端开发

```bash
# 安装依赖
cd web
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

### 后端开发

```bash
# 运行（开发模式）
go run . serve

# 构建
go build -o soliton-gen .
```

### 更新静态文件

修改前端代码后：

```bash
cd web
npm run build
cd ..
rm -rf server/static/*
cp -r web/dist/* server/static/
go build -o soliton-gen .
```

## 🎯 技术栈

**后端：**
- Go 1.22+
- Gin Web Framework
- Cobra CLI
- embed（静态文件嵌入）

**前端：**
- Vue 3（Composition API）
- TypeScript
- Vite
- Vue Router

## 📝 更新日志

### v1.1.0 (2026-01-05)

**新增功能：**
- ✅ 字段备注功能（GUI + CLI）
- ✅ 字段排序功能（↑↓ 按钮）
- ✅ 枚举字段编辑支持
- ✅ 完整的领域删除（清理所有相关文件）
- ✅ 生成后自动运行 go mod tidy
- ✅ Dashboard 手动更新依赖按钮

### v1.0.0 (2026-01-04)

**新增功能：**
- ✅ Web GUI 可视化界面
- ✅ 项目初始化向导
- ✅ 领域编辑器（可视化字段配置）
- ✅ 服务编辑器（可视化方法配置）
- ✅ 代码预览功能
- ✅ 中英双语界面
- ✅ 详细的操作提示和使用指南
- ✅ 软删除支持
- ✅ 分页查询支持
- ✅ 错误码常量

**技术改进：**
- ✅ 核心逻辑重构到 `core` 包
- ✅ RESTful API 设计
- ✅ 静态文件嵌入到二进制
- ✅ 单文件部署

## 🤝 贡献

欢迎提交 Issue 和 PR！

**前端开发：**
- 组件位于 `web/src/views/`
- 遵循 Vue 3 Composition API 规范

**后端开发：**
- API 处理器位于 `server/handlers/`
- 核心逻辑位于 `core/`
- 模板位于 `core/templates_*.go`

## 📄 许可证

与 Soliton-Go 框架相同。
