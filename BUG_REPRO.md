# BUG_REPRO: context 取消/超时未向下游传播

## Bug 是什么
`opsContext` 改用 `context.Background()` 丢掉父 ctx，`opsDelay` 无视 `ctx.Done`，`Create/Search/Get` 用 Background 调存储。调用方取消或超时后，后台操作仍继续执行、延迟等待不停。

## 如何触发
- 传入已取消/已到期的 ctx 调用 `Transition/Create/Search/Get`；
- 取消后等待 `opsDelay`，不会立即返回。

## 错误信息
定向测试断言取消后应立即返回 `ctx.Err()`，坏环境返回 nil/照常成功：
```text
已取消的 ctx 下 opsDelay 应返回错误，实际 nil（耗时 301ms）
父 ctx 已取消时 Transition 应返回错误
```
