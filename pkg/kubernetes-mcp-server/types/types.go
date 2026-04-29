package types

// EventsListInput 为 events_list 工具输入。
type EventsListInput struct {
	Namespace     string `json:"namespace,omitempty" jsonschema:"Namespace to list events from; empty means all namespaces"`
	FieldSelector string `json:"fieldSelector,omitempty" jsonschema:"Field selector used to filter events"`
	Limit         *int64 `json:"limit,omitempty" jsonschema:"Maximum number of events to return"`
}

// EventsListOutput 为 events_list 工具输出。
type EventsListOutput struct {
	Events     []EventInfo `json:"events"`
	TotalCount int         `json:"totalCount"`
}

type EventInfo struct {
	Name           string `json:"name"`
	Namespace      string `json:"namespace"`
	Type           string `json:"type"`
	Reason         string `json:"reason"`
	Message        string `json:"message"`
	InvolvedObject string `json:"involvedObject"`
	LastTimestamp  string `json:"lastTimestamp,omitempty"`
}

// NamespacesListInput 为 namespaces_list 工具输入。
type NamespacesListInput struct{}

// NamespacesListOutput 为 namespaces_list 工具输出。
type NamespacesListOutput struct {
	Namespaces []NamespaceInfo `json:"namespaces"`
	TotalCount int             `json:"totalCount"`
}

type NamespaceInfo struct {
	Name   string `json:"name"`
	Status string `json:"status,omitempty"`
}

// NodesLogInput 为 nodes_log 工具输入。
type NodesLogInput struct {
	Name         string `json:"name" jsonschema:"Node name"`
	SinceSeconds *int64 `json:"sinceSeconds,omitempty" jsonschema:"Only return logs newer than this many seconds"`
	TailLines    *int64 `json:"tailLines,omitempty" jsonschema:"Number of lines from the end of the logs to return"`
	LimitBytes   *int64 `json:"limitBytes,omitempty" jsonschema:"Maximum number of log bytes to return"`
	Timestamps   bool   `json:"timestamps,omitempty" jsonschema:"Include timestamps on each log line"`
}

// NodesLogOutput 为 nodes_log 工具输出。
type NodesLogOutput struct {
	Node string `json:"node"`
	Log  string `json:"log"`
}

// NodesStatsSummaryInput 为 nodes_stats_summary 工具输入。
type NodesStatsSummaryInput struct {
	Name string `json:"name,omitempty" jsonschema:"Node name; empty means all nodes summary"`
}

// NodesStatsSummaryOutput 为 nodes_stats_summary 工具输出。
type NodesStatsSummaryOutput struct {
	Summary map[string]any `json:"summary"`
}

// NodesTopInput 为 nodes_top 工具输入。
type NodesTopInput struct {
	Name string `json:"name,omitempty" jsonschema:"Node name; empty means all nodes"`
}

// NodesTopOutput 为 nodes_top 工具输出。
type NodesTopOutput struct {
	Nodes []NodeTopInfo `json:"nodes"`
}

type NodeTopInfo struct {
	Name   string `json:"name"`
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// PodsListInput 为 pods_list 工具输入。
type PodsListInput struct {
	LabelSelector string `json:"labelSelector,omitempty" jsonschema:"Label selector used to filter pods"`
	FieldSelector string `json:"fieldSelector,omitempty" jsonschema:"Field selector used to filter pods"`
	Limit         *int64 `json:"limit,omitempty" jsonschema:"Maximum number of pods to return"`
}

// PodsListOutput 为 pods_list 工具输出。
type PodsListOutput struct {
	Pods       []PodInfo `json:"pods"`
	TotalCount int       `json:"totalCount"`
}

// PodsListInNamespaceInput 为 pods_list_in_namespace 工具输入。
type PodsListInNamespaceInput struct {
	Namespace     string `json:"namespace" jsonschema:"Namespace name"`
	LabelSelector string `json:"labelSelector,omitempty" jsonschema:"Label selector used to filter pods"`
	FieldSelector string `json:"fieldSelector,omitempty" jsonschema:"Field selector used to filter pods"`
	Limit         *int64 `json:"limit,omitempty" jsonschema:"Maximum number of pods to return"`
}

// PodsListInNamespaceOutput 为 pods_list_in_namespace 工具输出。
type PodsListInNamespaceOutput struct {
	Pods       []PodInfo `json:"pods"`
	TotalCount int       `json:"totalCount"`
}

type PodInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	NodeName  string `json:"nodeName,omitempty"`
	Phase     string `json:"phase,omitempty"`
}

// PodsGetInput 为 pods_get 工具输入。
type PodsGetInput struct {
	Namespace string `json:"namespace" jsonschema:"Namespace name"`
	Name      string `json:"name" jsonschema:"Pod name"`
}

// PodsGetOutput 为 pods_get 工具输出。
type PodsGetOutput struct {
	Pod map[string]any `json:"pod"`
}

// PodsLogInput 为 pods_log 工具输入。
type PodsLogInput struct {
	Namespace    string `json:"namespace" jsonschema:"Namespace name"`
	Name         string `json:"name" jsonschema:"Pod name"`
	Container    string `json:"container,omitempty" jsonschema:"Container name"`
	Previous     bool   `json:"previous,omitempty" jsonschema:"Return logs for the previous terminated container instance"`
	SinceSeconds *int64 `json:"sinceSeconds,omitempty" jsonschema:"Only return logs newer than this many seconds"`
	TailLines    *int64 `json:"tailLines,omitempty" jsonschema:"Number of lines from the end of the logs to return"`
	LimitBytes   *int64 `json:"limitBytes,omitempty" jsonschema:"Maximum number of log bytes to return"`
	Timestamps   bool   `json:"timestamps,omitempty" jsonschema:"Include timestamps on each log line"`
}

// PodsLogOutput 为 pods_log 工具输出。
type PodsLogOutput struct {
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
	Log       string `json:"log"`
}

// PodsTopInput 为 pods_top 工具输入。
type PodsTopInput struct {
	Namespace string `json:"namespace,omitempty" jsonschema:"Namespace name; empty means all namespaces"`
	Name      string `json:"name,omitempty" jsonschema:"Pod name; empty means all pods"`
}

// PodsTopOutput 为 pods_top 工具输出。
type PodsTopOutput struct {
	Pods []PodTopInfo `json:"pods"`
}

type PodTopInfo struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	CPU        string `json:"cpu"`
	Memory     string `json:"memory"`
	Containers int    `json:"containers,omitempty"`
}

// ResourcesListInput 为 resources_list 工具输入。
type ResourcesListInput struct {
	APIVersion    string `json:"apiVersion" jsonschema:"Kubernetes apiVersion, for example apps/v1"`
	Kind          string `json:"kind" jsonschema:"Kubernetes kind, for example Deployment"`
	Namespace     string `json:"namespace,omitempty" jsonschema:"Namespace name; empty for cluster scoped resources"`
	LabelSelector string `json:"labelSelector,omitempty" jsonschema:"Label selector used to filter resources"`
	FieldSelector string `json:"fieldSelector,omitempty" jsonschema:"Field selector used to filter resources"`
	Limit         *int64 `json:"limit,omitempty" jsonschema:"Maximum number of resources to return"`
}

// ResourcesListOutput 为 resources_list 工具输出。
type ResourcesListOutput struct {
	Resources  []map[string]any `json:"resources"`
	TotalCount int              `json:"totalCount"`
}

// ResourcesGetInput 为 resources_get 工具输入。
type ResourcesGetInput struct {
	APIVersion string `json:"apiVersion" jsonschema:"Kubernetes apiVersion, for example apps/v1"`
	Kind       string `json:"kind" jsonschema:"Kubernetes kind, for example Deployment"`
	Name       string `json:"name" jsonschema:"Resource name"`
	Namespace  string `json:"namespace,omitempty" jsonschema:"Namespace name for namespaced resources"`
}

// ResourcesGetOutput 为 resources_get 工具输出。
type ResourcesGetOutput struct {
	Resource map[string]any `json:"resource"`
}
