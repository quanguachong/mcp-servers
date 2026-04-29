package kubernetesmcpserver

import (
	_ "k8s.io/api/core/v1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/kubernetes"
	_ "k8s.io/metrics/pkg/client/clientset/versioned"
	_ "k8s.io/utils/ptr"
	_ "sigs.k8s.io/yaml"
)
