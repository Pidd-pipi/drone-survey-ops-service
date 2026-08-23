# BUG_REPRO: 响应契约跨层不一致（中间件/健康检查）

## Bug 是什么
`opsEnterpriseMiddleware` 丢失 `X-Operations-Domain` 等响应头，`opsJSON` 不设 JSON Content-Type，`opsAllowed` 只认第一个允许方法，健康检查返回错误服务名 `ops`。

## 如何触发
- 真实 HTTP 请求任意接口检查响应头；
- 请求 `/health` 查看响应体。

## 错误信息
```text
缺少 X-Operations-Domain 头: ""
opsJSON 应设置 JSON Content-Type，实际 ""
POST 应被允许
健康检查服务名错误: "ops"
```
