# Tasks

- [ ] Task 1: 修复 NotifyAll 锁层级混乱
  - [ ] SubTask 1.1: 修改 NotifyAll() 先复制 bots 列表，释放 bm.mu 后再遍历通知

- [ ] Task 2: 修复 RemoveBot 锁内阻塞
  - [ ] SubTask 2.1: 修改 RemoveBot() 先从列表移除 bot 并释放锁，再调用 b.Stop()

- [ ] Task 3: 优化 Run() 循环逻辑
  - [ ] SubTask 3.1: 用 if-else 替代 if + continue 结构

- [ ] Task 4: 内联 insertAt 到 RequeueAt
  - [ ] SubTask 4.1: 删除 insertAt 包级函数，将逻辑内联到 RequeueAt 方法

- [ ] Task 5: 将 PENDING 排序移入 Controller.GetStatus()
  - [ ] SubTask 5.1: 在 GetStatus() 中对 pending 列表按 VIP 优先排序
  - [ ] SubTask 5.2: 移除 cmd/main.go 中的 sort.Slice 排序逻辑

- [ ] Task 6: 优化 Order.String() 用数组索引替代 switch
  - [ ] SubTask 6.1: 定义 typeStrs 和 statusStrs 数组，用索引替代 switch

- [ ] Task 7: 运行验证
  - [ ] SubTask 7.1: go test ./... -race -count=1 全部通过
  - [ ] SubTask 7.2: go vet ./... 无问题
  - [ ] SubTask 7.3: go build ./... 编译通过

# Task Dependencies
- Task 1 和 Task 2 都修改 bot.go，需串行执行
- Task 3 修改 bot.go，依赖 Task 1 和 Task 2
- Task 4 修改 queue.go，独立于 Task 1-3
- Task 5 修改 controller.go 和 main.go，独立于 Task 1-4
- Task 6 修改 order.go，独立于其他任务
- Task 7 依赖所有前置任务

# Parallelizable Work
- Task 4, Task 5, Task 6 可并行执行
