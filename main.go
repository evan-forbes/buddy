package main

import (
	"context"
	"os"

	"github.com/evan-forbes/buddy/cmd/abigen"
	"github.com/evan-forbes/buddy/cmd/solc"
	"github.com/evan-forbes/wand"
)

func main() {
	ctx := wand.NewDefaultContext(context.Background())
	spells := map[string]wand.Spell{
		"abigen": &abigen.Spell{},
		"solc":   &solc.Spell{},
	}
	wand.Run(ctx, spells, os.Args[1:])
}

/* TODO:

finish event and unpackers
save to file
finish cli tool


*/
