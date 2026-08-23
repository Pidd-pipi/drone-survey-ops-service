# BUG_REPRO: 规则标签池共享导致重复装载累积

## Bug 是什么
`opsRule0502/0602/0702/0802` 从各自包级标签池 `append` 累积标签，重复装载后 `RequiredLabels` 越攒越多（首次 4 个，再次 7/8 个），共享底层数组跨调用串场。

## 如何触发
- 重复调用 `opsRule0502()`（或 0602/0702/0802）。

## 错误信息
```text
OPS-0502 重复装载后标签越攒越多: 7: [site operator evidence site operator evidence reviewed]
```
