# MCP Servers Monorepo

这个仓库存储了多个 MCP Servers：

- `mongodb`：MongoDB 查询能力（`count`、`find`、`list-databases`、`list-collections`）
- `http`：HTTP 请求能力（`send_http_request`，支持多种认证方式）

## 推荐目录结构

当前结构可直接使用：

```text
mcp-server/
├── README.md
├── mongodb/
│   ├── README.md
│   ├── go.mod
│   ├── main.go
│   ├── go.sum
│   └── pkg/
└── http/
    ├── README.md
    ├── go.mod
    ├── go.sum
    ├── cmd/http-requests/main.go
    └── internal/
```

> 建议保持两个 Server 各自独立的 `go.mod`，避免依赖耦合。

## 快速开始

### 1) 安装两个 MCP Server（go install）

使用远程包路径安装:

```bash
go install github.com/quanguachong/mcp-servers/mongodb@latest
go install github.com/quanguachong/mcp-servers/http/cmd/http-requests@latest
```

安装后可执行文件会放到 `$GOBIN`（未设置时通常为 `$HOME/go/bin`）：

- `mongodb-mcp-server`
- `http-requests`

### 2) 配置 MCP Config

#### Cursor（`~/.cursor/mcp.json`）

编辑 `~/.cursor/mcp.json`，在 `mcpServers` 中增加（或合并）以下配置：

```json
{
  "mcpServers": {
    "mongodb-mcp-server": {
      "command": "mongodb-mcp-server",
      "env": {
        "MONGODB_URI": "mongodb://<your-uri>"
      }
    },
    "http-requests": {
      "command": "http-requests"
    }
  }
}
```

如果 Cursor 找不到命令，可改为绝对路径，例如 `/Users/<you>/go/bin/http-requests`。

更多参数和 tool 说明见：[mongodb/README.md](./mongodb/README.md) 与 [http/README.md](./mongodb/README.md)。
