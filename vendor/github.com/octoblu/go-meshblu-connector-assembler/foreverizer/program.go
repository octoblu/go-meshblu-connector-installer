package foreverizer

import (
	"os"

	"github.com/kardianos/service"
)

// Program inteface that is dumb and not real
type Program struct {
}

// Start service but not really
func (prg *Program) Start(srv service.Service) error {
	os.Exit(0)
	return nil
}

// Stop service but not really
func (prg *Program) Stop(srv service.Service) error {
	return nil
}
