package book

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/common"

	"github.com/pkg/errors"
)

// TODO: delete this package in favor of just using accounts? combine at the very least

type Book map[string]common.Address

func (b *Book) Write(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrapf(
			err,
			"Could not open or make address book file: %s",
			filename,
		)
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	err = enc.Encode(b)
	if err != nil {
		return errors.Wrapf(err, "could not write book as %s", filename)
	}
	return nil
}

// func (b *Book) MarshalJSON() ([]byte, error) {
// 	buf := bytes.NewBufferString("{")
// 	for name, addr := range *b {
// 		jsonVal, err := json.Marshal(addr.Hex())
// 		if err != nil {
// 			return nil, err
// 		}
// 		buf.WriteString(fmt.Sprintf(`"%s": "%s`))
// 	}
// }
