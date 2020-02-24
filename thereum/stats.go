package thereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

// StatsAPI is one of the *Thereum wrappers in BackendAPI. Specifically to provide the
// statistics portion of the api
type StatsAPI struct {
	back *Thereum
}

func (s *StatsAPI) ProtocolVersion() int {
	return 1
}

func (s *StatsAPI) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return s.back.config.GasPrice, nil
}

func (s *StatsAPI) ExtRPCEnabled() bool {
	return true
}

func (s *StatsAPI) RPCGasCap() *big.Int {
	return s.back.config.GasCap
}

func (s *StatsAPI) ChainConfig() *params.ChainConfig {
	return s.back.chainConfig
}

func (s *StatsAPI) CurrentBlock() *types.Block {
	return s.back.blockchain.CurrentBlock()
}
