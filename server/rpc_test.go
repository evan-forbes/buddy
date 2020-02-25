package server

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/evan-forbes/buddy/sim"
)

func defaultGenesisAlloc() core.GenesisAlloc {
	accnts := sim.NewAccounts("owner", "Alice", "Bob", "Celine", "Doug", "Erin", "Frank")
	return accnts.Genesis()
}

// func TestServer(t *testing.T) {
// 	var wg sync.WaitGroup
// 	mngr := cmd.NewManager(context.Background(), &wg)

// 	go mngr.Listen()
// 	Boot(&node.Config{
// 		HTTPHost: "127.0.0.1",
// 		HTTPPort: 8537,
// 		WSHost: "127.0.0.1",
// 		WSPort: 8538,
// 	},)

// 	<-mngr.Done()
// }
