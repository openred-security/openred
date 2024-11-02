package main

import (
	"fmt"
	"log"
	"os"

	"openred/openred/plugin"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: openred <command> [plugin]")
	}

	manager := plugin.NewManager()
	catalogDir := "./catalog" // Adjust as necessary

	// Load plugins from the catalog directory
	if err := manager.LoadPlugins(catalogDir); err != nil {
		log.Fatalf("Error loading plugins: %v", err)
	}

	command := os.Args[1]

	switch command {
	case "list":
		// List all loaded plugins
		manager.ListPlugins()

	case "download":
		// Download a specific plugin
		if len(os.Args) < 3 {
			log.Fatal("Usage: openred download <plugin>")
		}
		pluginID := os.Args[2]
		err := manager.DownloadAndDecompressPlugin(pluginID, catalogDir)
		if err != nil {
			log.Fatalf("Failed to download plugin %s: %v", pluginID, err)
		}
		fmt.Printf("Plugin %s downloaded successfully\n", pluginID)

	case "run":
		// Run a specific plugin, downloading if necessary
		if len(os.Args) < 3 {
			log.Fatal("Usage: openred run <plugin>")
		}
		pluginID := os.Args[2]
		pluginInstance, exists := manager.GetPlugin(pluginID)
		if !exists {
			log.Fatalf("Plugin %s not found", pluginID)
		}

		if err := pluginInstance.RunPlugin(); err != nil {
			log.Fatalf("Failed to run plugin %s: %v", pluginID, err)
		}
		fmt.Printf("Plugin %s executed successfully\n", pluginID)

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
