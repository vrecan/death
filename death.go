package death

//Manage the death of your application.

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Death manages the death of your application.
type Death struct {
	wg          *sync.WaitGroup
	sigChannel  chan os.Signal
	callChannel chan struct{}
	timeout     time.Duration
	log         Logger
}

// closer is a wrapper to the struct we are going to close with metadata
// to help with debuging close.
type closer struct {
	Index   int
	C       io.Closer
	Name    string
	PKGPath string
	Err     error
}

// NewDeath Create Death with the signals you want to die from.
func NewDeath(signals ...os.Signal) (death *Death) {
	death = &Death{timeout: 10 * time.Second,
		sigChannel:  make(chan os.Signal, 1),
		callChannel: make(chan struct{}, 1),
		wg:          &sync.WaitGroup{},
		log:         DefaultLogger()}
	signal.Notify(death.sigChannel, signals...)
	death.wg.Add(1)
	go death.listenForSignal()
	return death
}

// SetTimeout Overrides the time death is willing to wait for a objects to be closed.
func (d *Death) SetTimeout(t time.Duration) *Death {
	d.timeout = t
	return d
}

// SetLogger Overrides the default logger (seelog)
func (d *Death) SetLogger(l Logger) *Death {
	d.log = l
	return d
}

// WaitForDeath wait for signal and then kill all items that need to die. If they fail to
// die when instructed we return an error
func (d *Death) WaitForDeath(closable ...io.Closer) (err error) {
	d.wg.Wait()
	d.log.Info("Shutdown started...")
	count := len(closable)
	d.log.Debug("Closing ", count, " objects")
	if count > 0 {
		return d.closeInMass(closable...)
	}
	return nil
}

// WaitForDeathWithFunc allows you to have a single function get called when it's time to
// kill your application.
func (d *Death) WaitForDeathWithFunc(f func()) {
	d.wg.Wait()
	d.log.Info("Shutdown started...")
	f()
}

// getPkgPath for an io closer.
func getPkgPath(c io.Closer) (name string, pkgPath string) {
	t := reflect.TypeOf(c)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name(), t.PkgPath()
}

// closeInMass Close all the objects at once and wait for them to finish with a channel. Return an
// error if you fail to close all the objects
func (d *Death) closeInMass(closable ...io.Closer) (err error) {

	count := len(closable)
	sentToClose := make(map[int]closer)
	//call close async
	doneClosers := make(chan closer, count)
	for i, c := range closable {
		name, pkgPath := getPkgPath(c)
		closer := closer{Index: i, C: c, Name: name, PKGPath: pkgPath}
		go d.closeObjects(closer, doneClosers)
		sentToClose[i] = closer
	}

	// wait on channel for notifications.
	timer := time.NewTimer(d.timeout)
	failedClosers := []closer{}
	for {
		select {
		case <-timer.C:
			s := "failed to close: "
			pkgs := []string{}
			for _, c := range sentToClose {
				pkgs = append(pkgs, fmt.Sprintf("%s/%s", c.PKGPath, c.Name))
				d.log.Error("Failed to close: ", c.PKGPath, "/", c.Name)
			}
			return fmt.Errorf("%s", fmt.Sprintf("%s %s", s, strings.Join(pkgs, ", ")))
		case closer := <-doneClosers:
			delete(sentToClose, closer.Index)
			count--
			if closer.Err != nil {
				failedClosers = append(failedClosers, closer)
			}

			d.log.Debug(count, " object(s) left")
			if count != 0 || len(sentToClose) != 0 {
				continue
			}

			if len(failedClosers) != 0 {
				errString := generateErrString(failedClosers)
				return fmt.Errorf("errors from closers: %s", errString)
			}

			return nil
		}
	}
}

// closeObjects and return a bool when finished on a channel.
func (d *Death) closeObjects(closer closer, done chan<- closer) {
	err := closer.C.Close()
	if err != nil {
		d.log.Error(err)
		closer.Err = err
	}
	done <- closer
}

// FallOnSword manually initiates the death process.
func (d *Death) FallOnSword() {
	select {
	case d.callChannel <- struct{}{}:
	default:
	}
}

// ListenForSignal Manage death of application by signal.
func (d *Death) listenForSignal() {
	defer d.wg.Done()
	for {
		select {
		case <-d.sigChannel:
			return
		case <-d.callChannel:
			return
		}
	}
}

// generateErrString generates a string containing a list of tuples of pkgname to error message
func generateErrString(failedClosers []closer) (errString string) {
	for i, fc := range failedClosers {
		if i == 0 {
			errString = fmt.Sprintf("%s/%s: %s", fc.PKGPath, fc.Name, fc.Err)
			continue
		}
		errString = fmt.Sprintf("%s, %s/%s: %s", errString, fc.PKGPath, fc.Name, fc.Err)
	}

	return errString
}
