package pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	opTimeout        = 30 * time.Second
	disconnectTimout = 5 * time.Second
	envMongoURI      = "MONGODB_URI"
)

func withMongoClient(ctx context.Context, fn func(context.Context, *mongo.Client) error) error {
	uri := os.Getenv(envMongoURI)
	if uri == "" {
		return fmt.Errorf("环境变量 %s 未设置", envMongoURI)
	}

	connectCtx, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()

	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("连接 MongoDB 失败: %w", err)
	}
	defer func() {
		dctx, dcancel := context.WithTimeout(context.Background(), disconnectTimout)
		defer dcancel()
		_ = client.Disconnect(dctx)
	}()

	if err := client.Ping(connectCtx, nil); err != nil {
		return fmt.Errorf("Ping 失败: %w", err)
	}

	return fn(connectCtx, client)
}
