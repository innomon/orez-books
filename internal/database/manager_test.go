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
