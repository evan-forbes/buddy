package abigen

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/evan-forbes/wand"
)

// Spell funnels arguments and flags to the bite command
type Spell struct{}

// Cast fullfills the wand.Spell interface
func (s *Spell) Cast(ctx wand.Context) {
	// set flags and arsg
	flags := ctx.Flags()
	args := ctx.Args()
	// check to see if path was provided
	path := args["abigen"]

	// use specific flag OR look for files in the working directory
	abiPath, has := flags["abi"]
	if !has {
		// search for abi file in current directory
		abiPath, has = findFile(path, "abi")
		if !has {
			fmt.Println("could not find abi file in working dir")
			return
		}
	}
	binPath, has := flags["bin"]
	if !has {
		// search for abi file in current directory
		binPath, has = findFile(path, "bin")
		if !has {
			fmt.Println("could not find bin file in working dir")
			return
		}
	}

	// load files
	abiFile, err := os.Open(abiPath)
	if err != nil {
		fmt.Println("Could not find file:", abiPath, err)
		return
	}
	contractABI, err := abi.JSON(abiFile)
	if err != nil {
		fmt.Println("could not open abi file:", err)
		return 
	}
	
}

// findFile returns the first file found with the provided type
func findFile(path, ext string) (string, bool) {
	if path == "" {
		path = "."
	}
	items, err := ioutil.ReadDir(path)
	if err != nil {
		return "", false
	}
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		if strings.Contains(item.Name(), ext) {
			return item.Name(), true
		}
	}
	return "", false
}
