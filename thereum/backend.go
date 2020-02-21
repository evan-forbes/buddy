package thereum

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
)

type Thereum struct {
	config   *Config
	engine   consensus.Engine
	db       ethdb.Database
	eventMux *event.TypeMux

	wg   *sync.WaitGroup
	ctx  context.Context
	lock sync.RWMutex // Protects the variadic fields (e.g. gas price)
}

func New(ctx context.Context, wg *sync.WaitGroup, config *Config) (*Thereum, error) {
	// ensure gas price is not nil nor zero
	if config.GasPrice == nil || config.GasPrice <= 0 {
		config.GasPrice = new(big.Int).Set(DefaultConfig.GasPrice)
	}
	// start a database

	// Assemble the Thereum object
	chainDb, err := OpenDatabase("chaindata", config.DatabaseCache, config.DatabaseHandles, config.Path)
	if err != nil {
		return nil, err
	}
	chainConfig, chainHash, err := core.SetupGenesisBlock(chainDb, config.Genesis)

	// create engine
	engine := clique.New(config.CliqueConfig, chainDb)

	ther := &Thereum{
		config:   config,
		engine:   engine,
		db:       chainDb,
		eventMux: &event.TypeMux{},
		ctx:      ctx,
		wg:       wg,
	}
	return ther, nil
}

func OpenDatabase(name string, cache int, handles int, path string) (ethdb.Database, error) {
	if path == "" {
		return rawdb.NewMemoryDatabase(), nil
	}
	return rawdb.NewLevelDBDatabase(name, cache, handles, path)
}
