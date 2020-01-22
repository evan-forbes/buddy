package parse

import "github.com/ethereum/go-ethereum/core/types"

type Parsed struct {
	Tag  string
	Data interface{}
}

// Unpacker wraps generated methods intended to parse a standardized go ethereum log
// into a go binded version. If returning an interface is inconvient for your uses, 
// see the Parse generated method, which is identical except it returns the native type
type Unpacker interface {
	Unpack(log types.Log) (parsed *Parsed, err error)
}

// Filter helps parse logs by topic
type Filter map[string]Unpacker

// Merge combines the provided filter with the caller
func (f Filter) Merge(second Filter) {
	for topic, unpkr := range second {
		f[topic] = unpkr
	}
}

func (f Filter) Mux(log types.Log) (*Parsed, error) {
	unpkr, has := f[log.Topics[0]]
	if !has {
		return nil, fmt.Errorf("No unpacker found for topic: %s", log.Topics[0])
	}
}

type XMuxer struct {

}

/*

Generate code to implement interfaces for logs as well as methods

dai, _ := dai.NewDai(common.HexToAddress("0x000..."), client)

dai.Watch(ctx context.Context, client bind.Backend) <-chan *ApprovalLog, <-chan *TransferLog, <-chan error

or

// this function simply makes the neccessary channels for you
approveChan, transferChan := dai.EventChannels() (returns <-chan *ApprovalLog, <-chan *TransferLog)

// then, for whatever logs I'm interested in
go dai.Watch(ctx, client, transferChan, errc)

// 

*/

