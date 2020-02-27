package thereum

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/bloombits"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

type FilterAPI struct {
	back *Thereum
}

func (f *FilterAPI) BloomStatus() (uint64, uint64) {
	log.Println("BloomStatus not supported")
	return 0, 0
}

func (f *FilterAPI) GetLogs(ctx context.Context, blockHash common.Hash) ([][]*types.Log, error) {
	receipts := f.back.blockchain.GetReceiptsByHash(blockHash)
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (f *FilterAPI) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	log.Println("service filter not supported")
}

func (f *FilterAPI) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return f.back.BlockChain().SubscribeLogsEvent(ch)
}

func (f *FilterAPI) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return f.back.BlockChain().SubscribeRemovedLogsEvent(ch)
}

type MiningAPI struct {
	back *Thereum
}

func (m *MiningAPI) StartMining(threads int) error {
	return nil
}

func (m *MiningAPI) StopMining() {
	// Update the thread count within the consensus engine
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := m.back.engine.(threaded); ok {
		th.SetThreads(-1)
	}
	// Stop the block creating itself
	m.back.miner.Stop()
}

func (m *MiningAPI) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return m.back.miner.SubscribePendingLogs(ch)

}

func (m *MiningAPI) IsMining() bool {
	return m.back.miner.Mining()
}
