package workflows

import (
	"log"

	"go.temporal.io/sdk/client"
)

var C client.Client
var err error
var DefaultHostPort = "49.235.242.124:7233"

func init() {
	// The client and worker are heavyweight objects that should be created once per process.
	C, err = client.NewClient(client.Options{
		HostPort: DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}

}
