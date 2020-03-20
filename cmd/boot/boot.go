package boot

import (
	"log"

	"github.com/evan-forbes/buddy/node"
	"github.com/urfave/cli/v2"
)

func Boot(c *cli.Context) error {

	// Create a network node to run protocols with the default values.
	stack, err := node.New(&node.DefaultConfig)
	if err != nil {
		log.Fatalf("Failed to create network node: %v", err)
	}
	defer stack.Close()
	return nil
}

var defNodeconf = &node.Config{
	Name:      "Thereum",
	UserIdent: "Founder",
	Version:   "0.0.1",
	DataDir:   ".",
	HTTPHost:  "127.0.0.1",
	HTTPPort:  420024,
	WSHost:    "127.0.0.1",
	WSPort:    420025,
}
