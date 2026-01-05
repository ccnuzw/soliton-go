# Soliton-Gen Web GUI 使用指南

## 快速开始

### 场景 1：创建新项目

```bash
# 1. 在空目录或任意工作目录启动 Web GUI
mkdir my-workspace
cd my-workspace
./soliton-gen serve

# 2. 访问 http://127.0.0.1:3000

# 3. 点击"初始化项目 Init Project"
#    - 填写项目名称（如 my-project）
#    - 填写模块名称（如 github.com/yourname/my-project）
#    - 点击"预览"查看将生成的文件
#    - 点击"创建项目"

# 4. 项目创建完成后，进入项目目录
cd my-project

# 5. 再次启动 GUI 进行后续开发
soliton-gen serve
```

### 场景 2：在现有项目中使用

```bash
# 在现有项目根目录（包含 go.mod 和 internal/）启动
cd /path/to/existing/project
./soliton-gen serve

# 访问 http://127.0.0.1:3000
# 直接使用"生成领域"、"领域增强"、"生成服务"、"迁移中心"等功能
```

### 自定义端口和主机

```bash
./soliton-gen serve --port 8080 --host 0.0.0.0
```

---

## 功能说明

### 1. Dashboard（首页）

首页提供项目概览和快速导航。

**功能特性：**
- 自动检测当前目录是否为有效的 Go 项目
- 显示项目模块路径
- 提供五个主要功能入口
- 新增 DDD 组件入口（Value Object / Spec / Policy / Event / Handler）
- 新增迁移中心入口（详细日志 / 历史记录 / 筛选导出）
- **使用指南**：点击展开查看快速开始步骤和使用提示

**项目检测：**
- ✓ 已检测到项目：显示模块路径，可以使用所有功能
- ! 未找到项目：提示需要在包含 go.mod 和 internal/ 目录的项目中运行

---

### 2. Init Project（初始化项目）

创建一个新的 Soliton-Go 项目骨架。

#### 使用流程

**步骤 1：配置 Configure**

点击"💡 配置说明"查看详细帮助，包含：
- 项目名称格式建议
- 模块名称格式说明
- 框架替换路径用途

表单字段：
- **项目名称 Project Name** *（必填）
  - 说明：项目目录名称，将创建此名称的文件夹
  - 示例：`my-awesome-project`
  - 建议：使用小写字母和连字符

- **模块名称 Module Name**（可选）
  - 说明：Go 模块的导入路径
  - 默认：`github.com/soliton-go/<项目名>`
  - 示例：`github.com/yourname/my-project`
  - 格式：通常为 `github.com/username/project`

- **框架替换路径 Framework Replace**（可选）
  - 说明：本地开发时使用，指向 soliton-go/framework 的路径
  - 示例：`../framework` 或 `/path/to/framework`
  - 用途：仅在本地开发 Soliton 框架时使用

**步骤 2：预览 Preview**
- 查看将要创建的文件列表
- 确认配置信息
- 文件状态标识：
  - `NEW` - 新建文件
  - `SKIP` - 跳过（已存在）

**步骤 3：完成 Done**
- 显示生成结果
- 提供下一步操作指引
- 包含示例命令

#### 生成的文件结构

```
<项目名>/
├── cmd/
│   ├── main.go              # 应用入口
│   └── migrate/
│       └── main.go          # 迁移入口
├── configs/
│   └── config.yaml          # 配置文件
├── internal/
│   ├── application/         # 应用层
│   ├── domain/              # 领域层
│   ├── infrastructure/      # 基础设施层
│   └── interfaces/          # 接口层
├── go.mod                   # Go 模块定义
├── Makefile                 # 构建脚本
├── README.md                # 项目说明
└── .gitignore               # Git 忽略规则
```

---

### 3. Domain Editor（领域编辑器）

生成 DDD 风格的领域模块，包含完整的 CRUD 功能。

#### 使用指南

点击"📖 使用指南"查看：
- 领域名称格式（PascalCase）
- 字段类型详细说明
- 自动注入功能说明

#### 基本配置

**领域名称 Domain Name** *（必填）
- 说明：实体名称，将生成对应的 Go 结构体
- 格式：使用 PascalCase
- 示例：`User`、`Order`、`Product`

#### 字段配置

点击"+ 添加字段"添加新字段，每个字段包含：

1. **字段名**
   - 示例：`username`、`email`、`status`
   - 格式：使用 snake_case

2. **字段类型**
   - `string` - 字符串 (varchar 255)
   - `text` - 长文本
   - `int` - 整数 (32-bit)
   - `int64` - 整数 (64-bit)
   - `float64` - 浮点数 (64-bit)
   - `decimal` - 精确小数 (10,2 精度，适合金额)
   - `bool` - 布尔值
   - `time` - 时间戳
   - `time?` - 可选时间戳
   - `date` - 日期 (无时间部分)
   - `date?` - 可选日期
   - `uuid` - UUID (带索引)
   - `json` - JSON 对象 (需 gorm.io/datatypes)
   - `jsonb` - PostgreSQL JSONB
   - `bytes` - 二进制数据
   - `enum` - 枚举（需要提供枚举值）

3. **枚举值**（仅当类型为 enum 时）
   - 格式：用 `|` 分隔
   - 示例：`active|inactive|banned`
   - 提示：Hover 在输入框上查看格式说明

#### 选项说明

- **✓ 启用软删除 Soft Delete**
  - 说明：启用后将添加 DeletedAt 字段，删除时标记而非真删除
  - 适用场景：需要保留历史数据

- **✓ 自动注入到 main.go**
  - 说明：自动在 main.go 中注册此模块
  - 推荐：勾选此项可省去手动配置

- **✓ 强制覆盖 Force**
  - 说明：覆盖已存在的文件
  - 警告：会覆盖现有代码，请谨慎使用

> **提示**：领域生成完成后会自动执行一次 `go mod tidy`，可在首页“更新依赖”卡片手动再次执行。

#### 高级选项

点击"高级选项 Advanced"展开：

- **自定义表名 Table Name**
  - 默认：名称的复数形式（如 User → users）
  - 用途：自定义数据库表名

- **自定义路由前缀 Route Base**
  - 默认：名称的复数形式（如 User → /users）
  - 用途：自定义 HTTP 路由前缀

#### 生成的文件

```
internal/
├── domain/<name>/
│   ├── <name>.go           # 实体定义
│   ├── repository.go       # 仓储接口
│   └── events.go           # 领域事件
├── application/<name>/
│   ├── commands.go         # 命令（写操作）
│   ├── queries.go          # 查询（读操作）
│   ├── dto.go              # 数据传输对象
│   └── module.go           # 模块注册
├── infrastructure/persistence/
│   └── <name>_repo.go      # 仓储实现
└── interfaces/http/
    └── <name>_handler.go   # HTTP 处理器
```

#### 示例：创建 User 领域

1. 领域名称：`User`
2. 添加字段：
   - `username` - `string`
   - `email` - `string`
   - `status` - `enum` - `active|inactive|banned`
   - `created_at` - `time`
3. 勾选"自动注入到 main.go"
4. 点击"预览 Preview"查看将生成的文件
5. 点击"生成 Generate"执行生成

---

### 4. DDD Enhancer（领域增强）

生成 Value Object（值对象）、Specification（领域规格）、Policy（领域策略）、Event 和 Handler。

#### 功能说明

- **Value Object**：用于表达不可变的业务概念，如金额、邮箱、地址
- **Specification**：用于封装业务规则判断，便于复用与组合
- **Policy**：用于表达业务策略或决策逻辑
- **Event**：领域事件，描述领域中已发生的业务事实
- **Handler**：事件处理器，订阅并处理领域事件

#### 使用流程

1. 选择已有领域（Domain）
2. 选择需要生成的组件：
   - **Value Object**：定义值对象字段
   - **Specification**：定义规格，Target 可留空表示 any
   - **Policy**：定义策略，Target 可留空表示 any
   - **Event & Handler**：支持组合生成（Event + Handler）
3. 设置可选参数：
   - Topic（留空自动生成）
   - Force（可选，覆盖已有文件）
4. 点击"预览"确认生成结果
5. 点击"生成"执行生成

#### 组合生成说明

- 可以仅生成 Event 或仅生成 Handler
- 同时勾选时会依次生成 Event 与 Handler
- Handler Topic 留空时会复用 Event Topic
- 默认注入 EventBus Provider 与事件注册

---

### 5. Service Editor（服务编辑器）

生成应用服务层，用于编排跨领域的业务逻辑。

#### 使用指南

点击"📖 使用指南"查看：
- 服务名称格式说明
- 方法定义规则
- 默认方法说明
- Service 用途解释

#### 基本配置

**服务名称 Service Name** *（必填）
- 说明：应用服务名称，用于跨领域业务逻辑
- 格式：使用 PascalCase，自动添加 "Service" 后缀
- 示例：`OrderService`、`PaymentService`

#### 方法配置

点击"+ 添加方法"添加新方法：

- **方法名**
  - 格式：使用 PascalCase
  - 示例：`CreateOrder`、`ProcessPayment`、`CancelOrder`
  - 每行一个方法名

- **默认方法**
  - 如果不填写方法，将自动生成：
    - `Create` - 创建
    - `Get` - 获取单个
    - `List` - 获取列表

#### 选项说明

- **✓ 强制覆盖 Force**
  - 说明：覆盖已存在的文件
  - 警告：会覆盖现有代码，请谨慎使用

#### 生成的文件

```
internal/application/<name>/
└── <name>_service.go       # 服务实现
```

#### 下一步 Next Steps

生成后需要：
1. 在服务结构体中注入所需的 Repository
2. 在每个方法中实现业务逻辑
3. 在 main.go 中注册服务

#### 示例：创建 Order 服务

1. 服务名称：`Order`（自动变为 OrderService）
2. 添加方法：
   - `CreateOrder`
   - `ProcessPayment`
   - `CancelOrder`
3. 点击"预览 Preview"查看代码
4. 点击"生成 Generate"执行生成

---

### 6. Migration Center（迁移中心）

用于执行数据库迁移，并查看详细日志与历史记录。

#### 功能说明

- ✅ 支持执行迁移命令并返回完整日志
- ✅ 支持执行前自动 `go mod tidy`
- ✅ 支持日志筛选（INFO / ERROR / SYSTEM / TIDY / MIGRATE）
- ✅ 支持复制与下载日志
- ✅ 支持历史记录查看与回溯

#### 使用流程

1. 打开「迁移中心 Migration Center」
2. 确认检测到当前项目路径
3. 选择是否先执行 `go mod tidy`
4. 设置迁移超时（默认 300 秒）
5. 点击「开始迁移」并确认执行
6. 查看日志与结果状态

#### 日志说明

- **SYSTEM**：系统级日志（启动/结束/命令）
- **TIDY**：依赖整理日志
- **MIGRATE**：迁移执行日志
- **INFO / ERROR**：信息与错误级别

> **命令选择规则**：优先执行 `cmd/migrate/main.go`，不存在时回退 `cmd/migrate.go`。

---

## 界面特性

### 中英双语

所有界面采用中英双语显示：
- 标题和主要功能：中文 + 英文
- 技术术语：保留英文（如 Repository、CQRS）
- 说明文字：以中文为主

### 操作提示

每个页面都提供详细的操作提示：

1. **可折叠的使用指南**
   - 位于页面顶部
   - 点击展开/收起
   - 包含详细的使用说明和示例

2. **提示图标 ⓘ**
   - 位于字段标签旁边
   - Hover 显示详细说明
   - 帮助理解字段用途

3. **Placeholder 示例**
   - 输入框中提供多个示例
   - 展示正确的格式
   - 降低使用门槛

4. **Checkbox 提示**
   - Hover 在复选框上显示说明
   - 解释选项的作用
   - 提供使用建议

### 预览功能

所有生成操作都支持预览：
- 查看将要创建的文件列表
- 确认配置信息
- 避免误操作

### 领域增强组件

DDD 增强页面支持：
- Value Object / Specification / Policy / Event / Handler
- Event 与 Handler 组合生成

### 迁移中心

迁移中心支持：
- 详细日志、筛选、复制与下载
- 迁移前自动 tidy
- 历史记录回溯

### 状态反馈

- **文件状态**：NEW、OVERWRITE、SKIP
- **成功提示**：✅ 已生成文件
- **错误提示**：❌ 错误信息
- **加载状态**：加载中...、生成中...

---

## 开发模式

如果需要修改前端代码：

### 启动开发环境

```bash
# 终端 1：启动后端
cd tools/generator
go run . serve --port 3000

# 终端 2：启动前端开发服务器
cd web
npm run dev
```

访问 http://localhost:5173（前端开发服务器会自动代理 API 请求到后端）

### 构建前端

修改前端代码后：

```bash
cd web
npm run build
cd ..
rm -rf server/static/*
cp -r web/dist/* server/static/
go build -o soliton-gen .
```

---

## 故障排除

### 问题：访问页面显示 "No Project Found"

**原因**：当前目录不是有效的 Go 项目

**解决**：
- 确保当前目录包含 `go.mod` 文件
- 确保存在 `internal/` 目录
- 或者使用 "Init Project" 功能创建新项目

### 问题：API 请求失败

**检查**：
1. 后端服务是否正常运行
2. 端口是否被占用
3. 浏览器控制台是否有错误信息

### 问题：生成的代码无法编译

**检查**：
1. 是否在正确的项目目录下运行
2. 字段配置是否正确
3. 枚举值格式是否正确（用 `|` 分隔）
4. 运行 `go mod tidy` 更新依赖

### 问题：静态资源加载失败（404）

**原因**：embed.FS 配置问题

**解决**：
- 确保 `server/server.go` 中使用 `//go:embed all:static`
- 重新构建：`go build -o soliton-gen .`

### 问题：页面显示但功能不工作

**检查**：
1. 浏览器控制台是否有 JavaScript 错误
2. 检查 API 响应是否正常
3. 清除浏览器缓存后重试

---

## 最佳实践

### 项目初始化

1. 在空目录中运行 `soliton-gen serve`
2. 使用 "Init Project" 创建项目骨架
3. 进入项目目录：`cd <项目名>`
4. 运行 `GOWORK=off go mod tidy` 下载依赖

### 领域生成

1. 先规划好领域模型和字段
2. 使用"预览"功能确认生成内容
3. 首次生成时勾选"自动注入到 main.go"
4. 后续修改使用"强制覆盖"选项

### 服务生成

1. 明确服务的职责和方法
2. 服务名称要能体现业务含义
3. 方法名使用动词开头（如 Create、Process、Cancel）
4. 生成后及时实现业务逻辑

### 代码维护

1. 生成代码后立即提交 Git
2. 手动修改的代码不要使用"强制覆盖"
3. 定期运行 `go mod tidy` 清理依赖
4. 使用 `make` 命令统一构建流程

---

## 技术栈

**后端**：
- Go 1.22+
- Gin Web Framework
- Cobra CLI
- embed (静态文件嵌入)

**前端**：
- Vue 3 (Composition API)
- TypeScript
- Vite
- Vue Router

**工具**：
- npm/npx
- go modules

---

## 更新日志

### v1.0.0 (2026-01-04)

**新增功能**：
- ✅ Web GUI 界面
- ✅ 项目初始化向导
- ✅ 可视化领域编辑器
- ✅ 服务生成器
- ✅ 代码预览功能
- ✅ 中英双语界面
- ✅ 详细的操作提示和使用指南
- ✅ 可折叠的帮助文档
- ✅ Tooltip 提示图标
- ✅ 丰富的示例和格式说明

**技术改进**：
- ✅ 核心逻辑重构到 `core` 包
- ✅ RESTful API 设计
- ✅ 静态文件嵌入到二进制
- ✅ 单文件部署

---

## 技术支持

如有问题，请查看：
- [实现计划](file:///Users/mac/.gemini/antigravity/brain/d073472f-baf0-4d20-bc49-70a176ee09ea/implementation_plan.md)
- [开发总结](file:///Users/mac/.gemini/antigravity/brain/d073472f-baf0-4d20-bc49-70a176ee09ea/walkthrough.md)

---

## 贡献指南

欢迎贡献代码和建议！

**前端开发**：
- 组件位于 `web/src/views/`
- 样式使用 CSS Variables
- 遵循 Vue 3 Composition API 规范

**后端开发**：
- API 处理器位于 `server/handlers/`
- 核心逻辑位于 `core/`
- 模板位于 `core/templates_*.go`

**文档更新**：
- 使用指南：`WEB_GUI_GUIDE.md`
- API 文档：待补充
- 代码注释：使用中文
