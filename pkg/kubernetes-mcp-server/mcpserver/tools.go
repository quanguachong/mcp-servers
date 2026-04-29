package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type commonFormatInput struct {
	Format string `json:"format,omitempty" jsonschema:"Output format: yaml or json (default yaml)"`
}

type eventsListInput struct {
	Format    string `json:"format,omitempty" jsonschema:"Output format: yaml or json (default yaml)"`
	Namespace string `json:"namespace,omitempty" jsonschema:"Namespace name; empty means all namespaces"`
}

type nodesLogInput struct {
	Format    string `json:"format,omitempty" jsonschema:"Output format: yaml or json (default yaml)"`
	Name      string `json:"name" jsonschema:"Node name"`
	TailLines *int64 `json:"tailLines,omitempty" jsonschema:"Number of lines from the end of the logs to return"`
}

type nodesStatsSummaryInput struct {
	Format string `json:"format,omitempty" jsonschema:"Output format: yaml or json (default yaml)"`
	Name   string `json:"name" jsonschema:"Node name"`
}

type podsListInNamespaceInput struct {
	Format    string `json:"format,omitempty" jsonschema:"Output format: yaml or json (default yaml)"`
	Namespace string `json:"namespace" jsonschema:"Namespace name"`
}

type podsGetInput struct {
	Format    string `json:"format,omitempty" jsonschema:"Output format: yaml or json (default yaml)"`
	Namespace string `json:"namespace" jsonschema:"Namespace name"`
	Name      string `json:"name" jsonschema:"Pod name"`
}

type podsLogInput struct {
	Format    string `json:"format,omitempty" jsonschema:"Output format: yaml or json (default yaml)"`
	Namespace string `json:"namespace" jsonschema:"Namespace name"`
	Name      string `json:"name" jsonschema:"Pod name"`
	Container string `json:"container,omitempty" jsonschema:"Container name"`
	TailLines *int64 `json:"tailLines,omitempty" jsonschema:"Number of lines from the end of the logs to return"`
}

type podsTopInput struct {
	Format    string `json:"format,omitempty" jsonschema:"Output format: yaml or json (default yaml)"`
	Namespace string `json:"namespace" jsonschema:"Namespace name"`
	Name      string `json:"name" jsonschema:"Pod name"`
}

type resourcesListInput struct {
	Format     string `json:"format,omitempty" jsonschema:"Output format: yaml or json (default yaml)"`
	APIVersion string `json:"apiVersion" jsonschema:"Kubernetes apiVersion, for example apps/v1"`
	Kind       string `json:"kind" jsonschema:"Kubernetes kind, for example Deployment"`
	Namespace  string `json:"namespace,omitempty" jsonschema:"Namespace name for namespaced resources"`
}

type resourcesGetInput struct {
	Format     string `json:"format,omitempty" jsonschema:"Output format: yaml or json (default yaml)"`
	APIVersion string `json:"apiVersion" jsonschema:"Kubernetes apiVersion, for example apps/v1"`
	Kind       string `json:"kind" jsonschema:"Kubernetes kind, for example Deployment"`
	Namespace  string `json:"namespace,omitempty" jsonschema:"Namespace name for namespaced resources"`
	Name       string `json:"name" jsonschema:"Resource name"`
}

func RegisterTools(server *mcp.Server, manager Manager) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "events_list",
		Description: "列出 Kubernetes 事件（只读）",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in eventsListInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := HandleEventsList(ctx, manager, in.Namespace, in.Format)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"output": out}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "namespaces_list",
		Description: "列出 Kubernetes 命名空间（只读）",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in commonFormatInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := HandleNamespacesList(ctx, manager, in.Format)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"output": out}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "nodes_log",
		Description: "读取节点日志（只读）",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in nodesLogInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := HandleNodesLog(ctx, manager, in.Name, valueOrZero(in.TailLines), in.Format)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"output": out}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "nodes_stats_summary",
		Description: "读取节点 stats summary（只读）",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in nodesStatsSummaryInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := HandleNodesStatsSummary(ctx, manager, in.Name, in.Format)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"output": out}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "nodes_top",
		Description: "读取节点资源使用情况（只读）",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in commonFormatInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := HandleNodesTop(ctx, manager, in.Format)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"output": out}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "pods_list",
		Description: "列出所有命名空间 Pod（只读）",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in commonFormatInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := HandlePodsList(ctx, manager, in.Format)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"output": out}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "pods_list_in_namespace",
		Description: "列出指定命名空间 Pod（只读）",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in podsListInNamespaceInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := HandlePodsListInNamespace(ctx, manager, in.Namespace, in.Format)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"output": out}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "pods_get",
		Description: "读取单个 Pod 详情（只读）",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in podsGetInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := HandlePodsGet(ctx, manager, in.Namespace, in.Name, in.Format)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"output": out}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "pods_log",
		Description: "读取 Pod 日志（只读）",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in podsLogInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := HandlePodsLog(ctx, manager, in.Namespace, in.Name, in.Container, valueOrZero(in.TailLines), in.Format)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"output": out}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "pods_top",
		Description: "读取 Pod 资源使用情况（只读）",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in podsTopInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := HandlePodsTop(ctx, manager, in.Namespace, in.Name, in.Format)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"output": out}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "resources_list",
		Description: "按 apiVersion/kind 列资源（只读）",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in resourcesListInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := HandleResourcesList(ctx, manager, in.APIVersion, in.Kind, in.Namespace, in.Format)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"output": out}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "resources_get",
		Description: "按 apiVersion/kind/name 取资源详情（只读）",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in resourcesGetInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := HandleResourcesGet(ctx, manager, in.APIVersion, in.Kind, in.Namespace, in.Name, in.Format)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"output": out}, nil
	})
}

func valueOrZero(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}
