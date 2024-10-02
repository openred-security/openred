package launcher

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"gopkg.in/yaml.v2"
)

type plugin struct {
	Name      string
	Path      string
	Arguments []string
}

var plugins = []plugin{}

type Config struct {
	Binary string   `yaml:"binary"`
	Args   []string `yaml:"args"`
}

func Init(ctx context.Context) {
	// Añadir el plugin "last_logs"
	Add(ctx, "last_logs", "./plugins/last_logs/config.yml")
}

func Add(ctx context.Context, name string, configPath string) {
	// Leer configuración desde el archivo YAML
	config := Config{}
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error leyendo config.yml:", err)
		return
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Println("Error en la configuración YAML:", err)
		return
	}

	plugin := plugin{
		Name:      name,
		Path:      config.Binary,
		Arguments: config.Args,
	}
	plugins = append(plugins, plugin)
}

func Plugins() []plugin {
	return plugins
}

func Run(ctx context.Context, plugin plugin) {
	cmd := exec.CommandContext(ctx, plugin.Path, plugin.Arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error ejecutando el plugin:", err)
	}
}
