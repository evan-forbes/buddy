package thereum

import (
	"context"
	"fmt"
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

func (f *FilterAPI) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	fmt.Println("pending logs are not yet supported")
	return nil
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
	return
}

func (m *MiningAPI) IsMining() bool {
	return true
}
