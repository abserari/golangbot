package main

import (
	"log"

	"github.com/abserari/golangbot/pkg/config"
	"github.com/abserari/golangbot/pkg/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	config.Init()
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{
		HostPort: "49.235.242.124:7233",
	})
	if err != nil {
		log.Println(err)
	}
	defer c.Close()
	// This worker hosts both Workflow and Activity functions

	w := worker.New(c, "test", worker.Options{})
	w.RegisterWorkflow(workflows.GetAndSendMessage)
	w.RegisterActivity(workflows.SendMessage)
	workflows.RssWorker()
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
	}

}
