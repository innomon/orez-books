package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Manager struct {
	db     *gorm.DB
	dbPath string
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Connect(dbPath string) error {
	// Close existing connection if any
	if m.db != nil {
		sqlDB, err := m.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	// modernc.org/sqlite requires a specific dialector setup for GORM if not using default cgo one
	// But gorm.io/driver/sqlite by default uses github.com/mattn/go-sqlite3 (Cgo)
	// To use modernc.org/sqlite (Cgo-free), we need to use a different driver or configure GORM.
	// For now, let's stick to the principle of Cgo-free if possible.
	
	// Actually, gorm.io/driver/sqlite can be used with other drivers by providing the opener.
	// However, many users use a dedicated modernc driver for GORM.
	// Let's use the standard "sqlite" driver for now for simplicity, but we can switch if needed.
	
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	m.db = db
	m.dbPath = dbPath
	return nil
}

func (m *Manager) CreateNewDatabase(dbPath string) error {
	if _, err := os.Stat(dbPath); err == nil {
		err := os.Remove(dbPath)
		if err != nil {
			return fmt.Errorf("failed to remove existing database: %w", err)
		}
	}

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	return m.Connect(dbPath)
}

func (m *Manager) Close() error {
	if m.db == nil {
		return nil
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Close()
	m.db = nil
	return err
}

// SingleValue represents the settings table in Frappe Books
type SingleValue struct {
	Name       string `gorm:"primaryKey"`
	Parent     string `gorm:"index"`
	Fieldname  string `gorm:"index"`
	Value      string
	Created    string
	Modified   string
	CreatedBy  string
	ModifiedBy string
}

func (SingleValue) TableName() string {
	return "SingleValue"
}

func (m *Manager) Migrate() error {
	if m.db == nil {
		return fmt.Errorf("database not connected")
	}

	// Basic migration for SingleValue table which is core to Frappe Books
	return m.db.AutoMigrate(&SingleValue{})
}

func (m *Manager) MigrateFromSchemas(countryCode string) error {
	loader := NewSchemaLoader()
	if err := loader.LoadSchemas(countryCode); err != nil {
		return err
	}

	for _, schema := range loader.SchemaMap {
		if schema.IsSingle {
			continue
		}

		// Create table dynamically using raw SQL if GORM AutoMigrate is too restrictive
		// Or use m.db.Table(schema.Name).AutoMigrate(...) if we can define a dynamic struct
		// For now, let's use a simpler approach: create table if not exists
		if err := m.createTableFromSchema(schema); err != nil {
			return fmt.Errorf("failed to create table %s: %w", schema.Name, err)
		}
	}

	return nil
}

func (m *Manager) createTableFromSchema(schema *Schema) error {
	var columns []string
	for _, field := range schema.Fields {
		if field.Computed {
			continue
		}

		columnType := m.mapFieldTypeToSQLite(field.Fieldtype)
		columnDef := fmt.Sprintf("\"%s\" %s", field.Fieldname, columnType)
		if field.Fieldname == "name" {
			columnDef += " PRIMARY KEY"
		}
		if field.Required && field.Fieldname != "name" {
			columnDef += " NOT NULL"
		}
		if field.Default != nil {
			// Basic default value handling
			columnDef += fmt.Sprintf(" DEFAULT '%v'", field.Default)
		}
		columns = append(columns, columnDef)
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\" (%s)", schema.Name, strings.Join(columns, ", "))
	return m.db.Exec(query).Error
}

func (m *Manager) mapFieldTypeToSQLite(ft FieldType) string {
	switch ft {
	case FieldTypeInt:
		return "INTEGER"
	case FieldTypeFloat, FieldTypeCurrency:
		return "REAL"
	case FieldTypeCheck:
		return "INTEGER" // 0 or 1
	case FieldTypeDate, FieldTypeDatetime:
		return "TEXT" // SQLite doesn't have native Date
	default:
		return "TEXT"
	}
}

func (m *Manager) GetDB() *gorm.DB {
	return m.db
}
