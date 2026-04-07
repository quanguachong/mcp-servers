package main

import (
	"context"
	"log"

	mongodbmcpserver "github.com/quanguachong/mcp-servers/pkg/mongodb-mcp-server"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mongodb-mcp-server",
		Version: "v1.0.0",
	}, nil)

	mongodbmcpserver.RegisterMongoTools(server)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
