package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/evan-forbes/buddy/cmd/abigen"
	"github.com/evan-forbes/buddy/cmd/boot"
)

// TODOs:
/*
 - 1) finish fullfilling the Backend interface for thereum.Thereum

 ### Longerterm
 - cut more out of node/other forked packages to cut out bloat
   specifically: the node.Start
*/

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "buddy"

	// abiFlags are flags for the subcommand abigen
	abiFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "abi, a",
			Value: ".",
			Usage: "path to abi (.json or .abi)",
			// Destination: &abiPath,
		},
		&cli.StringFlag{
			Name:  "bin, b",
			Value: ".",
			Usage: "path to contract binary (usually a .bin)",
			// Destination: &binPath,
		},
		&cli.StringFlag{
			Name:  "type, t",
			Value: "",
			Usage: "specify the main type",
			// Destination: &tp,
		},
		&cli.StringFlag{
			Name:  "pkg, p",
			Value: "",
			Usage: "specify the package name",
			// Destination: &tp,
		},
		&cli.StringFlag{
			Name:  "out, o",
			Value: "",
			Usage: "specify the output file name (default = type_gen.go",
			// Destination: &tp,
		},
	}

	bootFlags := []cli.Flag{}

	// subcommands
	app.Commands = []*cli.Command{
		{
			Name:   "abigen",
			Usage:  "generate interface and mock friendly go bindings",
			Action: abigen.Cast,
			Flags:  abiFlags,
		},
		{
			Name:   "boot",
			Usage:  "generate interface and mock friendly go bindings",
			Action: boot.Boot,
			Flags:  bootFlags,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
