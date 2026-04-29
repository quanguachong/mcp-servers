package kubernetes

import (
	"testing"

	meta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestSelectKubeconfigPathOrder(t *testing.T) {
	const home = "/home/tester"
	defaultPath := home + "/.kube/config"

	t.Run("flag 优先级最高", func(t *testing.T) {
		path, ok := selectKubeconfigPath("/tmp/flag", "/tmp/env", home, func(string) bool { return true })
		if !ok || path != "/tmp/flag" {
			t.Fatalf("期望使用 flag 路径, got ok=%v path=%q", ok, path)
		}
	})

	t.Run("env 次优先", func(t *testing.T) {
		path, ok := selectKubeconfigPath("", "/tmp/missing:/tmp/env2", home, func(p string) bool { return p == "/tmp/env2" })
		if !ok || path != "/tmp/env2" {
			t.Fatalf("期望跳过不存在路径并使用可用 env 路径, got ok=%v path=%q", ok, path)
		}
	})

	t.Run("home 默认配置再次之", func(t *testing.T) {
		path, ok := selectKubeconfigPath("", "", home, func(p string) bool { return p == defaultPath })
		if !ok || path != defaultPath {
			t.Fatalf("期望使用 home 默认路径, got ok=%v path=%q", ok, path)
		}
	})

	t.Run("最终回退 incluster", func(t *testing.T) {
		path, ok := selectKubeconfigPath("", "", home, func(string) bool { return false })
		if ok || path != "" {
			t.Fatalf("期望回退 incluster, got ok=%v path=%q", ok, path)
		}
	})
}

func TestNamespaceOrDefault(t *testing.T) {
	m := &Manager{namespace: "team-a"}

	if got := m.NamespaceOrDefault("custom"); got != "custom" {
		t.Fatalf("期望 custom, got %q", got)
	}
	if got := m.NamespaceOrDefault(""); got != "team-a" {
		t.Fatalf("期望 team-a, got %q", got)
	}

	m.namespace = ""
	if got := m.NamespaceOrDefault(""); got != "default" {
		t.Fatalf("期望 default, got %q", got)
	}
}

func TestResolveGVR(t *testing.T) {
	mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{{Version: "v1"}})
	mapper.AddSpecific(
		schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
		schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
		schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pod"},
		meta.RESTScopeNamespace,
	)

	m := &Manager{RESTMapper: mapper}

	t.Run("可直接解析完整参数", func(t *testing.T) {
		gvr, err := m.ResolveGVR("deployments.v1.apps")
		if err != nil {
			t.Fatalf("不应报错: %v", err)
		}
		if gvr.Resource != "deployments" || gvr.Version != "v1" || gvr.Group != "apps" {
			t.Fatalf("解析结果异常: %+v", gvr)
		}
	})

	t.Run("可通过 mapper 解析短资源名", func(t *testing.T) {
		gvr, err := m.ResolveGVR("pods")
		if err != nil {
			t.Fatalf("不应报错: %v", err)
		}
		if gvr.Resource != "pods" || gvr.Version != "v1" {
			t.Fatalf("解析结果异常: %+v", gvr)
		}
	})
}
