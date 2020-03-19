package boot

import (
	"log"

	"github.com/evan-forbes/buddy/node"
	"github.com/urfave/cli/v2"
)

func Boot(c *cli.Context) error {
	// Create a network node to run protocols with the default values.
	stack, err := node.New(&node.Config{})
	if err != nil {
		log.Fatalf("Failed to create network node: %v", err)
	}
	defer stack.Close()
	return nil
}
