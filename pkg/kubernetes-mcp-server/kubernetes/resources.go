package kubernetes

import (
	"context"
	"fmt"
	"strings"

	meta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// ResourcesList 根据资源类型列出 Kubernetes 资源。
func ResourcesList(
	ctx context.Context,
	dyn dynamic.Interface,
	mapper meta.RESTMapper,
	resource string,
	namespace string,
) (*unstructured.UnstructuredList, error) {
	if dyn == nil {
		return nil, fmt.Errorf("dynamic client 未初始化")
	}
	gvr, scopeName, err := resolveResourceAndScope(mapper, resource)
	if err != nil {
		return nil, err
	}
	resourceClient, err := dynamicResourceClient(dyn, gvr, scopeName, namespace)
	if err != nil {
		return nil, err
	}
	return resourceClient.List(ctx, metav1.ListOptions{})
}

// ResourcesGet 根据资源类型读取单个 Kubernetes 资源。
func ResourcesGet(
	ctx context.Context,
	dyn dynamic.Interface,
	mapper meta.RESTMapper,
	resource string,
	namespace string,
	name string,
) (*unstructured.Unstructured, error) {
	if dyn == nil {
		return nil, fmt.Errorf("dynamic client 未初始化")
	}
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("资源名称不能为空")
	}
	gvr, scopeName, err := resolveResourceAndScope(mapper, resource)
	if err != nil {
		return nil, err
	}
	resourceClient, err := dynamicResourceClient(dyn, gvr, scopeName, namespace)
	if err != nil {
		return nil, err
	}
	return resourceClient.Get(ctx, name, metav1.GetOptions{})
}

func resolveResourceAndScope(mapper meta.RESTMapper, resource string) (schema.GroupVersionResource, meta.RESTScopeName, error) {
	if dynErr := validateDynamicMapper(mapper); dynErr != nil {
		return schema.GroupVersionResource{}, "", dynErr
	}
	resource = strings.TrimSpace(resource)
	if resource == "" {
		return schema.GroupVersionResource{}, "", fmt.Errorf("资源名不能为空")
	}

	parsed, _ := schema.ParseResourceArg(resource)
	if parsed == nil {
		parsed = &schema.GroupVersionResource{Resource: resource}
	}
	gvr, err := mapper.ResourceFor(*parsed)
	if err != nil {
		return schema.GroupVersionResource{}, "", fmt.Errorf("解析资源 %q 失败: %w", resource, err)
	}

	gvk, err := mapper.KindFor(gvr)
	if err != nil {
		return schema.GroupVersionResource{}, "", fmt.Errorf("解析资源 %q 的 GVK 失败: %w", resource, err)
	}
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, "", fmt.Errorf("解析资源 %q 的作用域失败: %w", resource, err)
	}
	return gvr, mapping.Scope.Name(), nil
}

func dynamicResourceClient(
	dyn dynamic.Interface,
	gvr schema.GroupVersionResource,
	scopeName meta.RESTScopeName,
	namespace string,
) (dynamic.ResourceInterface, error) {
	if scopeName == meta.RESTScopeNameNamespace {
		namespace = strings.TrimSpace(namespace)
		if namespace == "" {
			return nil, fmt.Errorf("namespaced 资源必须指定 namespace")
		}
		return dyn.Resource(gvr).Namespace(namespace), nil
	}
	return dyn.Resource(gvr), nil
}

func validateDynamicMapper(mapper meta.RESTMapper) error {
	if mapper == nil {
		return fmt.Errorf("RESTMapper 未初始化")
	}
	return nil
}

