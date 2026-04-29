package kubernetes

import (
	"encoding/json"

	"sigs.k8s.io/yaml"
)

// MarshalYaml 将任意对象序列化为 YAML 字节。
func MarshalYaml(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

// MarshalJSONPretty 将任意对象序列化为缩进格式的 JSON 字节。
func MarshalJSONPretty(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}
