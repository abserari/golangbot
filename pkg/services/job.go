package services

import (
	"github.com/abserari/golangbot/pkg/services/activities"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func registerWorker() {
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	defer c.Close()
	// This worker hosts both Workflow and Activity functions

	w := worker.New(c, "", worker.Options{})
	w.RegisterActivity(activities.SendMessage)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
	}

}
