# http-requests MCP Server

一个基于 Go 的 MCP `stdio` Server，提供按 HTTP Method 拆分的请求工具。

## 功能

- 默认提供 `GET/POST/PUT` 三个工具：
  - `get_http_request`
  - `post_http_request`
  - `put_http_request`
- 可选提供扩展工具（编译时开启）：
  - `delete_http_request`
  - `patch_http_request`
- 支持认证方式：
  - `bearer`
  - `api_key`（header/query）
  - `aksk_hmac`（HMAC-SHA256）
- 返回结构化响应：状态码、响应头、响应体、耗时、最终 URL。

## 快速开始

### 1) 安装（go install）

使用远程包路径安装：

```bash
go install github.com/quanguachong/mcp-servers/cmd/http-requests@latest
```

安装后命令通常位于 `$HOME/go/bin/http-requests`（或 `$GOBIN/http-requests`）。

### 2) 配置 Cursor（`~/.cursor/mcp.json`）

在 `~/.cursor/mcp.json` 的 `mcpServers` 中增加：

```json
{
  "mcpServers": {
    "http-requests": {
      "command": "http-requests"
    }
  }
}
```

如果 PATH 未包含 Go bin 目录，请将 `command` 改为绝对路径，例如 `/Users/<you>/go/bin/http-requests`。

## 构建行为

默认构建不会包含 `DELETE/PATCH` tool：

```bash
go build ./cmd/http-requests
```

如需启用 `DELETE/PATCH` tool，请使用 build tag：

```bash
go build -tags http_requests_extended_methods ./cmd/http-requests
```

## Tool 参数与返回

各 tool 用途说明：

- `get_http_request`
- `post_http_request`
- `put_http_request`
- `delete_http_request`（仅扩展构建）
- `patch_http_request`（仅扩展构建）

除 tool 名称和固定 Method 外，参数与返回结构一致。

### 参数

主要参数：

- `url`：请求地址（必填）
- `headers`：请求头 map
- `query`：查询参数 map
- `body`：字符串或 JSON 对象
- `timeout_ms`：超时（毫秒）
- `auth`：认证配置

返回字段：

- `status_code`
- `headers`
- `body`
- `body_base64`
- `latency_ms`
- `final_url`

## 认证示例

### Bearer

```json
{
  "tool": "get_http_request",
  "url": "https://httpbin.org/get",
  "auth": {
    "type": "bearer",
    "bearer": {
      "token": "YOUR_TOKEN"
    }
  }
}
```

### API Key (Header)

```json
{
  "tool": "get_http_request",
  "url": "https://httpbin.org/get",
  "auth": {
    "type": "api_key",
    "api_key": {
      "key": "X-API-Key",
      "value": "YOUR_API_KEY",
      "in": "header"
    }
  }
}
```

### AK/SK HMAC

```json
{
  "tool": "post_http_request",
  "url": "https://example.com/v1/resource?a=1",
  "body": {
    "name": "demo"
  },
  "auth": {
    "type": "aksk_hmac",
    "aksk_hmac": {
      "access_key": "YOUR_AK",
      "secret_key": "YOUR_SK"
    }
  }
}
```
