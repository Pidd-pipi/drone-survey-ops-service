# BUG_REPRO: 审计事件返回值被覆盖 + 生命周期资源泄漏

## Bug 是什么
`OpsAudit.Add` 的 `defer` 覆盖命名返回值导致永远返回零值事件；`Clear` 用 `[:0]` 保留底层数组不释放内存；`serveHTTP` 的 errCh 无缓冲造成关停后监听协程永久阻塞泄漏；`requestIDMiddleware` 非原子自增在并发下重复。

## 如何触发
- 写入审计事件后读取返回对象；
- 服务反复启动/关停（SIGTERM），goroutine 数只增不减；
- 并发请求查看 `X-Request-ID`。

## 错误信息
```text
Add 返回的事件 ID 为空
Clear 后底层数组未释放，cap=69
关停后 goroutine 数增长: before=7 after=11
请求 ID 重复
```
