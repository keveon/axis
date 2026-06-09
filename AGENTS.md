# Axis — AI Agent 协作指南

> **Axis** 是一个声明式物联网（IoT）数据平台，灵感来源于 Kubernetes 控制平面模型。
> 通过声明式 YAML 定义所有资源，实现物联网数据管道的全生命周期管理。

---

## 1. 项目概览

Axis 是一个开源 IoT 数据平台，核心理念是将 Kubernetes 的声明式控制平面模型应用于物联网数据管道。

### 设计哲学

- **声明式 API**：所有资源通过 CRD 风格的 YAML 定义，系统自动驱动至期望状态
- **控制平面与数据平面分离**：控制平面基于 etcd，数据平面可选 PostgreSQL/TimescaleDB
- **单一二进制，三种角色**：`axis` 二进制通过子命令扮演不同角色
- **组件即插件**：数据处理组件（Listener、Collector、Parser、Sink、Relay）是独立二进制，使用 Axis SDK 与控制平面交互

### 核心二进制：`axis`

| 角色 | 子命令 | 职责 |
|------|--------|------|
| **Server** | `axis server` | API Server + Scheduler + Controller Manager，控制平面核心 |
| **Agent** | `axis agent` | 机器注册、systemd 服务管理，部署在数据采集节点 |
| **Worker** | `axis worker` | 通用存储工作进程，处理数据持久化 |

### 数据处理组件（独立二进制，使用 Axis SDK）

| 组件 | 职责 |
|------|------|
| **Listener** | 监听网络端口，接收传入的数据流（如 TCP/UDP/HTTP） |
| **Collector** | 主动采集数据（如 SNMP 轮询、HTTP 拉取） |
| **Parser** | 解析原始数据为结构化格式（如 JSON 解析、协议解码） |
| **Sink** | 数据输出目标（如写入数据库、转发到外部系统） |
| **Relay** | 数据中继与转发，在组件间传递数据 |

---

## 2. 系统架构

```
┌─────────────────────────────────────────────────┐
│                  控制平面 (Control Plane)          │
│                                                  │
│   ┌──────────┐  ┌───────────┐  ┌──────────────┐ │
│   │ API Server│  │ Scheduler │  │Controller Mgr│ │
│   └────┬─────┘  └─────┬─────┘  └──────┬───────┘ │
│        │              │               │          │
│        └──────────────┼───────────────┘          │
│                       │                          │
│                   ┌───┴───┐                      │
│                   │  etcd  │  ← 所有资源状态存储   │
│                   └───────┘                      │
└──────────────────────┬──────────────────────────┘
                       │ Watch / Status Update
         ┌─────────────┼──────────────────┐
         │             │                  │
   ┌─────┴─────┐ ┌────┴─────┐   ┌────────┴────────┐
   │   Agent   │ │   Agent  │   │   Agent         │
   │ (Node 1)  │ │ (Node 2) │   │   (Node N)      │
   └─────┬─────┘ └────┬─────┘   └────────┬────────┘
         │            │                   │
    ┌────┴────┐  ┌────┴────┐        ┌────┴────┐
    │Listener │  │Collector│        │ Parser  │ ...
    │Collector│  │ Parser  │        │  Sink   │
    │  Sink   │  │  Sink   │        │         │
    └─────────┘  └─────────┘        └─────────┘
         │            │                   │
         └────────────┼──────────────────┘
                      │
              ┌───────┴───────┐
              │ 数据平面 (Data Plane) │
              │  PostgreSQL     │
              │  TimescaleDB    │
              └───────────────┘
```

### 关键设计决策

1. **etcd 作为控制平面存储**：所有资源对象的期望状态和实际状态都存储在 etcd 中
2. **可选数据平面数据库**：PostgreSQL（配合 TimescaleDB 扩展）用于时序数据持久化，非必需
3. **组件完全解耦**：每个数据处理组件是独立进程，通过 Axis SDK 与控制平面通信，不直接依赖核心代码
4. **核心保持最小化**：`axis` 核心只负责控制平面逻辑，不含任何数据处理实现

---

## 3. 资源类型

所有资源均通过声明式 YAML 定义，遵循 Kubernetes CRD 风格。

| 资源类型 | 说明 | 作用域 |
|----------|------|--------|
| **Protocol** | 协议定义，描述数据传输协议的编解码规则 | 集群全局 |
| **Pipeline** | 数据管道定义，描述数据从采集到存储的完整处理流程 | 集群全局 |
| **Store** | 存储配置，定义后端存储连接与 schema | 集群全局 |
| **Listener** | 监听器实例，部署于指定节点，接收传入数据流 | 节点级 |
| **Collector** | 采集器实例，部署于指定节点，主动拉取数据 | 节点级 |
| **Parser** | 解析器实例，部署于指定节点，解析原始数据 | 节点级 |
| **Sink** | 输出实例，部署于指定节点，将处理后的数据写入目标 | 节点级 |
| **Relay** | 中继实例，部署于指定节点，转发数据到其他组件 | 节点级 |

### 资源 YAML 示例结构

```yaml
apiVersion: axis.io/v1alpha1
kind: Listener
metadata:
  name: snmp-trap-listener
  namespace: default
spec:
  nodeSelector:
    matchLabels:
      role: collector
  protocol: snmp-v2c
  ports:
    - name: trap
      port: 162
      protocol: UDP
```

---

## 4. 目录结构

```
axis/
├── cmd/                    # 入口点
│   ├── axis/              # 主二进制（server / agent / worker）
│   │   └── main.go
│   └── components/        # 数据处理组件二进制
│       ├── listener/
│       ├── collector/
│       ├── parser/
│       ├── sink/
│       └── relay/
├── pkg/                    # 核心库代码
│   ├── api/               # API 定义（proto / OpenAPI）
│   ├── apis/              # API 资源类型定义（类似 k8s apiserver 的 types）
│   ├── server/            # Server 角色：API Server、Scheduler、Controller Manager
│   ├── agent/             # Agent 角色：机器注册、systemd 管理
│   ├── storeworker/       # Worker 角色实现
│   ├── controller/        # 控制器实现（调谐循环）
│   ├── scheduler/         # 调度器实现
│   └── client/            # 控制平面客户端（Axis SDK 的一部分）
├── sdk/                    # Axis SDK — 供组件开发者使用
│   └── go.mod             # 独立的 Go module
├── api/                    # API 相关文件（proto、OpenAPI spec）
├── deploy/                 # 部署配置
├── docs/                   # 文档
├── examples/               # 示例 YAML 和使用案例
├── go.mod                  # Go module 定义
├── go.sum
├── Makefile                # 构建、测试、Lint 入口
├── AGENTS.md               # AI Agent 协作指南（本文件）
└── README.md               # 项目说明
```

> **注意**：当前仓库处于早期阶段，目录结构将随开发进度逐步完善。上述结构为目标布局，开发时请参照此规划。

---

## 5. 严格开发规则

以下规则 **必须严格遵守**，任何 PR 或代码变更均需符合这些要求。

### 5.1 语言与工具链

| 项目 | 要求 |
|------|------|
| Go 版本 | **Go 1.22+**（当前 go.mod 使用 1.26.4） |
| 代码格式 | 必须通过 `gofmt` 和 `go vet` |
| Module 路径 | `github.com/keveon/axis` |
| 依赖管理 | Go Modules（`go.mod`） |

### 5.2 提交规范

使用 **Conventional Commits** 格式：

```
<type>(<scope>): <description>

[可选正文]

[可选 footer]
```

**允许的 type**：

| type | 用途 |
|------|------|
| `feat` | 新功能 |
| `fix` | Bug 修复 |
| `docs` | 文档变更 |
| `refactor` | 重构（不改变外部行为） |
| `chore` | 构建、工具、依赖等杂项变更 |
| `test` | 测试相关 |
| `ci` | CI/CD 配置变更 |

**示例**：
```
feat(server): add watch endpoint for Pipeline resources
fix(agent): handle systemd service restart timeout
docs(AGENTS.md): update resource type descriptions
refactor(scheduler): extract node scoring logic
```

### 5.3 代码质量

- ✅ 所有新代码 **必须包含测试**
- ✅ `gofmt` 格式化后无差异
- ✅ `go vet` 检查通过
- ✅ 所有测试通过：`go test ./...`
- ❌ 代码中 **禁止** 出现密钥、密码、token 等敏感信息
- ❌ 禁止引入任何公司相关的内容或内部代号

### 5.4 架构原则

- **核心最小化**：`axis` 核心二进制只包含控制平面逻辑
- **组件即插件**：数据处理组件（Listener、Collector、Parser、Sink、Relay）是独立二进制，通过 Axis SDK 通信
- **不硬编码组件**：核心代码中不应包含任何具体协议或数据格式的实现细节
- **SDK 独立**：Axis SDK 作为独立 module，组件开发者可以只依赖 SDK

---

## 6. 构建与测试命令

### 构建

```bash
# 构建主二进制
go build -o bin/axis ./cmd/axis

# 构建所有组件
go build ./cmd/components/...

# 构建特定组件
go build -o bin/listener ./cmd/components/listener
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试（带详细输出）
go test -v ./pkg/server/...

# 运行测试并显示覆盖率
go test -cover ./...
```

### 代码检查

```bash
# 格式检查（如有差异则失败）
gofmt -l .

# 静态分析
go vet ./...
```

### 常用开发流程

```bash
# 1. 格式化代码
gofmt -w .

# 2. 静态检查
go vet ./...

# 3. 运行测试
go test ./...

# 4. 构建
go build -o bin/axis ./cmd/axis
```

---

## 7. PR 与协作规范

### 创建 Pull Request

1. **分支命名**：
   - 功能分支：`feat/<简短描述>`
   - 修复分支：`fix/<简短描述>`
   - 文档分支：`docs/<简短描述>`
   - 重构分支：`refactor/<简短描述>`

2. **PR 标题格式**：遵循 Conventional Commits 格式，如：
   - `feat(server): implement Pipeline reconciliation loop`
   - `fix(agent): resolve node registration race condition`

3. **PR 描述模板**：
   ```markdown
   ## 变更内容
   简要描述本次变更的内容和目的。

   ## 变更类型
   - [ ] 新功能 (feat)
   - [ ] Bug 修复 (fix)
   - [ ] 文档 (docs)
   - [ ] 重构 (refactor)
   - [ ] 其他

   ## 测试
   - [ ] 已添加/更新相关测试
   - [ ] 所有测试通过 (`go test ./...`)
   - [ ] `gofmt` 和 `go vet` 通过

   ## 自查清单
   - [ ] 代码中无敏感信息
   - [ ] 核心代码保持最小化，未硬编码组件逻辑
   - [ ] 符合 Conventional Commits 规范
   ```

### Code Review 重点

- 控制平面逻辑是否属于核心职责（还是应放在组件中）
- 是否遵循声明式 API 设计模式
- 资源状态管理是否正确（期望状态 vs 实际状态）
- 错误处理和重试逻辑是否合理
- 测试覆盖率是否充分

---

## 8. 给 AI Agent 的特别说明

### 修改代码前

1. **先理解上下文**：阅读相关文件，理解当前架构和代码风格
2. **确认影响范围**：判断修改是否会影响其他模块
3. **遵守架构边界**：区分核心代码和组件代码，不要在核心中实现数据处理逻辑

### 编写代码时

1. **声明式思维**：资源管理应采用 "声明期望状态 → 控制器驱动至期望状态" 的模式
2. **最小变更原则**：每次变更应聚焦于单一目的
3. **保持一致性**：遵循项目已有的代码风格和命名规范

### 提交代码时

1. 使用 Conventional Commits 格式
2. 确保所有测试通过
3. 如涉及架构变更，在提交信息中详细说明理由

---

## 9. 参考资源

| 资源 | 链接 |
|------|------|
| 项目仓库 | [github.com/keveon/axis](https://github.com/keveon/axis) |
| Go Module | `github.com/keveon/axis` |
| Go 文档 | [go.dev/doc](https://go.dev/doc) |
| Conventional Commits | [conventionalcommits.org](https://www.conventionalcommits.org/) |

---

*本文档是 Axis 项目的 AI Agent 协作指南。如有疑问，请参考项目仓库中的其他文档或在 Issue 中讨论。*
