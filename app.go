package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"orez-books/pkg/config"
	"orez-books/pkg/database"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	dbManager *database.Manager
	config    *config.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	cfg, _ := config.NewConfig("OrezBooks")
	return &App{
		dbManager: database.NewManager(),
		config:    cfg,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// CheckDbAccess checks if a file is readable and writable
func (a *App) CheckDbAccess(filePath string) bool {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		if os.IsPermission(err) {
			return false
		}
		return false
	}
	file.Close()
	return true
}

// GetDbDefaultPath returns the default path for a new database
func (a *App) GetDbDefaultPath(companyName string) string {
	// In Wails, we don't have app.getPath('documents') directly like Electron.
	// We can use os.UserConfigDir or os.UserHomeDir.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "." // Fallback to current directory
	}

	// Follow Frappe Books logic: use Documents or UserData. We'll use Documents if available.
	docsDir := filepath.Join(homeDir, "Documents")
	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		configDir, err := os.UserConfigDir()
		if err == nil {
			docsDir = configDir
		} else {
			docsDir = homeDir
		}
	}

	dbsPath := filepath.Join(docsDir, "Frappe Books")
	backupPath := filepath.Join(dbsPath, "backups")
	
	_ = os.MkdirAll(backupPath, 0755)

	dbFilePath := filepath.Join(dbsPath, fmt.Sprintf("%s.books.db", companyName))

	if _, err := os.Stat(dbFilePath); err == nil {
		// File exists, ask user via dialog
		res, err := wailsRuntime.MessageDialog(a.ctx, wailsRuntime.MessageDialogOptions{
			Type:    wailsRuntime.QuestionDialog,
			Title:   "File Exists",
			Message: "Filename already exists. Do you want to overwrite the existing file or create a new one?",
			Buttons: []string{"Overwrite", "New"},
		})

		if err == nil && res == "New" {
			timestamp := "new" // In a real scenario, format current time
			dbFilePath = filepath.Join(dbsPath, fmt.Sprintf("%s_%s.books.db", companyName, timestamp))
			
			wailsRuntime.MessageDialog(a.ctx, wailsRuntime.MessageDialogOptions{
				Type:    wailsRuntime.InfoDialog,
				Title:   "New File",
				Message: fmt.Sprintf("New file: %s", filepath.Base(dbFilePath)),
			})
		}
	}

	return dbFilePath
}

// GetOpenFilePath opens a file selection dialog
func (a *App) GetOpenFilePath(title string) string {
	path, _ := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: title,
	})
	return path
}

// GetSaveFilePath opens a file save dialog
func (a *App) GetSaveFilePath(title string, defaultFilename string) string {
	path, _ := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultFilename,
	})
	return path
}

// ShowError displays an error dialog
func (a *App) ShowError(title string, content string) {
	wailsRuntime.MessageDialog(a.ctx, wailsRuntime.MessageDialogOptions{
		Type:    wailsRuntime.ErrorDialog,
		Title:   title,
		Message: content,
	})
}

// SaveData saves string data to a file
func (a *App) SaveData(data string, savePath string) error {
	return os.WriteFile(savePath, []byte(data), 0644)
}

// DeleteFile deletes a file
func (a *App) DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// GetEnv returns environment information
func (a *App) GetEnv() map[string]interface{} {
	return map[string]interface{}{
		"isDevelopment": true, // We should detect this based on build tags
		"platform":      runtime.GOOS,
		"version":       "0.36.0",
	}
}

// OpenExternal opens a link in the default browser
func (a *App) OpenExternal(link string) {
	wailsRuntime.BrowserOpenURL(a.ctx, link)
}

// MinimizeMainWindow minimizes the window
func (a *App) MinimizeMainWindow() {
	wailsRuntime.WindowMinimise(a.ctx)
}

// MaximizeMainWindow maximizes or unmaximizes the window
func (a *App) MaximizeMainWindow() {
	// Wails currently has WindowToggleMaximise
	wailsRuntime.WindowToggleMaximise(a.ctx)
}

// CloseMainWindow closes the application
func (a *App) CloseMainWindow() {
	wailsRuntime.Quit(a.ctx)
}

// DbCreate creates a new database
func (a *App) DbCreate(dbPath string, countryCode string) error {
	// CountryCode logic will be handled later, for now just create the DB
	err := a.dbManager.CreateNewDatabase(dbPath)
	if err != nil {
		return err
	}
	return a.dbManager.MigrateFromSchemas(countryCode)
}

// DbConnect connects to an existing database
func (a *App) DbConnect(dbPath string) error {
	err := a.dbManager.Connect(dbPath)
	if err != nil {
		return err
	}
	// We might want to pass countryCode here too if we store it in config
	return a.dbManager.MigrateFromSchemas("-")
}

// DbClose closes the current database connection
func (a *App) DbClose() error {
	return a.dbManager.Close()
}

// ConfigGet returns a config value
func (a *App) ConfigGet(key string) interface{} {
	return a.config.Get(key)
}

// ConfigSet sets a config value
func (a *App) ConfigSet(key string, value interface{}) error {
	return a.config.Set(key, value)
}

// ConfigDelete deletes a config value
func (a *App) ConfigDelete(key string) error {
	return a.config.Delete(key)
}

// ConfigGetAll returns all config values
func (a *App) ConfigGetAll() map[string]interface{} {
	return a.config.GetAll()
}

