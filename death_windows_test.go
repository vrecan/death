package death

import (
	"bytes"
	log "github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestDeath(t *testing.T) {

	defer log.Flush()

	Convey("Validate death happens cleanly on windows with ctrl-c event", t, func() {
		// create source file
		const source = `
package main
import (
	"syscall"
	"github.com/vrecan/death"
)
func main() {
	death := death.NewDeath(syscall.SIGINT)
	death.WaitForDeath()
}
`
		tmp, err := ioutil.TempDir("", "TestCtrlBreak")
		if err != nil {
			t.Fatal("TempDir failed: ", err)
		}
		defer os.RemoveAll(tmp)

		// write ctrlbreak.go
		name := filepath.Join(tmp, "ctlbreak")
		src := name + ".go"
		f, err := os.Create(src)
		if err != nil {
			t.Fatalf("Failed to create %v: %v", src, err)
		}
		defer f.Close()
		f.Write([]byte(source))

		// compile it
		exe := name + ".exe"
		defer os.Remove(exe)
		o, err := exec.Command("go", "build", "-o", exe, src).CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to compile: %v\n%v", err, string(o))
		}

		// run it
		cmd := exec.Command(exe)
		var b bytes.Buffer
		cmd.Stdout = &b
		cmd.Stderr = &b
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}
		err = cmd.Start()
		if err != nil {
			t.Fatalf("Start failed: %v", err)
		}
		go func() {
			time.Sleep(1 * time.Second)
			sendCtrlBreak(t, cmd.Process.Pid)
		}()
		err = cmd.Wait()
		if err != nil {
			t.Fatalf("Program exited with error: %v\n%v", err, string(b.Bytes()))
		}
	})
}

func sendCtrlBreak(t *testing.T, pid int) {
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		t.Fatalf("LoadDLL: %v\n", e)
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		t.Fatalf("FindProc: %v\n", e)
	}
	r, _, e := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		t.Fatalf("GenerateConsoleCtrlEvent: %v\n", e)
	}
}
