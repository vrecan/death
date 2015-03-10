package death

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

type SIGNAL string

const (
	TERMINATED SIGNAL = "terminated"
	INTERRUPT  SIGNAL = "interrupt"
	UNKNOWN    SIGNAL = "unknown"
)

type Death struct {
	wg           sync.WaitGroup
	sigChannel   chan os.Signal
	deathChannel chan string
	deathSignals map[SIGNAL]bool
}

func NewDeath() (death *Death) {
	death = &Death{}
	death.sigChannel = make(chan os.Signal, 10)
	death.deathChannel = make(chan string, 10)
	signal.Notify(death.sigChannel)
	go listenForSignal(death.sigChannel, death.deathChannel)
	return death
}

func (d *Death) SetDeathSignals(signals ...SIGNAL) {
	for _, sig := range signals {
		d.deathSignals[sig] = true
	}
}

func (d *Death) RemoveDeathSignals(signals ...SIGNAL) {
	for _, sig := range signals {
		delete(d.deathSignals, sig)
	}
}

func getSignal(sig string) SIGNAL {
	switch sig {
	case string(TERMINATED):
		{
			return TERMINATED
		}
	case string(INTERRUPT):
		{
			return INTERRUPT
		}
	default:
		{
			fmt.Println("Unknown signal received: ", sig)
			return UNKNOWN
		}
	}
}

func (d *Death) listenForDeath() {
	for sig := range d.deathChannel {
		signal := getSignal(sig)
		_, ok := d.deathSignals[signal]
		if ok {
			return
		}
	}
}

//Wait for death
func (d *Death) WaitForDeath() {
	d.listenForDeath()
}

//Manage death of application by signal
func listenForSignal(c <-chan os.Signal, death chan string) {
	for sig := range c {
		death <- sig.String()
	}
}

func (d *Death) Close() {
	if nil != d {
		close(d.sigChannel)
		close(d.deathChannel)
	}
}
