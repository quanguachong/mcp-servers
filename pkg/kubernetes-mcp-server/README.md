# kubernetes-mcp-server

`kubernetes-mcp-server` 提供 Kubernetes 集群只读查询能力，默认仅包含 12 个只读 tools，不包含任何写入/破坏性操作。

## Tools（12 个只读）

- `events_list`：列出事件
- `namespaces_list`：列出命名空间
- `nodes_log`：读取节点日志
- `nodes_stats_summary`：读取节点统计摘要
- `nodes_top`：读取节点资源使用情况
- `pods_list`：列出全部命名空间 Pod
- `pods_list_in_namespace`：列出指定命名空间 Pod
- `pods_get`：读取单个 Pod 详情
- `pods_log`：读取 Pod 日志
- `pods_top`：读取 Pod 资源使用情况
- `resources_list`：按 `apiVersion/kind` 列资源
- `resources_get`：按 `apiVersion/kind/name` 取资源详情

## 显式排除项

以下写入或破坏性操作不提供：

- `pods_exec`
- `pods_delete`
- `pods_run`
- `resources_delete`
- `resources_scale`
- `resources_create_or_update`

## kubeconfig 加载顺序

按以下顺序加载，命中即停止继续回退：

1. `--kubeconfig` 命令行参数
2. 环境变量 `KUBECONFIG`
3. 默认路径 `~/.kube/config`
4. 集群内配置 `InClusterConfig()`

## Cursor 配置示例

在 `~/.cursor/mcp.json` 中增加：

```json
{
  "mcpServers": {
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

如未传 `--kubeconfig`，服务会按上面的加载顺序自动解析 kubeconfig。
