package mcpserver

import "context"

// HandleNamespacesList 调用 manager 查询命名空间并返回 YAML/JSON 文本。
func HandleNamespacesList(ctx context.Context, manager Manager, format string) (string, error) {
	data, err := manager.ListNamespaces(ctx)
	if err != nil {
		return "", err
	}
	return renderOutput(data, format)
}
