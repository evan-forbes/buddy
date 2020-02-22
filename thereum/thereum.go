package thereum

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
)

// Thereum holds all datastructures needed for blockchain functionality.
// Various wrappers are used to implement further functionality.
type Thereum struct {
	config   *Config
	engine   consensus.Engine
	db       ethdb.Database
	eventMux *event.TypeMux

	wg   *sync.WaitGroup
	ctx  context.Context
	lock sync.RWMutex // Protects the variadic fields (e.g. gas price)
}

// New issues a new PoA Thereum blockchain based on the provided config.
func New(ctx context.Context, wg *sync.WaitGroup, config *Config) (*Thereum, error) {
	// ensure gas price is not nil nor zero
	if config.GasPrice == nil || config.GasPrice.Cmp(big.NewInt(0)) < 1 {
		config.GasPrice = new(big.Int).Set(DefaultConfig.GasPrice)
	}
	// open/start a database
	chainDb, err := openDatabase("chaindata", config.DatabaseCache, config.DatabaseHandles, config.Path)
	if err != nil {
		return nil, err
	}
	// setup a genesis block based on the config
	_, _, err = core.SetupGenesisBlock(chainDb, config.Genesis)
	if err != nil {
		return nil, errors.Wrap(err, "Could not setup genesis block during Thereum initialization")
	}
	// create PoA consensus engine
	engine := clique.New(config.CliqueConfig, chainDb)

	// construct the Thereum object
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

func openDatabase(name string, cache int, handles int, path string) (ethdb.Database, error) {
	if path == "" {
		return rawdb.NewMemoryDatabase(), nil
	}
	return rawdb.NewLevelDBDatabase(name, cache, handles, path)
}

type Backend struct {
}

func NewBackend() *Backend {
	// make thereum object
	// wrap with backend
	// create a gas price oracle and add it
}
