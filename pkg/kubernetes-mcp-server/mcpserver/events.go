package mcpserver

import "context"

// HandleEventsList 调用 manager 查询事件并返回 YAML/JSON 文本。
func HandleEventsList(ctx context.Context, manager Manager, namespace, format string) (string, error) {
	data, err := manager.ListEvents(ctx, namespace)
	if err != nil {
		return "", err
	}
	return renderOutput(data, format)
}
