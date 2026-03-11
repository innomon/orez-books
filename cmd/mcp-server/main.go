package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"orez-books/pkg/database"
)

func main() {
	dbPath := flag.String("db", "", "Path to the Orez Books SQLite database")
	transport := flag.String("transport", "stdio", "Transport to use (stdio or sse)")
	addr := flag.String("addr", "localhost:8080", "Address for SSE transport")
	flag.Parse()

	if *dbPath == "" {
		fmt.Fprintf(os.Stderr, "Error: -db flag is required\n")
		os.Exit(1)
	}

	// Initialize database manager
	dbManager := database.NewManager()
	if err := dbManager.Connect(*dbPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n")
		os.Exit(1)
	}
	defer dbManager.Close()

	// Initialize MCP server
	mcpServer := NewMCPServer(dbManager)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run MCP server
	if err := mcpServer.Run(ctx, *transport, *addr); err != nil {
		fmt.Fprintf(os.Stderr, "Error running MCP server: %v\n")
		os.Exit(1)
	}
}
