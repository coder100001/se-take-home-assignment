# Tasks

## Phase 1: 项目初始化
- [x] Task 1: 初始化 Go 模块
  - [x] SubTask 1.1: 创建 go.mod 文件，定义模块名称和 Go 版本
  - [x] SubTask 1.2: 验证 Go 环境配置正确

## Phase 2: 核心数据结构
- [x] Task 2: 实现订单结构和常量定义
  - [x] SubTask 2.1: 创建 internal/order/order.go
  - [x] SubTask 2.2: 定义 OrderType 常量（Normal, VIP）
  - [x] SubTask 2.3: 定义 OrderStatus 常量（Pending, Processing, Complete）
  - [x] SubTask 2.4: 定义 Order 结构体（ID, Type, Status, CreatedAt, Description）
  - [x] SubTask 2.5: 实现 OrderManager 结构体和方法（ID 生成器, 状态查询）
  - [x] SubTask 2.6: 明确 OrderManager 职责：管理订单元数据和全局 ID 计数器；OrderQueue 负责排序和存取

- [x] Task 3: 实现优先级队列
  - [x] SubTask 3.1: 创建 internal/queue/priority_queue.go
  - [x] SubTask 3.2: 实现 OrderQueue 结构体（VIPOrders, NormalOrders, mutex）
  - [x] SubTask 3.3: 实现 Enqueue 方法（根据订单类型插入对应队列）
  - [x] SubTask 3.4: 实现 Dequeue 方法（优先返回 VIP 订单头）
  - [x] SubTask 3.5: 实现 RequeueAt 方法（订单回队到原索引位置，保持优先级顺序）
  - [x] SubTask 3.6: 实现 IsEmpty 方法
  - [x] SubTask 3.7: 实现 PendingOrders 方法（返回状态快照用于显示）

## Phase 3: 机器人管理
- [x] Task 4: 实现机器人结构和管理器
  - [x] SubTask 4.1: 创建 internal/bot/bot.go
  - [x] SubTask 4.2: 定义 BotStatus 常量（IDLE, Processing）
  - [x] SubTask 4.3: 定义 Bot 结构体（ID, Status, CurrentOrder, StopChan）
  - [x] SubTask 4.4: 实现 Bot.Run 方法（goroutine 主循环）
  - [x] SubTask 4.5: 实现 Bot.Process 方法（处理订单 10 秒）
  - [x] SubTask 4.6: 实现 Bot.Stop 方法（停止 goroutine）
  - [x] SubTask 4.7: 实现 BotManager 结构体和方法（AddBot, RemoveBot）

## Phase 4: 控制器
- [x] Task 5: 实现控制器
  - [x] SubTask 5.1: 创建 internal/controller/controller.go
  - [x] SubTask 5.2: 实现 Controller 结构体（OrderManager, BotManager, Queue）
  - [x] SubTask 5.3: 实现 NewController 构造函数
  - [x] SubTask 5.4: 实现 CreateOrder 方法（normal/vip）
  - [x] SubTask 5.5: 实现 AddBot 方法
  - [x] SubTask 5.6: 实现 RemoveBot 方法
  - [x] SubTask 5.7: 实现 GetStatus 方法
  - [x] SubTask 5.8: 实现 notifyIdleBots 方法（通知空闲机器人）

## Phase 5: CLI 入口
- [x] Task 6: 实现 CLI 交互
  - [x] SubTask 6.1: 创建 cmd/main.go
  - [x] SubTask 6.2: 实现命令解析循环（bufio.Scanner）
  - [x] SubTask 6.3: 实现 normal 命令处理
  - [x] SubTask 6.4: 实现 vip 命令处理
  - [x] SubTask 6.5: 实现 +bot 命令处理
  - [x] SubTask 6.6: 实现 -bot 命令处理
  - [x] SubTask 6.7: 实现 status 命令处理
  - [x] SubTask 6.8: 实现 quit 命令处理
  - [x] SubTask 6.9: 实现 help 命令和启动帮助提示
  - [x] SubTask 6.10: 实现时间戳格式化（HH:MM:SS）
  - [x] SubTask 6.11: 设计 status 命令输出格式（PENDING/PROCESSING/COMPLETE 状态树 + 机器人状态）
  - [x] SubTask 6.12: CLI 通过 stdout 输出所有日志，通过 shell `>` 重定向到 result.txt（不直接写入文件）

## Phase 6: 测试
- [x] Task 7: 实现单元测试
  - [x] SubTask 7.1: 创建 internal/queue/priority_queue_test.go
  - [x] SubTask 7.2: 测试队列入队出队顺序
  - [x] SubTask 7.3: 测试 VIP 优先级
  - [x] SubTask 7.4: 测试空队列处理
  - [x] SubTask 7.4b: 测试 RequeueAt 回队到原位置
  - [x] SubTask 7.5: 测试并发安全（使用 -race 标志）
  
  - [x] SubTask 7.6: 创建 internal/bot/bot_test.go
  - [x] SubTask 7.7: 测试机器人创建
  - [x] SubTask 7.8: 测试订单处理
  - [x] SubTask 7.8b: 测试 Bot.Run 主循环（拾取订单、空闲等待、停止信号）
  - [x] SubTask 7.9: 测试机器人停止
  - [x] SubTask 7.10: 测试空闲状态
  
  - [x] SubTask 7.11: 创建 internal/controller/controller_test.go
  - [x] SubTask 7.11b: 测试 Controller.CreateOrder 单独逻辑
  - [x] SubTask 7.11c: 测试 Controller.AddBot 单独逻辑
  - [x] SubTask 7.11d: 测试 Controller.RemoveBot 单独逻辑（含回队）
  - [x] SubTask 7.12: 测试完整流程（创建订单→处理→完成）
  - [x] SubTask 7.13: 测试增减机器人
  - [x] SubTask 7.14: 测试 VIP 优先级
  - [x] SubTask 7.15: 测试并发场景

## Phase 7: 脚本更新
- [x] Task 8: 更新构建、测试和运行脚本
  - [x] SubTask 8.1: 更新 scripts/build.sh（go build 命令）
  - [x] SubTask 8.2: 更新 scripts/test.sh（go test 命令）
  - [x] SubTask 8.3: 更新 scripts/run.sh（执行 CLI 并输出到 result.txt）

## Phase 8: 验证和文档
- [x] Task 9: 最终验证
  - [x] SubTask 9.1: 运行所有测试，确保覆盖率 >= 80%
  - [x] SubTask 9.2: 运行构建脚本，确保编译成功
  - [x] SubTask 9.3: 运行 CLI，验证所有功能正常
  - [x] SubTask 9.4: 验证 result.txt 格式正确，包含时间戳
  - [x] SubTask 9.5: 验证 GitHub Actions 工作流通过（本地验证通过，需 PR 后确认）

# Task Dependencies
- Task 2 (订单结构) 无依赖，可并行
- Task 3 (优先级队列) 依赖 Task 2（需使用 Order 结构）
- Task 4 (机器人管理) 依赖 Task 2 和 Task 3（需使用 Order 和 Queue）
- Task 5 (控制器) 依赖 Task 2, Task 3, Task 4（协调所有模块）
- Task 6 (CLI 入口) 依赖 Task 5（需使用 Controller 接口）
- Task 7 (测试) 可与 Task 2-6 并行开发（TDD 方式）
- Task 8 (脚本更新) 依赖 Task 5 和 Task 6（需确认构建和运行命令）
- Task 9 (验证) 依赖所有前置任务

# Task Changes Impact
- Task 2 变更 → 影响 Task 3, Task 4, Task 5
- Task 3 变更 → 影响 Task 4, Task 5
- Task 4 变更 → 影响 Task 5, Task 6
- Task 5 变更 → 影响 Task 6

# Parallelizable Work
- Task 2 和 Task 7 的测试部分可以并行开发（TDD）
- Task 3 和 Task 4 可以在 Task 2 完成后并行开发
