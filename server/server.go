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
	return nil
}

func Constructor(ctx *node.ServiceContext) (node.Service, error) {
	return &Server{}, nil
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
