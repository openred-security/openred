// Package plugin provides functionality for managing multiple plugins
// defined in a catalog directory with each plugin configured by a YAML file.
package plugin

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Manager is responsible for managing multiple plugins loaded from a catalog directory.
type Manager struct {
	plugins map[string]*Plugin // Stores plugins loaded from the catalog
	mu      sync.RWMutex
}

// NewManager creates a new Manager instance.
func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]*Plugin),
	}
}

// LoadPlugins reads all plugin configurations from the specified catalog directory.
// Each plugin is expected to have its configuration defined in a config.yml file.
//
// Parameters:
// - catalogDir: string - the path to the catalog directory containing plugin subdirectories.
//
// Returns:
// - error - an error if there is a problem reading any config.yml file in the catalog directory.
func (m *Manager) LoadPlugins(catalogDir string) error {
	files, err := ioutil.ReadDir(catalogDir)
	if err != nil {
		return fmt.Errorf("failed to read catalog directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			pluginConfigPath := filepath.Join(catalogDir, file.Name(), "config.yml")
			if _, err := os.Stat(pluginConfigPath); err == nil {
				plugin, err := LoadPlugin(pluginConfigPath)
				if err != nil {
					log.Printf("failed to load plugin from %s: %v", pluginConfigPath, err)
					continue
				}
				m.plugins[plugin.ID] = plugin
			} else {
				log.Printf("config.yml not found in directory: %s", file.Name())
			}
		}
	}

	return nil
}

// GetPlugin retrieves a plugin by its ID.
func (m *Manager) GetPlugin(pluginID string) (*Plugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	plugin, exists := m.plugins[pluginID]
	return plugin, exists
}

// ListPlugins prints the names and versions of all loaded plugins in a formatted way.
func (m *Manager) ListPlugins() {
	fmt.Println("Loaded Plugins:")
	for _, plugin := range m.plugins {
		fmt.Printf("%s - %s v%s\n", plugin.ID, plugin.Name, plugin.Version)
	}
}

// DownloadAndDecompressPlugin downloads and decompresses the plugin's archive or binary file
// to the plugin's subdirectory in the catalog directory.
//
// Parameters:
// - pluginID: string - the ID of the plugin to download and decompress.
// - catalogDir: string - the path to the catalog directory where plugins are stored.
//
// Returns:
// - error - an error if the plugin is not found, or if there is an issue during download or decompression.
func (m *Manager) DownloadAndDecompressPlugin(pluginID, catalogDir string) error {
	plugin, exists := m.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin with ID %s not found", pluginID)
	}

	pluginDir := filepath.Join(catalogDir, pluginID)
	archivePath, err := plugin.DownloadPlugin(pluginDir)
	if err != nil {
		return fmt.Errorf("failed to download plugin %s: %w", pluginID, err)
	}

	// Decompress if the downloaded file is an archive
	if filepath.Ext(archivePath) != "" {
		err = plugin.Decompress(archivePath, pluginDir)
		if err != nil {
			return fmt.Errorf("failed to decompress plugin %s: %w", pluginID, err)
		}
	}

	fmt.Printf("Plugin %s downloaded and decompressed successfully\n", plugin.Name)
	return nil
}
