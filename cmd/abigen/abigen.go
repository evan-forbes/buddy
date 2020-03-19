package abigen

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/evan-forbes/buddy/bind"
	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

// TODO: recursively generate over a directory

// Cast runs the abigen command
func Cast(ctx *cli.Context) error {
	////// check for flags //////
	// fail if no package is specified
	if ctx.String("pkg") == "" {
		return errors.New("No package declared. Use flag --pkg or -p")
	}
	// set value of type (default is whatever pkg is set to)
	var tp string
	if ctx.String("type") == "" {
		tp = ctx.String("pkg")
	}
	// use defaults for abi and bin file locations (".")
	abiPath := ctx.String("abi")
	binPath := ctx.String("bin")
	path := "."
	if ctx.NArg() > 0 {
		path = ctx.Args().Get(1)
	}
	// try to find paths to abi and bin files if not specified.
	if abiPath == "." {
		newABIPath, has := findFile(path, ".abi")
		if !has {
			return errors.Errorf("Could not find a .abi file in %s", path)
		}
		abiPath = newABIPath
	}
	if binPath == "." {
		binPath, _ = findFile(path, ".bin")
	}
	// load the files
	jsonABI, hexBin, err := openFiles(abiPath, binPath)
	if err != nil {
		return errors.Wrapf(err, "Problem loading files in abi path: %s bin path: %s", abiPath, binPath)
	}
	// generate bindings
	code, err := bind.Bind(
		[]string{tp},
		[]string{jsonABI},
		[]string{hexBin},
		ctx.String("pkg"),
	)
	// write to file
	filename := fmt.Sprintf("%s_gen.go", tp)
	if ctx.String("out") != "" {
		filename = ctx.String("out")
	}
	err = ioutil.WriteFile(filename, []byte(code), 0644)
	if err != nil {
		return err
	}

	return nil
}

// open files
func openFiles(abiPath, binPath string) (a, b string, err error) {
	b, err = loadBin(binPath)
	if err != nil {
		return a, b, err
	}
	a, err = loadABI(abiPath)
	if err != nil {
		return a, b, err
	}
	return a, b, nil
}

// load the binary into a string, if path is not viable ignore
func loadBin(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	// load bin
	rawBin, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "Could read bin file at path: %s", path)
	}
	return string(rawBin), nil
}

// load abi file into a string, throw error if path not viable
func loadABI(path string) (string, error) {
	if path == "" {
		return "", errors.New("could not load abi json")
	}
	// load abi
	rawABI, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "Could read abi file: %s", path)
	}
	return string(rawABI), nil
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

// Recycling bin
// // set flags and args
// flags := ctx.Flags()
// args := ctx.Args()
// // check to see if path was provided
// path := args["abigen"]

// if path != "" !strings.Contains(path, "--") {
// 	// load all bins abi and types
// }

// // use specific flag for paths OR look for files in the working directory
// abiPath, has := flags["abi"]
// if !has {
// 	// search for abi file in current directory
// 	abiPath, has = findFile(path, "abi")
// 	if !has {
// 		fmt.Println("could not find abi file \n(use flag --abi= or ensure you have a abi in the working dir)")
// 		return
// 	}
// }
// binPath, has := flags["bin"]
// if !has {
// 	// search for bin file in current directory
// 	binPath, _ = findFile(path, "bin")
// }

// // Load files
// a, err := loadABI(abiPath)
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
// hexBin := loadBin(binPath)
