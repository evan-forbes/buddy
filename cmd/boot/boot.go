package boot

import (
	"context"
	"log"
	"sync"

	"github.com/evan-forbes/buddy/cmd"
	"github.com/evan-forbes/buddy/node"
	"github.com/evan-forbes/buddy/thereum"
	"github.com/urfave/cli/v2"
)

func Boot(c *cli.Context) error {
	// create a waitgroup and listen for ctrl + c to cancel
	var wg sync.WaitGroup
	mngr := cmd.NewManager(context.Background(), &wg)
	go mngr.Listen()

	// Create a network node to run protocols with the default values.
	stack, err := node.New(&node.DefaultConfig)
	if err != nil {
		log.Fatalf("Failed to create network node: %v", err)
	}
	defer stack.Close()

	// define the constructor
	constructor := func(ctx *node.ServiceContext) (node.Service, error) {
		return thereum.New(mngr.Ctx, &wg, &thereum.DefaultConfig)
	}

	err = stack.Register(constructor)
	if err != nil {
		return err
	}
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
