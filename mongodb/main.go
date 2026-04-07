package main

import (
	"context"
	"log"

	"mongodb-mcp-server/pkg"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mongodb-mcp-server",
		Version: "v1.0.0",
	}, nil)

	pkg.RegisterMongoTools(server)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
