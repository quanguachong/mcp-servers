package mongodbmcpserver

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
)

var objectIDCall = regexp.MustCompile(`ObjectId\s*\(\s*"([a-fA-F0-9]{24})"\s*\)`)

func normalizeSelector(selector string) string {
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
