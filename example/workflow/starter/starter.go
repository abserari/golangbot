package main

import (
	"context"
	"log"

	"github.com/abserari/golangbot/pkg/workflows"
	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.NewClient(client.Options{
		HostPort: "49.235.242.124:7233",
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "transfer-money-workflow",
		TaskQueue: "test",
	}

	var hello = "hello"

	we, err := c.ExecuteWorkflow(context.Background(), options, workflows.GetAndSendMessage, hello)
	if err != nil {
		log.Fatalln("error starting TransferMoney workflow", err)
	}

	var werr error
	err = we.Get(context.Background(), &werr)
	if err != nil {
		log.Println(err)

	}

}
