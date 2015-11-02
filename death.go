package death

import (
	LOG "github.com/cihub/seelog"
	"io"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"time"
)

type Death struct {
	wg         *sync.WaitGroup
	sigChannel chan os.Signal
	timeout    time.Duration
	log        Logger
}

var empty struct{}

//Logger interface to log.
type Logger interface {
	Error(v ...interface{}) error
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{}) error
}

type closer struct {
	C       io.Closer
	Name    string
	PKGPath string
}

//Create Death with the signals you want to die from.
func NewDeath(signals ...os.Signal) (death *Death) {
	death = &Death{timeout: 10 * time.Second,
		sigChannel: make(chan os.Signal, 1),
		wg:         &sync.WaitGroup{},
		log:        LOG.Current}
	signal.Notify(death.sigChannel, signals...)
	death.wg.Add(1)
	go death.listenForSignal(death.sigChannel)
	return death
}

//Override the time death is willing to wait for a objects to be closed.
func (d *Death) setTimeout(t time.Duration) {
	d.timeout = t
}

//setLogger Override the default logger (seelog)
//logger methods must be thread safe.
func (d *Death) setLogger(l Logger) {
	d.log = l
}

//Wait for death and then kill all items that need to die.
func (d *Death) WaitForDeath(closable ...io.Closer) {
	d.wg.Wait()
	d.log.Info("Shutdown started...")
	count := len(closable)
	d.log.Debug("Closing ", count, " objects")
	if count > 0 {
		d.closeInMass(closable...)
	}
}

//Close all the objects at once and wait forr them to finish with a channel.
func (d *Death) closeInMass(closable ...io.Closer) {

	count := len(closable)
	sentToClose := make(map[closer]struct{})
	//call close async
	doneClosers := make(chan closer, count)
	for _, c := range closable {
		elem := reflect.TypeOf(c).Elem()
		closer := closer{C: c, Name: elem.Name(), PKGPath: elem.PkgPath()}
		go d.closeObjects(closer, doneClosers)
		sentToClose[closer] = empty
	}

	//wait on channel for notifications.

	timer := time.NewTimer(d.timeout)
	for {
		select {
		case <-timer.C:
			d.log.Warn(count, " object(s) remaining but timer expired.")
			for c, _ := range sentToClose {
				d.log.Error("Failed to close: ", c.PKGPath, "/", c.Name)
			}
			return
		case closer := <-doneClosers:
			delete(sentToClose, closer)
			count--
			d.log.Debug(count, " object(s) left")
			d.log.Debug("Closers: ", sentToClose)
			if count == 0 && len(sentToClose) == 0 {
				d.log.Debug("Finished closing objects")
				return
			}
		}
	}
}

//closeObjects and return a bool when finished on a channel.
func (d *Death) closeObjects(closer closer, done chan<- closer) {
	err := closer.C.Close()
	if nil != err {
		d.log.Error(err)
	}
	done <- closer
}

//ListenForSignal Manage death of application by signal.
func (d *Death) listenForSignal(c <-chan os.Signal) {
	defer d.wg.Done()
	for range c {
		return
	}
}
