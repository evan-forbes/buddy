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

/*

for the type of code generation that I want, I will need

each of the event log structs and some unpacking method

a wrapper around bind.BoundContract along with
a wrapper methods for each public contract method


*/
