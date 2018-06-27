package death

import (
	"testing"

	log "github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetPkgPath(t *testing.T) {
	defer log.Flush()

	Convey("Give pkgPath a ptr", t, func() {
		c := &Closer{}
		name, pkgPath := getPkgPath(c)
		So(name, ShouldEqual, "Closer")
		So(pkgPath, ShouldEqual, "github.com/vrecan/death")

	})

	Convey("Give pkgPath a interface", t, func() {
		var closable Closable
		closable = Closer{}
		name, pkgPath := getPkgPath(closable)
		So(name, ShouldEqual, "Closer")
		So(pkgPath, ShouldEqual, "github.com/vrecan/death")
	})

	Convey("Give pkgPath a copy", t, func() {
		c := Closer{}
		name, pkgPath := getPkgPath(c)
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
