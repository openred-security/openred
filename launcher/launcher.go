package launcher

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"gopkg.in/yaml.v3"
)

type PluginConfig struct {
	Name   string `yaml:"name"`
	Binary string `yaml:"binary"`
}

type Plugin struct {
	Name   string
	Binary string
}

func LoadPlugins(pluginDir string) ([]Plugin, error) {
	var plugins []Plugin

	err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process directories at the first level (the plugin directories)
		if info.IsDir() && path != pluginDir {
			configPath := filepath.Join(path, "config.yml")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				return nil // Skip if config.yml doesn't exist
			}

			// Read the config.yml file
			configFile, err := os.ReadFile(configPath)
			if err != nil {
				return err
			}

			var pluginConfig PluginConfig
			err = yaml.Unmarshal(configFile, &pluginConfig)
			if err != nil {
				return err
			}

			plugin := Plugin{
				Name:   pluginConfig.Name,
				Binary: filepath.Join(path, pluginConfig.Binary),
			}

			plugins = append(plugins, plugin)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return plugins, nil
}

func RunPlugin(plugin Plugin) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, plugin.Binary)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the plugin binary
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run plugin: %s, error: %v, stderr: %s", plugin.Name, err, stderr.String())
	}

	return stdout.String(), nil
}
