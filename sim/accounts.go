package sim

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Account represents a singular wallet
type Account struct {
	Name    string            `json:"name"`
	Address common.Address    `json:"address"`
	PrivKey *ecdsa.PrivateKey `json:"private_key"`
	Balance *big.Int          `json:"balance"`
	TxOpts  *bind.TransactOpts
}

// IncrNonce increases the nonce by plus (default of 1 if plus == nil)
func (a *Account) IncrNonce(plus *big.Int) {
	if plus == nil {
		plus = new(big.Int).SetInt64(1)
	}
	a.TxOpts.Nonce.Add(a.TxOpts.Nonce, plus)
}

// SendETH signs a transaction and sends it to the client sending the provided amount of ETH
// to the provided address
func (a *Account) SendETH(client bind.ContractBackend, addr common.Address, amount *big.Int) (string, error) {
	tx := types.NewTransaction(
		a.TxOpts.Nonce.Uint64(), addr, amount,
		a.TxOpts.GasLimit,
		a.TxOpts.GasPrice,
		[]byte{},
	)
	tx, err := a.Sign(tx)
	if err != nil {
		return "", err
	}
	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// Sign uses info in Account a to sign the provided transaction
func (a *Account) Sign(tx *types.Transaction) (*types.Transaction, error) {
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), a.PrivKey)
	if err != nil {
		return nil, err
	}
	a.IncrNonce(nil)
	return signedTx, nil
}

// NewAccount issues a new account with a freshly generated private key
func NewAccount(name string, bal *big.Int) (*Account, error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	topt := bind.NewKeyedTransactor(priv)

	return &Account{
		Name:    name,
		Address: topt.From,
		TxOpts:  topt,
		Balance: bal,
	}, nil
}

// Accounts connects simple names (or any string) to a transactor
type Accounts map[string]*Account

// NewAccounts generates private keys and author accounts for testing
func NewAccounts(names ...string) Accounts {
	out := make(map[string]*Account)
	for _, name := range names {
		acc, err := NewAccount(name, new(big.Int))
		if err != nil {
			return nil
		}
		out[name] = acc
	}
	return Accounts(out)
}

// Genesis converts transactors to Genesis Accounts with 100 ETH allocations
func (ta Accounts) Genesis() core.GenesisAlloc {
	out := make(core.GenesisAlloc)
	for _, acc := range ta {
		out[acc.TxOpts.From] = core.GenesisAccount{Balance: acc.Balance}
	}
	return out
}

// Render formats account info for easy printing
func (ta *Accounts) Render() string {
	var out []string
	for name, acc := range *ta {
		s := fmt.Sprintf("%s: %s", name, acc.TxOpts.From.Hex())
		out = append(out, s)
	}
	return strings.Join(out, "\n")
}

// SetGasPrice uses the provided client to suggest a gas price and sets it
// for all accounts.
func (ta *Accounts) SetGasPrice(gasPrice *big.Int) error {
	for _, acc := range *ta {
		acc.TxOpts.GasPrice = gasPrice
	}
	return nil
}

// SetNonce uses the provided client fetch nonce for each account and sets it
// for all accounts.
func (ta *Accounts) SetNonce(back bind.ContractBackend) error {
	for _, acc := range *ta {
		nonce, err := back.PendingNonceAt(context.Background(), acc.TxOpts.From)
		if err != nil {
			return err
		}
		acc.TxOpts.Nonce = new(big.Int).SetUint64(nonce)
	}
	return nil
}

// newBlankAuth generates a new private key and creates an authenticated
// transactor with that key
func newBlankAuth() (*bind.TransactOpts, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return bind.NewKeyedTransactor(privateKey), nil
}
