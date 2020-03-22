package thereum

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/bloombits"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/evan-forbes/buddy/ethapi"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
)

// Thereum holds all datastructures needed for blockchain functionality.
// Various wrappers are used to implement further functionality.
type Thereum struct {
	config        *Config
	engine        consensus.Engine
	db            ethdb.Database
	blockchain    *core.BlockChain
	txPool        *core.TxPool
	chainConfig   *params.ChainConfig
	eventMux      *event.TypeMux
	miner         *miner.Miner
	bloom         *core.ChainIndexer
	bloomRequests chan chan *bloombits.Retrieval
	accountMngr   *accounts.Manager

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
		config:        config,
		engine:        engine,
		db:            chainDb,
		chainConfig:   chainConfig,
		eventMux:      &event.TypeMux{},
		ctx:           ctx,
		wg:            wg,
		bloom:         eth.NewBloomIndexer(chainDb, params.BloomBitsBlocks, params.BloomConfirms),
		bloomRequests: make(chan chan *bloombits.Retrieval),
		accountMngr:   accounts.NewManager(&accounts.Config{}),
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

// create methods to follow the node.Service interface

func (t *Thereum) Protocols() []p2p.Protocol { return []p2p.Protocol{} }

func (t *Thereum) APIs() (out []rpc.API) {
	// add the api functionality to the base thereum obj
	back := NewBackendAPI(t)
	apis := ethapi.GetAPIs(back)

	// Add any APIs exposed explicitly by the consensus engine
	apis = append(apis, t.engine.APIs(t.BlockChain())...)

	// Add local APIs
	apis = append(apis, []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(back, false),
			Public:    true,
		},
	}...)

	return apis
}

func (t *Thereum) Start(serv *p2p.Server) error {
	t.wg.Add(1)
	t.startBloomHandlers(params.BloomBitsBlocks)
	return nil
}

func (t *Thereum) Stop() error {
	t.bloom.Close()
	t.blockchain.Stop()
	t.engine.Close()
	t.txPool.Stop()
	t.miner.Stop()
	t.eventMux.Stop()
	t.db.Close()
	t.wg.Done()
	return nil
}

// The following apis have been removed for now
// {
// 	Namespace: "eth",
// 	Version:   "1.0",
// 	Service:   downloader.NewPublicDownloaderAPI(t.protocolManager.downloader, s.eventMux),
// 	Public:    true,
// },
// {
// 	Namespace: "debug",
// 	Version:   "1.0",
// 	Service:   NewPublicDebugAPI(t),
// 	Public:    true,
// }, {
// 	Namespace: "debug",
// 	Version:   "1.0",
// 	Service:   NewPrivateDebugAPI(t),
// },
// {
// 	Namespace: "eth",
// 	Version:   "1.0",
// 	Service:   NewPublicEthereumAPI(t),
// 	Public:    true,
// },
// {
// 	Namespace: "net",
// 	Version:   "1.0",
// 	Service:   t.netRPCService,
// 	Public:    true,
// },
// {
// 	Namespace: "eth",
// 	Version:   "1.0",
// 	Service:   NewPublicMinerAPI(t),
// 	Public:    true,
// },
// {
// 	Namespace: "miner",
// 	Version:   "1.0",
// 	Service:   NewPrivateMinerAPI(t),
// 	Public:    false,
// },

// * defo add later

// {
// 	Namespace: "admin",
// 	Version:   "1.0",
// 	Service:   NewPrivateAdminAPI(t),
// },
// // PrivateAdminAPI is the collection of Ethereum full node-related APIs
// // exposed over the private admin endpoint.
// type PrivateAdminAPI struct {
// 	eth *Ethereum
// }

// // NewPrivateAdminAPI creates a new API definition for the full node private
// // admin methods of the Ethereum service.
// func NewPrivateAdminAPI(eth *Ethereum) *PrivateAdminAPI {
// 	return &PrivateAdminAPI{eth: eth}
// }

// // ExportChain exports the current blockchain into a local file,
// // or a range of blocks if first and last are non-nil
// func (api *PrivateAdminAPI) ExportChain(file string, first *uint64, last *uint64) (bool, error) {
// 	if first == nil && last != nil {
// 		return false, errors.New("last cannot be specified without first")
// 	}
// 	if first != nil && last == nil {
// 		head := api.t.BlockChain().CurrentHeader().Number.Uint64()
// 		last = &head
// 	}
// 	if _, err := os.Stat(file); err == nil {
// 		// File already exists. Allowing overwrite could be a DoS vecotor,
// 		// since the 'file' may point to arbitrary paths on the drive
// 		return false, errors.New("location would overwrite an existing file")
// 	}
// 	// Make sure we can create the file to export into
// 	out, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
// 	if err != nil {
// 		return false, err
// 	}
// 	defer out.Close()

// 	var writer io.Writer = out
// 	if strings.HasSuffix(file, ".gz") {
// 		writer = gzip.NewWriter(writer)
// 		defer writer.(*gzip.Writer).Close()
// 	}

// 	// Export the blockchain
// 	if first != nil {
// 		if err := api.eth.BlockChain().ExportN(writer, *first, *last); err != nil {
// 			return false, err
// 		}
// 	} else if err := api.eth.BlockChain().Export(writer); err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// func hasAllBlocks(chain *core.BlockChain, bs []*types.Block) bool {
// 	for _, b := range bs {
// 		if !chain.HasBlock(b.Hash(), b.NumberU64()) {
// 			return false
// 		}
// 	}

// 	return true
// }

// // ImportChain imports a blockchain from a local file.
// func (api *PrivateAdminAPI) ImportChain(file string) (bool, error) {
// 	// Make sure the can access the file to import
// 	in, err := os.Open(file)
// 	if err != nil {
// 		return false, err
// 	}
// 	defer in.Close()

// 	var reader io.Reader = in
// 	if strings.HasSuffix(file, ".gz") {
// 		if reader, err = gzip.NewReader(reader); err != nil {
// 			return false, err
// 		}
// 	}

// 	// Run actual the import in pre-configured batches
// 	stream := rlp.NewStream(reader, 0)

// 	blocks, index := make([]*types.Block, 0, 2500), 0
// 	for batch := 0; ; batch++ {
// 		// Load a batch of blocks from the input file
// 		for len(blocks) < cap(blocks) {
// 			block := new(types.Block)
// 			if err := stream.Decode(block); err == io.EOF {
// 				break
// 			} else if err != nil {
// 				return false, fmt.Errorf("block %d: failed to parse: %v", index, err)
// 			}
// 			blocks = append(blocks, block)
// 			index++
// 		}
// 		if len(blocks) == 0 {
// 			break
// 		}

// 		if hasAllBlocks(api.eth.BlockChain(), blocks) {
// 			blocks = blocks[:0]
// 			continue
// 		}
// 		// Import the batch and reset the buffer
// 		if _, err := api.eth.BlockChain().InsertChain(blocks); err != nil {
// 			return false, fmt.Errorf("batch %d: failed to insert: %v", batch, err)
// 		}
// 		blocks = blocks[:0]
// 	}
// 	return true, nil
// }
