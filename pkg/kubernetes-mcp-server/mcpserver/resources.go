package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"
)

const (
	// OutputFormatYAML 表示输出 YAML 文本。
	OutputFormatYAML = "yaml"
	// OutputFormatJSON 表示输出 JSON 文本。
	OutputFormatJSON = "json"
)

// Manager 定义 handler 依赖的只读 Kubernetes 查询能力。
type Manager interface {
	ListEvents(ctx context.Context, namespace string) (any, error)
	ListNamespaces(ctx context.Context) (any, error)
	NodesLog(ctx context.Context, nodeName string, tailLines int64) (any, error)
	NodesStatsSummary(ctx context.Context, nodeName string) (any, error)
	NodesTop(ctx context.Context) (any, error)
	PodsList(ctx context.Context) (any, error)
	PodsListInNamespace(ctx context.Context, namespace string) (any, error)
	PodsGet(ctx context.Context, namespace, name string) (any, error)
	PodsLog(ctx context.Context, namespace, name, container string, tailLines int64) (any, error)
	PodsTop(ctx context.Context, namespace, name string) (any, error)
	ResourcesList(ctx context.Context, apiVersion, kind, namespace string) (any, error)
	ResourcesGet(ctx context.Context, apiVersion, kind, namespace, name string) (any, error)
}

// HandleResourcesList 调用 manager 列表查询并返回 YAML/JSON 文本。
func HandleResourcesList(ctx context.Context, manager Manager, apiVersion, kind, namespace, format string) (string, error) {
	data, err := manager.ResourcesList(ctx, apiVersion, kind, namespace)
	if err != nil {
		return "", err
	}
	return renderOutput(data, format)
}

// HandleResourcesGet 调用 manager 详情查询并返回 YAML/JSON 文本。
func HandleResourcesGet(ctx context.Context, manager Manager, apiVersion, kind, namespace, name, format string) (string, error) {
	data, err := manager.ResourcesGet(ctx, apiVersion, kind, namespace, name)
	if err != nil {
		return "", err
	}
	return renderOutput(data, format)
}

func renderOutput(data any, format string) (string, error) {
	switch normalizeFormat(format) {
	case OutputFormatYAML:
		b, err := yaml.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("marshal yaml failed: %w", err)
		}
		return string(b), nil
	case OutputFormatJSON:
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "", fmt.Errorf("marshal json failed: %w", err)
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
}

func normalizeFormat(format string) string {
	if format == "" {
		return OutputFormatYAML
	}
	return strings.ToLower(format)
}
