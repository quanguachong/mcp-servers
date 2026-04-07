# mongodb-mcp-server

基于官方 Go MCP SDK 的 MongoDB MCP Server，支持旧版本的 MongoDB（例如: 4.0.12）。

## 支持的 tools

- `count`
- `find`
- `list-databases`
- `list-collections`

## 依赖

- Go 1.21+
- 可访问的 MongoDB 实例（已在 4.0.12 目标场景考虑兼容）

## 快速开始

### 1) 安装（go install）

使用远程包路径安装：

```bash
go install github.com/quanguachong/mcp-servers/cmd/mongodb-mcp-server@latest
```

安装后命令通常位于 `$HOME/go/bin/mongodb-mcp-server`（或 `$GOBIN/mongodb-mcp-server`）。

### 2) 配置 Cursor（`~/.cursor/mcp.json`）

在 `~/.cursor/mcp.json` 的 `mcpServers` 中增加：

```json
{
  "mcpServers": {
    "mongodb-mcp-server": {
      "command": "mongodb-mcp-server",
      "env": {
        "MONGODB_URI": "mongodb://<your-uri>"
      }
    }
  }
}
```

如果 PATH 未包含 Go bin 目录，请将 `command` 改为绝对路径，例如 `/Users/<you>/go/bin/mongodb-mcp-server`。

MCP Server 使用 stdio transport 运行（由 MCP Host 拉起该进程并通过 stdin/stdout 通信）。

## 环境变量

| 变量名        | 说明                     |
| ------------- | ------------------------ |
| `MONGODB_URI` | 必填。MongoDB 连接 URI。 |

## Tool 参数与返回

各 tool 用途说明：

- `count`：按条件统计文档数量，适合先评估数据规模。
- `find`：按过滤条件查询文档，支持投影、排序、条数限制和响应大小限制。
- `list-databases`：列出当前连接可见的数据库及基础统计信息。
- `list-collections`：列出指定数据库下的集合名称及数量。

### 1) `count`

参数：

```json
{
  "database": "test",
  "collection": "users",
  "query": { "status": 1 }
}
```

返回：

```json
{
  "count": 123
}
```

### 2) `find`

参数（示例）：

```json
{
  "database": "test",
  "collection": "users",
  "filter": { "status": 1 },
  "projection": { "name": 1, "_id": 0 },
  "sort": { "createdAt": -1 },
  "limit": 10,
  "responseBytesLimit": 1048576
}
```

- `limit` 默认值为 `10`。
- `responseBytesLimit` 默认值为 `1048576`（1MB），超限会返回错误。

返回：

```json
{
  "documents": [{ "name": "alice" }]
}
```

### 3) `list-databases`

参数：

```json
{}
```

返回：

```json
{
  "databases": [{ "name": "admin", "size": 12345 }],
  "totalCount": 1
}
```

### 4) `list-collections`

参数：

```json
{
  "database": "test"
}
```

返回：

```json
{
  "collections": [{ "name": "users" }],
  "totalCount": 1
}
```
