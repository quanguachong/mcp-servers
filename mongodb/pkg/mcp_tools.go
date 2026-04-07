package pkg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultFindLimit int64 = 10
const defaultResponseBytesLimit int64 = 1048576

type countInput struct {
	Database   string         `json:"database" jsonschema:"Database name"`
	Collection string         `json:"collection" jsonschema:"Collection name"`
	Query      map[string]any `json:"query,omitempty" jsonschema:"Filter query"`
}

type countOutput struct {
	Count int64 `json:"count"`
}

type findInput struct {
	Database           string         `json:"database" jsonschema:"Database name"`
	Collection         string         `json:"collection" jsonschema:"Collection name"`
	Filter             map[string]any `json:"filter,omitempty" jsonschema:"Query filter"`
	Projection         map[string]any `json:"projection,omitempty" jsonschema:"Projection document"`
	Limit              *int64         `json:"limit,omitempty" jsonschema:"Maximum number of documents to return"`
	Sort               map[string]any `json:"sort,omitempty" jsonschema:"Sort document"`
	ResponseBytesLimit *int64         `json:"responseBytesLimit,omitempty" jsonschema:"Maximum response size in bytes"`
}

type findOutput struct {
	Documents []bson.M `json:"documents"`
}

type listDatabasesOutput struct {
	Databases  []databaseInfo `json:"databases"`
	TotalCount int            `json:"totalCount"`
}

type databaseInfo struct {
	Name string  `json:"name"`
	Size float64 `json:"size"`
}

type listCollectionsInput struct {
	Database string `json:"database" jsonschema:"Database name"`
}

type listCollectionsOutput struct {
	Collections []collectionInfo `json:"collections"`
	TotalCount  int              `json:"totalCount"`
}

type collectionInfo struct {
	Name string `json:"name"`
}

func RegisterMongoTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "count",
		Description: "Gets the number of documents in a MongoDB collection",
	}, handleCount)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "find",
		Description: "Run a find query against a MongoDB collection",
	}, handleFind)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list-databases",
		Description: "List all databases for a MongoDB connection",
	}, handleListDatabases)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list-collections",
		Description: "List all collections for a given database",
	}, handleListCollections)
}

func handleCount(ctx context.Context, _ *mcp.CallToolRequest, in countInput) (*mcp.CallToolResult, countOutput, error) {
	if in.Database == "" || in.Collection == "" {
		return nil, countOutput{}, fmt.Errorf("database 和 collection 为必填参数")
	}
	query := bson.M{}
	if in.Query != nil {
		query = bson.M(in.Query)
	}

	var out countOutput
	err := withMongoClient(ctx, func(opCtx context.Context, client *mongo.Client) error {
		coll := client.Database(in.Database).Collection(in.Collection)
		count, err := coll.CountDocuments(opCtx, query)
		if err != nil {
			return fmt.Errorf("执行 count 失败: %w", err)
		}
		out.Count = count
		return nil
	})
	return nil, out, err
}

func handleFind(ctx context.Context, _ *mcp.CallToolRequest, in findInput) (*mcp.CallToolResult, findOutput, error) {
	if in.Database == "" || in.Collection == "" {
		return nil, findOutput{}, fmt.Errorf("database 和 collection 为必填参数")
	}

	limit, responseBytesLimit := applyFindDefaults(in.Limit, in.ResponseBytesLimit)
	if limit < 0 {
		return nil, findOutput{}, fmt.Errorf("limit 不能小于 0")
	}
	if responseBytesLimit <= 0 {
		return nil, findOutput{}, fmt.Errorf("responseBytesLimit 必须大于 0")
	}

	filter := bson.M{}
	if in.Filter != nil {
		filter = bson.M(in.Filter)
	}

	findOpts := options.Find()
	if limit > 0 {
		findOpts.SetLimit(limit)
	}
	if in.Projection != nil {
		findOpts.SetProjection(bson.M(in.Projection))
	}
	if in.Sort != nil {
		findOpts.SetSort(bson.M(in.Sort))
	}

	out := findOutput{Documents: []bson.M{}}
	err := withMongoClient(ctx, func(opCtx context.Context, client *mongo.Client) error {
		coll := client.Database(in.Database).Collection(in.Collection)
		cursor, err := coll.Find(opCtx, filter, findOpts)
		if err != nil {
			return fmt.Errorf("执行 find 失败: %w", err)
		}
		defer cursor.Close(opCtx)

		if err := cursor.All(opCtx, &out.Documents); err != nil {
			return fmt.Errorf("读取 find 结果失败: %w", err)
		}
		if out.Documents == nil {
			out.Documents = []bson.M{}
		}
		return nil
	})
	if err != nil {
		return nil, findOutput{}, err
	}

	if err := ensureResponseBytesWithinLimit(out.Documents, responseBytesLimit); err != nil {
		return nil, findOutput{}, err
	}
	return nil, out, nil
}

func handleListDatabases(ctx context.Context, _ *mcp.CallToolRequest, _ map[string]any) (*mcp.CallToolResult, listDatabasesOutput, error) {
	out := listDatabasesOutput{Databases: []databaseInfo{}}
	err := withMongoClient(ctx, func(opCtx context.Context, client *mongo.Client) error {
		result, err := client.ListDatabases(opCtx, bson.D{})
		if err != nil {
			return fmt.Errorf("列出数据库失败: %w", err)
		}
		out.Databases = make([]databaseInfo, 0, len(result.Databases))
		for _, db := range result.Databases {
			out.Databases = append(out.Databases, databaseInfo{
				Name: db.Name,
				Size: float64(db.SizeOnDisk),
			})
		}
		out.TotalCount = len(out.Databases)
		return nil
	})
	return nil, out, err
}

func handleListCollections(ctx context.Context, _ *mcp.CallToolRequest, in listCollectionsInput) (*mcp.CallToolResult, listCollectionsOutput, error) {
	if in.Database == "" {
		return nil, listCollectionsOutput{}, fmt.Errorf("database 为必填参数")
	}

	out := listCollectionsOutput{Collections: []collectionInfo{}}
	err := withMongoClient(ctx, func(opCtx context.Context, client *mongo.Client) error {
		names, err := client.Database(in.Database).ListCollectionNames(opCtx, bson.D{})
		if err != nil {
			return fmt.Errorf("列出集合失败: %w", err)
		}
		out.Collections = make([]collectionInfo, 0, len(names))
		for _, name := range names {
			out.Collections = append(out.Collections, collectionInfo{Name: name})
		}
		out.TotalCount = len(out.Collections)
		return nil
	})
	return nil, out, err
}

func applyFindDefaults(limit *int64, responseBytesLimit *int64) (int64, int64) {
	finalLimit := defaultFindLimit
	if limit != nil {
		finalLimit = *limit
	}

	finalBytesLimit := defaultResponseBytesLimit
	if responseBytesLimit != nil {
		finalBytesLimit = *responseBytesLimit
	}
	return finalLimit, finalBytesLimit
}

func ensureResponseBytesWithinLimit(documents []bson.M, responseBytesLimit int64) error {
	data, err := json.Marshal(documents)
	if err != nil {
		return fmt.Errorf("序列化 find 结果失败: %w", err)
	}
	if int64(len(data)) > responseBytesLimit {
		return fmt.Errorf("find 结果大小 %d 字节超出 responseBytesLimit=%d，请缩小查询范围或降低返回字段", len(data), responseBytesLimit)
	}
	return nil
}
