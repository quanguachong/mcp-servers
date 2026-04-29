package kubernetes

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestEventsListNilClient(t *testing.T) {
	_, err := EventsList(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error when client is nil")
	}
}

func TestEventsList(t *testing.T) {
	client := fake.NewSimpleClientset(
		&corev1.Event{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns-b",
				Name:      "event-2",
			},
			Reason:  "Scaled",
			Message: "scaled deployment",
			Type:    "Normal",
			Count:   2,
		},
		&corev1.Event{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns-a",
				Name:      "event-1",
			},
			Reason:  "Failed",
			Message: "image pull failed",
			Type:    "Warning",
			Count:   1,
		},
	)

	got, err := EventsList(context.Background(), client)
	if err != nil {
		t.Fatalf("EventsList returned error: %v", err)
	}
	if got.TotalCount != 2 {
		t.Fatalf("unexpected total count: %d", got.TotalCount)
	}
	if len(got.Items) != 2 {
		t.Fatalf("unexpected items length: %d", len(got.Items))
	}
	if got.Items[0].Namespace != "ns-a" || got.Items[0].Name != "event-1" {
		t.Fatalf("items should be sorted by namespace/name, got first=%+v", got.Items[0])
	}
	if got.Items[0].Reason != "Failed" || got.Items[0].Type != "Warning" || got.Items[0].Count != 1 {
		t.Fatalf("unexpected first item fields: %+v", got.Items[0])
	}
}
