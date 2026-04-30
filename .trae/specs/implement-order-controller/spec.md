# McDonald's 订单控制器 Spec

## Why
McDonald's 在 COVID-19 期间需要构建自动化烹饪机器人系统来减少人力成本并提高效率。作为软件工程师，需要创建一个订单控制器来处理订单控制流程，支持 VIP 优先级、动态机器人管理和并发处理。

## What Changes
- 新增交互式 CLI 订单控制系统
- 新增订单管理模块（支持普通订单和 VIP 订单）
- 新增机器人管理模块（支持动态增减机器人）
- 新增优先级队列（VIP > Normal）
- 新增并发处理机制（多个机器人同时处理订单）
- 修改构建脚本、测试脚本和运行脚本
- 新增单元测试和集成测试

## Impact
- Affected specs: 无（新系统）
- Affected code: 
  - 新增 `cmd/main.go` - CLI 入口
  - 新增 `internal/order/order.go` - 订单管理
  - 新增 `internal/queue/priority_queue.go` - 优先级队列
  - 新增 `internal/bot/bot.go` - 机器人管理
  - 新增 `internal/controller/controller.go` - 控制器
  - 新增测试文件
  - 修改 `scripts/build.sh`, `scripts/test.sh`, `scripts/run.sh`

## ADDED Requirements

### Requirement: 订单创建功能
系统应提供订单创建功能，支持普通订单和 VIP 订单。

#### Scenario: 创建普通订单
- **WHEN** 用户输入 `normal <description>` 命令
- **THEN** 系统创建一个普通订单，订单号唯一且递增，状态为 PENDING，显示在 PENDING 区域

#### Scenario: 创建 VIP 订单
- **WHEN** 用户输入 `vip <description>` 命令
- **THEN** 系统创建一个 VIP 订单，订单号唯一且递增，状态为 PENDING，排在所有普通订单之前，但在现有 VIP 订单之后

#### Scenario: 订单号唯一性
- **WHEN** 创建多个订单（普通或 VIP）
- **THEN** 每个订单的编号唯一且递增

### Requirement: 机器人管理功能
系统应提供机器人管理功能，支持动态增减机器人。

#### Scenario: 增加机器人
- **WHEN** 用户输入 `+bot` 命令
- **THEN** 系统创建一个新的机器人，机器人立即尝试从 PENDING 区域获取订单并处理

#### Scenario: 减少机器人
- **WHEN** 用户输入 `-bot` 命令
- **THEN** 系统销毁最新的机器人，如果机器人正在处理订单，订单返回到原位置（保持 VIP/Normal 优先级）

#### Scenario: 机器人处理订单
- **WHEN** 机器人获取到一个订单
- **THEN** 机器人处理该订单 10 秒，订单状态变为 PROCESSING，处理完成后订单状态变为 COMPLETE 并移至 COMPLETE 区域

#### Scenario: 机器人空闲状态
- **WHEN** PENDING 区域没有订单
- **THEN** 机器人保持 IDLE 状态，等待新订单到来

### Requirement: 优先级队列
系统应实现优先级队列，确保 VIP 订单优先处理。

#### Scenario: VIP 优先级
- **WHEN** PENDING 区域同时存在 VIP 订单和普通订单
- **THEN** 机器人优先处理 VIP 订单，VIP 队列为空时才处理普通订单

#### Scenario: VIP 订单排序
- **WHEN** 创建多个 VIP 订单
- **THEN** VIP 订单按照 FIFO（先进先出）顺序处理

#### Scenario: 普通订单排序
- **WHEN** 创建多个普通订单
- **THEN** 普通订单按照 FIFO（先进先出）顺序处理

### Requirement: 并发安全
系统应确保并发安全，防止订单重复处理或丢失。

#### Scenario: 多机器人并发处理
- **WHEN** 多个机器人同时从队列获取订单
- **THEN** 每个订单只被一个机器人处理一次，无竞态条件

#### Scenario: 机器人销毁时订单回退
- **WHEN** 销毁正在处理订单的机器人
- **THEN** 订单安全地返回到队列原位置，保持优先级

#### Scenario: 订单回退保持精确位置
- **GIVEN** PENDING 区域队列顺序为 [VIP#1, Normal#1, Normal#2]
- **WHEN** Normal#1 被机器人拾取并开始处理
- **AND** 该机器人被销毁
- **THEN** Normal#1 回到 Normal 队列的原始位置（即 Normal#2 之前），队列变为 [VIP#1, Normal#1, Normal#2]

#### Scenario: 快速连续增减机器人
- **WHEN** 用户快速连续输入 `+bot` 和 `-bot`
- **THEN** 系统正确处理每个操作，无竞态条件或状态不一致

#### Scenario: 机器人数量为 0 时批量创建订单
- **WHEN** 机器人数量为 0 时创建大量订单
- **AND** 随后添加一个机器人
- **THEN** 新机器人在创建后立即开始处理 PENDING 中的最高优先级订单

### Requirement: CLI 交互
系统应提供交互式 CLI 界面。

#### Scenario: 命令解析
- **WHEN** 用户输入命令
- **THEN** 系统正确解析并执行相应操作

#### Scenario: 无效命令处理
- **WHEN** 用户输入无效命令
- **THEN** 系统提示"未知命令"，显示可用命令列表帮助

#### Scenario: 启动帮助提示
- **WHEN** 程序启动
- **THEN** 系统显示可用命令帮助信息，方便新用户快速上手

#### Scenario: 状态查看
- **WHEN** 用户输入 `status` 命令
- **THEN** 系统显示所有订单状态（PENDING, PROCESSING, COMPLETE）和机器人状态（IDLE, PROCESSING）

#### Scenario: 退出程序
- **WHEN** 用户输入 `quit` 命令
- **THEN** 系统优雅退出

### Requirement: 输出格式
系统应按照要求格式输出结果。

#### Scenario: 时间戳格式
- **WHEN** 系统输出任何信息
- **THEN** 输出包含 HH:MM:SS 格式的时间戳

#### Scenario: result.txt 文件
- **WHEN** 程序运行完成
- **THEN** `scripts/result.txt` 文件包含有意义的输出，且包含时间戳

### Requirement: 测试覆盖
系统应具备充分的测试覆盖。

#### Scenario: 单元测试
- **WHEN** 运行 `go test ./... -v -race`
- **THEN** 所有测试通过，覆盖率 >= 80%，无竞态条件

#### Scenario: 集成测试
- **WHEN** 运行 `scripts/test.sh`
- **THEN** 所有测试通过

#### Scenario: 构建验证
- **WHEN** 运行 `scripts/build.sh`
- **THEN** 成功编译生成可执行文件

## MODIFIED Requirements
无（新系统）

## REMOVED Requirements
无（新系统）
