package launcher

import (
	"context"
	"openred/openred-agent/process"
	"fmt"
	"os"
	"os/exec"
)

type plugin struct {
	name      string
	path      string
	arguments []string
}

var plugins = []plugin{}

func Init(ctx context.Context) {
	current_path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	Add(ctx, "sample", current_path+"/plugins/sample/sample", []string{" "})
}

func Add(ctx context.Context, name string, path string, arguments []string) {
	plugin := plugin{
		name:      name,
		path:      path,
		arguments: arguments,
	}
	plugins = append(plugins, plugin)
}

func Run(ctx context.Context, plugin plugin) {

	proc, err := process.Start(
		plugin.path,
		process.WithContext(ctx),
		process.WithArgs(plugin.arguments),
		process.WithCmdOptions(func(c *exec.Cmd) error {
			//c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			return nil
		}))
	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
	}

	resChan := make(chan *os.ProcessState)
	go func() {
		procState, _ := proc.Process.Wait()
		resChan <- procState
	}()

}

func RunAll(ctx context.Context) {

	for _, plugin := range plugins {
		Run(ctx, plugin)
	}

}
