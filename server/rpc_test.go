package server

import (
	"testing"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/node"
	"github.com/evan-forbes/buddy/sim"
)

func defaultGenesisAlloc() core.GenesisAlloc {
	accnts := sim.NewAccounts("owner", "Alice", "Bob", "Celine", "Doug", "Erin", "Frank")
	return accnts.Genesis()
}

func TestServer(t *testing.T) {
	Boot(&node.Config{})
}
