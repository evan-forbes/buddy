package sim

// func TestSendEth(t *testing.T) {
// 	back := NewBackend(uint64(4712388))
// 	alice := back.Accounts["Alice"]
// 	bobAddr := back.Accounts["Bob"].From

// }

// func TestBackendLife(t *testing.T) {
// 	var wg sync.WaitGroup
// 	mngr := cmd.NewManager(context.Background(), &wg)

// 	go mngr.Listen()

// 	back := NewBackend(uint64(4712388))
// 	err := back.SetGasPrice()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer back.Close()
// 	err = back.SetNonce()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	headerch := make(chan *types.Header)

// 	sub, err := back.SubscribeNewHead(mngr.Ctx, headerch)

// 	go func() {
// 		for {
// 			select {
// 			case head, ok := <-headerch:
// 				if !ok {
// 					return
// 				}
// 				fmt.Println(head.Hash().Hex())
// 			case err := <-sub.Err():
// 				t.Error(err)
// 				return
// 			}
// 		}
// 	}()

// 	go func() {
// 		for {
// 			select {
// 			case <-mngr.Ctx.Done():
// 				return
// 			default:
// 				time.Sleep(time.Second)
// 				back.Commit()
// 			}
// 		}
// 	}()
// 	<-mngr.Done()
// }
