package kubernetes

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNamespacesListNilClient(t *testing.T) {
	_, err := NamespacesList(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error when client is nil")
	}
}

func TestNamespacesList(t *testing.T) {
	client := fake.NewSimpleClientset(
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "z-ns"},
			Status:     corev1.NamespaceStatus{Phase: corev1.NamespaceTerminating},
		},
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "a-ns"},
			Status:     corev1.NamespaceStatus{Phase: corev1.NamespaceActive},
		},
	)

	got, err := NamespacesList(context.Background(), client)
	if err != nil {
		t.Fatalf("NamespacesList returned error: %v", err)
	}
	if got.TotalCount != 2 {
		t.Fatalf("unexpected total count: %d", got.TotalCount)
	}
	if len(got.Items) != 2 {
		t.Fatalf("unexpected items length: %d", len(got.Items))
	}
	if got.Items[0].Name != "a-ns" || got.Items[0].Phase != string(corev1.NamespaceActive) {
		t.Fatalf("unexpected first item: %+v", got.Items[0])
	}
}
