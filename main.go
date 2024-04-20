package main

import (
	"context"
	"openred/openred-agent/launcher"
	"fmt"
	"time"
)

func main() {

	// Global context
	ctx := context.Background()

	launcher.Init(ctx)
	launcher.RunAll(ctx)

	for {
		fmt.Println("Still here..")
		time.Sleep(time.Second * 10)
	}

}
