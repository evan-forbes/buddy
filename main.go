package main

import (
	"log"
	"os"

	"github.com/evan-forbes/buddy/cmd/abigen"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []*cli.Command{
		{
			Name:   "abigen",
			Usage:  "generate interface and mock friendly go bindings",
			Action: abigen.Cast,
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
