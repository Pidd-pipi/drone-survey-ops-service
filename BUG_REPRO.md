# BUG_REPRO: 运维存储读路径丢锁 + Labels 引用逃逸

## Bug 是什么
`OpsStore.Get/List` 读路径丢失 RLock，且返回记录不再 Clone（Labels map 与内部存储共享）；`Put` 直接保存调用方 Labels 引用。并发读写发生 data race，调用方改返回记录的标签会反向污染库里数据。

## 如何触发
- 并发调用 `List/Get` 与 `Put/Update`（竞态检测必报 DATA RACE）；
- 读取记录后修改其 Labels，再读同一记录可见篡改。

## 错误信息
```text
WARNING: DATA RACE
Write at 0x00c000122f00 by goroutine 14:
  runtime.mapaccess2_faststr()
  drone-survey-ops-service.(*OpsStore).Put()
      .../backend/ops_store.go:58 +0x204
Previous read at 0x00c000122f00 by goroutine 9:
  drone-survey-ops-service.(*OpsStore).List()
      .../backend/ops_store.go:40 +0x78
```
