package kubernetes

import (
	"context"
	"fmt"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NamespaceItem 表示单个命名空间的基础信息。
type NamespaceItem struct {
	Name  string `json:"name" yaml:"name"`
	Phase string `json:"phase" yaml:"phase"`
}

// NamespacesListResult 表示命名空间列表查询结果。
type NamespacesListResult struct {
	Items      []NamespaceItem `json:"items" yaml:"items"`
	TotalCount int             `json:"totalCount" yaml:"totalCount"`
}

// NamespacesList 只读列出集群中的命名空间。
func NamespacesList(ctx context.Context, client kubernetes.Interface) (*NamespacesListResult, error) {
	if client == nil {
		return nil, fmt.Errorf("kubernetes client is nil")
	}

	list, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list namespaces failed: %w", err)
	}

	items := make([]NamespaceItem, 0, len(list.Items))
	for _, it := range list.Items {
		items = append(items, NamespaceItem{
			Name:  it.Name,
			Phase: string(it.Status.Phase),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return &NamespacesListResult{
		Items:      items,
		TotalCount: len(items),
	}, nil
}
