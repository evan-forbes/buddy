package thereum

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

type Config struct {
	Prune           bool          `json:"prune"`
	GasPrice        *big.Int      `json:"gas_price"`
	TxPoolConfig    TxPoolConfig  `json:"tx_pool_config"`
	CliqueConfig    CliqueConfig  `json:"clique_config"`
	TrieCleanCache  int           `json:"trie_clean_cache"`
	TrieDirtyCache  int           `json:"trie_dirty_cache"`
	TrieTimeout     time.Duration `json:"trie_timeout"`
	DatabaseHandles int           `json:"database_handles"`
	DatabaseCache   int           `json:"database_cache"`
	DatabaseFreezer string        `json:"database_freezer"`
	InMemDB         bool          `json:"in_memory_db"`
	Path            string
}

var DefaultConfig = Config{
	Prune:          false,
	GasPrice:       big.NewInt(params.GWei),
	TxPoolConfig:   DefaultTxPoolConfig,
	CliqueConfig:   DefaultCliqueConfig,
	DatabaseCache:  512,
	TrieCleanCache: 256,
	TrieDirtyCache: 256,
	TrieTimeout:    60 * time.Minute,
	Path:           "./chaindata",
}

// TxPoolConfig are the configuration parameters of the transaction pool.
type TxPoolConfig struct {
	Locals    []common.Address // Addresses that should be treated by default as local
	NoLocals  bool             // Whether local transaction handling should be disabled
	Journal   string           // Journal of local transactions to survive node restarts
	Rejournal time.Duration    // Time interval to regenerate the local transaction journal

	PriceLimit uint64 // Minimum gas price to enforce for acceptance into the pool
	PriceBump  uint64 // Minimum price bump percentage to replace an already existing transaction (nonce)

	AccountSlots uint64 // Number of executable transaction slots guaranteed per account
	GlobalSlots  uint64 // Maximum number of executable transaction slots for all accounts
	AccountQueue uint64 // Maximum number of non-executable transaction slots permitted per account
	GlobalQueue  uint64 // Maximum number of non-executable transaction slots for all accounts

	Lifetime time.Duration // Maximum amount of time non-executable transaction are queued
}

// DefaultTxPoolConfig contains the default configurations for the transaction
// pool.
var DefaultTxPoolConfig = TxPoolConfig{
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

// CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
type CliqueConfig struct {
	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}

var DefaultCliqueConfig = CliqueConfig{
	Period: 1,
	Epoch:  2000,
}
