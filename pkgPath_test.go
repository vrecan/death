package death

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetPkgPath(t *testing.T) {

	Convey("Give pkgPath a ptr", t, func() {
		c := &Closer{}
		name, pkgPath := getPkgPath(c)
		So(name, ShouldEqual, "Closer")
		So(pkgPath, ShouldEqual, "github.com/vrecan/death/v3")

	})

	Convey("Give pkgPath a interface", t, func() {
		var closable Closable
		closable = Closer{}
		name, pkgPath := getPkgPath(closable)
		So(name, ShouldEqual, "Closer")
		So(pkgPath, ShouldEqual, "github.com/vrecan/death/v3")
	})

	Convey("Give pkgPath a copy", t, func() {
		c := Closer{}
		name, pkgPath := getPkgPath(c)
		So(name, ShouldEqual, "Closer")
		So(pkgPath, ShouldEqual, "github.com/vrecan/death/v3")
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
