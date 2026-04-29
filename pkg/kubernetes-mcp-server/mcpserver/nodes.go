package mcpserver

import "context"

// HandleNodesLog 调用 manager 查询节点日志并返回 YAML/JSON 文本。
func HandleNodesLog(ctx context.Context, manager Manager, nodeName string, tailLines int64, format string) (string, error) {
	data, err := manager.NodesLog(ctx, nodeName, tailLines)
	if err != nil {
		return "", err
	}
	return renderOutput(data, format)
}

// HandleNodesStatsSummary 调用 manager 查询节点统计摘要并返回 YAML/JSON 文本。
func HandleNodesStatsSummary(ctx context.Context, manager Manager, nodeName, format string) (string, error) {
	data, err := manager.NodesStatsSummary(ctx, nodeName)
	if err != nil {
		return "", err
	}
	return renderOutput(data, format)
}

// HandleNodesTop 调用 manager 查询节点资源使用并返回 YAML/JSON 文本。
func HandleNodesTop(ctx context.Context, manager Manager, format string) (string, error) {
	data, err := manager.NodesTop(ctx)
	if err != nil {
		return "", err
	}
	return renderOutput(data, format)
}
