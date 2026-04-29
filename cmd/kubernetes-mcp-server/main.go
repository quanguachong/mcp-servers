package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/quanguachong/mcp-servers/pkg/kubernetes-mcp-server/kubernetes"
	"github.com/quanguachong/mcp-servers/pkg/kubernetes-mcp-server/mcpserver"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type managerAdapter struct {
	manager *kubernetes.Manager
}

func (a *managerAdapter) ListEvents(ctx context.Context, namespace string) (any, error) {
	data, err := kubernetes.EventsList(ctx, a.manager.Clientset)
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		return data, nil
	}
	filtered := make([]kubernetes.EventItem, 0, len(data.Items))
	for _, item := range data.Items {
		if item.Namespace == namespace {
			filtered = append(filtered, item)
		}
	}
	return &kubernetes.EventsListResult{
		Items:      filtered,
		TotalCount: len(filtered),
	}, nil
}

func (a *managerAdapter) ListNamespaces(ctx context.Context) (any, error) {
	return kubernetes.NamespacesList(ctx, a.manager.Clientset)
}

func (a *managerAdapter) NodesLog(ctx context.Context, nodeName string, tailLines int64) (any, error) {
	query := url.Values{}
	if tailLines > 0 {
		query.Set("tailLines", fmt.Sprintf("%d", tailLines))
	}
	data, err := kubernetes.NodesLog(ctx, a.manager.Clientset, nodeName, query)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"node": nodeName,
		"log":  string(data),
	}, nil
}

func (a *managerAdapter) NodesStatsSummary(ctx context.Context, nodeName string) (any, error) {
	data, err := kubernetes.NodesStatsSummary(ctx, a.manager.Clientset, nodeName)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"node":    nodeName,
		"summary": string(data),
	}, nil
}

func (a *managerAdapter) NodesTop(ctx context.Context) (any, error) {
	return kubernetes.NodesTop(ctx, a.manager.Metrics)
}

func (a *managerAdapter) PodsList(ctx context.Context) (any, error) {
	return kubernetes.PodsListInAllNamespaces(ctx, a.manager.Clientset)
}

func (a *managerAdapter) PodsListInNamespace(ctx context.Context, namespace string) (any, error) {
	return kubernetes.PodsListInNamespace(ctx, a.manager.Clientset, a.manager.NamespaceOrDefault(namespace))
}

func (a *managerAdapter) PodsGet(ctx context.Context, namespace, name string) (any, error) {
	return kubernetes.PodsGet(ctx, a.manager.Clientset, a.manager.NamespaceOrDefault(namespace), name)
}

func (a *managerAdapter) PodsLog(ctx context.Context, namespace, name, container string, tailLines int64) (any, error) {
	ns := a.manager.NamespaceOrDefault(namespace)
	opts := kubernetes.BuildPodLogOptions(tailLines)
	return kubernetes.PodsLog(ctx, a.manager.Clientset, ns, name, container, opts)
}

func (a *managerAdapter) PodsTop(ctx context.Context, namespace, name string) (any, error) {
	return kubernetes.PodsTop(ctx, a.manager.Metrics, a.manager.NamespaceOrDefault(namespace), name)
}

func (a *managerAdapter) ResourcesList(ctx context.Context, apiVersion, kind, namespace string) (any, error) {
	resource, err := resolveResourceArg(a.manager, apiVersion, kind)
	if err != nil {
		return nil, err
	}
	return kubernetes.ResourcesList(ctx, a.manager.Dynamic, a.manager.RESTMapper, resource, namespace)
}

func (a *managerAdapter) ResourcesGet(ctx context.Context, apiVersion, kind, namespace, name string) (any, error) {
	resource, err := resolveResourceArg(a.manager, apiVersion, kind)
	if err != nil {
		return nil, err
	}
	return kubernetes.ResourcesGet(ctx, a.manager.Dynamic, a.manager.RESTMapper, resource, namespace, name)
}

func resolveResourceArg(manager *kubernetes.Manager, apiVersion, kind string) (string, error) {
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return "", fmt.Errorf("apiVersion 无效: %w", err)
	}
	mapping, err := manager.RESTMapper.RESTMapping(schema.GroupKind{Group: gv.Group, Kind: kind}, gv.Version)
	if err != nil {
		return "", fmt.Errorf("解析资源映射失败: %w", err)
	}
	return mapping.Resource.String(), nil
}

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to kubeconfig file (optional)")
	namespace := flag.String("namespace", "", "Default namespace (optional)")
	flag.Parse()

	manager, err := kubernetes.NewManager(context.Background(), kubernetes.Options{
		Kubeconfig: *kubeconfig,
		Namespace:  *namespace,
	})
	if err != nil {
		log.Fatal(err)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "kubernetes-mcp-server",
		Version: "v1.0.0",
	}, nil)

	mcpserver.RegisterTools(server, &managerAdapter{manager: manager})

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
