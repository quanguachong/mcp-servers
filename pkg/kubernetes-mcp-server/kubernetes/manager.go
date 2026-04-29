package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	meta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
)

// Options 用于控制 manager 初始化参数。
type Options struct {
	Kubeconfig string
	Namespace  string
}

// Manager 聚合 Kubernetes 只读访问所需客户端。
type Manager struct {
	Clientset  kubernetes.Interface
	Dynamic    dynamic.Interface
	Discovery  discovery.DiscoveryInterface
	Metrics    metricsclient.Interface
	RESTMapper meta.RESTMapper
	namespace  string
}

// NewManager 按固定顺序加载配置并初始化只读客户端集合。
func NewManager(ctx context.Context, opts Options) (*Manager, error) {
	cfg, err := loadRESTConfig(opts.Kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化 clientset 失败: %w", err)
	}
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化 dynamic client 失败: %w", err)
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化 discovery client 失败: %w", err)
	}
	metricsClient, err := metricsclient.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化 metrics client 失败: %w", err)
	}

	groupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return nil, fmt.Errorf("加载 API 资源清单失败: %w", err)
	}

	_ = ctx // 预留上下文参数给后续调用链扩展。
	return &Manager{
		Clientset:  clientset,
		Dynamic:    dynamicClient,
		Discovery:  discoveryClient,
		Metrics:    metricsClient,
		RESTMapper: restmapper.NewDiscoveryRESTMapper(groupResources),
		namespace:  strings.TrimSpace(opts.Namespace),
	}, nil
}

func loadRESTConfig(flagPath string) (*rest.Config, error) {
	path, ok := selectKubeconfigPath(
		flagPath,
		os.Getenv("KUBECONFIG"),
		homedir.HomeDir(),
		func(p string) bool {
			_, err := os.Stat(p)
			return err == nil
		},
	)
	if ok {
		cfg, err := clientcmd.BuildConfigFromFlags("", path)
		if err != nil {
			return nil, fmt.Errorf("加载 kubeconfig 失败(%s): %w", path, err)
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("加载集群内配置失败: %w", err)
	}
	return cfg, nil
}

// selectKubeconfigPath 返回 kubeconfig 路径；false 表示应回退 InCluster。
func selectKubeconfigPath(flagPath, envValue, homeDir string, exists func(string) bool) (string, bool) {
	if p := strings.TrimSpace(flagPath); p != "" {
		return p, true
	}

	for _, p := range filepath.SplitList(envValue) {
		if p = strings.TrimSpace(p); p != "" {
			if exists(p) {
				return p, true
			}
		}
	}

	if homeDir != "" {
		p := filepath.Join(homeDir, ".kube", "config")
		if exists(p) {
			return p, true
		}
	}

	return "", false
}

// NamespaceOrDefault 返回指定或默认命名空间。
func (m *Manager) NamespaceOrDefault(ns string) string {
	if v := strings.TrimSpace(ns); v != "" {
		return v
	}
	if v := strings.TrimSpace(m.namespace); v != "" {
		return v
	}
	return "default"
}

// ResolveGVR 将资源参数解析为完整 GVR。
func (m *Manager) ResolveGVR(resource string) (schema.GroupVersionResource, error) {
	resource = strings.TrimSpace(resource)
	if resource == "" {
		return schema.GroupVersionResource{}, fmt.Errorf("资源名不能为空")
	}

	if parsed, _ := schema.ParseResourceArg(resource); parsed != nil {
		return *parsed, nil
	}
	if m.RESTMapper == nil {
		return schema.GroupVersionResource{}, fmt.Errorf("RESTMapper 未初始化")
	}

	gvr, err := m.RESTMapper.ResourceFor(schema.GroupVersionResource{Resource: resource})
	if err != nil {
		return schema.GroupVersionResource{}, fmt.Errorf("解析资源 %q 失败: %w", resource, err)
	}
	return gvr, nil
}
