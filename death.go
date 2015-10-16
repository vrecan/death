package death

import (
	log "github.com/cihub/seelog"
	"os"
	"os/signal"
	"sync"
	"time"
	"io"
)

type Death struct {
	wg         sync.WaitGroup
	sigChannel chan os.Signal
	timeout    time.Duration
}

//Create Death with the signals you want to die from.
func NewDeath(signals ...os.Signal) (death *Death) {
	death = &Death{timeout: 10 * time.Second}
	death.sigChannel = make(chan os.Signal, 1)
	signal.Notify(death.sigChannel, signals...)
	death.wg.Add(1)
	go death.listenForSignal(death.sigChannel)
	return death
}

//Override the time death is willing to wait for a objects to be closed.
func (d *Death) setTimeout(t time.Duration) {
	d.timeout = t
}

//Wait for death and then kill all items that need to die.
func (d *Death) WaitForDeath(closable ...io.Closer) {
	d.wg.Wait()
	log.Info("Shutdown started...")
	count := len(closable)
	log.Debug("Closing ", count, " objects")
	if count > 0 {
		d.closeInMass(closable...)
	}
}

//Close all the objects at once and wait forr them to finish with a channel.
func (d *Death) closeInMass(closable ...io.Closer) {
	count := len(closable)
	//call close async
	done := make(chan bool, count)
	for _, c := range closable {
		go d.closeObjects(c, done)
	}

	//wait on channel for notifications.

	timer := time.NewTimer(d.timeout)
	for {
		select {
		case <-timer.C:
			log.Warn(count, " object(s) remaining but timer expired.")
			return
		case <-done:
			count--
			log.Debug(count, " object(s) left")
			if count == 0 {
				log.Debug("Finished closing objects")
				return
			}
		}
	}
}

//Close objects and return a bool when finished on a channel.
func (d *Death) closeObjects(c io.Closer, done chan<- bool) {
	c.Close()
	done <- true
}

//Manage death of application by signal.
func (d *Death) listenForSignal(c <-chan os.Signal) {
	defer d.wg.Done()
	for range c {
		return
	}
}
