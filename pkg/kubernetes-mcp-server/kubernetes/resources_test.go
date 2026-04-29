package kubernetes

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/dynamic/fake"
)

func TestResourcesList_NamespacedResourceRequiresNamespace(t *testing.T) {
	dyn := fake.NewSimpleDynamicClient(runtime.NewScheme())
	mapper := newTestRESTMapper()

	_, err := ResourcesList(context.Background(), dyn, mapper, "pods", "")
	if err == nil {
		t.Fatalf("expected error when namespace is empty for namespaced resource")
	}
}

func TestResourcesGet_ClusterResourceIgnoresNamespace(t *testing.T) {
	dyn := fake.NewSimpleDynamicClient(runtime.NewScheme(), &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Namespace",
			"metadata": map[string]any{
				"name": "default",
			},
		},
	})
	mapper := newTestRESTMapper()

	got, err := ResourcesGet(context.Background(), dyn, mapper, "namespaces", "ignored", "default")
	if err != nil {
		t.Fatalf("ResourcesGet returned error: %v", err)
	}
	if got.GetName() != "default" {
		t.Fatalf("unexpected name: %s", got.GetName())
	}
}

func newTestRESTMapper() meta.RESTMapper {
	mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{
		{Version: "v1"},
	})
	mapper.AddSpecific(
		schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
		schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
		schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
		meta.RESTScopeNamespace,
	)
	mapper.AddSpecific(
		schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"},
		schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"},
		schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespace"},
		meta.RESTScopeRoot,
	)
	return mapper
}

var _ = metav1.ListOptions{}
