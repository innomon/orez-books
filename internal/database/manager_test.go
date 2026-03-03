package database

import (
	"os"
	"testing"
)

func TestManager_CreateAndConnect(t *testing.T) {
	m := NewManager()
	dbPath := "test_books.db"
	defer os.Remove(dbPath)

	// Test Create
	err := m.CreateNewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("Database file was not created")
	}

	// Test Migrate
	err = m.Migrate()
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Test MigrateFromSchemas
	err = m.MigrateFromSchemas("-")
	if err != nil {
		t.Fatalf("Failed to migrate from schemas: %v", err)
	}

	// Verify that a table from schemas exists (e.g., Account)
	var count int64
	err = m.db.Raw("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='Account'").Scan(&count).Error
	if err != nil {
		t.Fatalf("Failed to check for Account table: %v", err)
	}
	if count == 0 {
		t.Errorf("Account table was not created")
	}

	// Test Insert into SingleValue (verify schema exists)
	sv := SingleValue{
		Name:      "test_setting",
		Parent:    "SystemSettings",
		Fieldname: "version",
		Value:     "0.36.0",
	}
	result := m.db.Create(&sv)
	if result.Error != nil {
		t.Fatalf("Failed to insert into SingleValue: %v", result.Error)
	}

	// Test Connect to existing
	err = m.Close()
	if err != nil {
		t.Fatalf("Failed to close database: %v", err)
	}

	err = m.Connect(dbPath)
	if err != nil {
		t.Fatalf("Failed to connect to existing database: %v", err)
	}

	var foundSv SingleValue
	m.db.First(&foundSv, "name = ?", "test_setting")
	if foundSv.Value != "0.36.0" {
		t.Errorf("Expected value 0.36.0, got %s", foundSv.Value)
	}

	m.Close()
}
