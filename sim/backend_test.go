package sim

import "testing"

func TestSendEth(t *testing.T) {
	back := NewBackend(uint64(4712388))
	alice := back.Accounts["Alice"]
	bobAddr := back.Accounts["Bob"].From

}
