package kubernetes

import (
	"context"
	"fmt"
	"io"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
	"k8s.io/utils/ptr"
)

const defaultContainerAnnotationKey = "kubectl.kubernetes.io/default-container"

// PodsListInAllNamespaces 列出所有命名空间中的 Pod。
func PodsListInAllNamespaces(ctx context.Context, client kubernetes.Interface) ([]corev1.Pod, error) {
	if client == nil {
		return nil, fmt.Errorf("kubernetes client 不能为空")
	}
	list, err := client.CoreV1().Pods(metav1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// PodsListInNamespace 列出指定命名空间中的 Pod。
func PodsListInNamespace(ctx context.Context, client kubernetes.Interface, namespace string) ([]corev1.Pod, error) {
	if client == nil {
		return nil, fmt.Errorf("kubernetes client 不能为空")
	}
	if strings.TrimSpace(namespace) == "" {
		return nil, fmt.Errorf("namespace 不能为空")
	}
	list, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// PodsGet 获取指定命名空间和名称的 Pod。
func PodsGet(ctx context.Context, client kubernetes.Interface, namespace string, name string) (*corev1.Pod, error) {
	if client == nil {
		return nil, fmt.Errorf("kubernetes client 不能为空")
	}
	if strings.TrimSpace(namespace) == "" {
		return nil, fmt.Errorf("namespace 不能为空")
	}
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("pod name 不能为空")
	}
	return client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
}

// PodsLog 读取 Pod 日志；当未显式指定 container 时，会按注解/单容器/首容器顺序自动解析。
func PodsLog(ctx context.Context, client kubernetes.Interface, namespace string, name string, container string, opts corev1.PodLogOptions) (string, error) {
	pod, err := PodsGet(ctx, client, namespace, name)
	if err != nil {
		return "", err
	}

	resolvedContainer, err := resolvePodContainer(*pod, container)
	if err != nil {
		return "", err
	}
	opts.Container = resolvedContainer

	stream, err := client.CoreV1().Pods(namespace).GetLogs(name, &opts).Stream(ctx)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	data, err := io.ReadAll(stream)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// PodsTop 获取单个 Pod 的资源使用情况。
func PodsTop(ctx context.Context, metricsClient metricsclientset.Interface, namespace string, name string) (*metricsv1beta1.PodMetrics, error) {
	if metricsClient == nil {
		return nil, fmt.Errorf("metrics client 不能为空")
	}
	if strings.TrimSpace(namespace) == "" {
		return nil, fmt.Errorf("namespace 不能为空")
	}
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("pod name 不能为空")
	}
	return metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(ctx, name, metav1.GetOptions{})
}

func resolvePodContainer(pod corev1.Pod, container string) (string, error) {
	if strings.TrimSpace(container) != "" {
		if podHasContainer(pod, container) {
			return container, nil
		}
		return "", fmt.Errorf("pod %s/%s 不存在容器 %q", pod.Namespace, pod.Name, container)
	}

	if defaultContainer := strings.TrimSpace(pod.Annotations[defaultContainerAnnotationKey]); defaultContainer != "" && podHasContainer(pod, defaultContainer) {
		return defaultContainer, nil
	}

	if len(pod.Spec.Containers) == 1 {
		return pod.Spec.Containers[0].Name, nil
	}
	if len(pod.Spec.Containers) > 1 {
		return pod.Spec.Containers[0].Name, nil
	}
	return "", fmt.Errorf("pod %s/%s 不包含可用容器", pod.Namespace, pod.Name)
}

func podHasContainer(pod corev1.Pod, name string) bool {
	for _, c := range pod.Spec.Containers {
		if c.Name == name {
			return true
		}
	}
	return false
}

// BuildPodLogOptions 构造只读日志查询选项。
func BuildPodLogOptions(tailLines int64) corev1.PodLogOptions {
	opts := corev1.PodLogOptions{}
	if tailLines > 0 {
		opts.TailLines = ptr.To(tailLines)
	}
	return opts
}
