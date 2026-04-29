package kubernetes

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	metricsfake "k8s.io/metrics/pkg/client/clientset/versioned/fake"
)

func TestResolvePodContainer(t *testing.T) {
	t.Run("优先使用 default-container 注解", func(t *testing.T) {
		pod := corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					defaultContainerAnnotationKey: "c2",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{Name: "c1"}, {Name: "c2"}},
			},
		}
		got, err := resolvePodContainer(pod, "")
		if err != nil {
			t.Fatalf("resolvePodContainer returned error: %v", err)
		}
		if got != "c2" {
			t.Fatalf("expected c2, got %s", got)
		}
	})

	t.Run("仅一个容器时自动选择", func(t *testing.T) {
		pod := corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{Name: "only"}},
			},
		}
		got, err := resolvePodContainer(pod, "")
		if err != nil {
			t.Fatalf("resolvePodContainer returned error: %v", err)
		}
		if got != "only" {
			t.Fatalf("expected only, got %s", got)
		}
	})

	t.Run("多容器且无注解时选择首容器", func(t *testing.T) {
		pod := corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{Name: "first"}, {Name: "second"}},
			},
		}
		got, err := resolvePodContainer(pod, "")
		if err != nil {
			t.Fatalf("resolvePodContainer returned error: %v", err)
		}
		if got != "first" {
			t.Fatalf("expected first, got %s", got)
		}
	})
}

func TestPodsListAndGet(t *testing.T) {
	client := fake.NewSimpleClientset(
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "p1", Namespace: "ns1"},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "p2", Namespace: "ns2"},
		},
	)

	all, err := PodsListInAllNamespaces(context.Background(), client)
	if err != nil {
		t.Fatalf("PodsListInAllNamespaces returned error: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 pods, got %d", len(all))
	}

	nsPods, err := PodsListInNamespace(context.Background(), client, "ns1")
	if err != nil {
		t.Fatalf("PodsListInNamespace returned error: %v", err)
	}
	if len(nsPods) != 1 || nsPods[0].Name != "p1" {
		t.Fatalf("unexpected namespace pods: %+v", nsPods)
	}

	got, err := PodsGet(context.Background(), client, "ns1", "p1")
	if err != nil {
		t.Fatalf("PodsGet returned error: %v", err)
	}
	if got.Name != "p1" {
		t.Fatalf("expected pod p1, got %s", got.Name)
	}
}

func TestPodsTop(t *testing.T) {
	_, err := PodsTop(context.Background(), nil, "ns1", "p1")
	if err == nil {
		t.Fatal("expected error when metrics client is nil")
	}

	_, err = PodsTop(context.Background(), metricsfake.NewSimpleClientset(), "", "p1")
	if err == nil {
		t.Fatal("expected error when namespace is empty")
	}

	_, err = PodsTop(context.Background(), metricsfake.NewSimpleClientset(), "ns1", "")
	if err == nil {
		t.Fatal("expected error when pod name is empty")
	}
}
