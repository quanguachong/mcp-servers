# http-requests MCP Server

一个基于 Go 的 MCP `stdio` Server，提供 `send_http_request` 工具，用于发送 HTTP 请求。

## 功能

- 支持所有 HTTP Method（`method` 字符串透传）。
- 支持认证方式：
  - `bearer`
  - `api_key`（header/query）
  - `aksk_hmac`（HMAC-SHA256）
- 返回结构化响应：状态码、响应头、响应体、耗时、最终 URL。

## 快速开始

### 1) 安装（go install）

使用远程包路径安装:

```bash
go install github.com/quanguachong/mcp-servers/http/cmd/http-requests@latest
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

## Tool 参数与返回

各 tool 用途说明：

- `send_http_request`：发送 HTTP 请求并返回结构化响应，支持自定义 headers/query/body、超时设置与多种认证方式。

### 1） `send_http_request`

主要参数：

- `url`：请求地址（必填）
- `method`：HTTP 方法（必填）
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
  "url": "https://httpbin.org/get",
  "method": "GET",
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
  "url": "https://httpbin.org/get",
  "method": "GET",
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
  "url": "https://example.com/v1/resource?a=1",
  "method": "POST",
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
