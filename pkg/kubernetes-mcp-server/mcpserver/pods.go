package mcpserver

import "context"

// HandlePodsList 调用 manager 查询所有命名空间 Pod 列表并返回 YAML/JSON 文本。
func HandlePodsList(ctx context.Context, manager Manager, format string) (string, error) {
	data, err := manager.PodsList(ctx)
	if err != nil {
		return "", err
	}
	return renderOutput(data, format)
}

// HandlePodsListInNamespace 调用 manager 查询指定命名空间 Pod 列表并返回 YAML/JSON 文本。
func HandlePodsListInNamespace(ctx context.Context, manager Manager, namespace, format string) (string, error) {
	data, err := manager.PodsListInNamespace(ctx, namespace)
	if err != nil {
		return "", err
	}
	return renderOutput(data, format)
}

// HandlePodsGet 调用 manager 查询 Pod 详情并返回 YAML/JSON 文本。
func HandlePodsGet(ctx context.Context, manager Manager, namespace, name, format string) (string, error) {
	data, err := manager.PodsGet(ctx, namespace, name)
	if err != nil {
		return "", err
	}
	return renderOutput(data, format)
}

// HandlePodsLog 调用 manager 查询 Pod 日志并返回 YAML/JSON 文本。
func HandlePodsLog(ctx context.Context, manager Manager, namespace, name, container string, tailLines int64, format string) (string, error) {
	data, err := manager.PodsLog(ctx, namespace, name, container, tailLines)
	if err != nil {
		return "", err
	}
	return renderOutput(data, format)
}

// HandlePodsTop 调用 manager 查询 Pod 资源使用并返回 YAML/JSON 文本。
func HandlePodsTop(ctx context.Context, manager Manager, namespace, name, format string) (string, error) {
	data, err := manager.PodsTop(ctx, namespace, name)
	if err != nil {
		return "", err
	}
	return renderOutput(data, format)
}
