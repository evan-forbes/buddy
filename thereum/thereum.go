package thereum

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
)

// Thereum holds all datastructures needed for blockchain functionality.
// Various wrappers are used to implement further functionality.
type Thereum struct {
	config      *Config
	engine      consensus.Engine
	db          ethdb.Database
	blockchain  *core.BlockChain
	txPool      *core.TxPool
	chainConfig *params.ChainConfig
	eventMux    *event.TypeMux
	miner       *miner.Miner
	bloom       *core.ChainIndexer

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
	chainDb, err := openDatabase(
		"chaindata",
		config.DatabaseCache,
		config.DatabaseHandles,
		config.Path,
	)
	if err != nil {
		return nil, err
	}
	// setup a genesis block based on the config
	chainConfig, _, err := core.SetupGenesisBlock(chainDb, config.Genesis)
	if err != nil {
		return nil, errors.Wrap(err, "Could not setup genesis block during Thereum initialization")
	}
	// create PoA consensus engine
	engine := clique.New(params.AllCliqueProtocolChanges.Clique, chainDb)

	// construct the Thereum object
	ther := &Thereum{
		config:      config,
		engine:      engine,
		db:          chainDb,
		chainConfig: chainConfig,
		eventMux:    &event.TypeMux{},
		ctx:         ctx,
		wg:          wg,
		bloom:       eth.NewBloomIndexer(chainDb, params.BloomBitsBlocks, params.BloomConfirms),
	}

	ther.blockchain, err = core.NewBlockChain(
		chainDb,
		config.CacheConfig,
		chainConfig,
		engine,
		vm.Config{},
		nil,
	)
	if err != nil {
		return nil, err
	}

	ther.bloom.Start(ther.blockchain)

	ther.txPool = core.NewTxPool(config.TxPoolConfig, chainConfig, ther.blockchain)
	ther.miner = miner.New(ther, config.MinerConfig, chainConfig, ther.eventMux, ther.engine, ther.IsLocalBlock)

	return ther, nil
}

func (t *Thereum) IsLocalBlock(block *types.Block) bool {
	author, err := t.engine.Author(block.Header())
	if err != nil {
		log.Warn("Failed to retrieve block author", "number", block.NumberU64(), "hash", block.Hash(), "err", err)
		return false
	}
	// Check whether the given address is etherbase.
	t.lock.RLock()
	etherbase := t.config.MinerConfig.Etherbase
	t.lock.RUnlock()
	if author == etherbase {
		return true
	}
	return false
}

func (t *Thereum) BlockChain() *core.BlockChain {
	return t.blockchain
}

func (t *Thereum) TxPool() *core.TxPool {
	return t.txPool
}

func openDatabase(name string, cache int, handles int, path string) (ethdb.Database, error) {
	if path == "" {
		return rawdb.NewMemoryDatabase(), nil
	}
	return rawdb.NewLevelDBDatabase(name, cache, handles, path)
}

// BackendAPI fullfills the backend.Backend (and ethapi.Backend) interface for plugging into the
// rpc backend managed by node.Node.
type BackendAPI struct {
	back *Thereum
	FilterAPI
	MiningAPI
	StatsAPI
	TxPoolAPI
	BlockchainAPI
}

func NewBackendAPI(t *Thereum) *BackendAPI {
	return &BackendAPI{
		back:          t,
		FilterAPI:     FilterAPI{back: t},
		MiningAPI:     MiningAPI{back: t},
		StatsAPI:      StatsAPI{back: t},
		TxPoolAPI:     TxPoolAPI{back: t},
		BlockchainAPI: BlockchainAPI{back: t},
	}
}
