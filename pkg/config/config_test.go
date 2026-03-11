package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	appName := "OrezBooksTest"
	c, err := NewConfig(appName)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	defer os.RemoveAll(filepath.Dir(c.filePath))

	// Test Set and Get
	err = c.Set("theme", "dark")
	if err != nil {
		t.Fatalf("Failed to set config: %v", err)
	}

	// Ensure file was created
	if _, err := os.Stat(c.filePath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", c.filePath)
	}

	val := c.Get("theme")
	if val != "dark" {
		t.Errorf("Expected dark, got %v", val)
	}

	// Test Persistence
	c2, err := NewConfig(appName)
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	val2 := c2.Get("theme")
	if val2 != "dark" {
		t.Errorf("Expected persistent value dark, got %v", val2)
	}

	// Test Delete
	err = c.Delete("theme")
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}

	if c.Get("theme") != nil {
		t.Errorf("Expected nil after delete, got %v", c.Get("theme"))
	}
}
