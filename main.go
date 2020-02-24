package main

import (
	"log"
	"os"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/evan-forbes/buddy/cmd/abigen"
)

// TODOs:
/*
 - 1) finish fullfilling the Backend interface for thereum.Thereum
 - 2) fork the rest of the node package to eliminate confusion/whatever

 ### Longerterm
 - cut more out of node/other forked packages to cut out bloat
*/

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
