package thereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

type Statistics struct {
	back *Thereum
}

func (s *Statistics) ProtocolVersion() int {
	return 1
}

func (s *Statistics) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return
}

func (s *Statistics) ExtRPCEnabled() bool {

}

func (s *Statistics) RPCGasCap() *big.Int {

}

func (s *Statistics) ChainConfig() *params.ChainConfig {

}

func (s *Statistics) CurrentBlock() *types.Block {

}
