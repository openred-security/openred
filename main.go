package main

import (
	"context"
	"openred/openred-agent/launcher"
	"openred/openred-agent/api"
	"openred/openred-agent/sender"
	"fmt"
	"time"
)

func main() {

	// Global context
	ctx := context.Background()
	

	sender.Sender()
	launcher.Init(ctx)
	launcher.RunAll(ctx)
	
	api.Start()

	for {
		fmt.Println("Still here..")
		time.Sleep(time.Second * 10)
	}

}
