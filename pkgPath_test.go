package death

import (
	log "github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetPkgPath(t *testing.T) {
	defer log.Flush()

	Convey("Give pkgPath a ptr", t, func() {
		c := &Closer{}
		name, pkgPath := GetPkgPath(c)
		So(name, ShouldEqual, "Closer")
		So(pkgPath, ShouldEqual, "github.com/vrecan/death")

	})

	Convey("Give pkgPath a interface", t, func() {
		var closable Closable
		closable = Closer{}
		name, pkgPath := GetPkgPath(closable)
		So(name, ShouldEqual, "Closer")
		So(pkgPath, ShouldEqual, "github.com/vrecan/death")
	})

	Convey("Give pkgPath a copy", t, func() {
		c := Closer{}
		name, pkgPath := GetPkgPath(c)
		So(name, ShouldEqual, "Closer")
		So(pkgPath, ShouldEqual, "github.com/vrecan/death")
	})
}

type Closable interface {
	Close() error
}
type Closer struct {
}

func (c Closer) Close() error {
	return nil
}
