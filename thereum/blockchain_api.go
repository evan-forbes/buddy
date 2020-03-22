package thereum

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

type BlockchainAPI struct {
	back *Thereum
}

func (b *BlockchainAPI) ChainDb() ethdb.Database {
	return b.back.db
}

func (b *BlockchainAPI) AccountManager() *accounts.Manager {
	return b.back.accountMngr
}

func (b *BlockchainAPI) Downloader() *downloader.Downloader {
	log.Println("Downloader not supported")
	return nil
}

func (b *BlockchainAPI) SetHead(number uint64) {
	log.Println("Setting Head not supported")
}

func (b *BlockchainAPI) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	return b.back.blockchain.GetHeaderByNumber(uint64(number)), nil
}

func (b *BlockchainAPI) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.back.blockchain.GetHeaderByHash(hash), nil
}

func (b *BlockchainAPI) HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.HeaderByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header := b.back.blockchain.GetHeaderByHash(hash)
		if header == nil {
			return nil, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.back.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return nil, errors.New("hash is not currently canonical")
		}
		return header, nil
	}
	return nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *BlockchainAPI) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	// Otherwise resolve and return the block
	if number == rpc.LatestBlockNumber {
		return b.back.blockchain.CurrentBlock(), nil
	}
	return b.back.blockchain.GetBlockByNumber(uint64(number)), nil
}

func (b *BlockchainAPI) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.back.blockchain.GetBlockByHash(hash), nil
}

func (b *BlockchainAPI) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.BlockByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header := b.back.blockchain.GetHeaderByHash(hash)
		if header == nil {
			return nil, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.back.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return nil, errors.New("hash is not currently canonical")
		}
		block := b.back.blockchain.GetBlock(hash, header.Number.Uint64())
		if block == nil {
			return nil, errors.New("header found, but block body is missing")
		}
		return block, nil
	}
	return nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *BlockchainAPI) StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, nil, err
	}
	if header == nil {
		return nil, nil, errors.New("header not found")
	}
	stateDb, err := b.back.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *BlockchainAPI) StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.StateAndHeaderByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header, err := b.HeaderByHash(ctx, hash)
		if err != nil {
			return nil, nil, err
		}
		if header == nil {
			return nil, nil, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.back.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return nil, nil, errors.New("hash is not currently canonical")
		}
		stateDb, err := b.back.BlockChain().StateAt(header.Root)
		return stateDb, header, err
	}
	return nil, nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *BlockchainAPI) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	return b.back.blockchain.GetReceiptsByHash(hash), nil
}

func (b *BlockchainAPI) GetTd(hash common.Hash) *big.Int {
	return b.back.blockchain.GetTdByHash(hash)
}

func (b *BlockchainAPI) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.back.BlockChain(), nil)
	return vm.NewEVM(context, state, b.back.blockchain.Config(), *b.back.blockchain.GetVMConfig()), vmError, nil
}

func (b *BlockchainAPI) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.back.BlockChain().SubscribeChainEvent(ch)
}

func (b *BlockchainAPI) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.back.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *BlockchainAPI) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.back.BlockChain().SubscribeChainSideEvent(ch)
}
