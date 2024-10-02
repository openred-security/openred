package main

import (
	"context"
	"flag"
	"fmt"
	"openred/openred/launcher"
	"os"
)

func main() {
	// Crear comandos y subcomandos
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	// Revisar argumentos
	if len(os.Args) < 2 {
		fmt.Println("Se necesita un subcomando como 'run'")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		runCmd.Parse(os.Args[2:])
		if runCmd.NArg() < 1 {
			fmt.Println("Se necesita el nombre del plugin para ejecutar. Ejemplo: openred run last_logs")
			os.Exit(1)
		}
		pluginName := runCmd.Arg(0)
		executePlugin(pluginName)
	default:
		fmt.Println("Comando no reconocido:", os.Args[1])
		os.Exit(1)
	}
}

// Función para ejecutar un plugin
func executePlugin(pluginName string) {
	// Crear un contexto global
	ctx := context.Background()

	// Inicializar los plugins
	launcher.Init(ctx)

	// Buscar y ejecutar el plugin específico
	for _, plugin := range launcher.Plugins() {
		if plugin.Name == pluginName {
			fmt.Println("Ejecutando plugin:", pluginName)
			launcher.Run(ctx, plugin)
			return
		}
	}

	fmt.Println("Plugin no encontrado:", pluginName)
}
