package thereum

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/params"
)

type Config struct {
	Prune           bool                 `json:"prune"`
	GasPrice        *big.Int             `json:"gas_price"`
	GasCap          *big.Int             `json:"gas_cap"`
	TxPoolConfig    core.TxPoolConfig    `json:"tx_pool_config"`
	CliqueConfig    *params.CliqueConfig `json:"clique_config"`
	CacheConfig     *core.CacheConfig    `json:"cache_config"`
	MinerConfig     *miner.Config        `json:"miner_config"`
	VMConfig        *vm.Config           `json:"vm_config"`
	TrieCleanCache  int                  `json:"trie_clean_cache"`
	TrieDirtyCache  int                  `json:"trie_dirty_cache"`
	TrieTimeout     time.Duration        `json:"trie_timeout"`
	DatabaseHandles int                  `json:"database_handles"`
	DatabaseCache   int                  `json:"database_cache"`
	DatabaseFreezer string               `json:"database_freezer"`
	Path            string               `json:"path_to_db"`
	Genesis         *core.Genesis        `json:"genesis"`
}

var DefaultConfig = Config{
	Prune:        false,
	GasPrice:     big.NewInt(params.GWei),
	GasCap:       big.NewInt(1000000000000),
	TxPoolConfig: DefaultTxPoolConfig,
	CliqueConfig: DefaultCliqueConfig,
	CacheConfig:  DefaultCacheConfig,
	Genesis:      DefaultGenesis,
	MinerConfig:  DefaultMinerConfig,

	// DatabaseCache:  512,
	// TrieCleanCache: 256,
	// TrieDirtyCache: 256,
	// TrieTimeout:    60 * time.Minute,
	Path: "./chaindata",
}

// DefaultTxPoolConfig contains the default configurations for the transaction
// pool.
var DefaultTxPoolConfig = core.TxPoolConfig{
	Journal:   "transactions.rlp",
	Rejournal: time.Hour,
	NoLocals:  true,

	PriceLimit: 1,
	PriceBump:  10,

	AccountSlots: 16,
	GlobalSlots:  4096,
	AccountQueue: 64,
	GlobalQueue:  1024,

	Lifetime: 3 * time.Hour,
}

// // CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
// type CliqueConfig struct {
// 	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
// 	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
// }

var DefaultCliqueConfig = &params.CliqueConfig{
	Period: 1,
	Epoch:  2000,
}

var DefaultGenesis = &core.Genesis{
	Config:     params.AllCliqueProtocolChanges,
	GasLimit:   1000000000000,
	Timestamp:  uint64(time.Now().Unix()),
	Difficulty: big.NewInt(0),
	Alloc:      core.GenesisAlloc{},
}

var DefaultCacheConfig = &core.CacheConfig{
	TrieCleanLimit:      256,
	TrieCleanNoPrefetch: false,
	TrieDirtyLimit:      256,
	TrieDirtyDisabled:   true,
	TrieTimeLimit:       time.Hour,
}
