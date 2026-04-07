package main

import (
	"log"
	"os"

	mcpserver "github.com/quanguachong/mcp-servers/http/internal/mcp"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"http-requests",
		"0.1.0",
		server.WithToolCapabilities(true),
	)
	mcpserver.RegisterTools(s)

	if err := server.ServeStdio(s); err != nil {
		log.New(os.Stderr, "", log.LstdFlags).Fatal(err)
	}
}
