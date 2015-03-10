package death

import (
	log "github.com/cihub/seelog"
	"os"
	"os/signal"
	"sync"
)

type Death struct {
	wg         sync.WaitGroup
	sigChannel chan os.Signal
}

func NewDeath(signals ...os.Signal) (death *Death) {
	death = &Death{}
	death.sigChannel = make(chan os.Signal, 1)
	signal.Notify(death.sigChannel, signals...)
	death.wg.Add(1)
	go death.listenForSignal(death.sigChannel)
	return death
}

type Closable interface {
	Close()
}

//Wait for death
func (d *Death) WaitForDeath(closable ...Closable) {
	d.wg.Wait()
	log.Info("Shutdown started...")
	log.Debug("Closing ", len(closable), " objects")
	for _, c := range closable {
		c.Close()
	}
}

//Manage death of application by signal
func (d *Death) listenForSignal(c <-chan os.Signal) {
	defer d.wg.Done()
	for _ = range c {
		return
	}
}

//Shutdown.
func (d *Death) Close() {
	if nil != d {
		close(d.sigChannel)
	}
}
