package death

import (
	// "fmt"
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
	death.sigChannel = make(chan os.Signal, 10)
	signal.Notify(death.sigChannel, signals...)
	return death
}

type Closable interface {
	Close()
}

//Wait for death
func (d *Death) WaitForDeath(closable []Closable) {
	d.wg.Wait()
	for _, c := range closable {
		c.Close()
	}
}

//Manage death of application by signal
func (d *Death) listenForSignal(c <-chan os.Signal, death chan string) {
	d.wg.Add(1)
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
