# Drone Survey Ops Service

标准库无人机测绘外业任务服务。提供 `GET /health`、`GET /api/missions` 和 `POST /api/missions/status`；POST 示例为 `{"id":"mission-241","status":"landed"}`，状态支持 `planned`、`flying`、`landed`、`aborted`。根路径是内置任务页面，默认监听 8080。启动命令为 `cd backend && go run .`。

## Directory Tree and API

```text
backend/       Go 模块、HTTP 服务和内嵌 web 静态资源
database/      数据库说明
output/        验证记录
README.md      项目说明
prompt.txt     任务提示
runtime_smoke.json  运行冒烟配置
```

健康检查：`GET /health`。API：`GET /api/missions`、`POST /api/missions/status`；页面入口：`GET /`。

## Verification

验证日期：2026-08-21。`gofmt -w *.go`、`go build ./...`、`go test ./...` 均成功。启动后真实验证 health、missions、合法 POST 均返回 200；非法状态返回 400，未知任务返回 404；首页和 `/app.js` 返回 200。验证完成后服务已关闭。

## Engineering Notes

无人机测绘流程代码按领域模型、校验、状态转换、并发安全存储、审计事件和 HTTP 生命周期分层。请求会保留请求标识并经过恢复与超时保护；状态写入使用版本校验，错误通过可识别的领域错误返回。

除现有接口回归测试外，项目还保留可复用的分页、过滤、策略、工作流和运行健康能力，便于后续扩展而不把业务规则堆积到处理器中。
