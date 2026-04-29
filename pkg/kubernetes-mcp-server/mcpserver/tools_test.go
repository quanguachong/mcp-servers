package mcpserver

import (
	"context"
	"slices"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type fakeManager struct {
	podsLogCalls []struct {
		namespace string
		name      string
		container string
		tailLines int64
	}
	resourcesGetCalls []struct {
		apiVersion string
		kind       string
		namespace  string
		name       string
	}
}

func (m *fakeManager) ListEvents(context.Context, string) (any, error) { return map[string]any{}, nil }
func (m *fakeManager) ListNamespaces(context.Context) (any, error)      { return map[string]any{}, nil }
func (m *fakeManager) NodesLog(context.Context, string, int64) (any, error) {
	return map[string]any{}, nil
}
func (m *fakeManager) NodesStatsSummary(context.Context, string) (any, error) {
	return map[string]any{}, nil
}
func (m *fakeManager) NodesTop(context.Context) (any, error) { return map[string]any{}, nil }
func (m *fakeManager) PodsList(context.Context) (any, error) { return map[string]any{}, nil }
func (m *fakeManager) PodsListInNamespace(context.Context, string) (any, error) {
	return map[string]any{}, nil
}
func (m *fakeManager) PodsGet(context.Context, string, string) (any, error) {
	return map[string]any{}, nil
}
func (m *fakeManager) PodsLog(_ context.Context, namespace, name, container string, tailLines int64) (any, error) {
	m.podsLogCalls = append(m.podsLogCalls, struct {
		namespace string
		name      string
		container string
		tailLines int64
	}{
		namespace: namespace,
		name:      name,
		container: container,
		tailLines: tailLines,
	})
	return map[string]any{"ok": true}, nil
}
func (m *fakeManager) PodsTop(context.Context, string, string) (any, error) {
	return map[string]any{}, nil
}
func (m *fakeManager) ResourcesList(context.Context, string, string, string) (any, error) {
	return map[string]any{}, nil
}
func (m *fakeManager) ResourcesGet(_ context.Context, apiVersion, kind, namespace, name string) (any, error) {
	m.resourcesGetCalls = append(m.resourcesGetCalls, struct {
		apiVersion string
		kind       string
		namespace  string
		name       string
	}{
		apiVersion: apiVersion,
		kind:       kind,
		namespace:  namespace,
		name:       name,
	})
	return map[string]any{"ok": true}, nil
}

func TestRegisterTools_BasicMappingAndParamPassthrough(t *testing.T) {
	ctx := context.Background()
	manager := &fakeManager{}
	server := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "v0.0.1"}, nil)
	RegisterTools(server, manager)

	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("连接服务端失败: %v", err)
	}
	defer serverSession.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0.0.1"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("连接客户端失败: %v", err)
	}
	defer clientSession.Close()

	tools, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("列出工具失败: %v", err)
	}

	gotNames := make([]string, 0, len(tools.Tools))
	for _, tool := range tools.Tools {
		gotNames = append(gotNames, tool.Name)
	}
	wantNames := []string{
		"events_list",
		"namespaces_list",
		"nodes_log",
		"nodes_stats_summary",
		"nodes_top",
		"pods_list",
		"pods_list_in_namespace",
		"pods_get",
		"pods_log",
		"pods_top",
		"resources_list",
		"resources_get",
	}
	for _, want := range wantNames {
		if !slices.Contains(gotNames, want) {
			t.Fatalf("工具未注册: %s, got=%v", want, gotNames)
		}
	}

	_, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "pods_log",
		Arguments: map[string]any{
			"namespace": "team-a",
			"name":      "pod-a",
			"container": "c1",
			"tailLines": 17,
			"format":    "json",
		},
	})
	if err != nil {
		t.Fatalf("调用 pods_log 失败: %v", err)
	}
	if len(manager.podsLogCalls) != 1 {
		t.Fatalf("pods_log 未透传到 manager, calls=%d", len(manager.podsLogCalls))
	}
	podsCall := manager.podsLogCalls[0]
	if podsCall.namespace != "team-a" || podsCall.name != "pod-a" || podsCall.container != "c1" || podsCall.tailLines != 17 {
		t.Fatalf("pods_log 参数透传异常: %+v", podsCall)
	}

	_, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "resources_get",
		Arguments: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"namespace":  "team-a",
			"name":       "demo",
			"format":     "yaml",
		},
	})
	if err != nil {
		t.Fatalf("调用 resources_get 失败: %v", err)
	}
	if len(manager.resourcesGetCalls) != 1 {
		t.Fatalf("resources_get 未透传到 manager, calls=%d", len(manager.resourcesGetCalls))
	}
	resCall := manager.resourcesGetCalls[0]
	if resCall.apiVersion != "apps/v1" || resCall.kind != "Deployment" || resCall.namespace != "team-a" || resCall.name != "demo" {
		t.Fatalf("resources_get 参数透传异常: %+v", resCall)
	}
}
