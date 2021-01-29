package process

import (
	"os"
	"testing"

	"golang.org/x/sys/windows"
)

func TestHandler(t *testing.T) {
	var hdlr prochdlr

	pid := uint32(os.Getpid())

	err := hdlr.open(pid, windows.PROCESS_QUERY_INFORMATION)
	if err != nil {
		t.Fatalf("Could not open process handler")
	}

	hdlr.close()
}
