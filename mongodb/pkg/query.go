package pkg

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
)

// objectIDCall 匹配 mongosh 风格的 ObjectId("24位十六进制")
var objectIDCall = regexp.MustCompile(`ObjectId\s*\(\s*"([a-fA-F0-9]{24})"\s*\)`)

// normalizeSelector 将 ObjectId("...") 转为 Canonical Extended JSON 的 {"$oid":"..."} 片段，便于后续 UnmarshalExtJSON。
func normalizeSelector(selector string) string {
	// $$ 在替换模板中表示字面量 $，否则 $oid 会被误解析
	return objectIDCall.ReplaceAllString(selector, `{"$$oid":"$1"}`)
}

func parseFilter(selector string) (bson.M, error) {
	normalized := normalizeSelector(selector)
	var filter bson.M
	if err := bson.UnmarshalExtJSON([]byte(normalized), true, &filter); err != nil {
		return nil, fmt.Errorf("解析 selector 失败: %w", err)
	}
	return filter, nil
}
