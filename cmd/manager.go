package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Manager is a simple context and waitgroup controller
type Manager struct {
	Cancel   context.CancelFunc
	Ctx      context.Context
	WG       *sync.WaitGroup
	Interupt chan os.Signal
	DoneChan chan struct{}
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
	if wg == nil {
		wg = &sync.WaitGroup{}
	}
	mnger := &Manager{
		Cancel:   cancel,
		Ctx:      ctx,
		WG:       wg,
		Interupt: make(chan os.Signal),
		DoneChan: make(chan struct{}, 1),
		// LogFile:  json.NewEncoder(f),
	}
	return mnger
}

// Listen watches for interuption via ctrl + C
func (m *Manager) Listen() {
	signal.Notify(m.Interupt, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-m.Interupt:
			m.Cancel()
		case <-m.Ctx.Done():
			m.WG.Wait()
			m.DoneChan <- struct{}{}
			os.Exit(0)
		}
	}
}

// Done wraps around the Manager's context's Done method to block until
//
func (m *Manager) Done() <-chan struct{} {
	return m.DoneChan
}

type Action int

const (
	SHUTDOWN Action = iota
	REBOOT
)

// Report is a simple error wrapper for logging
type Report struct {
	Action Action
	Data   interface{}
	error
}

// Error fulffils the error interface
func (r *Report) Error() string {
	return r.Error()
}

// HandleReports can manages error handling for reports
func (m *Manager) HandleReports(errc <-chan error) {
	for err := range errc {
		rep, ok := err.(*Report)
		if !ok {
			log.Println(err)
			continue
		}
		switch rep.Action {
		case SHUTDOWN:
			m.Cancel()
			return
		case REBOOT:
			m.Cancel()
			// m.BootFunc()
		}
	}
}

// CancelOn funs the context's cancel function upon recieving any error
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

// Config contains settings data
type Config struct {
	ArangoPass string `json:"arango_pass"`
	ArangoUser string `json:"arango_user"`
	Workers    int    `json:"workers"`
	RPCAddress string `json:"rpc_address"`
	DBAddress  string `json:"db_address"`
	Col        string `json:"col"`
	DB         string `json:"db"`
}

// LoadConfig reads and parses the configuration file
func LoadConfig() (Config, error) {
	var out Config
	jsonFile, err := ioutil.ReadFile("/home/evan/.creds/mesh.json")
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(jsonFile, &out)
	if err != nil {
		return out, err
	}
	return out, nil
}

// LoadHexAddressMap reads and parses the configuration file
func LoadHexAddressMap() (map[string]string, error) {
	out := make(map[string]string)
	jsonFile, err := ioutil.ReadFile("/home/evan/.creds/HexAddressMap.json")
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(jsonFile, &out)
	if err != nil {
		return out, err
	}
	return out, nil
}
