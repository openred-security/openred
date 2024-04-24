package main

import (
	"context"
	"fmt"
	"openred/openred-agent/api"
	"openred/openred-agent/launcher"
	"time"
)

func main() {

	// Global context
	ctx := context.Background()

	launcher.Init(ctx)
	//launcher.RunAll(ctx)

	api.Start()

	for {
		fmt.Println("Still here..")
		time.Sleep(time.Second * 10)
	}

}
