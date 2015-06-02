// +build linux bsd darwin

package death

import (
	log "github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestDeath(t *testing.T) {
	defer log.Flush()

	Convey("Validate death happens cleanly", t, func() {
		death := NewDeath(syscall.SIGTERM)
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
		death.WaitForDeath()

	})

	Convey("Validate death happens with other signals", t, func() {
		death := NewDeath(syscall.SIGHUP)
		closeMe := &CloseMe{}
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		death.WaitForDeath(closeMe)
		So(closeMe.Closed, ShouldEqual, 1)
	})

	Convey("Validate death gives up after timeout", t, func() {
		death := NewDeath(syscall.SIGHUP)
		death.setTimeout(10 * time.Millisecond)
		neverClose := &neverClose{}
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		death.WaitForDeath(neverClose)

	})

}

type neverClose struct {
}

func (n *neverClose) Close() {
	time.Sleep(2 * time.Minute)
}

type CloseMe struct {
	Closed int
}

func (c *CloseMe) Close() {
	c.Closed++
}
