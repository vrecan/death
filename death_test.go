package death

import (
	log "github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"syscall"
	"testing"
)

func TestDeath(t *testing.T) {
	defer log.Flush()

	Convey("Validate death happens cleanly", t, func() {
		death := NewDeath(syscall.SIGTERM)
		defer death.Close()
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
		death.WaitForDeath()

	})

	Convey("Validate death happens with other signals", t, func() {
		death := NewDeath(syscall.SIGHUP)
		closeMe := &CloseMe{}
		defer death.Close()
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		death.WaitForDeath(closeMe)
		So(closeMe.Closed, ShouldEqual, 1)
	})

	// Convey("Validate death happens cleanly with closable", t, func() {

	// })

}

type CloseMe struct {
	Closed int
}

func (c *CloseMe) Close() {
	c.Closed++
}
