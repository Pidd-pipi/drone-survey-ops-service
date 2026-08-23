# BUG_REPRO: 规则分组/标签共享切片累积

## Bug 是什么
`opsRules03/04` 复用包级切片直接 append 不清零，重复装载后条数越攒越多；`opsRule0302/0406` 的标签从共享池 append 累积，重复装载后标签翻倍。

## 如何触发
- 重复调用 `opsRules03()/opsRules04()`；
- 重复调用 `opsRule0302()/opsRule0406()`。

## 错误信息
```text
分组 03 重复装载后长度增长: 16
OPS-0302 重复装载后标签越攒越多: 7: [site operator evidence site operator evidence reviewed]
```
