package server

import (
	"testing"

	"github.com/ethereum/go-ethereum/node"
	"github.com/evan-forbes/buddy/sim"
)

func TestServer(t *testing.T) {
	accnts := sim.NewAccounts("owner", "Alice", "Bob", "Celine", "Doug", "Erin", "Frank")
	genesisAlloc := accnts.Genesis()
	s := &Server{
		Backend: sim.NewSimulatedBackend(genesisAlloc, uint64(4712388)),
	}
	Boot(&node.Config{})
}
