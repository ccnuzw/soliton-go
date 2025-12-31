# Soliton-Gen 代码生成器使用指南

**Soliton-Gen** 是 Soliton-Go 框架配套的命令行工具 (CLI)，旨在自动化生成符合 DDD 规范的样板代码，减少重复劳动，并确保项目结构的一致性。

---

## 1. 🛠 安装与编译

代码生成器的源码位于 `tools/generator` 目录下。在使用前，您需要将其编译为可执行文件。

### 1.1 编译
在项目根目录下执行：

```bash
cd tools/generator
go build -o soliton-gen main.go
```
执行成功后，当前目录下会生成一个名为 `soliton-gen` (Windows 下为 `soliton-gen.exe`) 的可执行文件。

### 1.2配置环境变量 (可选)
为了方便在任何目录下使用，您可以将 `tools/generator` 目录添加到系统的 `PATH` 环境变量中，或者将编译好的 `soliton-gen` 移动到 `/usr/local/bin` 等全局路径下。

---

## 2. 💻 命令详解

目前 `soliton-gen` 支持基于 Cobra 的子命令体系。

### 2.1 查看帮助
运行不带参数的命令或 `-h` 查看帮助信息：
```bash
./soliton-gen -h
```

### 2.2 生成领域对象 (`domain`)
该命令用于生成 DDD 领域层的核心文件，包括聚合根 (Entity) 和仓储接口 (Repository Interface)。

**语法**:
```bash
./soliton-gen domain [EntityName]
```

**参数**:
*   `EntityName`: 实体名称（首字母建议大写），例如 `Order`, `Product`, `User`。

**示例**:
```bash
./soliton-gen domain Product
```

Running this command will create **5 files** across Domain and Infrastructure layers:

```text
# Domain Layer (Interfaces)
Created internal/domain/Product/Product.go
Created internal/domain/Product/repository.go
Created internal/domain/Product/mapper.go

# Infrastructure Layer (Implementations)
Created internal/infrastructure/persistence/Product_repo.go
Created internal/infrastructure/persistence/Product_mapper.go
```

**生成文件结构 (预览)**:

1.  **Domain Layer**:
    *   `Product.go`: 聚合根定义 (包含 `ProductID` 强类型)。
    *   `repository.go`: 仓储接口定义。
    *   `mapper.go`: SQLMapper 接口定义。
2.  **Infrastructure Layer**:
    *   `Product_repo.go`: 自动实现 Repository 接口 (GORM).
    *   `Product_mapper.go`: 自动实现 Mapper 接口 (GORM).

---

## 3. ⚙️ 二次开发与定制

目前的 `soliton-gen` 提供了一个基于 `spf13/cobra` 的 CLI 骨架。为了适应您的具体业务需求（例如公司特定的代码规范、版权头等），您可能需要对其进行定制。

### 3.1 目录结构
```text
tools/generator/
├── go.mod          # 依赖定义
├── main.go         # 入口文件
└── cmd/            # 命令实现目录
    ├── root.go     # 根命令 (rootCmd)
    └── domain.go   # domain 子命令实现
```

### 3.2 添加新模板
如果您希望生成 `Infrastructure` 层代码（如 GORM 实现），可以参考 `cmd/domain.go` 添加一个新的 `infra.go` 命令：

1.  **复制 `cmd/domain.go`** 为 `cmd/infra.go`。
2.  **修改 Command Use**: `Use: "infra [name]"`。
3.  **编写生成逻辑**: 使用 Go 的 `text/template` 库来渲染 `.go` 文件内容，并使用 `os.WriteFile` 写入到 `application/internal/infrastructure/persistence` 目录。

### 3.3 模板示例 (Go Template)
建议在 `tools/generator/templates` 下管理您的 `.tmpl` 文件，例如：

```go
package {{.PackageName}}

import "github.com/soliton-go/framework/ddd"

type {{.EntityName}} struct {
    ddd.BaseAggregateRoot
    ID string `gorm:"primaryKey"`
}
```

然后在 `domain.go` 中解析并渲染该模板。

---

## 4. 🛡 安全生成 (覆盖保护)

生成器内置了 **覆盖保护机制 (Lock)**，确保您手写的代码不会被意外覆盖。

*   **机制**: 在生成文件前，工具会检查目标路径是否存在文件。
*   **行为**:
    *   如果文件**不存在**：正常生成，输出 `[CREATED]`。
    *   如果文件**已存在**：自动跳过，输出 `[LOCK] Skipping`。
*   **如何强制重新生成**: 如果您确实需要重置文件，请手动删除旧文件后再运行命令。

## 5. 最佳实践

1.  **先设计，后生成**: 在运行生成器前，先规划好聚合根的名字。
2.  **增量生成**: 利用内置的覆盖保护，您可以放心地在项目中多次运行生成命令，补充缺失的文件而不破坏已有逻辑。
3.  **结合 IDE**: 生成代码后，建议使用 IDE (如 VS Code, GoLand) 的格式化工具 (`gofmt`) 美化代码。

---

> **注意**: 当前版本的生成器主要用于演示 CLI 结构和工作流。在实际生产使用中，建议将 `fmt.Println` 替换为真实的 `text/template` 渲染和文件 I/O 操作。
