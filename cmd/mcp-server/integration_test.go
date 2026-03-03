package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"orez-books/internal/database"
)

func TestMCP_Integration(t *testing.T) {
	// 1. Setup temporary database
	dbPath := filepath.Join(t.TempDir(), "test_books.db")
	dbManager := database.NewManager()
	if err := dbManager.CreateNewDatabase(dbPath); err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
	if err := dbManager.MigrateFromSchemas("-"); err != nil {
		t.Fatalf("Failed to migrate schemas: %v", err)
	}
	
	// Add a test account
	now := time.Now().Format("2006-01-02 15:04:05")
	dbManager.GetDB().Table("Account").Create(map[string]interface{}{
		"name":        "Test Asset Account",
		"rootType":    "Asset",
		"accountType": "Cash",
		"isGroup":     0,
		"submitted":   0,
		"cancelled":   0,
		"lft":         1,
		"rgt":         2,
		"created":     now,
		"modified":    now,
		"createdBy":   "test",
		"modifiedBy":  "test",
	})
	dbManager.Close()

	// 2. Build the MCP server binary
	binPath := filepath.Join(t.TempDir(), "orez-books-mcp")
	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Build failed: %v\n%s", err, string(out))
	}

	// 3. Start the MCP server process
	cmd := exec.Command(binPath, "-db", dbPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer cmd.Process.Kill()

	// Use a scanner to read the JSON-RPC output line by line
	scanner := bufio.NewScanner(stdout)

	// 4. Helper to send JSON-RPC
	sendRequest := func(method string, params interface{}) {
		req := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  method,
			"params":  params,
		}
		data, _ := json.Marshal(req)
		stdin.Write(data)
		stdin.Write([]byte("\n"))
	}

	// 5. Test Initialize
	sendRequest("initialize", map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo":      map[string]interface{}{"name": "test-client", "version": "1.0.0"},
	})

	if !scanner.Scan() {
		t.Fatalf("Failed to read initialize response")
	}
	line := scanner.Bytes()
	if !bytes.Contains(line, []byte("orez-books-mcp")) {
		t.Errorf("Unexpected initialize response: %s", string(line))
	}

	// 6. Test List Tools
	sendRequest("tools/list", map[string]interface{}{})
	if !scanner.Scan() {
		t.Fatalf("Failed to read list tools response")
	}
	line = scanner.Bytes()
	if !bytes.Contains(line, []byte("create_invoice")) || !bytes.Contains(line, []byte("list_accounts")) {
		t.Errorf("Missing expected tools in list: %s", string(line))
	}

	// 7. Test Call Tool (list_accounts)
	sendRequest("tools/call", map[string]interface{}{
		"name": "list_accounts",
		"arguments": map[string]interface{}{
			"root_type": "Asset",
		},
	})
	
	if !scanner.Scan() {
		t.Fatalf("Failed to read call tool response")
	}
	line = scanner.Bytes()
	if !bytes.Contains(line, []byte("Test Asset Account")) {
		t.Errorf("Tool response did not contain test account: %s", string(line))
	}
}
