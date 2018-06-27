// +build linux bsd darwin

package death

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
)

type Unhashable map[string]interface{}

func (u Unhashable) Close() error {
	return nil
}

func TestDeath(t *testing.T) {
	defer seelog.Flush()

	Convey("Validate death handles unhashable types", t, func() {
		u := make(Unhashable)
		death := NewDeath(syscall.SIGTERM)
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
		err := death.WaitForDeath(u)
		So(err, ShouldBeNil)
	})

	Convey("Validate death happens cleanly", t, func() {
		death := NewDeath(syscall.SIGTERM)
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
		err := death.WaitForDeath()
		So(err, ShouldBeNil)
	})

	Convey("Validate death happens with other signals", t, func() {
		death := NewDeath(syscall.SIGHUP)
		closeMe := &CloseMe{}
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		err := death.WaitForDeath(closeMe)
		So(err, ShouldBeNil)
		So(closeMe.Closed, ShouldEqual, 1)
	})

	Convey("Validate death happens with a manual call", t, func() {
		death := NewDeath(syscall.SIGHUP)
		closeMe := &CloseMe{}
		death.FallOnSword()
		err := death.WaitForDeath(closeMe)
		So(err, ShouldBeNil)
		So(closeMe.Closed, ShouldEqual, 1)
	})

	Convey("Validate multiple sword falls do not block", t, func() {
		death := NewDeath(syscall.SIGHUP)
		closeMe := &CloseMe{}
		death.FallOnSword()
		death.FallOnSword()
		err := death.WaitForDeath(closeMe)
		So(err, ShouldBeNil)
		So(closeMe.Closed, ShouldEqual, 1)
	})

	Convey("Validate multiple sword falls do not block even after we have exited waitForDeath", t, func() {
		death := NewDeath(syscall.SIGHUP)
		closeMe := &CloseMe{}
		death.FallOnSword()
		death.FallOnSword()
		err := death.WaitForDeath(closeMe)
		So(err, ShouldBeNil)
		death.FallOnSword()
		death.FallOnSword()
		So(closeMe.Closed, ShouldEqual, 1)
	})

	Convey("Validate death gives up after timeout", t, func() {
		death := NewDeath(syscall.SIGHUP)
		death.SetTimeout(10 * time.Millisecond)
		neverClose := &neverClose{}
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		err := death.WaitForDeath(neverClose)
		So(err, ShouldNotBeNil)
	})

	Convey("Validate death uses new logger", t, func() {
		death := NewDeath(syscall.SIGHUP)
		closeMe := &CloseMe{}
		logger := &MockLogger{}
		death.SetLogger(logger)

		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		err := death.WaitForDeath(closeMe)
		So(err, ShouldBeNil)
		So(closeMe.Closed, ShouldEqual, 1)
		So(logger.Logs, ShouldNotBeEmpty)
	})

	Convey("Close multiple things with one that fails the timer", t, func() {
		death := NewDeath(syscall.SIGHUP)
		death.SetTimeout(10 * time.Millisecond)
		neverClose := &neverClose{}
		closeMe := &CloseMe{}
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		err := death.WaitForDeath(neverClose, closeMe)
		So(err, ShouldNotBeNil)
		So(closeMe.Closed, ShouldEqual, 1)
	})

	Convey("Close with anonymous function", t, func() {
		death := NewDeath(syscall.SIGHUP)
		death.SetTimeout(5 * time.Millisecond)
		closeMe := &CloseMe{}
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		death.WaitForDeathWithFunc(func() {
			closeMe.Close()
			So(true, ShouldBeTrue)
		})
		So(closeMe.Closed, ShouldEqual, 1)
	})

	Convey("Validate death errors when closer returns error", t, func() {
		death := NewDeath(syscall.SIGHUP)
		killMe := &KillMe{}
		death.FallOnSword()
		err := death.WaitForDeath(killMe)
		So(err, ShouldNotBeNil)
	})

}

func TestGenerateErrString(t *testing.T) {
	Convey("Generate for multiple errors", t, func() {
		closers := []closer{
			closer{
				Err:     fmt.Errorf("error 1"),
				Name:    "foo",
				PKGPath: "my/pkg",
			},
			closer{
				Err:     fmt.Errorf("error 2"),
				Name:    "bar",
				PKGPath: "my/otherpkg",
			},
		}

		expected := "my/pkg/foo: error 1, my/otherpkg/bar: error 2"
		actual := generateErrString(closers)

		So(actual, ShouldEqual, expected)
	})

	Convey("Generate for single error", t, func() {
		closers := []closer{
			closer{
				Err:     fmt.Errorf("error 1"),
				Name:    "foo",
				PKGPath: "my/pkg",
			},
		}

		expected := "my/pkg/foo: error 1"
		actual := generateErrString(closers)

		So(actual, ShouldEqual, expected)
	})
}

type MockLogger struct {
	Logs []interface{}
}

func (l *MockLogger) Info(v ...interface{}) {
	for _, log := range v {
		l.Logs = append(l.Logs, log)
	}
}

func (l *MockLogger) Debug(v ...interface{}) {
	for _, log := range v {
		l.Logs = append(l.Logs, log)
	}
}

func (l *MockLogger) Error(v ...interface{}) {
	for _, log := range v {
		l.Logs = append(l.Logs, log)
	}
}

func (l *MockLogger) Warn(v ...interface{}) {
	for _, log := range v {
		l.Logs = append(l.Logs, log)
	}
}

type neverClose struct {
}

func (n *neverClose) Close() error {
	time.Sleep(2 * time.Minute)
	return nil
}

// CloseMe returns nil from close
type CloseMe struct {
	Closed int
}

func (c *CloseMe) Close() error {
	c.Closed++
	return nil
}

// KillMe returns an error from close
type KillMe struct{}

func (c *KillMe) Close() error {
	return fmt.Errorf("error from closer")
}
