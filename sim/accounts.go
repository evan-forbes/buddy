package sim

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

// Account represents a singular wallet
type Account struct {
	Name         string            `json:"name"`
	Address      common.Address    `json:"address"`
	PrivKey      *ecdsa.PrivateKey `json:"private_key"`
	Balance      *big.Int          `json:"balance"`
	TransactOpts *bind.TransactOpts
}

// NewAccount issues a new account with a freshly generated private key
func NewAccount(name string, bal *big.Int) (*Account, error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	topt := bind.NewKeyedTransactor(priv)

	return &Account{
		Name:        name,
		Address:     topt.From,
		TransactOps: topt,
		Balance:     bal,
	}, nil
}

// Accounts connects simple names (or any string) to a transactor
type Accounts map[string]*Account

// NewAccounts generates private keys and author accounts for testing
func NewAccounts(names ...string) Accounts {
	out := make(map[string]*bind.TransactOpts)
	for _, name := range names {
		acc, err := NewAccount(name, new(big.Int))
		if err != nil {
			return nil, errors.Wrap(err, "Could not generate account")
		}
		out[name] = acc
	}
	return Accounts(out)
}

// Genesis converts transactors to Genesis Accounts with 100 ETH allocations
func (ta Accounts) genesis() core.GenesisAlloc {
	out := make(core.GenesisAlloc)
	for _, acc := range ta {
		out[acc.TransactOpts.From] = core.GenesisAccount{Balance: acc.Balance}
	}
	return out
}

// Render formats account info for easy printing
func (ta *Accounts) Render() string {
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
