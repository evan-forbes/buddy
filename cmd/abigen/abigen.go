package abigen

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"

	"github.com/evan-forbes/buddy/bind"
	"github.com/evan-forbes/wand"
)

// Spell funnels arguments and flags to the bite command
type Spell struct{}

// TODOS:
// - load files properly and check local path for many bins and abis
// - generate code maybe modify to generate into files by type
// - figure out how abigen inputs types
// - write to file
// - generate code and fix bugs

// Cast fullfills the wand.Spell interface
func (s *Spell) Cast(ctx wand.Context) {
	// set flags and args
	flags := ctx.Flags()
	args := ctx.Args()
	// check to see if path was provided
	path := args["abigen"]

	if path != "" {
		// load all bins abi and types
	}

	// use specific flag for paths OR look for files in the working directory
	abiPath, has := flags["abi"]
	if !has {
		// search for abi file in current directory
		abiPath, has = findFile(path, "abi")
		if !has {
			fmt.Println("could not find abi file \n(use flag --abi= or ensure you have a abi in the working dir)")
			return
		}
	}
	binPath, has := flags["bin"]
	if !has {
		// search for bin file in current directory
		binPath, _ = findFile(path, "bin")
	}

	// Load files
	a, err := loadABI(abiPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	hexBin := loadBin(binPath)

	bind.Bind()

}

func loadBin(path string) string {
	if path == "" {
		return ""
	}
	// load bin
	rawBin, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could read bin file:", path, err)
		return ""
	}
	return rawBin
}

func loadABI(path string) (*abi.ABI, error) {
	abiFile, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not open abi file: %s", path)
	}
	contractABI, err := abi.JSON(abiFile)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse abi file: %s", path)
	}
	return &contractABI, nil
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
