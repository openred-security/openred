package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"openred/openred/launcher"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: openred [command] [plugin]")
		os.Exit(1)
	}

	pluginsDir := filepath.Join(".", "plugins")

	// Load all plugins from the plugins directory
	plugins, err := launcher.LoadPlugins(pluginsDir)
	if err != nil {
		log.Fatalf("Error loading plugins: %v", err)
	}

	command := os.Args[1]

	switch command {
	case "plugin", "plugins":
		if len(os.Args) < 3 || os.Args[2] != "list" {
			fmt.Println("Usage: openred plugin list")
			os.Exit(1)
		}

		// List all loaded plugins
		for _, plugin := range plugins {
			fmt.Printf("Plugin: %s, Binary: %s\n", plugin.Name, plugin.Binary)
		}

	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Usage: openred run [plugin_name]")
			os.Exit(1)
		}

		pluginName := os.Args[2]

		// Find the plugin by name
		var pluginToRun *launcher.Plugin
		for _, plugin := range plugins {
			if plugin.Name == pluginName {
				pluginToRun = &plugin
				break
			}
		}

		if pluginToRun == nil {
			fmt.Printf("Plugin '%s' not found.\n", pluginName)
			os.Exit(1)
		}

		// Run the plugin
		output, err := launcher.RunPlugin(*pluginToRun)
		if err != nil {
			log.Fatalf("Error running plugin: %v", err)
		}

		fmt.Printf("Output from plugin '%s':\n%s", pluginToRun.Name, output)

	default:
		fmt.Println("Unknown command. Usage: openred [plugin|run] [plugin_name]")
		os.Exit(1)
	}
}
