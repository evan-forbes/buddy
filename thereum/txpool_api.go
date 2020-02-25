package thereum

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

type TxPoolAPI struct {
	back *Thereum
}

func (t *TxPoolAPI) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return nil
}

func (t *TxPoolAPI) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, error) {
	tx, blockHash, blockNumber, index := rawdb.ReadTransaction(t.back.db, txHash)
	return tx, blockHash, blockNumber, index, nil
}

func (t *TxPoolAPI) GetPoolTransactions() (types.Transactions, error) {
	pending, err := t.back.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (t *TxPoolAPI) GetPoolTransaction(txHash common.Hash) *types.Transaction {
	return t.back.txPool.Get(txHash)
}

func (t *TxPoolAPI) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return t.back.txPool.Nonce(addr), nil
}

func (t *TxPoolAPI) Stats() (pending int, queued int) {
	return t.back.txPool.Stats()
}

func (t *TxPoolAPI) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return t.back.txPool.Content()
}

func (t *TxPoolAPI) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return t.back.txPool.SubscribeNewTxsEvent(ch)
}

func (t *TxPoolAPI) TxPool() *core.TxPool {
	return t.back.txPool
}
