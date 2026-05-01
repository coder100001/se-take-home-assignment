# 代码优化修复 Spec

## Why
项目代码评审发现多处潜在死锁风险、锁层级混乱和代码清晰度问题，需要修复以确保系统安全性和可维护性。

## What Changes
- 修复 `NotifyAll()` 锁层级混乱：先复制 bots 列表再遍历通知
- 修复 `RemoveBot()` 锁内阻塞：先释放 `bm.mu` 再调用 `b.Stop()`
- 优化 `Run()` 循环逻辑：用 `if-else` 替代 `if + continue`
- 将 `insertAt` 包级函数内联到 `RequeueAt` 方法中
- 将 PENDING 排序逻辑移入 `Controller.GetStatus()`
- 优化 `Order.String()` 用数组索引替代 switch

## Impact
- Affected specs: implement-order-controller
- Affected code: internal/bot/bot.go, internal/queue/priority_queue.go, internal/controller/controller.go, internal/order/order.go, cmd/main.go

## ADDED Requirements

### Requirement: 锁层级安全
系统应确保锁获取顺序一致，避免潜在死锁。

#### Scenario: NotifyAll 不持有 BotManager 锁
- **WHEN** 调用 `NotifyAll()` 通知空闲机器人
- **THEN** 不在持有 `BotManager.mu` 的同时获取 `Bot.mu`

#### Scenario: RemoveBot 不在锁内阻塞等待
- **WHEN** 调用 `RemoveBot()` 移除机器人
- **THEN** 先从列表移除并释放 `BotManager.mu`，再调用 `b.Stop()` 等待 goroutine 退出

## MODIFIED Requirements

### Requirement: Bot.Run() 循环逻辑
`Run()` 方法应使用 `if-else` 结构替代 `if + continue`，减少认知负担。

### Requirement: Controller.GetStatus() 返回已排序数据
`GetStatus()` 返回的 pending 列表应按 VIP 优先排序，CLI 层不再做排序。

### Requirement: Order.String() 使用数组索引
`String()` 方法应使用 `[...]string` 数组索引替代 switch 语句。

## REMOVED Requirements

### Requirement: insertAt 包级函数
**Reason**: 该函数仅被 `RequeueAt` 使用，应内联消除不必要的抽象
**Migration**: 将 insertAt 逻辑直接写入 `RequeueAt` 方法
