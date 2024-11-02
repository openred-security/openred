// Package plugin provides functionality for loading, downloading, and decompressing plugins defined by a YAML configuration file.
package plugin

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"openred/openred/process"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archiver/v3"
	"gopkg.in/yaml.v2"
)

// Plugin represents the structure of a plugin as defined in a YAML configuration file.
// It includes metadata about the plugin and download URLs for various OS/architecture combinations.
type Plugin struct {
	ID          string                       `yaml:"id"`           // Unique identifier for the plugin
	Name        string                       `yaml:"name"`         // Name of the plugin
	Description string                       `yaml:"description"`  // Description of the plugin's functionality
	Site        string                       `yaml:"site"`         // URL to the plugin's main site
	Docs        string                       `yaml:"docs"`         // URL to the plugin's documentation
	Target      string                       `yaml:"target"`       // Target environment for the plugin (e.g., "engine", "console", "agent")
	License     string                       `yaml:"license"`      // License under which the plugin is distributed
	Category    string                       `yaml:"category"`     // Category for organizing plugins
	Version     string                       `yaml:"version"`      // Version of the plugin
	CatalogURL  string                       `yaml:"catalog_url"`  // URL to the plugin's catalog page
	DownloadURL map[string]map[string]string `yaml:"download_url"` // URLs for downloading the plugin per OS/architecture
	BinaryName  string                       `yaml:"binary_name"`  // Name of the binary to execute
	Env         []string                     `yaml:"env"`          // Environment variables for the plugin
	Args        []string                     `yaml:"args"`         // Arguments to pass to the plugin binary
	BinaryPath  string                       `yaml:"-"`            // Path where the binary will be located (not in YAML)
}

// LoadPlugin reads a YAML file from the specified path and unmarshals it into a Plugin struct.
func LoadPlugin(yamlFile string) (*Plugin, error) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var plugin Plugin
	if err := yaml.Unmarshal(data, &plugin); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Set the BinaryPath based on the directory of the YAML file
	pluginDir := filepath.Dir(yamlFile)
	plugin.BinaryPath = filepath.Join(pluginDir, "run", plugin.BinaryName)

	absBinaryPath, err := filepath.Abs(plugin.BinaryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to determine absolute path for binary: %w", err)
	}
	plugin.BinaryPath = absBinaryPath

	return &plugin, nil
}

// GetDownloadURL determines the appropriate download URL based on the current operating system and architecture.
func (p *Plugin) GetDownloadURL() (string, error) {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	if osName == "darwin" {
		osName = "macos"
	}
	urls, ok := p.DownloadURL[osName]
	if !ok {
		return "", fmt.Errorf("no download URL available for OS: %s", osName)
	}

	downloadURL, ok := urls[arch]
	if !ok {
		return "", fmt.Errorf("no download URL available for architecture: %s", arch)
	}

	return downloadURL, nil
}

// DownloadPlugin downloads the plugin binary or archive based on its download URL.
func (p *Plugin) DownloadPlugin(destinationDir string) (string, error) {
	url, err := p.GetDownloadURL()
	if err != nil {
		return "", err
	}

	// Ensure the 'run' directory exists
	runDir := filepath.Join(destinationDir, "run")
	if err := os.MkdirAll(runDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create run directory: %w", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download plugin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status code: %d", resp.StatusCode)
	}

	fileName := filepath.Base(url)
	destinationPath := filepath.Join(runDir, fileName) // Save in the run directory
	file, err := os.Create(destinationPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Printf("Plugin downloaded to %s\n", destinationPath)
	return destinationPath, nil
}

// Decompress extracts the contents of the downloaded plugin file if it's in a supported archive format.
func (p *Plugin) Decompress(archivePath, destinationDir string) error {
	ext := strings.ToLower(filepath.Ext(archivePath))
	var decompressor archiver.Unarchiver

	switch ext {
	case ".zip":
		decompressor = archiver.NewZip()
	case ".tar":
		decompressor = archiver.NewTar()
	case ".gz", ".tgz":
		decompressor = archiver.NewTarGz()
	default:
		return fmt.Errorf("unsupported archive format: %s", ext)
	}

	err := decompressor.Unarchive(archivePath, destinationDir)
	if err != nil {
		return fmt.Errorf("failed to decompress plugin: %w", err)
	}

	fmt.Printf("Plugin decompressed to %s\n", destinationDir)
	return nil
}

// RunPlugin starts the plugin binary as defined in its configuration.
func (p *Plugin) RunPlugin() error {
	// Check if the plugin binary is already downloaded
	if _, err := os.Stat(p.BinaryPath); os.IsNotExist(err) {
		log.Printf("Binary for plugin %s not found, downloading...", p.Name)
		if _, err := p.DownloadPlugin(filepath.Dir(p.BinaryPath)); err != nil {
			return fmt.Errorf("failed to download plugin: %w", err)
		}
	}

	// Start the plugin process
	proc, err := process.Start(p.BinaryPath,
		process.WithContext(context.Background()),
		process.WithArgs(p.Args),
		process.WithEnv(p.Env),
	)
	if err != nil {
		return fmt.Errorf("failed to start plugin %s: %w", p.Name, err)
	}
	defer proc.Process.Kill()

	// Capture plugin output
	go func() {
		io.Copy(os.Stdout, proc.Stdout)
	}()
	go func() {
		io.Copy(os.Stderr, proc.Stderr)
	}()

	// Wait for the process to complete
	state := <-proc.Wait()
	log.Printf("Plugin %s exited with status %v", p.Name, state)
	return nil
}
