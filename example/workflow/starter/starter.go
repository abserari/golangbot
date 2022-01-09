package main

import (
	"context"
	"log"

	"github.com/abserari/golangbot/pkg/workflows"
	"github.com/google/uuid"
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

	// options := client.StartWorkflowOptions{
	// 	ID:        "transfer-money-workflow",
	// 	TaskQueue: "test",
	// }

	// var hello = "hello"

	// we, err := c.ExecuteWorkflow(context.Background(), options, workflows.GetAndSendMessage, hello)
	// if err != nil {
	// 	log.Fatalln("error starting TransferMoney workflow", err)
	// }

	// var werr error
	// err = we.Get(context.Background(), &werr)
	// if err != nil {
	// 	log.Println(err)

	// }

	startcron(c)

}

func startcron(c client.Client) {
	// This workflow ID can be user business logic identifier as well.
	workflowID := "cron_" + uuid.New().String()
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "cron",
		// CronSchedule: "@daily",
		CronSchedule: "* * * * *",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.CronRssFeedWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

}
