package kubernetes

import (
	"context"
	"fmt"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// EventItem 表示单条 Kubernetes Event 的精简信息。
type EventItem struct {
	Namespace string `json:"namespace" yaml:"namespace"`
	Name      string `json:"name" yaml:"name"`
	Reason    string `json:"reason,omitempty" yaml:"reason,omitempty"`
	Message   string `json:"message,omitempty" yaml:"message,omitempty"`
	Type      string `json:"type,omitempty" yaml:"type,omitempty"`
	Count     int32  `json:"count" yaml:"count"`
}

// EventsListResult 表示 Event 列表查询结果。
type EventsListResult struct {
	Items      []EventItem `json:"items" yaml:"items"`
	TotalCount int         `json:"totalCount" yaml:"totalCount"`
}

// EventsList 只读列出集群中的 Event。
func EventsList(ctx context.Context, client kubernetes.Interface) (*EventsListResult, error) {
	if client == nil {
		return nil, fmt.Errorf("kubernetes client is nil")
	}

	list, err := client.CoreV1().Events("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list events failed: %w", err)
	}

	items := make([]EventItem, 0, len(list.Items))
	for _, it := range list.Items {
		items = append(items, EventItem{
			Namespace: it.Namespace,
			Name:      it.Name,
			Reason:    it.Reason,
			Message:   it.Message,
			Type:      it.Type,
			Count:     it.Count,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return items[i].Name < items[j].Name
		}
		return items[i].Namespace < items[j].Namespace
	})

	return &EventsListResult{
		Items:      items,
		TotalCount: len(items),
	}, nil
}
