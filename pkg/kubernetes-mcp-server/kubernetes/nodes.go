package kubernetes

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
)

// NodeTopItem 表示单个节点的资源使用摘要。
type NodeTopItem struct {
	Name        string `json:"name" yaml:"name"`
	CPU         string `json:"cpu" yaml:"cpu"`
	Memory      string `json:"memory" yaml:"memory"`
	Timestamp   string `json:"timestamp" yaml:"timestamp"`
	Window      string `json:"window" yaml:"window"`
	NodeMetrics any    `json:"nodeMetrics,omitempty" yaml:"nodeMetrics,omitempty"`
}

// NodesTopResult 表示节点 top 查询结果及 metrics API 可用性。
type NodesTopResult struct {
	Available bool          `json:"available" yaml:"available"`
	Reason    string        `json:"reason,omitempty" yaml:"reason,omitempty"`
	Items     []NodeTopItem `json:"items,omitempty" yaml:"items,omitempty"`
}

// NodesLog 通过 kubelet proxy 读取节点日志，返回原始字节数据。
// 仅执行只读 GET 请求，不会修改集群状态。
func NodesLog(ctx context.Context, client kubernetes.Interface, nodeName string, query url.Values) ([]byte, error) {
	if client == nil {
		return nil, fmt.Errorf("kubernetes client is nil")
	}
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}

	req := client.CoreV1().RESTClient().Get().
		Resource("nodes").
		Name(nodeName).
		SubResource("proxy", "logs")
	if len(query) > 0 {
		req = req.VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec)
		for k, values := range query {
			for _, v := range values {
				req = req.Param(k, v)
			}
		}
	}
	return req.Do(ctx).Raw()
}

// NodesStatsSummary 通过 kubelet proxy 读取节点 stats/summary 数据。
// 返回 kubelet 原始响应，便于上层按 JSON/YAML 自行序列化。
func NodesStatsSummary(ctx context.Context, client kubernetes.Interface, nodeName string) ([]byte, error) {
	if client == nil {
		return nil, fmt.Errorf("kubernetes client is nil")
	}
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}

	return client.CoreV1().RESTClient().Get().
		Resource("nodes").
		Name(nodeName).
		SubResource("proxy", "stats", "summary").
		Do(ctx).
		Raw()
}

// NodesTop 查询 metrics.k8s.io 的 NodeMetrics，提供节点资源使用信息。
// 当 metrics API 不可用时，不返回错误中断，而是给出可用性状态与原因。
func NodesTop(ctx context.Context, metricsCli metricsclient.Interface) (*NodesTopResult, error) {
	if metricsCli == nil {
		return &NodesTopResult{
			Available: false,
			Reason:    "metrics client is nil",
		}, nil
	}

	list, err := metricsCli.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return &NodesTopResult{
			Available: false,
			Reason:    err.Error(),
		}, nil
	}

	items := make([]NodeTopItem, 0, len(list.Items))
	for _, it := range list.Items {
		items = append(items, NodeTopItem{
			Name:        it.Name,
			CPU:         it.Usage.Cpu().String(),
			Memory:      it.Usage.Memory().String(),
			Timestamp:   it.Timestamp.Time.Format("2006-01-02T15:04:05Z07:00"),
			Window:      strconv.FormatInt(int64(it.Window.Duration.Seconds()), 10) + "s",
			NodeMetrics: it,
		})
	}

	// 结果按节点名排序，保证输出稳定，便于测试与对比。
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return &NodesTopResult{
		Available: true,
		Items:     items,
	}, nil
}
