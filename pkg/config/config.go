package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	mu       sync.RWMutex
	filePath string
	data     map[string]interface{}
}

func NewConfig(appName string) (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(configDir, appName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	filePath := filepath.Join(appDir, "config.json")
	c := &Config{
		filePath: filePath,
		data:     make(map[string]interface{}),
	}

	if err := c.load(); err != nil {
		// If file doesn't exist, it's fine, we start with empty data
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return c, nil
}

func (c *Config) load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(c.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&c.data)
}

func (c *Config) save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	file, err := os.Create(c.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c.data)
}

func (c *Config) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

func (c *Config) Set(key string, value interface{}) error {
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
	return c.save()
}

func (c *Config) Delete(key string) error {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
	return c.save()
}

func (c *Config) GetAll() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Return a copy to avoid external mutation
	copy := make(map[string]interface{})
	for k, v := range c.data {
		copy[k] = v
	}
	return copy
}
