# MCP Servers Monorepo

这个仓库存储了多个 MCP Servers：

- `mongodb`：MongoDB 查询能力（`count`、`find`、`list-databases`、`list-collections`）
- `http`：HTTP 请求能力（`get/post/put_http_request` 默认可用，支持多种认证方式；`delete/patch_http_request` 需扩展构建）
- `kubernetes`：Kubernetes 只读查询能力（12 个只读 tools；不包含 `pods_exec/pods_delete/pods_run/resources_delete/resources_scale/resources_create_or_update`）

## 推荐目录结构

当前结构可直接使用：

```text
mcp-server/
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   ├── http-requests/
│   │   └── main.go
│   ├── kubernetes-mcp-server/
│   │   └── main.go
│   └── mongodb-mcp-server/
│       └── main.go
└── pkg/
    ├── http-requests/
    │   ├── README.md
    │   ├── auth/
    │   ├── httpclient/
    │   ├── mcpserver/
    │   └── types/
    ├── kubernetes-mcp-server/
    │   ├── README.md
    │   └── *.go
    └── mongodb-mcp-server/
        ├── README.md
        └── *.go
```

> 当前仓库使用根目录单一 `go.mod` 管理依赖，`cmd/` 放启动入口，`pkg/` 放两个 MCP 的业务逻辑。

## 快速开始

### 1) 安装三个 MCP Server（go install）

使用远程包路径安装:

```bash
go install github.com/quanguachong/mcp-servers/cmd/mongodb-mcp-server@latest
go install github.com/quanguachong/mcp-servers/cmd/http-requests@latest
go install github.com/quanguachong/mcp-servers/cmd/kubernetes-mcp-server@latest
```

安装后可执行文件会放到 `$GOBIN`（未设置时通常为 `$HOME/go/bin`）：

- `mongodb-mcp-server`
- `http-requests`
- `kubernetes-mcp-server`

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
    },
    "kubernetes-mcp-server": {
      "command": "kubernetes-mcp-server",
      "args": [
        "--kubeconfig",
        "/Users/<you>/.kube/config"
      ]
    }
  }
}
```

如果 Cursor 找不到命令，可改为绝对路径，例如 `/Users/<you>/go/bin/http-requests`。

更多参数和 tool 说明见：

- [pkg/mongodb-mcp-server/README.md](./pkg/mongodb-mcp-server/README.md)
- [pkg/http-requests/README.md](./pkg/http-requests/README.md)
- [pkg/kubernetes-mcp-server/README.md](./pkg/kubernetes-mcp-server/README.md)
