package thereum

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/core/rawdb"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
)

type Thereum struct {
	engine   consensus.Engine
	db       ethdb.Database
	eventMux *event.TypeMux

	wg   *sync.WaitGroup
	lock sync.RWMutex // Protects the variadic fields (e.g. gas price)
}

func New(ctx context.Context, wg *sync.WaitGroup, config *Config) (*Thereum, error) {
	// ensure gas price is not nil nor zero
	if config.GasPrice == nil || config.GasPrice <= 0 {
		config.GasPrice = new(big.Int).Set(DefaultConfig.GasPrice)
	}
	// start a database

	// Assemble the Thereum object
	chainDb, err := ctx.OpenDatabaseWithFreezer("chaindata", config.DatabaseCache, config.DatabaseHandles, config.DatabaseFreezer, path)
	if err != nil {
		return nil, err
	}
}

func getDB(config *Config) (ethdb.Database, error) {
	if config.InMemDB {
		return rawdb.NewMemoryDatabase(), nil
	}

}

func OpenDatabase(name string, cache int, handles int, namespace string) (ethdb.Database, error) {
	if ctx.config.DataDir == "" {
		return rawdb.NewMemoryDatabase(), nil
	}
	return rawdb.NewLevelDBDatabase(ctx.config.ResolvePath(name), cache, handles, namespace)
}
