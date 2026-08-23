# Drone Survey Ops Backend

Go 模块：`drone-survey-ops-service`。

```bash
go test ./...
go build ./...
go run .
```

入口为 `main.go`，默认监听 `:8080`，可用 `PORT` 调整。健康检查为 `GET /health`；API 为 `GET /api/missions` 和 `POST /api/missions/status`。
