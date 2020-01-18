package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Manager struct {
	Cancel   context.CancelFunc
	Ctx      context.Context
	WG       *sync.WaitGroup
	Interupt chan os.Signal
	// LogFile  *json.Encoder
}

// NewManager return a pointer to a setup manager.
func NewManager(superctx context.Context, wg *sync.WaitGroup) *Manager {
	ctx, cancel := context.WithCancel(superctx)
	// f, err := os.OpenFile("file.Log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	fmt.Println("problem creating log file")
	// 	cancel()
	// }
	mnger := &Manager{
		Cancel:   cancel,
		Ctx:      ctx,
		WG:       wg,
		Interupt: make(chan os.Signal),
		// LogFile:  json.NewEncoder(f),
	}
	return mnger
}

func (m *Manager) Listen() {
	signal.Notify(m.Interupt, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-m.Interupt:
			m.Cancel()
		case <-m.Ctx.Done():
			m.WG.Wait()
			os.Exit(0)
		}
	}
}

type Report struct {
	Data interface{}
	error
}

func (r *Report) Error() string {
	return r.Error()
}

func (m *Manager) CancelOn(errc <-chan error) {
	for err := range errc {
		if err != nil {
			log.Println(err, "\nsending cancels")
			m.Cancel()
		}
	}
}

// LogErr writes any and every log directly to file
func (m *Manager) LogErr(errc <-chan error) {
	for err := range errc {
		if err != nil {
			log.Println(err)
		}
	}
}
