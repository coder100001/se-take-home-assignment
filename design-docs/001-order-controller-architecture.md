# Design Doc 001: 订单控制器架构设计

## 元数据
- **编号**: 001
- **标题**: 订单控制器架构设计
- **状态**: draft
- **创建日期**: 2026-04-30
- **最后更新**: 2026-04-30
- **关联任务**: McDonald's Order Controller Implementation
- **复杂度级别**: L2
- **前置 Design Doc**: 无

## 1. 背景与动机

### 为什么需要这个改动？
McDonald's 在 COVID-19 期间需要构建自动化烹饪机器人系统来减少人力成本并提高效率。作为软件工程师，需要创建一个订单控制器来处理订单控制流程。

### 当前系统的痛点是什么？
- 这是一个全新的系统，没有现有实现
- 需要支持 VIP 优先级处理
- 需要支持动态增减机器人数量
- 需要处理并发场景（多个机器人同时处理订单）

### 业务/技术驱动因素
- **业务驱动**: 提高订单处理效率，减少人力成本
- **技术驱动**: 需要可靠的并发处理机制，确保订单不丢失、不重复处理

## 2. 调研与现状分析

### 2.1 现有实现
- 当前项目为空模板，只有基本的脚本框架
- GitHub Actions 工作流已配置，支持 Go 和 Node.js
- 需要从零开始实现所有功能

### 2.2 业界实践

#### 订单管理系统参考
1. **消息队列模式** (Kafka/RabbitMQ)
   - 优点: 解耦、可靠、可扩展
   - 缺点: 对于原型过于复杂
   
2. **优先级队列** (Priority Queue)
   - 优点: 简单、高效、适合内存处理
   - 缺点: 需要自己实现并发安全
   
3. **工作池模式** (Worker Pool)
   - 优点: 控制并发数、资源复用
   - 缺点: 需要处理任务分发逻辑

#### 并发模型参考
1. **Go goroutine + channel**
   - 天然支持并发
   - 轻量级协程
   - channel 通信安全
   
2. **Node.js async/await + Promise**
   - 单线程事件循环
   - 需要手动处理并发控制
   - 适合 I/O 密集型

### 2.3 技术约束
- **语言限制**: 必须使用 Go 或 Node.js
- **环境限制**: 需要在 GitHub Actions 中运行
- **输出要求**: 必须输出到 result.txt，包含时间戳
- **时间限制**: 建议在 30 分钟内完成（原型）
- **数据持久化**: 不需要，纯内存处理

## 3. 可选方案

### 方案 A: Go + Goroutine + 优先级队列

#### 描述
使用 Go 语言实现，每个机器人作为一个独立的 goroutine 运行。订单使用优先级队列管理，VIP 订单排在普通订单之前。

#### 架构设计
```
┌─────────────────────────────────────────┐
│           OrderController               │
├─────────────────────────────────────────┤
│  ┌──────────────┐  ┌────────────────┐  │
│  │ OrderManager │  │  BotManager    │  │
│  │  - Queue     │  │  - Bots[]      │  │
│  │  - Counter   │  │  - AddBot()    │  │
│  └──────────────┘  │  - RemoveBot() │  │
│         │          └────────────────┘  │
│         │                   │           │
│         ▼                   ▼           │
│  ┌──────────────────────────────────┐  │
│  │      Priority Queue               │  │
│  │  [VIP1, VIP2] [Normal1, Normal2] │  │
│  └──────────────────────────────────┘  │
│         │                              │
│         ▼                              │
│  ┌──────────────────────────────────┐  │
│  │    Bot Goroutines                │  │
│  │  Bot1 (goroutine) - Processing   │  │
│  │  Bot2 (goroutine) - IDLE         │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

#### 数据结构
```go
type OrderType int
const (
    Normal OrderType = iota
    VIP
)

type OrderStatus int
const (
    Pending OrderStatus = iota
    Processing
    Complete
)

type Order struct {
    ID        int
    Type      OrderType
    Status    OrderStatus
    CreatedAt time.Time
}

type Bot struct {
    ID           int
    Status       BotStatus
    CurrentOrder *Order
    StopChan     chan bool
}

type OrderQueue struct {
    VIPOrders    []*Order
    NormalOrders []*Order
    mu           sync.Mutex
}
```

#### 核心流程
1. **创建订单**: 
   - 生成唯一 ID
   - 根据类型插入队列（VIP 插入 VIP 队列尾部，Normal 插入 Normal 队列尾部）
   - 通知空闲机器人

2. **增加机器人**:
   - 创建 Bot 结构
   - 启动 goroutine
   - 立即尝试获取订单

3. **减少机器人**:
   - 标记为停止
   - 如果正在处理订单，将订单放回队列原位置
   - 等待 goroutine 结束

4. **机器人处理流程**:
   ```go
   func (b *Bot) Run() {
       for {
           select {
           case <-b.StopChan:
               return
           default:
               order := queue.Pop() // 获取最高优先级订单
               if order != nil {
                   b.Process(order)
               } else {
                   time.Sleep(100 * time.Millisecond) // 空闲等待
               }
           }
       }
   }
   ```

#### 优点
- ✅ Go 原生并发支持，goroutine 轻量级
- ✅ 类型安全，编译时检查
- ✅ 性能好，适合 CPU 密集型任务
- ✅ 标准库丰富，无需第三方依赖
- ✅ 编译为单一二进制文件，部署简单

#### 缺点
- ⚠️ 需要手动处理并发同步（mutex）
- ⚠️ 学习曲线相对陡峭（如果团队不熟悉 Go）

#### 工作量
- **预估**: 中等
- **文件数**: 4-5 个
- **代码行**: 300-400 行

---

### 方案 B: Node.js + async/await + 事件驱动

#### 描述
使用 Node.js 实现，采用事件驱动模型。每个机器人作为一个异步任务运行，使用事件发射器进行通信。

#### 架构设计
```
┌─────────────────────────────────────────┐
│           OrderController               │
├─────────────────────────────────────────┤
│  ┌──────────────┐  ┌────────────────┐  │
│  │ OrderManager │  │  BotManager    │  │
│  │  - Queue     │  │  - Bots[]      │  │
│  │  - Counter   │  │  - addBot()    │  │
│  └──────────────┘  │  - removeBot() │  │
│         │          └────────────────┘  │
│         │                   │           │
│         ▼                   ▼           │
│  ┌──────────────────────────────────┐  │
│  │      Event Emitter                │  │
│  │  - 'orderCreated'                │  │
│  │  - 'orderCompleted'              │  │
│  └──────────────────────────────────┘  │
│         │                              │
│         ▼                              │
│  ┌──────────────────────────────────┐  │
│  │    Bot Async Functions            │  │
│  │  Bot1 (async) - Processing        │  │
│  │  Bot2 (async) - IDLE              │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

#### 数据结构
```javascript
class Order {
    constructor(id, type) {
        this.id = id;
        this.type = type; // 'VIP' or 'Normal'
        this.status = 'PENDING';
        this.createdAt = new Date();
    }
}

class Bot {
    constructor(id) {
        this.id = id;
        this.status = 'IDLE';
        this.currentOrder = null;
        this.stopFlag = false;
    }
    
    async process(order) {
        this.status = 'PROCESSING';
        this.currentOrder = order;
        order.status = 'PROCESSING';
        
        await sleep(10000); // 10秒处理时间
        
        order.status = 'COMPLETE';
        this.status = 'IDLE';
        this.currentOrder = null;
    }
}
```

#### 核心流程
1. **创建订单**: 
   - 生成唯一 ID
   - 插入队列
   - 触发 'orderCreated' 事件

2. **增加机器人**:
   - 创建 Bot 实例
   - 启动异步处理循环

3. **减少机器人**:
   - 设置停止标志
   - 如果正在处理，将订单放回队列

4. **机器人处理流程**:
   ```javascript
   async run() {
       while (!this.stopFlag) {
           const order = queue.pop();
           if (order) {
               await this.process(order);
           } else {
               await sleep(100); // 空闲等待
           }
       }
   }
   ```

#### 优点
- ✅ 开发速度快，代码简洁
- ✅ 异步 I/O 天然支持
- ✅ 适合快速原型开发
- ✅ JSON 处理方便

#### 缺点
- ⚠️ 单线程，无法利用多核
- ⚠️ 需要小心处理异步错误
- ⚠️ 并发控制需要额外库支持
- ⚠️ 类型安全需要 TypeScript

#### 工作量
- **预估**: 小
- **文件数**: 3-4 个
- **代码行**: 200-300 行

---

### 方案 C: Go + Channel + Worker Pool

#### 描述
使用 Go 的 channel 实现 worker pool 模式。订单通过 channel 分发给机器人，机器人从 channel 中消费订单。

#### 架构设计
```
┌─────────────────────────────────────────┐
│           OrderController               │
├─────────────────────────────────────────┤
│  ┌──────────────────────────────────┐  │
│  │      Order Channel                │  │
│  │  (buffered, priority based)      │  │
│  └──────────────────────────────────┘  │
│         │                              │
│         ▼                              │
│  ┌──────────────────────────────────┐  │
│  │    Worker Pool                    │  │
│  │  Bot1 (worker) - Processing      │  │
│  │  Bot2 (worker) - IDLE            │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

#### 优点
- ✅ channel 天然线程安全
- ✅ 优雅的并发模型
- ✅ 易于扩展

#### 缺点
- ⚠️ channel 不支持优先级，需要额外逻辑
- ⚠️ 机器人销毁时，channel 中的订单难以回退
- ⚠️ 复杂度较高

#### 工作量
- **预估**: 大
- **文件数**: 5-6 个
- **代码行**: 400-500 行

---

## 4. 决策

### 选定方案
**方案 A: Go + Goroutine + 优先级队列**

### 决策理由
1. **并发支持**: Go 的 goroutine 天然适合处理多个机器人并发场景
2. **类型安全**: 编译时类型检查，减少运行时错误
3. **性能**: 编译型语言，性能优于 Node.js
4. **简单性**: 相比方案 C，优先级队列更容易实现 VIP 优先逻辑
5. **部署**: 编译为单一二进制文件，符合 GitHub Actions 要求
6. **可维护性**: 代码结构清晰，易于测试

### 权衡取舍
- **放弃了 Node.js 的快速开发优势**，换取了更好的并发控制和性能
- **放弃了 channel 的优雅性**，换取了更简单的优先级队列实现
- **需要手动处理 mutex 同步**，但这是可控的复杂度

## 5. 影响范围

### 影响的模块/包
- 新增 `cmd/main.go` - CLI 入口
- 新增 `internal/order/` - 订单管理
- 新增 `internal/bot/` - 机器人管理
- 新增 `internal/queue/` - 优先级队列
- 修改 `scripts/build.sh` - 编译脚本
- 修改 `scripts/test.sh` - 测试脚本
- 修改 `scripts/run.sh` - 运行脚本

### 影响的接口
- 新增 CLI 命令接口（交互式）
- 新增订单管理 API
- 新增机器人管理 API

### 影响的配置
- 无配置文件，纯命令行交互

### 影响的部署
- GitHub Actions 工作流已配置，无需修改
- 输出文件 `result.txt` 格式需要符合要求

## 6. 开放问题

- [ ] CLI 交互方式：使用 `bufio.Scanner` 还是 `cobra` 库？
  - **倾向**: 使用标准库 `bufio.Scanner`，避免引入第三方依赖
  
- [ ] 时间戳格式：使用 `time.Now().Format("15:04:05")` 还是其他格式？
  - **倾向**: 使用 `15:04:05` 格式，符合 HH:MM:SS 要求
  
- [ ] 测试策略：单元测试覆盖率目标？
  - **倾向**: >= 80%，核心逻辑 100%

- [ ] 错误处理：如何处理异常情况？
  - **倾向**: 使用 Go 的 error 返回值，记录日志

---

**文档状态**: approved
**下一步**: 等待评审通过后进入 Phase 1 (PLAN.md)
