package sim

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
)

// TODO: delete and make it something that doesn't just install arbitrary defaults

// Backend manages the simulated backend and accounts for testing
type Backend struct {
	Accounts Accounts
	*backends.SimulatedBackend
}

// NewBackend generates a new simulated backend with 7 generated accounts. Typically "owner"
// is used for deploying infrastructural contracts
func NewBackend(gasLim uint64) *Backend {
	accnts := NewAccounts("owner", "Alice", "Bob", "Celine", "Doug", "Erin", "Frank")
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
		auth.TxOpts.GasPrice = gasPrice
	}
	return nil
}

// SetNonce uses the provided client fetch nonce for each account and sets it
// for all accounts.
func (back *Backend) SetNonce() error {
	for _, auth := range back.Accounts {
		nonce, err := back.PendingNonceAt(context.Background(), auth.TxOpts.From)
		if err != nil {
			return err
		}
		auth.TxOpts.Nonce = new(big.Int).SetUint64(nonce)
	}
	return nil
}
