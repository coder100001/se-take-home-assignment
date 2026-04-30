# 功能实现计划: McDonald's 订单控制器

## 1. 需求概述
- **功能名称**: McDonald's 自动化烹饪机器人订单控制器
- **目标**: 实现一个交互式 CLI 订单控制系统，支持订单创建（普通/VIP）、机器人管理（增/减）、优先级队列处理
- **优先级**: P0
- **复杂度级别**: L2 (中等)
- **Design Doc**: [001 订单控制器架构设计](design-docs/001-order-controller-architecture.md)
- **技术选型**: Go (Golang) + Goroutine + 优先级队列

## 2. 变更范围

### 2.1 新增文件
| 文件路径 | 用途 | 行数预估 |
|---------|------|---------|
| `cmd/main.go` | CLI 交互式入口，命令解析 | ~80 |
| `internal/order/order.go` | 订单结构与订单管理器 | ~80 |
| `internal/queue/priority_queue.go` | 优先级队列（VIP > Normal） | ~80 |
| `internal/bot/bot.go` | 机器人结构与机器人管理器 | ~120 |
| `internal/controller/controller.go` | 控制器：协调订单和机器人 | ~80 |
| `internal/controller/controller_test.go` | 控制器单元测试 | ~100 |
| `internal/queue/priority_queue_test.go` | 队列单元测试 | ~60 |
| `internal/bot/bot_test.go` | 机器人单元测试 | ~80 |
| `go.mod` | Go 模块定义 | ~10 |

### 2.2 修改文件
| 文件路径 | 变更类型 | 影响范围 |
|---------|---------|---------|
| `scripts/build.sh` | 填充构建命令 | 编译脚本 |
| `scripts/test.sh` | 填充测试命令 | 测试脚本 |
| `scripts/run.sh` | 填充运行命令 | 运行脚本 |

### 2.3 删除文件
无

## 3. 逻辑设计

### 3.1 CLI 交互命令
```
命令列表:
  normal <description>   - 创建普通订单
  vip <description>      - 创建 VIP 订单
  +bot                   - 增加一个机器人
  -bot                   - 减少一个机器人
  status                 - 查看当前状态
  quit                   - 退出
```

### 3.2 核心数据流
```
用户输入命令
    │
    ├─ normal/vip order → OrderManager.AddOrder(type, desc)
    │       └── OrderQueue.Enqueue(order) ── 根据类型插入队列
    │       └── Notify idle bots ── 如果有空闲机器人，触发处理
    │
    ├─ +bot → BotManager.AddBot()
    │       └── 创建 Bot goroutine
    │       └── Bot.Pickup() ── 立即尝试从队列获取订单
    │
    ├─ -bot → BotManager.RemoveBot()
    │       └── 停止最新的机器人
    │       └── 如果机器人在处理中，订单回队
    │
    └─ status → 打印所有订单状态和机器人状态
```

### 3.3 订单队列优先级规则
```
队列结构:
  VIP队列:    [VIP#1, VIP#2, ...]  ← 新VIP插入尾部
  Normal队列: [Normal#1, Normal#2, ...]  ← 新Normal插入尾部

出队规则:
  1. 优先从 VIP 队列头部取
  2. VIP 队列为空时从 Normal 队列头部取
```

### 3.4 机器人状态机
```
IDLE ──→ PICK_ORDER ──→ PROCESSING (10秒) ──→ COMPLETE ──→ IDLE
  ↑                                                          │
  └──────────────────────────────────────────────────────────┘
   (无订单时保持IDLE，等待新订单通知)
```

### 3.5 边界条件
| 场景 | 输入 | 预期输出 | 处理方式 |
|-----|------|---------|---------|
| -Bot 减少正在处理的机器人 | `-bot` 命令 | 销毁最新bot，处理中的订单回队列原位置 | 发送 stop 信号，等待当前处理完成后回队 |
| -Bot 减少 IDLE 机器人 | `-bot` 命令 | 直接销毁最新bot | 立即终止 goroutine |
| 无订单时 +Bot | `+bot` 命令 | 创建 bot，变为 IDLE | 等待新订单通知 |
| 订单被多次处理 | 并发场景 | 每个订单只被一个 bot 处理一次 | 使用 mutex 保护出队操作 |
| VIP 订单优先级 | 混合订单 | VIP 始终先于 Normal | 双队列 + 优先级出队逻辑 |

## 4. 包结构设计

```go
se-take-home-assignment/
├── cmd/
│   └── main.go              // CLI 入口, 交互式命令循环
├── internal/
│   ├── controller/
│   │   ├── controller.go     // 控制器: 协调订单和机器人
│   │   └── controller_test.go
│   ├── order/
│   │   └── order.go          // 订单结构, 订单类型, 状态
│   ├── queue/
│   │   ├── priority_queue.go // 优先级队列 (VIP + Normal)
│   │   └── priority_queue_test.go
│   └── bot/
│       ├── bot.go            // 机器人结构, 生命周期管理
│       └── bot_test.go
├── scripts/
│   ├── build.sh
│   ├── run.sh
│   ├── test.sh
│   └── result.txt
├── go.mod
├── PLAN.md
└── README.md
```

## 5. 依赖分析

### 5.1 内部依赖
- `cmd/main.go` → `internal/controller`
- `internal/controller` → `internal/order`, `internal/queue`, `internal/bot`
- `internal/queue` → 无内部依赖（纯数据结构）
- `internal/order` → 无内部依赖（纯数据结构）
- `internal/bot` → `internal/order`, `internal/queue`

### 5.2 外部依赖
- 标准库 only：`fmt`, `bufio`, `os`, `strings`, `sync`, `time`, `strconv`
- 零第三方依赖

### 5.3 循环依赖风险
- [x] 无风险（依赖方向单向：main → controller → order/queue/bot）

## 6. 测试策略

### 6.1 单元测试
- 覆盖率目标: >= 80%
- `queue/priority_queue_test.go`: 队列入队出队顺序、VIP 优先级、空队列、并发安全
- `bot/bot_test.go`: 机器人创建、处理订单、取消处理、空闲状态
- `controller/controller_test.go`: 完整流程场景（创建订单→处理→完成，增减机器人）

### 6.2 集成测试
- `scripts/test.sh` 运行 `go test ./... -v -race`
- 使用 `-race` 检测竞态条件

### 6.3 E2E 验证
- `scripts/run.sh` 执行 CLI 并验证 `result.txt` 输出格式
- 输出必须包含 HH:MM:SS 时间戳

## 7. 验收标准
- [ ] `normal` / `vip` 命令能正确创建订单并显示在 PENDING 区域
- [ ] VIP 订单优先级高于 Normal 订单
- [ ] `+bot` 创建机器人并立即处理 PENDING 订单
- [ ] 机器人处理一个订单需 10 秒，完成后移到 COMPLETE
- [ ] `-bot` 销毁最新机器人，处理中的订单回队
- [ ] 无订单时机器人 IDLE，新订单到来后自动处理
- [ ] 订单编号唯一且递增
- [ ] 所有输出包含 HH:MM:SS 时间戳
- [ ] `result.txt` 内容符合要求
- [ ] 所有测试通过（`go test ./... -v -race`）

## 8. 回滚策略
- [x] 纯内存实现，无数据库迁移
- [x] Git revert 可回滚所有代码变更

---

**计划状态**: 待评审
**创建时间**: 2026-04-30
