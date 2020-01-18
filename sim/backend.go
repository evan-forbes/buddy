package sim

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestAccounts connects simple names (or any string) to a transactor
type TestAccounts map[string]*bind.TransactOpts

// NewTestAccounts generates private keys and author accounts for testing
func NewTestAccounts(names ...string) TestAccounts {
	out := make(map[string]*bind.TransactOpts)
	for _, name := range names {
		auth, err := NewBlankAuth()
		if err != nil {
			fmt.Println("Could not generate author: ", err)
		}
		out[name] = auth
	}
	return TestAccounts(out)
}

// Genesis converts transactors to Genesis Accounts with 100 ETH allocations
func (ta TestAccounts) Genesis() core.GenesisAlloc {
	out := make(core.GenesisAlloc)
	for _, auth := range ta {
		bal, _ := new(big.Int).SetString("10000000000000000000", 10)
		out[auth.From] = core.GenesisAccount{Balance: bal}
	}
	return out
}

// Render formats account info for easy printing
func (ta *TestAccounts) Render() string {
	var out []string
	for name, auth := range *ta {
		s := fmt.Sprintf("%s: %s", name, auth.From.Hex())
		out = append(out, s)
	}
	return strings.Join(out, "\n")
}

// NewBlankAuth generates a new private key and creates an authenticated
// transactor with that key
func NewBlankAuth() (*bind.TransactOpts, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return bind.NewKeyedTransactor(privateKey), nil
}

// Backend manages the simulated backend and accounts for testing
type Backend struct {
	Accounts TestAccounts
	*backends.SimulatedBackend
}

// NewBackend generates a new simulated backend with 7 generated accounts. Typically "owner"
// is used for deploying infrastructural contracts
func NewBackend(gasLim uint64) *Backend {
	accnts := NewTestAccounts("owner", "Alice", "Bob", "Celine", "Doug", "Erin", "Frank")
	genesisAlloc := accnts.Genesis()
	return &Backend{
		Accounts:         accnts,
		SimulatedBackend: backends.NewSimulatedBackend(genesisAlloc, gasLim),
	}
}

// SetGasPrice uses the provided client to suggest a gas price and sets it
// for all accounts.
func (back *Backend) SetGasPrice() error {
	gasPrice, err := back.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	for _, auth := range back.Accounts {
		auth.GasPrice = gasPrice
	}
	return nil
}

// SetNonce uses the provided client to suggest a gas price and sets it
// for all accounts.
func (back *Backend) SetNonce() error {
	for _, auth := range back.Accounts {
		nonce, err := back.PendingNonceAt(context.Background(), auth.From)
		if err != nil {
			return err
		}
		auth.Nonce = nonce
	}
	return nil
}
