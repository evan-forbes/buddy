package server

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/evan-forbes/buddy/sim"
	"github.com/pkg/errors"
)

// I need to plug the simulated backend interface into the interfaces that it
// implements, and then plug those things into the node rpc example below

// todo: find out all apis included in eth/backend.go ethereum.APIs() and implement/test one by connected it below
// 		 - make a client that connects to the local host or whatever ip is being used.

type Server struct {
	Backend *sim.SimulatedBackend
}

func (s *Server) Protocols() []p2p.Protocol { return nil }

func (s *Server) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.Backend, false),
			Public:    true,
		},
	}
}

func (s *Server) Start(p *p2p.Server) error {
	return nil
}

func (s *Server) Stop() error {
	s.Backend.Close()
	s.Backend.ChainDb().Close()
	return nil
}

func Constructor(ctx *node.ServiceContext) (node.Service, error) {
	accnts := sim.NewAccounts("owner", "Alice", "Bob", "Celine", "Doug", "Erin", "Frank")
	genesisAlloc := accnts.Genesis()
	s := &Server{
		Backend: sim.NewSimulatedBackend(genesisAlloc, uint64(4712388)),
	}
	return s, nil
}

func Boot(config *node.Config) error {
	stack, err := node.New(config)
	if err != nil {
		return errors.Wrapf(err, "Could not boot server %s", config.Name)
	}
	defer stack.Close()
	if err = stack.Register(Constructor); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	// Boot up the entire protocol stack, do a restart and terminate
	if err = stack.Start(); err != nil {
		log.Fatalf("Failed to start the protocol stack: %v", err)
	}
	time.Sleep(time.Second * 8)
	if err = stack.Stop(); err != nil {
		log.Fatalf("Failed to stop the protocol stack: %v", err)
	}
	return nil
}
