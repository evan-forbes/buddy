package main

import (
	"log"
	"os"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/evan-forbes/buddy/cmd/abigen"
)

// I need to find out how to design services, like the example node
// and I need to get rid of the services that I don't need.

// those services connect throught the websocket and http server in 
// the node instance. the node instance uses the startRPC method to
// boot up operations. The startRPC method asks for []rpc.API to start

// still need to figure out how to connect the simulated blockchain to the local rpc

// implement loading and unloading of a chain using cmdutils.ImportChain

// Planned Commands:
//  - mock
//  - sim
//  - abigen

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "buddy"

	// abiFlags are flags for the subcommand abigen
	abiFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "abi, a",
			Value: ".",
			Usage: "path to abi (.json or .abi)",
			// Destination: &abiPath,
		},
		cli.StringFlag{
			Name:  "bin, b",
			Value: ".",
			Usage: "path to contract binary (usually a .bin)",
			// Destination: &binPath,
		},
		cli.StringFlag{
			Name:  "type, t",
			Value: "",
			Usage: "specify the main type",
			// Destination: &tp,
		},
		cli.StringFlag{
			Name:  "pkg, p",
			Value: "",
			Usage: "specify the package name",
			// Destination: &tp,
		},
		cli.StringFlag{
			Name:  "out, o",
			Value: "",
			Usage: "specify the output file name (default = type_gen.go",
			// Destination: &tp,
		},
	}

	// subcommands
	app.Commands = []cli.Command{
		{
			Name:   "abigen",
			Usage:  "generate interface and mock friendly go bindings",
			Action: abigen.Cast,
			Flags:  abiFlags,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

/* TODO:

finish event and unpackers
save to file
finish cli tool

	ctx := wand.NewDefaultContext(context.Background())
	spells := map[string]wand.Spell{
		"abigen": &abigen.Spell{},
		"solc":   &solc.Spell{},
	}
	wand.Run(ctx, spells, os.Args[1:])
*/
