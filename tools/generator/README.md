# Soliton-Gen 代码生成器

Soliton-Gen 是 Soliton-Go 框架的代码生成工具，支持命令行和 Web GUI 两种使用方式。

## 🚀 快速开始

### 编译

```bash
cd tools/generator
go build -o soliton-gen .
```

### 使用 Web GUI（推荐）

```bash
# 启动 Web 界面
./soliton-gen serve

# 自定义端口
./soliton-gen serve --port 8080 --host 0.0.0.0
```

访问 http://127.0.0.1:3000 即可使用可视化界面。

### 使用命令行

```bash
# 初始化项目
./soliton-gen init my-project

# 生成领域模块
./soliton-gen domain User --fields "username,email,status:enum(active|inactive)" --wire

# 生成应用服务
./soliton-gen service Order --methods "CreateOrder,ProcessPayment" --wire
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
- 支持多种字段类型（string、int、enum 等）
- 枚举值可视化配置
- 软删除选项
- 自动注入到 main.go
- 代码预览功能

### 4. Service Editor（服务编辑器）
- 可视化方法配置
- 默认方法生成
- 代码预览功能

## 🔌 命令列表

| 命令 | 说明 | 示例 |
|------|------|------|
| `init <name>` | 初始化新项目 | `soliton-gen init my-project` |
| `domain <name>` | 生成领域模块 | `soliton-gen domain User --fields "username,email"` |
| `service <name>` | 生成应用服务 | `soliton-gen service Order --methods "CreateOrder"` |
| `serve` | 启动 Web GUI | `soliton-gen serve --port 3000` |

### Domain 命令参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `--fields` | 字段列表 | `--fields "username,email,age:int"` |
| `--table` | 自定义表名 | `--table "sys_users"` |
| `--route` | 自定义路由 | `--route "/v1/users"` |
| `--soft-delete` | 启用软删除 | `--soft-delete` |
| `--wire` | 自动注入到 main.go | `--wire` |
| `--force` | 强制覆盖 | `--force` |

### Service 命令参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `--methods` | 方法列表 | `--methods "Create,Update,Delete"` |
| `--force` | 强制覆盖 | `--force` |

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
